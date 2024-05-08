package room

import (
	"github.com/Diarkis/diarkis/user"
)

/*
DeleteRoom deletes the target room regardless of the room being not empty and sends notification message to its members.

	[IMPORTANT] This function does NOT work if the room is not on the same server.

	[NOTE] All the members of the room will be automatically kicked out of the room.
	[NOTE] Uses mutex lock internally.

Parameters

	roomID   - Target room ID to delete.
	userData - User that performs the deletion of the room.
	ver      - Command version to send the notification as.
	cmd      - Command ID to send the notification as.
	message  - Message byte array to be sent to the members.
*/
func DeleteRoom(roomID string, userData *user.User, ver uint8, cmd uint16, message []byte) {
}
