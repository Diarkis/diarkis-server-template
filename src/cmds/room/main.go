package roomcmds

import (
	"encoding/binary"
	"{0}/lib/meshCmds"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/roomSupport"
	"github.com/Diarkis/diarkis/user"
	"strconv"
)

const ver uint8 = 3
const onRoomOwnerChangeCmd uint16 = 100

var logger = log.New("room")

func Expose() {
	// we broadcast user ID to room members when the client leaves unexpectedly
	// The message will be propagated by OnMemberLeave even on the client
	room.SetOnDiscardCustomMessage(onDiscardCustomMessage)
	// Sends a custom broadcast to all members when the room owner changes
	room.SetOnRoomOwnerChange(onRoomOwnerChange)
	// When a user joins a room, the user receives the room owner ID by message
	room.AfterJoinRoomCmd(afterJoinRoom)
	// expose commands
	room.ExposeCommands()
	roomSupport.ExposeCommands()
	// set up custom mesh commands
	meshCmds.Setup()
}

func onDiscardCustomMessage(roomID string, userID string, sid string) []byte {
	logger.Debug("OnDiscardCustomMessage roomID:%v userID:%v sid:%v", roomID, userID, sid)
	// UE4 sample client uses uin64 as the data type for userID
	conv, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return []byte(userID)
	}
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(conv))
	logger.Debug("OnDiscardCustomMessage message user ID %v bytes %v", conv, bytes)
	return bytes
}

func onRoomOwnerChange(params interface{}) {
	data := params.(map[string]string)
	roomID := data["roomID"]
	ownerID := data["ownerID"]
	syncRoomOwnerID(roomID, ownerID)
}

func afterJoinRoom(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		// user is not in the room
		next(nil)
		return
	}
	ownerID := room.GetRoomOwnerID(roomID)
	if ownerID == "" {
		// there is no owner yet...
		next(nil)
		return
	}
	syncRoomOwnerID(roomID, ownerID)
}

func syncRoomOwnerID(roomID string, ownerID string) {
	logger.Debug("OnRoomOwnerChange => broadcast the new owner ID %s to room %s", ownerID, roomID)
	message := []byte(ownerID)
	room.Announce(roomID, room.GetMemberIDs(roomID), ver, onRoomOwnerChangeCmd, message, true)
}
