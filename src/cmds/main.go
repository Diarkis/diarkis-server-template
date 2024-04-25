package cmds

import (
	customcmds "github.com/Diarkis/diarkis-server-template/cmds/custom"
	dmcmds "github.com/Diarkis/diarkis-server-template/cmds/dm"
	fieldcmds "github.com/Diarkis/diarkis-server-template/cmds/field"
	groupcmds "github.com/Diarkis/diarkis-server-template/cmds/group"
	httpcmds "github.com/Diarkis/diarkis-server-template/cmds/http"
	matchmakercmds "github.com/Diarkis/diarkis-server-template/cmds/matchmaker"
	p2pcmds "github.com/Diarkis/diarkis-server-template/cmds/p2p"
	roomcmds "github.com/Diarkis/diarkis-server-template/cmds/room"
	sessioncmds "github.com/Diarkis/diarkis-server-template/cmds/session"
)

var rootpath string

func Setup(rpath string) {
	rootpath = rpath
}

func ExposeHTTP() {
	httpcmds.Expose(rootpath)
	dmcmds.Expose(rootpath)
	fieldcmds.Expose(rootpath)
}

func ExposeServer() {
	roomcmds.Expose()
	groupcmds.Expose(rootpath)
	fieldcmds.Expose(rootpath)
	dmcmds.Expose(rootpath)
	matchmakercmds.Expose(rootpath)
	p2pcmds.Expose()
	customcmds.Expose()
	sessioncmds.Expose()
}
