package roomcmds

import (
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/roomSupport"
)

var logger = log.New("room")

func Expose() {
	room.ExposeCommands()
	roomSupport.ExposeCommands()
	// we broadcast user ID to room members when the client leaves unexpectedly
	// The message will be propagated by OnMemberLeave even on the client
	room.SetOnDiscardCustomMessage(onDiscardCustomMessage)
}

func onDiscardCustomMessage(roomID string, userID string) []byte {
	logger.Debug("OnDiscardCustomMessage roomID:%v userID:%v", roomID, userID)
	return []byte(userID)
}
