package field

import (
	"github.com/Diarkis/diarkis/user"
)

/*
ExposeCommands exposes commands to the client to work with Field package
*/
func ExposeCommands() {
}

/*
BeforeSyncCmd registers a command to be executed before field sync: Must be called before ExposeCommands()

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeSyncCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterSyncCmd registers a command to be executed after field sync: Must be called before ExposeCommands()

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterSyncCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeDisappearCmd registers a command to be executed before field disappear: Must be called before ExposeCommands()

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeDisappearCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterDisappearCmd registers a command to be executed after field disappear: Must be called before ExposeCommands()

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterDisappearCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeLeaveCmd registers a command to be executed before field leave: Must be called before ExposeCommands()

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeLeaveCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterLeaveCmd registers a command to be executed after field leave: Must be called before ExposeCommands()

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterLeaveCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}
