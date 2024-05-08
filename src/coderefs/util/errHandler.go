package util

/*
ErrData returns either a UTF8 string message as a byte array or structured error data as a byte array.

	[IMPORT] This is meant to be sent to the client.

In order to use structured error data, use DIARKIS_USE_STRUCT_ERR env.
*/
func ErrData(message string, code uint16) []byte {
	return nil
}
