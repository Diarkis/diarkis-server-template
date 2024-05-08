package group

/*
IsNewGroupError returns true if the given error is or contains a GroupNewError.
*/
func IsNewGroupError(err error) bool {
	return false
}

/*
IsJoinError returns true if the given error is or contains a GroupJoinError.
*/
func IsJoinError(err error) bool {
	return false
}

/*
IsLeaveError returns true if the given error is or contains a GroupLeaveError.
*/
func IsLeaveError(err error) bool {
	return false
}
