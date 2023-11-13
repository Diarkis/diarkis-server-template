package dmcmds

import (
	"fmt"
	"github.com/Diarkis/diarkis/dm"
)

func Expose(rootpath string) {
	// rootpath is defined in cmds/main.go
	dm.Setup(fmt.Sprintf("%s/configs/shared/dm.json", rootpath))
	dm.ExposeCommands()
}
