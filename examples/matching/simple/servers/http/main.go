// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/matching"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	// For this example we do not need to configure any of the Diarkis modules.
	// We solely rely on the base-HTTP server itself.
	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json"},
	})

	exposeMatching()

	diarkisexec.SetupDiarkisHTTPServer("configs/http/main.json")
	diarkisexec.StartDiarkis()
}

func exposeMatching() {
	{
		levelMatchProfile := make(map[string]int)
		levelMatchProfile["level"] = 10 // level
		matching.Define("LevelMatch", levelMatchProfile)
	}

	{
		levelMatchProfile := make(map[string]int)
		levelMatchProfile["level"] = 1 // level
		matching.Define("LevelMatchExact", levelMatchProfile)
	}
}
