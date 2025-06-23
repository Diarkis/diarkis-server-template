// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis-server-template/examples/dive/user-online-status/cmds"
	"github.com/Diarkis/diarkis/diarkisexec"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Dive: &diarkisexec.Options{ConfigPath: "configs/shared/dive.json", ExposeCommands: true},
		// We configure the room module in order to test the online status
		// when a user is in a room.
		Room: &diarkisexec.Options{ExposeCommands: true},
	})

	cmds.SetupUDP()

	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")

	diarkisexec.StartDiarkis()
}
