package dmcmds

import (
	"fmt"
	"github.com/Diarkis/diarkis/dm"
)

func Expose(rootpath string) {
	// rootpath is defined in cmds/main.go
	dm.Setup()
	dm.ExposeCommands()
}
