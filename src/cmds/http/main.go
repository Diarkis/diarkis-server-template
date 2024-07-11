// © 2019-2024 Diarkis Inc. All rights reserved.

package httpcmds

import (
	"github.com/Diarkis/diarkis/groupsupport" // if you do not need group module, comment out this line
	"github.com/Diarkis/diarkis/roomsupport"  // if you do not need room module, comment out this line
	"github.com/Diarkis/diarkis/server/http"
)

func Expose() {
	// this is a sample HTTP request handler
	http.Get("/hello", handleHello)
	// if you do not need HTTP-based match making comment out this line
	exposeMatchMaker()
	// if you do not need room module, comment out this line
	roomsupport.DefineRoomSupport()
	// if you do not need group module, comment out this line
	groupsupport.DefineGroupSupport()
	// custom room operations
	exposeRoom()
	// custom user online status
	exposeOnlineStatus()
}

func handleHello(res *http.Response, req *http.Request, params *http.Params, next func(error)) {
	res.Respond("Hello", http.Ok)
	next(nil)
}
