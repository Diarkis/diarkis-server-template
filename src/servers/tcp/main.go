// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis-server-template/cmds"
	"github.com/Diarkis/diarkis/diarkisexec"
)

func main() {
	logConfigPath := "/configs/shared/log.json"
	meshConfigPath := "/configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Room:       &diarkisexec.Options{ExposeCommands: true},
		Group:      &diarkisexec.Options{ConfigPath: "/configs/shared/group.json", ExposeCommands: true},
		Dive:       &diarkisexec.Options{ConfigPath: "/configs/shared/dive.json", ExposeCommands: true},
		Field:      &diarkisexec.Options{ConfigPath: "/configs/shared/field.json", ExposeCommands: true},
		DM:         &diarkisexec.Options{ConfigPath: "/configs/shared/dm.json", ExposeCommands: true},
		MatchMaker: &diarkisexec.Options{ConfigPath: "/configs/shared/matching.json", ExposeCommands: true},
		Session:    &diarkisexec.Options{ConfigPath: "/configs/shared/session.json", ExposeCommands: true},
	})

	cmds.SetupTCP()

	diarkisexec.SetupDiarkisTCPServer("/configs/tcp/main.json")

	diarkisexec.StartDiarkis()
}
