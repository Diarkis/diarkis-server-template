// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"errors"

	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/server/http"
)

func main() {
	logConfigPath := "/configs/shared/log.json"
	meshConfigPath := "/configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Dive:       &diarkisexec.Options{ConfigPath: "/configs/shared/dive.json", ExposeCommands: true},
		Field:      &diarkisexec.Options{ConfigPath: "/configs/shared/field.json", ExposeCommands: true},
		DM:         &diarkisexec.Options{ConfigPath: "/configs/shared/dm.json", ExposeCommands: true},
		MatchMaker: &diarkisexec.Options{ConfigPath: "/configs/shared/matching.json", ExposeCommands: true},
	})

	expose()

	diarkisexec.SetupDiarkisHTTPServer("/configs/http/main.json")

	diarkisexec.StartDiarkis()
}

func expose() {
	http.Post("/echo", handleEcho)
}

func handleEcho(res *http.Response, req *http.Request, params *http.Params, next func(error)) {
	if req.JSONBody == nil {
		err := errors.New("expect JSON body")
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}

	message, ok := req.JSONBody["message"].(string)
	if !ok {
		err := errors.New("missing parameter message")
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}

	enc, err := json.Marshal(map[string]any{"echo": message})
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}

	res.SetHeader("Content-Type", "application/json")
	res.SendBytes(enc, http.Ok)
	// move on to other handlers
	next(nil)
}
