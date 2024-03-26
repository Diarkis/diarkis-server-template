package dmcmds

import (
	"fmt"

	"github.com/Diarkis/diarkis/dm"
	"github.com/Diarkis/diarkis/user"
)

func Expose(rootpath string) {
	// rootpath is defined in cmds/main.go
	dm.Setup(fmt.Sprintf("%s/configs/shared/dm.json", rootpath))

	// If you do not need to send a disconnect message automatically, comment out the following line
	dm.SetOnUserDisconnect(func(disconnectedUser *user.User, peerUserID string) []byte {
		return []byte(fmt.Sprintf("User %s is gone", disconnectedUser.ID))
	})

	dm.ExposeCommands()
}
