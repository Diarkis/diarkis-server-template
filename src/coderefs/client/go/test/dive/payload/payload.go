package divepayload

/*
DivePayload represents data structure of test payload
*/
type DivePayload struct {
	Key    string
	Value  string
	Values []string
	TTL    int64
	From   int
	To     int
}

/*
Pack encodes the given DivePayload struct into byte array
+----------+---------+------------+
| key size | 1 byte  |            |
| key      | N bytes |            |
| ttl      | 8 bytes | Big Endian |
| value    | N bytes |            |
+----------+---------+------------+
*/
func Pack(m *DivePayload, multipleValues bool) []byte {
	return nil
}

/*
Unpack decodes the given byte array to DivePayload struct
*/
func Unpack(payload []byte, multipleValues bool) (*DivePayload, error) {
	return nil, nil
}
