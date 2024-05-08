package dm

import (
	"github.com/Diarkis/diarkis/user"
)

/*
SendData represents internally used data
*/
type SendData struct {
	MyUserID   string `json:"myUserID"`
	YourUserID string `json:"yourUserID"`
	Ver        uint8  `json:"ver"`
	Cmd        uint16 `json:"cmd"`
	Message    []byte `json:"message"`
}

/*
FindData represents internally used data
*/
type FindData struct {
	MyUserID   string `json:"myUserID"`
	YourUserID string `json:"yourUserID"`
	MyNodeAddr string `json:"myNodeAddr"`
	Ver        uint8  `json:"ver"`
	Cmd        uint16 `json:"cmd"`
	Message    []byte `json:"message"`
}

/*
Setup must be called in the server main function to setup dm module before diarkis.Start

	[IMPORTANT] DM module uses Dive module internally as its distributed memory storage.
*/
func Setup(confpath string) {
}

/*
SetOnClearCache assigns a callback to be invoked when a cache (connection) between two user has been cleared (deleted).

The callback is invoked on the user that stayed with the connection.

This callback is useful when you need to notify the other user when the user is gone from the connection unexpectedly.

	[IMPORTANT] The callbacks will NOT be invoked if ClearCacheByUserID with message byte array
	            because if the message is not empty the message will be sent to the user to notify that the connection cache has been cleared.

Parameters

	cb - Callback
	     func(clearedUserID string, userData *user.User)
	       clearedUserID - User ID of the user that left the connection.
	       userData      - User that stayed with the connection.
*/
func SetOnClearCache(cb func(clearedUserID string, userData *user.User)) {
}

/*
SetOnReceiveMessage assigns a callback to be invoked when a user receives a message from another user.

If the callback returns false, the receiving end user declines to receive the message and the message will not be delivered.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This function is not goroutine safe and it recommended to use this function only when the server process starts.

Parameters

	cb - Callback function to be invoked on every message sent.
	     cb func(ver uint8, cmd uint16, message []byte, receiver *user.User) bool
	     ver      - Direct message command ver.
	     cmd      - Direct message command ID.
	     message  - Direct message byte array.
	     senderID - The user ID of the sender.
	     receiver - Receiver user data.
*/
func SetOnReceiveMessage(cb func(ver uint8, cmd uint16, message []byte, senderID string, receiver *user.User) bool) bool {
	return false
}

/*
SetOnMessage assigns a callback to be invoked when a user receives a message, but just before the server sends the message to the user client.

The callback allows you to change the message byte array by returning a modified
or completely changed message byte array to be sent to the client.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] Multiple callbacks maybe assigned and invoked in the order of assignments.
	            The last callback's returned message byte array will be sent to the user client.
	[IMPORTANT] If the callback returns nil, it will not change the message byte array.
	[IMPORTANT] If the callback returns an empty byte array, it will change the message byte array.
	            If the message is an empty byte array, the message will NOT be sent to the user client.
*/
func SetOnMessage(cb func(ver uint8, cmd uint16, message []byte, senderID string, receiver *user.User) []byte) bool {
	return false
}

/*
SetOnUserDisconnect assigns a callback to be invoked
when the user is disconnected to dynamically create a disconnection DM message for the peer user.

The returned byte array is the message to be sent to the peer user.
*/
func SetOnUserDisconnect(cb func(disconnectedUserData *user.User, peerUserID string) []byte) bool {
	return false
}

/*
SendMessageByUserID sends a message to another user by their user ID.

	[IMPORTANT] It may take some time (order of seconds) to initially send a message to another user client as the number of servers increase.

	[NOTE] If the message byte array is empty, the message will not be sent.

Error Cases:

	+-------------------+-------------------------------------------------------------------------------------------+
	| Error             | Reason                                                                                    |
	+-------------------+-------------------------------------------------------------------------------------------+
	| Reliable time out | Internal server communication time out. Most likely caused by server load being to high.  |
	| Message is nil    | Message byte array is nil and cannot proceed with the operation.                          |
	| Nodes not found   | There are no servers in the cluster to search for user. Likely caused by corrupt cluster. |
	| User not found    | Either the target user is not found or no longer connected.                               |
	+-------------------+-------------------------------------------------------------------------------------------+

Parameters

	userData   - Sender user data: sender
	yourUserID - Target user ID  : recipient
	ver        - Server push command version
	cmd        - Server push command ID
	message    - Message byte array: If empty, no message will be sent
	cb         - Callback invoked on message being sent
*/
func SendMessageByUserID(userData *user.User, yourUserID string, ver uint8, cmd uint16, message []byte, cb func(error)) {
}

/*
ClearCacheByUserID deletes the locally stored and remotely stored communication channel cache
and sends a message to the peer user client for one last time.

This terminates the direct message communication with the peer user client and sends a message to the other.

	[IMPORTANT] If message size is 0 byte or a nil, it will not send a last message to the peer client.
	[IMPORTANT] If there is no cache to delete or no peer user client found, it will fail silently.

Parameters

	userData   - User that is deleting the local cache and terminating the communication with the peer user.
	yourUserID - Peer user ID of the cache to delete.
	ver        - One last message ver.
	cmd        - One last message command ID.
	message    - One last message byte array.
*/
func ClearCacheByUserID(userData *user.User, yourUserID string, ver uint8, cmd uint16, message []byte) {
}
