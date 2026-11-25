// © 2019-2024 Diarkis Inc. All rights reserved.

package roomcmds

import (
	"github.com/Diarkis/diarkis-server-template/lib/meshCmds"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/room"
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
	// When a new room is created (by create room command or by random join), we setup on join event
	// to sync owner ID by message. SetOnCreate is invoked exactly once per room, so the on join callback
	// is never registered more than once for the same room.
	room.SetOnCreate(setupOnJoinCallback)
	// When a new room is created by create room command, we sync owner ID to the creator
	room.AfterCreateRoomCmd(afterCreateRoom)
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
	ownerID := room.GetRoomOwnerID(roomID)
	syncRoomOwnerID(roomID, ownerID)
	next(nil)
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
