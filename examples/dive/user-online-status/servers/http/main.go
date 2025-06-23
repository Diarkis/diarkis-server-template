// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis-server-template/examples/dive/user-online-status/cmds"
	"github.com/Diarkis/diarkis/diarkisexec"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	// In order to use dive to store the user online status, we need to
	// enable the module on setup.
	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Dive: &diarkisexec.Options{ConfigPath: "configs/shared/dive.json", ExposeCommands: true},
	})

	cmds.SetupHTTP()

	diarkisexec.SetupDiarkisHTTPServer("configs/http/main.json")

	diarkisexec.StartDiarkis()
}
