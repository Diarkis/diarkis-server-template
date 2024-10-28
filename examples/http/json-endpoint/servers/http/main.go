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

	// For this example we do not need to configure any of the diarkis module.
	// We solely rely on the http server.
	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{})

	expose()

	diarkisexec.SetupDiarkisHTTPServer("/configs/http/main.json")

	diarkisexec.StartDiarkis()
}

func expose() {
	// Expose a simple handler on POST /echo
	http.Post("/echo", handleEcho)
}

func handleEcho(res *http.Response, req *http.Request, params *http.Params, next func(error)) {
	// diarkis automatically parse JSON request body if the content-type
	// is set to application/json.
	if req.JSONBody == nil {
		err := errors.New("expect JSON body")
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}

	// If the message property is missing or if its type is not string
	// return an error.
	message, ok := req.JSONBody["message"].(string)
	if !ok {
		err := errors.New("missing/invalid parameter message")
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}

	// Create JSON response
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
