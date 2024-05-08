package room

/*
IsNewRoomError returns true if the given error is or contains a RoomNewError.
*/
func IsNewRoomError(err error) bool {
	return false
}

/*
IsJoinError returns true if the given error is or contains a RoomJoinError.
*/
func IsJoinError(err error) bool {
	return false
}

/*
IsLeaveError returns true if the given error is or contains a RoomLeaveError.
*/
func IsLeaveError(err error) bool {
	return false
}

/*
IsReserveError returns true if the given error is or contains a RoomReserveError.
*/
func IsReserveError(err error) bool {
	return false
}

/*
IsCancelReserveError returns true if the given error is or contains a RoomCancelReserveError.
*/
func IsCancelReserveError(err error) bool {
	return false
}

/*
IsStateUpdateError returns true if the given error is or contains a RoomStateUpdateError.
*/
func IsStateUpdateError(err error) bool {
	return false
}

/*
IsStateUpsertError returns true if the given error is or contains a RoomStateUpsertError.
*/
func IsStateUpsertError(err error) bool {
	return false
}

/*
IsFailureError returns true if the given error is or contains a generic RoomFailureError.
*/
func IsFailureError(err error) bool {
	return false
}
