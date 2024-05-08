package server

import (
	"github.com/Diarkis/diarkis/user"
)

/*
IsDebugCommandsEnabled returns true if debug commands are enabled
*/
func IsDebugCommandsEnabled() bool {
	return false
}

/*
SetDebugCommandCallback assigns a callback for a debug command with the given debug command ID
*/
func SetDebugCommandCallback(cmd uint16, cb func(*user.User, []byte, func(error, []byte))) {
}

/*
ExposeDebugCommands exposes debug commands to the client.
Server package comes with debug-only commands for tests and debugging.

Available Debug Commands:

▷ Set Server State Online

Sets the server to "ONLINE" state.

ONLINE State - The server accepts new client connections execute all diarkis functions.

	version:          0
	command:          900
	Required payload: empty byte array
	Response:         none

▷ Set Server State Taken

Sets the server to "TAKEN" state.

TAKEN State - The server does **NOT** accept new client connections, but able to execute all diarkis functions.

	version:          0
	command:          901
	Required payload: empty byte array
	Response:         none

▷ Set Server State Offline

Sets the server to "OFFLINE" state.

OFFLINE State - The server does **NOT** accept new client connections, and can **NOT** create new rooms,
but able to join rooms and execute all other diarkis functions.

	version:          0
	command:          902
	Required payload: empty byte array
	Response:         none

▷ Server Terminate

Stops the server.
Diarkis server waits for the connected clients

	version:           0
	command:          903
	Required payload: empty byte array
	Response:         none

▷ Room Data Dump

Dumps the room data to the client.

	version:          0
	command:          904
	Required payload: empty byte array
	Response:         byte array encoded room data dump as JSON
*/
func ExposeDebugCommands() {
}

/*
UserDataDump implies a payload format for the dump user data command
*/
type UserDataDump struct {
	ID      string                 `json:"ID"`
	SID     string                 `json:"SID"`
	Data    map[string]interface{} `json:"Data"`
	Latency int64                  `json:"Latency"`
}
