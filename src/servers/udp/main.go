package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis-server-template/cmds"
)

func main() {
	logConfigPath := "/configs/shared/log.json"
	meshConfigPath := "/configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Room:       &diarkisexec.Options{ ExposeCommands: true },
		P2P:        &diarkisexec.Options{ ExposeCommands: true },
		Group:      &diarkisexec.Options{ ConfigPath: "/configs/shared/group.json", ExposeCommands: true },
		Dive:       &diarkisexec.Options{ ConfigPath: "/configs/shared/dive.json", ExposeCommands: true },
		Field:      &diarkisexec.Options{ ConfigPath: "/configs/shared/field.json", ExposeCommands: true },
		DM:         &diarkisexec.Options{ ConfigPath: "/configs/shared/dm.json", ExposeCommands: true },
		MatchMaker: &diarkisexec.Options{ ConfigPath: "/configs/shared/matching.json", ExposeCommands: true },
		Session:    &diarkisexec.Options{ ConfigPath: "/configs/shared/session.json", ExposeCommands: true },
	})

	cmds.SetupUDP()

	diarkisexec.SetupDiarkisUDPServer("/configs/udp/main.json")

	diarkisexec.StartDiarkis()
}
