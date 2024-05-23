package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"

	"github.com/Diarkis/diarkis-server-template/cmds"
)

func main() {
	logConfigPath := "/configs/shared/log.json"
	meshConfigPath := "/configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Room:       &diarkisexec.Options{},
		Group:      &diarkisexec.Options{},
		Dive:       &diarkisexec.Options{ConfigPath: "/configs/shared/dive.json", ExposeCommands: true},
		Field:      &diarkisexec.Options{ConfigPath: "/configs/shared/field.json", ExposeCommands: true},
		DM:         &diarkisexec.Options{ConfigPath: "/configs/shared/dm.json", ExposeCommands: true},
		MatchMaker: &diarkisexec.Options{ConfigPath: "/configs/shared/matching.json", ExposeCommands: true},
	})

	cmds.SetupHTTP()

	diarkisexec.SetupDiarkisHTTPServer("/configs/http/main.json")

	diarkisexec.StartDiarkis()
}
