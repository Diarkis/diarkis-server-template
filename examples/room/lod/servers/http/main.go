// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	// For this example we do not need to configure any of the Diarkis modules.
	// We solely rely on the base-HTTP server itself.
	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Room:       &diarkisexec.Options{ConfigPath: "configs/shared/room.json"},
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json"},
		Dive:       &diarkisexec.Options{ConfigPath: "configs/shared/dive.json"},
	})

	diarkisexec.SetupDiarkisHTTPServer("configs/http/main.json")
	diarkisexec.StartDiarkis()
}
