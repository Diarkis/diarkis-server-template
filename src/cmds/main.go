// © 2019-2024 Diarkis Inc. All rights reserved.

package cmds

import (
	"github.com/Diarkis/diarkis"
	customcmds "github.com/Diarkis/diarkis-server-template/cmds/custom"
	dmcmds "github.com/Diarkis/diarkis-server-template/cmds/dm"
	httpcmds "github.com/Diarkis/diarkis-server-template/cmds/http"
	matchmakercmds "github.com/Diarkis/diarkis-server-template/cmds/matchmaker"
	roomcmds "github.com/Diarkis/diarkis-server-template/cmds/room"
	"github.com/Diarkis/diarkis/server"
)

func SetupUDP() {
	dmcmds.Setup()
	matchmakercmds.Setup()
	roomcmds.Setup()
	customcmds.Expose()
	exposeDebugCommands()
}

func SetupTCP() {
	dmcmds.Setup()
	matchmakercmds.Setup()
	roomcmds.Setup()
	customcmds.Expose()
	exposeDebugCommands()
}

func SetupHTTP() {
	httpcmds.Expose()
}

// exposeDebugCommands makes the Diarkis debug commands available to the clients.
// For the list of the debug commands, refer to server.ExposeDebugCommands.
//
// WARNING: Do not make the debug commands available in a production environment.
func exposeDebugCommands() {
	diarkis.OnReady(func(next func(error)) {
		server.ExposeDebugCommands()
		next(nil)
	})
}
