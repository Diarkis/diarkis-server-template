package cmds

import (
	httpcmds "{0}/cmds/http"
	roomcmds "{0}/cmds/room"
	groupcmds "{0}/cmds/group"
	fieldcmds "{0}/cmds/field"
	dmcmds "{0}/cmds/dm"
	matchmakercmds "{0}/cmds/matchmaker"
	customcmds "{0}/cmds/custom"
)

var rootpath string

func Setup(rpath string) {
	rootpath = rpath
}

func ExposeHTTP() {
	httpcmds.Expose(rootpath)
}

func ExposeServer() {
	roomcmds.Expose()
	groupcmds.Expose(rootpath)
	fieldcmds.Expose(rootpath)
	dmcmds.Expose(rootpath)
	matchmakercmds.Expose(rootpath)
	customcmds.Expose()
}
