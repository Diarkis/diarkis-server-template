package httpcmds

// rootpath is defined in cmds/main.go

import (
	"fmt"
	"errors"
	"github.com/Diarkis/diarkis-server-template/lib/meshCmds"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/server/http"
	"strings"
)

func exposeRoom() {
	// IMPORTANT: This API does NOT work with custom server type.
	// serverType={UDP|TCP}
	// maxMembers={number of allowed members}
	// ttl={TTL of empty room in seconds}
	// interval={broadcast buffer interval}
	http.Post("/room/create/:serverType/:maxMembers/:ttl/:interval", createRoom)
}

func createRoom(resp *http.Response, req *http.Request, params *http.Params, next func(error)) {
	serverType, err := params.GetAsString("serverType")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	serverType = strings.ToUpper(serverType)
	if serverType != server.UDPType && serverType != server.TCPType {
		resp.Respond(fmt.Sprintf("Invalid server type %s", serverType), http.Bad)
		next(errors.New(fmt.Sprintf("Invalid server type %s", serverType)))
		return
	}
	maxMembers, err := params.GetAsInt("maxMembers")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	ttl, err := params.GetAsInt64("ttl")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	interval, err := params.GetAsInt64("interval")
	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	// execute internal RPC command to create a new room remotely
	meshCmds.CreateRemoteRoom(serverType, maxMembers, ttl, interval, func(err error, roomID string) {
		if err != nil {
			resp.Respond(err.Error(), http.Bad)
			next(err)
			return
		}
		resp.Respond(roomID, http.Ok)
		next(nil)
	})
}
