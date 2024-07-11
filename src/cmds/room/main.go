// © 2019-2024 Diarkis Inc. All rights reserved.

package roomcmds

import (
	"github.com/Diarkis/diarkis-server-template/lib/meshCmds"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/roomsupport"
	"github.com/Diarkis/diarkis/user"
)

const ver uint8 = 3
const onRoomOwnerChangeCmd uint16 = 100

var logger = log.New("room")

func Setup() {
	// we broadcast user ID to room members when the client leaves unexpectedly
	// The message will be propagated by OnMemberLeave even on the client
	room.SetOnDiscardCustomMessage(onDiscardCustomMessage)
	// Sends a custom broadcast to all members when the room owner changes
	room.SetOnRoomOwnerChange(onRoomOwnerChange)
	// When a new room is create, we setup on join event to sync owner ID by message
	room.AfterCreateRoomCmd(afterCreateRoom)
	// When a user performs random join and the server creates a new room, we setup on join event
	// to sync owner ID by message
	roomsupport.AfterRandomRoomCmd(afterRandomJoin)
	// set up custom mesh commands
	meshCmds.Setup()
}

func onDiscardCustomMessage(roomID string, userID string, sid string) []byte {
	logger.Debug("OnDiscardCustomMessage roomID:%v userID:%v sid:%v", roomID, userID, sid)
	return []byte(userID)
}

func onRoomOwnerChange(roomID string, ownerID string) {
	syncRoomOwnerID(roomID, ownerID)
}

func afterCreateRoom(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		// user is not in the room
		next(nil)
		return
	}
	setupOnJoinCallback(roomID)
	ownerID := room.GetRoomOwnerID(roomID)
	syncRoomOwnerID(roomID, ownerID)
}

func afterRandomJoin(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		// user is not in the room
		next(nil)
		return
	}
	// first byte is a flag to tell use either create or join
	if payload[0] != 0x00 {
		// it was join
		next(nil)
		return
	}
	setupOnJoinCallback(roomID)
}

func setupOnJoinCallback(roomID string) {
	room.SetOnJoinCompleteByID(roomID, func(rid string, ud *user.User) {
		ownerID := room.GetRoomOwnerID(roomID)
		if ownerID == "" {
			// there is no owner yet...
			return
		}
		syncRoomOwnerID(roomID, ownerID)
	})
}

func syncRoomOwnerID(roomID string, ownerID string) {
	logger.Debug("OnRoomOwnerChange => broadcast the new owner ID %s to room %s", ownerID, roomID)
	message := []byte(ownerID)
	room.Announce(roomID, room.GetMemberIDs(roomID), ver, onRoomOwnerChangeCmd, message, true)
}
