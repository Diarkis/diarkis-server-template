package mesh

/*
IsSendError returns bool if the given error is or contains Send error.

The Send error is used by Send, USend, SendMany, USendMany.
*/
func IsSendError(err error) bool {
	return false
}

/*
IsRequestError returns true if the given error is or contains a SendRequest error.
*/
func IsRequestError(err error) bool {
	return false
}

/*
IsRequestTimeoutError returns true if the given error is or contains a SendRequest time out error.
*/
func IsRequestTimeoutError(err error) bool {
	return false
}

/*
IsUnhealthyError returns true if the given error is or contains mesh health check error.
*/
func IsUnhealthyError(err error) bool {
	return false
}
