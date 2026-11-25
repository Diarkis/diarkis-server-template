// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/matching"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json"},
		Dive:       &diarkisexec.Options{ConfigPath: "configs/shared/dive.json"},
		DM:         &diarkisexec.Options{ConfigPath: "configs/shared/dm.json"},
	})

	exposeMatching()

	diarkisexec.SetupDiarkisHTTPServer("configs/http/main.json")
	diarkisexec.StartDiarkis()
}

func exposeMatching() {
	// Test-use matching profiles that do not have any criteria:
	matching.Define("team", map[string]int{"level": 1})
	matching.Define("battle", map[string]int{"level": 1})
}
