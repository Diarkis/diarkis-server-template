package util

/*
NewError creates a new error with given error code, message and stack trace.

The error message is a formatted string with variables.

# The same syntax as fmt.Sprintf

Example:

	errMessage := "Some error occurred: VariableOne:%!v(MISSING) VariableTwo:%!s(MISSING)"

	var1 := 100

	var2 := "Something bad happened"

	err := NewError(errMessage, var1, var2)

Parameters

	message - Error message as a string or an error.
	...vars - optional variables for the error message.
*/
func NewError(message interface{}, vars ...interface{}) error {
	return nil
}

/*
StackError adds another error to the given error to create a stack of multiple errors.

	[IMPORTANT] Errors should be formatted errors created by Diarkis' util.NewError

Parameters

	err  - Anchor error for the other errors to join.
	       If the error given is nil, the function returns nil.
	errs - Optional errors to stack on the anchor error.
	       If nil is given, the nil will be ignored.
*/
func StackError(err error, errs ...error) error {
	return nil
}

/*
RemoveErrorStackTrace removes error stack traces from the formatted error created by Diarkis' util.NewError.
*/
func RemoveErrorStackTrace(message string) string {
	return ""
}

/*
FlattenErrorStackTrace replaces error stack traces' line breaks and convert them to a tab.
*/
func FlattenErrorStackTrace(message string) string {
	return ""
}
