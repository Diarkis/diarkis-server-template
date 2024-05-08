package server

import (
	"github.com/Diarkis/diarkis/server/tcp"
	"github.com/Diarkis/diarkis/server/udp"
	"github.com/Diarkis/diarkis/user"
)

const Ok = 1
const Bad = 4
const Err = 5
const TCPType = tcp.Type
const UDPType = udp.Type

/*
SetupAsConnector Initialize the server as a Connector server for an external real-time server.
This is meant to be used with an external real-time server such as Unreal Engine Dedicated server etc.
Diarkis itself will not communicate wit the clients, but focuses on auto-scaling and horizontal scaling on K8s.

	path string - This is the absolute path of the configuration file to read.
*/
func SetupAsConnector(path string) {
}

/*
SetupAsTCPServer Initialize as a TCP server.
Load configurations from the file path given.
Pass an empty string, if there is no need for configurations.

	path string - This is the absolute path of the configuration file to read.
*/
func SetupAsTCPServer(path string) {
}

/*
SetupAsUDPServer Initialize as a UDP server.
Load configurations from the file path given.
Pass an empty string, if there is no need for configurations.

	path string - This is the absolute path of the configuration file to read.
*/
func SetupAsUDPServer(path string) {
}

/*
IsP2PEnabled returns true if UDP server enabled P2P by its configuration
*/
func IsP2PEnabled() bool {
	return false
}

/*
SetupServerForGCP sets public endpoint in GCP environment
Must be called in diarkis.OnReady() as a callback
*/
func SetupServerForGCP() error {
	return nil
}

/*
SetupServerForAWS sets public endpoint in AWS environment
Must be called in diarkis.OnReady() as a callback
*/
func SetupServerForAWS() error {
	return nil
}

/*
SetupServerForAzure sets public endpoint in Microsoft Azure environment
Must be called in diarkis.OnReady() as a callback
*/
func SetupServerForAzure() error {
	return nil
}

/*
SetupServerForAlibaba sets public endpoint in Alibaba Cloud environment
Must be called in diarkis.OnReady() as a callback
*/
func SetupServerForAlibaba() error {
	return nil
}

/*
SetupServerForTencent sets public endpoint in Tencent Cloud environment
Must be called in diarkis.OnReady() as a callback
*/
func SetupServerForTencent() error {
	return nil
}

/*
SetupServerForLinode sets public endpoint in Linode environment
Must be called in diarkis.OnReady() as a callback
*/
func SetupServerForLinode() error {
	return nil
}

/*
SetupGenericCloudServer sets public endpoint in a generic way
Must be called in diarkis.OnReady() as a callback
*/
func SetupGenericCloudServer() error {
	return nil
}

/*
SetPublicEndPoint sets public end point to be sent to the client to be used to connect

	addr string - This is the public endpoint for the server to be registered and used by the clients.

This is usually used internally.
*/
func SetPublicEndPoint(addr string) {
}

/*
GetEndPoint returns the end point address the server is bound with
*/
func GetEndPoint() string {
	return ""
}

/*
IsTCP Returns true if the server is setup as TCP server
*/
func IsTCP() bool {
	return false
}

/*
IsUDP Returns true if the server is setup as UDP server
*/
func IsUDP() bool {
	return false
}

/*
OnDisconnect Registers a callback on connection termination for TCP/UDP
*/
func OnDisconnect(callback func(string, *user.User)) {
}

/*
MarkServerAsTakenIf flags the server as "TAKEN" if the callback returns true.
The callback is invoked every 2 second and if the callback returns true,
the server will be marked as "ONLINE".

	[IMPORTANT] If the server is marked as "OFFLINE", this function will be ignored.

	[NOTE] This function is executed every 2 seconds so there will be race condition and it is not precise.

TAKEN - When the server is marked as "TAKEN", the server will NOT accept new user connections.

Example: The example code below uses CCU of the node to control TAKEN <==> ONLINE

	server.MarkServerAsTakenIf(func() bool {
		// user package can tell you the CCU of the node
		ccu := user.GetCCU()
		if ccu >= maxAllowedCCU {
			// this will mark the server as TAKEN
			return true
		}
		// this will mark the server as ONLINE
		return false
	})
*/
func MarkServerAsTakenIf(callback func() bool) {
}

/*
Online flags the server to be online and will accept new users
*/
func Online() {
}

/*
Taken flags the server to be taken and will not accept new users.
*/
func Taken() {
}

/*
Offline flags the server to be offline (will be shutdown) and will not accept new users
*/
func Offline() {
}

/*
SendPM sends a direct message to another user - The sent private message is a reliable delivery.

# Deprecated

This function has been deprecated and will be removed in the future version without a warning.
Use DM module instead.

	nodeAddr string - This is the node (server) internal address of the target user.
	sid      string - This is the sid (session ID) of the target user.
	ver      uint8  - This is the command version to be used with the message sent.
	cmd      uint16 - This is the command version to be used with the message sent.
	message  []byte - This is the message data byte array.
*/
func SendPM(nodeAddr string, sid string, ver uint8, cmd uint16, message []byte) {
}

/*
HookAllCommands Registers a packet handler as a hook to all commands

	handler func(userData *user.User, next func(error)) - Function to be invoked on command hook.

NOTE: The second argument next func(error) must be called at the end of handler to move the operation to next.
*/
func HookAllCommands(handler func(userData *user.User, next func(error))) {
}

/*
IsCommandDuplicate returns true if ver and command ID are used elsewhere when this function is invoked

	ver uint8  - Command version to be checked for duplication.
	cmd uint16 - Command ID to be checked for duplication.
*/
func IsCommandDuplicate(ver uint8, commandID uint16) bool {
	return false
}

/*
Command has been deprecated: Use HandleCommand instead for other servers.

# Deprecated

This function has been deprecated and will be removed in the future version without a warning.

	ver     uint8  - Command version of the handler.
	cmd     uint16 - Command ID of the handler.
	handler func(userData *user.User, next func(error))

Returns an error if it fails to assign a handler.

NOTE: If you assign multiple handles with the same ver and cmd, all handlers will be executed in the order of assignment.
*/
func Command(ver uint8, cmd uint16, handler func(userData *user.User, next func(error))) {
}

/*
HandleCommand registers a command handler function for TCP and UDP/RUDP server.

Returns an error if it fails to assign a handler.

	[NOTE] If you assign multiple handles with the same ver and cmd,
	       all handlers will be executed in the order of assignment.

Handler callback's next function:

The next function that the handler callback receives must be invoked at the end of all operations within the handler.

This ensures the handling of a command is completed and Diarkis server proceeds to the next handler of the same command ID or another command ID.

By blocking or not invoking the next function, Diarkis server will pause the handling of all commands for the user
and eventually the user client will disconnect because it cannot handle echo (keep alive) command that every user client needs to stay connected to the server.

Parameters

	ver     - Command version of the handler.
	cmd     - Command ID of the handler.
	handler - Handler function to be executed when the server receives the command.
	          func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))

Handler Callback Parameters

	ver     - Command ver sent from the client.
	cmd     - Command ID sent from the client.
	payload - Command payload sent from the client.
	userData - User data representing the client that sent the command.
	next    - The function to signal the command to move on to the next operation of the same ver and command ID.
	          This function must be called at the end of all operations in the callback.
	          If you pass an error, it will not proceed to the next operations of the same command ver and ID.
*/
func HandleCommand(ver uint8, cmd uint16, handler func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))) error {
	return nil
}

/*
OnKeepAlive assigns a callback function on TCP heartbeat or UDP echo

	handler - Function to be invoked on keep alive message from the clients.

	[NOTE] next func(error) must be called at the end of handler's operation.
*/
func OnKeepAlive(handler func(userData *user.User, next func(error))) {
}
