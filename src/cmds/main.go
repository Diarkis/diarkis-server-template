package cmds

import (
	customcmds "github.com/Diarkis/diarkis-server-template/cmds/custom"
	httpcmds "github.com/Diarkis/diarkis-server-template/cmds/http"
	dmcmds "github.com/Diarkis/diarkis-server-template/cmds/dm"
	matchmakercmds "github.com/Diarkis/diarkis-server-template/cmds/matchmaker"
	roomcmds "github.com/Diarkis/diarkis-server-template/cmds/room"
	"github.com/Diarkis/diarkis-server-template/lib/onlinestatus"
)

func SetupUDP() {
	dmcmds.Setup()
	matchmakercmds.Setup()
	roomcmds.Setup()
	onlinestatus.Setup()
	customcmds.Expose()
}

func SetupTCP() {
	dmcmds.Setup()
	matchmakercmds.Setup()
	roomcmds.Setup()
	onlinestatus.Setup()
	customcmds.Expose()
}

func SetupHTTP() {
	httpcmds.Expose()
	onlinestatus.Setup()
}
