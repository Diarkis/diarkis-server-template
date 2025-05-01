// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/matching"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	// For this example we do not need to configure any of the diarkis module.
	// We solely rely on the http server.
	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json"},
	})

	exposeMatching()

	diarkisexec.SetupDiarkisHTTPServer("configs/http/main.json")

	diarkisexec.StartDiarkis()
}

func exposeMatching() {
	// Test matching profile that does not have any criteria.
	matching.Define("All", map[string]int{
		"level": 1,
	})
}
