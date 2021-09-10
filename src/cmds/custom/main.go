package customcmds

import (
	"encoding/binary"
	"github.com/Diarkis/diarkis/log"
	dpayload "{0}/lib/payload"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

const customVer = 2
const helloCmdID = 10
const pushCmdID = 11
const matchmakerAdd = 100
const matchmakerRm = 101
const matchmakerSearch = 102
const matchmakerComplete = 103 // sent when room is full

var logger = log.New("CUSTOM")

func Expose() {
	server.HandleCommand(customVer, helloCmdID, helloCmd)
	server.HandleCommand(customVer, helloCmdID, afterHelloCmd)
	server.HandleCommand(customVer, pushCmdID, pushCmd)
	server.HandleCommand(customVer, matchmakerAdd, addToMatchMaker)
	server.HandleCommand(customVer, matchmakerSearch, searchMatchMaker)
}

func helloCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Hello command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	// if this is executed as UDP, reliable = true means sending the packet as RRUDP
	reliable := true
	// we send a response back to the client with the byte array sent from the client
	userData.ServerRespond(payload, ver, cmd, server.Ok, reliable)
	// move on to the next command handler if there is any
	next(nil)
}

func afterHelloCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("This is executed after Hello command has been handled")
	next(nil)
}

func pushCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Push command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	// if this is executed as UDP, reliable = true means sending the packet as RRUDP
	reliable := true
	// we send a push packet to the client that sent the data to this command
	userData.ServerPush(ver, cmd, payload, reliable)
	// move on to the next command handler if there is any
	next(nil)
}

// add to matching and create a room
func addToMatchMaker(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	mmAdd := dpayload.UnpackMMAdd(payload)
	if mmAdd == nil {
		userData.ServerRespond([]byte("Invalid payload"), ver, cmd, server.Bad, true)
		return
	}
	// raise Room.OnCreate event
	maxMembers := 10
	allowEmpty := false
	join := true
	ttl := int64(120)
	interval := int64(0)
	roomID, err := room.NewRoom(userData, maxMembers, allowEmpty, join, ttl, interval)
	if err != nil {
		userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Bad, true)
		next(err)
		return
	}
	created := room.GetCreatedTime(roomID)
	timeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(timeBytes, uint32(created))
	roomIDBytes := []byte(roomID)
	bytes := append(timeBytes, roomIDBytes...)
	userData.ServerRespond(bytes, util.CmdBuiltInVer, util.CmdCreateRoom, server.Ok, true)
	metadata := make(map[string]interface{})
	metadata["serialized"] = mmAdd.Metadata
	metadata["roomID"] = roomID
	metadata["maxMembers"] = maxMembers
	matching.Add(mmAdd.ID, mmAdd.UID, mmAdd.Props, metadata, mmAdd.TTL, 2)
	// response for matchmaking
	userData.ServerRespond([]byte("OK"), ver, cmd, server.Ok, true)
	next(nil)
}

func searchMatchMaker(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	mmSearch := dpayload.UnpackMMSearch(payload)
	if mmSearch == nil {
		userData.ServerRespond([]byte("Invalid payload"), ver, cmd, server.Bad, true)
		return
	}
	howmany := 10
	matching.Search(mmSearch.IDs, mmSearch.Props, howmany, func(err error, results []interface{}) {
		if err != nil {
			userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Bad, true)
			return
		}
		logger.Sys("MatchMaker search results %v", results)
		if len(results) == 0 {
			userData.ServerRespond([]byte("Mathing not found"), ver, cmd, server.Bad, true)
			return
		}
		list := make([]map[string]interface{}, len(results))
		operations := make([]func(func(error)), len(results))
		index := 0
		joinedRoomID := ""
		ct := int64(0)
		isRoomFull := false
		done := func(err error) {
			if err != nil {
				userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Bad, true)
			}
			if joinedRoomID == "" {
				userData.ServerRespond([]byte("No room found to join"), ver, cmd, server.Bad, true)
				return
			}
			timeBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(timeBytes, uint32(ct))
			roomIDBytes := []byte(joinedRoomID)
			bytes := append(timeBytes, roomIDBytes...)
			userData.ServerRespond(bytes, util.CmdBuiltInVer, util.CmdRandJoinRoom, server.Ok, true)
			userData.ServerRespond([]byte(joinedRoomID), ver, cmd, server.Ok, true)
			if isRoomFull {
				userData.ServerPush(ver, cmd, []byte(joinedRoomID), true)
			}
			next(err)
		}
		join := func(moveon func(error)) {
			item := list[index]
			roomID := item["roomID"].(string)
			logger.Sys("Try to join room %v", roomID)
			room.Join(roomID, userData, util.CmdBuiltInVer, util.CmdJoinRoom, []byte(userData.ID),
			func(err error, memberIDs []string, ownerID string, createdTime int64, props map[string]interface{}) {
				if err != nil {
					// try the next room
					index++
					moveon(nil)
					return
				}
				logger.Sys("Successfully joined a from by MatchMaker %v", roomID)
				joinedRoomID = roomID
				ct = createdTime
				isRoomFull = len(memberIDs) == item["maxMembers"].(int)
				// joined successfully
				done(nil)
			})
		}
		for i, v := range results {
			list[i] = v.(map[string]interface{})
			operations[i] = join
		}
		util.Waterfall(operations, done)
	})
}
