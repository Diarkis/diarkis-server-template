package cmds

import (
	customcmds "github.com/Diarkis/diarkis-server-template/cmds/custom"
	dmcmds "github.com/Diarkis/diarkis-server-template/cmds/dm"
	httpcmds "github.com/Diarkis/diarkis-server-template/cmds/http"
	matchmakercmds "github.com/Diarkis/diarkis-server-template/cmds/matchmaker"
	roomcmds "github.com/Diarkis/diarkis-server-template/cmds/room"
)

func Setup() {
	dmcmds.Setup()
	matchmakercmds.Setup()
	roomcmds.Setup()
}

func SetupHTTP() {
	httpcmds.Expose()
}

func ExposeCustomCommands() {
	customcmds.Expose()
}
