package roomsupport

import (
	"github.com/Diarkis/diarkis/user"
)

/*
DefineRoomSupport required to use roomsupport.ExposeCommands - This function MUST be called in HTTP server ONLY
*/
func DefineRoomSupport() {
}

/*
ExposeCommands exposes optional commands for room.

[IMPORTANT] room.ExposeCommands must be called BEFORE calling this function

[IMPORTANT] roomsupport.DefineRoomSupport() MUST be called in HTTP server for this function to work

This function exposes the following commands to the server:

  - Random Room Join
  - Register Room
  - Find Registered Room By Type
  - Sync and Update Objects
  - Chat Sync and Log
  - Start P2P
*/
func ExposeCommands() {
}

/*
BeforeRandomRoomCmd registers a callback function to be executed before random join room command:
Must be called before ExposeCommands
*/
func BeforeRandomRoomCmd(callback func(uint8, uint16, []byte, *user.User, func(error))) {
}

/*
AfterRandomRoomCmd registers a callback function to be executed before random join room command:
Must be called before ExposeCommands
*/
func AfterRandomRoomCmd(callback func(uint8, uint16, []byte, *user.User, func(error))) {
}

/*
IsRandomRoomCreated returns true if the given payload byte array of RandomRoomJoin has created a room.
This function is meant to be used with AfterRandomRoomCmd.
*/
func IsRandomRoomCreated(payload []byte) bool {
	return false
}
