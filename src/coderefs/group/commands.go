package group

import (
	"github.com/Diarkis/diarkis/user"
)

/*
BeforeCreateGroupCmd registers a callback function to be executed before create group command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeCreateGroupCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterCreateGroupCmd registers a callback function to be executed after create group command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterCreateGroupCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeJoinGroupCmd registers a callback function to be executed before join group command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeJoinGroupCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterJoinGroupCmd registers a callback function to be executed after join group command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterJoinGroupCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeLeaveGroupCmd registers a callback function to be executed before leave group command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeLeaveGroupCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterLeaveGroupCmd registers a callback function to be executed after leave group command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterLeaveGroupCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeBroadcastGroupCmd registers a callback function to be executed before broadcast group command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeBroadcastGroupCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterBroadcastGroupCmd registers a callback function to be executed after broadcast group command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterBroadcastGroupCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
ExposeCommands exposes commands to the client to work with the Group package

The following commands will be exposed to the client:

  - Create Group
  - Join Group
  - Leave Group
  - Broadcast To Group
*/
func ExposeCommands() {
}
