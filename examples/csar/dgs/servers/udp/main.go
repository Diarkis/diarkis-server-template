// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis-server-template/examples/csar/dgs/cmds"
	"github.com/Diarkis/diarkis/diarkisexec"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		DGS:  &diarkisexec.Options{ConfigPath: "configs/shared/dgs.json"},
		Room: &diarkisexec.Options{ConfigPath: "configs/shared/room.json", ExposeCommands: true},
		P2P:  &diarkisexec.Options{ExposeCommands: true},
		Dive: &diarkisexec.Options{ConfigPath: "configs/shared/dive.json", ExposeCommands: true},
	})

	cmds.SetupUDP()

	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")

	diarkisexec.StartDiarkis()
}
