package httpcmds

import (
	dmcmds "github.com/Diarkis/diarkis-server-template/cmds/dm"

	"github.com/Diarkis/diarkis/groupsupport" // if you do not need group module, comment out this line
	"github.com/Diarkis/diarkis/roomsupport"  // if you do not need room module, comment out this line
	"github.com/Diarkis/diarkis/server/http"
)

func Expose(rootpath string) {
	// this is a sample HTTP request handler
	http.Get("/hello", handleHello)
	// if you do not need HTTP-based match making comment out this line
	exposeMatchMaker(rootpath)
	// if you do not need room module, comment out this line
	roomsupport.DefineRoomSupport()
	// if you do not need group module, comment out this line
	groupsupport.DefineGroupSupport()
	// custom room operations
	exposeRoom(rootpath)
	// if you do not need dm module, comment out this line
	dmcmds.Expose(rootpath)
}

func handleHello(res *http.Response, req *http.Request, params *http.Params, next func(error)) {
	res.Respond("Hello", http.Ok)
	next(nil)
}
