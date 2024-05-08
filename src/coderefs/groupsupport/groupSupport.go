package groupsupport

import (
	"github.com/Diarkis/diarkis/user"
)

/*
DefineGroupSupport required to use groupsupport.ExposeCommands - This function MUST be called in HTTP server ONLY
*/
func DefineGroupSupport() {
}

/*
ExposeCommands exposes optional commands for group - group.ExposeCommands must be called BEFORE calling this function
groupsupport.DefineGroupSupport() MUST be called in HTTP server for this function to work
*/
func ExposeCommands() {
}

/*
BeforeRandomGroupCmd registers a callback function to be executed before random join group command:
Must be called before ExposeCommands
*/
func BeforeRandomGroupCmd(callback func(uint8, uint16, []byte, *user.User, func(error))) {
}

/*
AfterRandomGroupCmd registers a callback function to be executed before random join group command:
Must be called before ExposeCommands
*/
func AfterRandomGroupCmd(callback func(uint8, uint16, []byte, *user.User, func(error))) {
}
