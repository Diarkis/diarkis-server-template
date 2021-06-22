package roomcmds

import (
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/roomSupport"
)

func Expose() {
	room.ExposeCommands()
	roomSupport.ExposeCommands()
}
