package matching

import (
	"github.com/Diarkis/diarkis/user"
)

/*
ExposeCommands exposes built-in commands to the client.
*/
func ExposeCommands() {
}

/*
SetCustomJoinCondition assigns a on join evaluation callback
to be called to evaluate if the user should join or not
*/
func SetCustomJoinCondition(callback func(string, *user.User) bool) {
}
