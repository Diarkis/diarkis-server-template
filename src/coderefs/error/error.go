package error

/*
Error is the custom error data
*/
type Error struct{}

/*
New creates a new error data
*/
func New(message string, code uint16) *Error {
	return nil
}

/*
ToError converts error bytes to an error data
*/
func ToError(bytes []byte) *Error {
	return nil
}

/*
Message returns the error message
*/
func (e *Error) Message() string {
	return ""
}

/*
Code returns the error code
*/
func (e *Error) Code() uint16 {
	return 0
}

/*
Bytes returns the byte array for the client
*/
func (e *Error) Bytes() []byte {
	return nil
}

/*
StackTrace return a string of stacktrace
*/
func (e *Error) StackTrace() string {
	return ""
}
