package room

import (
	"github.com/Diarkis/diarkis/user"
)

/*
NewRoom Creates a new room

	[NOTE] Uses mutex lock internally.
	[NOTE] This function creates a new room on the same server as it is invoked.

Error Cases:

	┌───────────────────────────────┬───────────────────────────────────────────────────────────────────────────────────────────┐
	│ Error                         │ Reason                                                                                    │
	╞═══════════════════════════════╪═══════════════════════════════════════════════════════════════════════════════════════════╡
	│ allowEmpty and join are false │ If allowEmpty and join are bot false, the room will be deleted immediately.               │
	├───────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────┤
	│ Server is offline             │ The server is in offline state (getting ready to shutdown)                                │
	│                               │ and creating new rooms is not allowed.                                                    │
	├───────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────┤
	│ Invalid TTL                   │ TTL must be greater than 0 second.                                                        │
	├───────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────┤
	│ User in another room          │ You may not join more than one room at a time.                                            │
	├───────────────────────────────┼───────────────────────────────────────────────────────────────────────────────────────────┤
	│ Max member is invalid         │ The range of max member is from 2 to 1000.                                                │
	└───────────────────────────────┴───────────────────────────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

	room.IsNewRoomError(err error) // Room module failed to create a new room

Parameters

	userData   - User data struct.
	             If you do not have user available, you may use a dummy user created by user.CreateBlankUser.
	             Using a dummy user allows you to create an empty room without having the access to a user.
	             With a dummy user, you must have allowEmpty=true, join=false.
	maxMembers - Maximum number of members allowed to join the room
	allowEmpty - If true, the room will not be deleted until TTL expires even if all members leave
	join       - If true, the creator user will join the room automatically
	ttl        - TTL of the room when it is empty in seconds. The value must be greater than 0. Minimum is 10 seconds
	interval   - Broadcast interval in milliseconds. Interval below 100ms will be forced to 0. If it's 0, no packet merge
*/
func NewRoom(userData *user.User, maxMembers int, allowEmpty bool, join bool, ttl int64, interval int64) (string, error) {
	return "", nil
}

/*
IsReserved returns true if the given uid (User ID) has a reservation with the given room.

[NOTE] Uses mutex lock internally.
*/
func IsReserved(roomID string, uid string) bool {
	return false
}

/*
Reserve the user may reserve a slot in a room so that it will be guaranteed to join the room later.

	[IMPORTANT] This function does not work if the room is not on the same server.

	[NOTE] Uses mutex lock internally.

Error Cases

	┌──────────────────────────────┬──────────────────────────────────────────────────────────────────────────────┐
	│ Error                        │ Reason                                                                       │
	╞══════════════════════════════╪══════════════════════════════════════════════════════════════════════════════╡
	│ memberIDs array is empty     │ No member IDs to make reservations with.                                     │
	├──────────────────────────────┼──────────────────────────────────────────────────────────────────────────────┤
	│ Room or room owner not found │ Cannot make a reservation if the room and/or its owner is missing.           │
	├──────────────────────────────┼──────────────────────────────────────────────────────────────────────────────┤
	│ Must be owner                │ If mustBeOwner is true, only the room owner is allowed to make reservations. │
	├──────────────────────────────┼──────────────────────────────────────────────────────────────────────────────┤
	│ Room not found               │ Cannot make a reservation if the room is missing.                            │
	├──────────────────────────────┼──────────────────────────────────────────────────────────────────────────────┤
	│ Exceeds max members          │ Make reservations that exceeds max members of the room is not allowed.       │
	├──────────────────────────────┼──────────────────────────────────────────────────────────────────────────────┤
	│ Reservation is full          │ There is no more room for reservations.                                      │
	└──────────────────────────────┴──────────────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

	room.IsReserveError(err error) // Failed to make a reservation of a room

Parameters

	roomID      - Target room ID to make a reservation.
	userData    - Owner user of the room to make a reservation.
	memberIDs   - An array of user IDs to make a reservation for.
	mustBeOwner - If true, userData must be the owner of the room.
*/
func Reserve(roomID string, userData *user.User, memberIDs []string, mustBeOwner bool) error {
	return nil
}

/*
CancelReservation removes reservation per members

	[IMPORTANT] This function does not work if the room is not on the same server.

	[NOTE] Uses mutex lock internally.

Error Cases

	┌───────────────────────────────┬────────────────────────────────────────────────────────────────────────────────┐
	│ Error                         │ Reason                                                                         │
	╞═══════════════════════════════╪════════════════════════════════════════════════════════════════════════════════╡
	│ Room and/or its owner missing │ It is not possible to cancel reservations if room and/or its owner is missing. │
	├───────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Must be owner                 │ If mustBeOwner is true, only the room owner is allowed to cancel reservations. │
	├───────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Room not found                │ Reservations cannot be canceled if the room is missing.                        │
	├───────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ No reservation                │ Cannot cancel reservations if there is no reservation.                         │
	└───────────────────────────────┴────────────────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

	room.IsCancelReserveError(err error) // Failed to cancel a reservation of a room

Parameters

	room ID     - Target room ID to cancel a reservation.
	userData    - Owner user to cancel the reservation.
	memberIDs   - An array of user IDs to cancel the reservation for.
	mustBeOwner - If true userData must be the room owner
*/
func CancelReservation(roomID string, userData *user.User, memberIDs []string, mustBeOwner bool) error {
	return nil
}

/*
CancelReservationRemote removes reservation per members

	[IMPORTANT] This function works even if the room is not on the same server.
	[IMPORTANT] The callback may be invoked on the server where the room is not held.
	            Do not attempt to use functions that works only if the room is on the same server.

	[NOTE] This function is asynchronous.

Error Cases

	┌───────────────────────────────┬────────────────────────────────────────────────────────────────────────────────┐
	│ Error                         │ Reason                                                                         │
	╞═══════════════════════════════╪════════════════════════════════════════════════════════════════════════════════╡
	│ Room and/or its owner missing │ It is not possible to cancel reservations if room and/or its owner is missing. │
	├───────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Must be owner                 │ If mustBeOwner is true, only the room owner is allowed to cancel reservations. │
	├───────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Room not found                │ Reservations cannot be canceled if the room is missing.                        │
	├───────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ No reservation                │ Cannot cancel reservations if there is no reservation.                         │
	└───────────────────────────────┴────────────────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

Besides the following errors, there are Mesh module errors as well.

	room.IsCancelReserveError(err error) // Failed to cancel a reservation of a room

Parameters

	roomID      - Target room ID to cancel a reservation.
	userData    - Owner user to cancel the reservation.
	memberIDs   - An array of user IDs to cancel the reservation for.
	mustBeOwner - If true userData must be the room owner.
	cb          - Callback function to be invoked when the cancel operation is completed.
*/
func CancelReservationRemote(roomID string, userData *user.User, memberIDs []string, mustBeOwner bool, cb func(error)) {
}
