// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"

	"github.com/Diarkis/diarkis-server-template/examples/csar/dgs/cmds"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Room: &diarkisexec.Options{},
		Dive: &diarkisexec.Options{ConfigPath: "configs/shared/dive.json", ExposeCommands: true},
	})

	cmds.SetupHTTP()

	diarkisexec.SetupDiarkisHTTPServer("configs/http/main.json")

	diarkisexec.StartDiarkis()
}
