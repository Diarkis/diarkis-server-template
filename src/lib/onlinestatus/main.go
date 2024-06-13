package onlinestatus

import (
	"github.com/Diarkis/diarkis/dive"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/session"
	"github.com/Diarkis/diarkis/td"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"

	"github.com/Diarkis/diarkis-server-template/lib/meshCmds"
)

const returnLimit = 10
const userStatusTTL = 10 // in seconds
const storageName = "OnlineStatus"

var storage *dive.Storage
var logger = log.New("STATUS")

// UserStatusData represents a user's online status data.
type UserStatusData struct {
	UID         string
	InRoom      bool
	SessionData []UserSessionData
}

// UserSessionData represents user's session status data.
type UserSessionData struct {
	Type uint8
	ID   string
}

var userStatus = td.DefineTransportData([]td.Property{
	td.Property{Name: "UID", Type: td.String},
	td.Property{Name: "SessionData", Type: td.Bytes},
	td.Property{Name: "InRoom", Type: td.Bool},
})

var userStatusList = td.DefineTransportData([]td.Property{
	td.Property{Name: "List", Type: td.BytesArray},
})

var uidList = td.DefineTransportData([]td.Property{
	td.Property{Name: "List", Type: td.StringArray},
})

var userSessionData = td.DefineTransportData([]td.Property{
	td.Property{Name: "Types", Type: td.Uint8Array},
	td.Property{Name: "IDs", Type: td.StringArray},
})

func Setup() {
	if storage != nil {
		// Setup has already been called...
		return
	}

	logger.Debug("Setting up online status lib")

	// setup callbacks on user becoming online and offline
	// we keep updating the status on keep alive
	// once the user disconnects, keep alive stops.
	// we will let the status be deleted by TTL expiration
	user.OnNew(setUserAsOnline)

	mesh.HandleRPC(meshCmds.GetOnlineStatusListCmd, handleGetUserStatusList)
}

func getStorage() *dive.Storage {
	if storage == nil {
		// storage is defined in configs/dive.json file
		storage = dive.GetStorageByName(storageName)
	}
	return storage
}

func setUserAsOnline(userData *user.User) {
	us := userStatus.New()
	us.SetAsString("UID", userData.ID)
	err := getStorage().SetEx(userData.ID, us.Pack(), userStatusTTL)
	if err != nil {
		logger.Error("Failed to set user as online: UID:%s Error:%v", userData.ID, err.Error())
		return
	}
	logger.Debug("Set user as online: UID:%s", userData.ID)
	// set up user status updater on every keep alive
	server.OnKeepAlive(updateUserStatus)
}

func updateUserStatus(userData *user.User, next func(error)) {
	roomID := room.GetRoomID(userData)
	sessionDataList := session.GetSessionDataByUser(userData)

	inRoom := false
	if len(roomID) > 0 {
		inRoom = true
	}

	usd := userSessionData.New()
	types := make([]uint8, len(sessionDataList))
	ids := make([]string, len(sessionDataList))
	for i, sd := range sessionDataList {
		types[i] = sd.Type
		ids[i] = sd.ID
	}
	usd.SetAsUint8Array("Types", types)
	usd.SetAsStringArray("IDs", ids)

	us := userStatus.New()
	us.SetAsString("UID", userData.ID)
	us.SetAsBool("InRoom", inRoom)
	us.SetAsBytes("SessionData", usd.Pack())

	packed := us.Pack()
	// update the remote status data
	err := getStorage().SetEx(userData.ID, packed, userStatusTTL)
	if err != nil {
		logger.Error("Failed to update user status: UID:%s Error:%v", userData.ID, err.Error())
	}

	// we must call next to make sure the reset of keep alive operations move forward
	next(nil)
}

// GetUserStatus returns a td.TransportData representing user status.
//
// Returned *td.TransportData has the following properties:
//
//	Name: "Online" Type: td.Bool - GetAsBool("Online")
func GetUserStatus(uid string) (*UserStatusData, error) {
	if getStorage() == nil {
		return nil, util.NewError("Storage must be setup")
	}

	res := getStorage().Get(uid)
	data, err := res.ToBytes()

	if err != nil {
		return nil, err
	}

	return getUserStatus(data), nil
}

func getUserStatus(data []byte) *UserStatusData {
	us := userStatus.New()
	us.Unpack(data)

	userStatusData := &UserStatusData{}
	uid, _ := us.GetAsString("UID")
	inRoom, err := us.GetAsBool("InRoom")

	if err != nil {
		inRoom = false
	}

	userStatusData.UID = uid
	userStatusData.InRoom = inRoom

	b, err := us.GetAsBytes("SessionData")

	if err == nil {
		usd := userSessionData.New()
		usd.Unpack(b)
		types, _ := usd.GetAsUint8Array("Types")
		ids, _ := usd.GetAsStringArray("IDs")
		userStatusData.SessionData = make([]UserSessionData, len(ids))
		for i, id := range ids {
			userStatusData.SessionData[i] = UserSessionData{
				Type: types[i],
				ID:   id,
			}
		}
	}

	return userStatusData
}

// GetUserStatusList returns a list of td.TransportData representing user status.
//
//	[IMPORTANT] The number of user status list to retrieve at a time is limited to 10.
//
// Returned *td.TransportData has the following properties:
//
//	Name: "Online" Type: td.Bool - GetAsBool("Online")
func GetUserStatusList(uids []string) ([]*UserStatusData, error) {
	if getStorage() == nil {
		return nil, util.NewError("Storage must be setup")
	}

	// check the limit
	if len(uids) > returnLimit {
		logger.Debug("Number UIDs exceeds the limit of %v and the array is trucated", returnLimit)
		return nil, util.NewError("User status list retrieval is up to %v, but %v given", returnLimit, len(uids))
	}

	// first, we group UIDs that resolve to the same mesh address
	// we do this because every getStorage().Get calls on remote server to retrieve the online status data
	// we send the user IDs that resolve to the same server in one remote server call.
	addrs := make(map[string][]string, len(uids))
	for _, uid := range uids {
		addr := getStorage().ResolveKey(uid)
		if _, ok := addrs[addr]; !ok {
			addrs[addr] = make([]string, 0)
		}
		addrs[addr] = append(addrs[addr], uid)
	}

	ret := make([]*UserStatusData, len(uids))
	i := 0

	// now we send grouped UIDs to the servers with the resolved address
	list := uidList.New()
	for addr, uids := range addrs {
		list.SetAsStringArray("List", uids)
		resp, err := mesh.SendRPC(meshCmds.GetOnlineStatusListCmd, addr, list.Pack())

		if err != nil {
			return nil, util.StackError(util.NewError("Failed to get user status list -> UIDs:%v", uids), err)
		}

		uslist := userStatusList.New()
		uslist.Unpack(resp)
		l, _ := uslist.GetAsBytesArray("List")
		for _, b := range l {
			ret[i] = getUserStatus(b)
			i++
		}
	}

	return ret[0:i], nil
}

func handleGetUserStatusList(payload []byte, sender string) ([]byte, error) {
	list := uidList.New()
	list.Unpack(payload)
	uids, _ := list.GetAsStringArray("List")

	ret := make([][]byte, len(uids))
	i := 0

	for _, uid := range uids {
		res := getStorage().Get(uid)
		b, err := res.ToBytes()

		// we allow errors because offline user status will not be found and returns as an error.
		if err != nil {
			logger.Debug("User status get returned an error -> UID:%v Error:%v", uid, err.Error())
			continue
		}

		ret[i] = b
		i++
	}

	uslist := userStatusList.New()
	uslist.SetAsBytesArray("List", ret[0:i])

	return uslist.Pack(), nil
}
