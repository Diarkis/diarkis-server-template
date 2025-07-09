// © 2019-2024 Diarkis Inc. All rights reserved.

package customcmds

import (
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"

	pufferDgs "github.com/Diarkis/diarkis-server-template/examples/csar/dgs/puffer/go/dgs"
)

var logger = log.New("CUSTOM")

func Expose() {
	// defined in dgs.go
	diarkisexec.SetServerCommandHandler(pufferDgs.InstanceCreateRequestVer, pufferDgs.InstanceCreateRequestCmd, dgsCreateSessionFromRoom)
	diarkisexec.SetServerCommandHandler(pufferDgs.InstanceBackfillRequestVer, pufferDgs.InstanceBackfillRequestCmd, dgsSessionBackfillFromRoom)
}
