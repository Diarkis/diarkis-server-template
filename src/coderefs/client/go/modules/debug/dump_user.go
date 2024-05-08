package debug

import (
	"github.com/Diarkis/diarkis/server"
)

/*
DumpUser dumps the user data content
*/
func (d *Debug) DumpUser() {
}

/*
OnDumpUserResponse assigns a callback to DumpUser
*/
func (d *Debug) OnDumpUserResponse(callback func(success bool, userData *server.UserDataDump)) {
}
