package uuid

/*
UUID Data structure of UUID
*/
type UUID struct {
	Bytes  []byte
	String string
}

/*
New Creates a new UUID struct
*/
func New() (*UUID, error) {
	return nil, nil
}

/*
GenUUIDBytes returns a UUID v4 byte array.
*/
func GenUUIDBytes() ([]byte, error) {
	return nil, nil
}

/*
GenUUIDString returns a UUID v4 string.
*/
func GenUUIDString() (string, error) {
	return "", nil
}

/*
Clear Resets internal values of UUID struct
*/
func (uuid *UUID) Clear() {
}

/*
FromString Converts string to UUID struct
*/
func FromString(uuidString string) (*UUID, error) {
	return nil, nil
}

/*
FromBytes Converts []byte to UUID struct
*/
func FromBytes(uuidBytes []byte) (*UUID, error) {
	return nil, nil
}
