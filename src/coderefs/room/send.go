package room

import (
	"github.com/Diarkis/diarkis/user"
)

/*
Broadcast Sends a broadcast message to the other members in the room

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] Multiple broadcast payloads maybe combined and the client must parse the combined payloads BytesToBytesList of Diarkis client SDK.
	[IMPORTANT] Built-in events such as OnMemberBroadcast handles the separation of combined payloads internally,
	            so the application does not need to handle it by using BytesToBytesList.

	[NOTE]      Uses mutex lock internally.

Parameters

	roomID     - Target room ID.
	senderUser - User to send broadcast.
	ver        - Command version to be used as broadcast message.
	cmd        - Command ID to be used as broadcast message.
	message    - Message byte array.
	reliable   - If true, UDP will be RUDP.
*/
func Broadcast(roomID string, senderUser *user.User, ver uint8, cmd uint16, message []byte, reliable bool) {
}

/*
Message Sends a message to selected members of the room

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] All message payloads have leading 4-byte-long header. The client must skip the first 4 bytes.
	[IMPORTANT] Built-in events such as OnMemberMessage handles the separation of combined payloads internally,
	            so the application does not need to handle it by using BytesToBytesList.

	[NOTE]      Uses mutex lock internally.

Parameters

	roomID     - Target room ID.
	memberIDs  - An array of member user IDs to send message to.
	senderUser - User to send broadcast.
	ver        - Command version to be used as broadcast message.
	cmd        - Command ID to be used as broadcast message.
	message    - Message byte array.
	reliable   - If true, UDP will be RUDP.
*/
func Message(roomID string, memberIDs []string, senderUser *user.User, ver uint8, cmd uint16, message []byte, reliable bool) {
}

/*
Announce Sends a message to selected members to the room without having a "sender" - This function must NOT be called by user

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] All announce payloads have leading 4-byte-long header. The client must skip the first 4 bytes.

	[NOTE]      Uses mutex lock internally.

Parameters

	roomID     - Target room ID.
	memberIDs  - An array of member user IDs to send message to.
	senderUser - User to send broadcast.
	ver        - Command version to be used as broadcast message.
	cmd        - Command ID to be used as broadcast message.
	message    - Message byte array.
	reliable   - If true, UDP will be RUDP.
*/
func Announce(roomID string, memberIDs []string, ver uint8, cmd uint16, message []byte, reliable bool) {
}

/*
Relay sends a message to all members of the room except for the sender.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] Unlike Broadcast, Message, and Announce, Relay does not combine multiple messages.
	            Therefore, the message sent does not have to be separated by BytesToBytesList on the client side.
	[IMPORTANT] Unlike Broadcast, Message, and Announce, Relay does not send the message to the sender user.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID   - Room ID of the room to send relay message to.
	sender   - Sender user. The sender user must be a member of the room.
	ver      - Command version of the message sent to the members of the room.
	cmd      - Command ID of the message sent to the members of the room.
	message  - Message to be sent to the members of the room.
	reliable - If true, the message will be sent as an RUDP if the network protocol is UDP.
	           If the network protocol used is TCP, this flag is ignored.
*/
func Relay(roomID string, sender *user.User, ver uint8, cmd uint16, message []byte, reliable bool) bool {
	return false
}

/*
RelayTo sends a message to selected members of the room except for the sender.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] Unlike Broadcast, Message, and Announce, RelayTo does not combine multiple messages.
	            Therefore, the message sent does not have to be separated by BytesToBytesList on the client side.
	[IMPORTANT] Unlike Broadcast, Message, and Announce, Relay does not send the message to the sender user.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID   - Room ID of the room to send relay message to.
	sender   - Sender user. The sender user must be a member of the room.
	members  - An array of target member user IDs to send relay message to.
	ver      - Command version of the message sent to the members of the room.
	cmd      - Command ID of the message sent to the members of the room.
	message  - Message to be sent to the members of the room.
	reliable - If true, the message will be sent as an RUDP if the network protocol is UDP.
	           If the network protocol used is TCP, this flag is ignored.
*/
func RelayTo(roomID string, sender *user.User, members []string, ver uint8, cmd uint16, message []byte, reliable bool) bool {
	return false
}

/*
SyncCreateTime sends room's created time (in seconds) to the selected (or all) members of the room

	[IMPORTANT] This function does NOT work if the room is not on the same server.

	[NOTE]      Uses mutex lock internally.

Parameters

	roomID    - Target room ID.
	ver       - Command version to be used for the message.
	cmd       - Command ID to be used for the message.
	memberIDs - An array of member user IDs to sync the room creation time.
*/
func SyncCreateTime(roomID string, ver uint8, cmd uint16, memberIDs []string) {
}
