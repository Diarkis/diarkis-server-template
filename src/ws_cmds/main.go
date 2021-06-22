package cmds

import (
	httpcmds "{0}/ws_cmds/http"
	roomcmds "{0}/ws_cmds/room"
	groupcmds "{0}/ws_cmds/group"
	fieldcmds "{0}/ws_cmds/field"
	customcmds "{0}/ws_cmds/custom"
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
	customcmds.Expose()
}
