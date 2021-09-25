package roomcmds

import (
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/roomSupport"
)

func Expose() {
	room.ExposeCommands()
	roomSupport.ExposeCommands()
	// we broadcast user ID to room members when the client leaves unexpectedly
	room.SetOnDiscardCustomMessage(onDiscardCustomMessage)
}

func onDiscardCustomMessage(roomID string, userID string) []byte {
	return []byte(userID)
}
