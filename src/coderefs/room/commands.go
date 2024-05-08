package room

import (
	"github.com/Diarkis/diarkis/user"
)

/*
BeforeCreateRoomCmd registers a callback function to be executed before create room command:

# Must be called before ExposeCommands

Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeCreateRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeUpdateRoomPropCmd registers a callback function to be executed before update room properties:

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
func BeforeUpdateRoomPropCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeGetRoomPropCmd registers a callback function to be executed before get room properties:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeGetRoomPropCmd(callback func(ver uint8, cmd uint16, payloa []byte, userData *user.User, next func(error))) {
}

/*
AfterCreateRoomCmd registers a callback function to be executed after create room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterCreateRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeJoinRoomCmd registers a callback function to be executed before join room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeJoinRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterJoinRoomCmd registers a callback function to be executed after join room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterJoinRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeLeaveRoomCmd registers a callback function to be executed before leave room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeLeaveRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterLeaveRoomCmd registers a callback function to be executed after leave room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterLeaveRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeBroadcastRoomCmd registers a callback function to be executed before broadcast room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeBroadcastRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterBroadcastRoomCmd registers a callback function to be executed after broadcast room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterBroadcastRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
BeforeMessageRoomCmd registers a callback function to be executed before message room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func BeforeMessageRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterMessageRoomCmd registers a callback function to be executed after message room command:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterMessageRoomCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterUpdateRoomPropCmd registers a callback function to be executed after update room properties:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterUpdateRoomPropCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
AfterGetRoomPropCmd registers a callback function to be executed after get room properties:

# Must be called before ExposeCommands

Parameters

	ver      - Command ver sent from the client.
	cmd      - Command ID sent from the client.
	payload  - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next     - The function to signal the command to move on to the next operation of the same ver and command ID.
	           This function must be called at the end of all operations in the callback.
	           If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func AfterGetRoomPropCmd(callback func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) {
}

/*
ExposeCommands Exposes TCP or UDP/RUDP commands to execute room functions from the client

This function exposes the following commands to the client:

  - Create Room
  - Join Room
  - Leave Room
  - Broadcast To Room
  - Send Message To Room Members
  - Get and Update Room Properties
  - Reserve Room
  - Cancel Room Reservation
  - Get Number of Room Members
  - Room Migration
  - Property Sync
  - Propagate Latency
*/
func ExposeCommands() {
}

/*
ExposePufferCommands Exposes TCP or UDP/RUDP puffer commands to execute room functions from the client

This function exposes the following commands to the client:

  - Create Room
  - Join Room
  - Leave Room
  - Broadcast To Room
  - Send Message To Room Members
  - Get and Update Room Properties
  - Reserve Room
  - Cancel Room Reservation
  - Get Number of Room Members
  - Room Migration
  - Property Sync
  - Propagate Latency
*/
func ExposePufferCommands() {
}
