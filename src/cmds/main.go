package cmds

import (
	customcmds "{0}/cmds/custom"
	dmcmds "{0}/cmds/dm"
	fieldcmds "{0}/cmds/field"
	groupcmds "{0}/cmds/group"
	httpcmds "{0}/cmds/http"
	matchmakercmds "{0}/cmds/matchmaker"
	p2pcmds "{0}/cmds/p2p"
	roomcmds "{0}/cmds/room"
	sessioncmds "{0}/cmds/session"
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
