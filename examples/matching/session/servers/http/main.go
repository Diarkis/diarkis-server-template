// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Dive: &diarkisexec.Options{ConfigPath: "configs/shared/dive.json"},
		DM:   &diarkisexec.Options{ConfigPath: "configs/shared/dm.json"},
	})
	diarkisexec.SetupDiarkisHTTPServer("configs/http/main.json")

	diarkisexec.StartDiarkis()
}
