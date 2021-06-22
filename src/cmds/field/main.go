package fieldcmds

import (
	"fmt"
	"github.com/Diarkis/diarkis/field"
)

func Expose(rootpath string) {
	// rootpath is defined in cmds/main.go
	field.Setup(fmt.Sprintf("%s/configs/shared/field.json", rootpath))
	field.ExposeCommands()
}
