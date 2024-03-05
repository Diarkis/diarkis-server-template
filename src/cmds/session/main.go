package sessioncmds

import (
	"github.com/Diarkis/diarkis/session"
)

func Expose() {
	session.Setup()
	session.ExposeCommands()
}
