package room

import (
	"github.com/Diarkis/diarkis/user"
)

/*
Leave Leaves from a room and notify the other members on leaving. If message is empty Broadcast will not be sent.

	[NOTE] Uses mutex lock internally.

	[IMPORTANT] This function does not work if the room is not on the same server.
	[IMPORTANT] When a member user of a room disconnects unexpectedly, the room will automatically detect the leave of
	            the user and sends on leave message to the reset of the member users.
	            The message contains room ID and the user ID that disconnected.
	            In order to have custom message for this behavior, use SetOnDiscardCustomMessage.

Error Cases

	┌──────────────────────────┬─────────────────────────────────────────────────────┐
	│ Error                    │ Reason                                              │
	╞══════════════════════════╪═════════════════════════════════════════════════════╡
	│ Room not found           │ Room to leave from is not found.                    │
	├──────────────────────────┼─────────────────────────────────────────────────────┤
	│ The user is not a member │ The user is not a member of the room to leave from. │
	└──────────────────────────┴─────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

	room.IsLeaveError(err error) // Failed to leave from a room

Parameters

	roomID string       - Target room ID to leave from.
	userData *user.User - User that will be leaving the room.
	ver uint8           - Command version to be used for message sent when leave is successful.
	cmd uint16          - Command ID to be used for message sent when leave is successful.
	message []byte      - Message byte array to be sent as successful message to other room members.
	                      If message is either nil or empty, the message will not be sent.
*/
func Leave(roomID string, userData *user.User, ver uint8, cmd uint16, message []byte) error {
	return nil
}
