package groupcmds

import (
	"fmt"
	"github.com/Diarkis/diarkis/group"
	"github.com/Diarkis/diarkis/groupsupport"
)

func Expose(rootpath string) {
	// rootpath is defined in cmds/main.go
	group.Setup(fmt.Sprintf("%s/configs/shared/group.json", rootpath))
	group.ExposeCommands()
	groupsupport.ExposeCommands()
}
