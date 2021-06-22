package groupcmds

import (
	"fmt"
	"github.com/Diarkis/diarkis/group"
	"github.com/Diarkis/diarkis/groupSupport"
)

func Expose(rootpath string) {
	// rootpath is defined in cmds/main.go
	group.Setup(fmt.Sprintf("%s/configs/shared/group.json", rootpath))
	group.ExposeCommands()
	groupSupport.ExposeCommands()
}
