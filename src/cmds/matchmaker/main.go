package roomcmds

import (
	"fmt"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
)

var logger = log.New("matching")

func Expose(rootpath string) {
	matching.Setup(fmt.Sprintf("%s/configs/shared/matching.json", rootpath))
	matching.ExposeCommands()
}
