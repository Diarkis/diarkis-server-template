// Code generated by Diarkis Puffer module: DO NOT EDIT.
//
// # Auto-generated by Diarkis Version 1.0.0
//
// - Maximum length of a string is 65535 bytes
// - Maximum length of a byte array is 65535 bytes
// - Maximum length of any array is 65535 elements
package onlinestatus

import "encoding/binary"
import "errors"
import "fmt"
import "strings"

// OnlineStatusRequestVer represents the ver of the protocol's command.
//
//	[NOTE] The value is optional and if ver is not given in the definition JSON, it will be 0.
const OnlineStatusRequestVer uint8 = 2

// OnlineStatusRequestCmd represents the command ID of the protocol's command ID.
//
//	[NOTE] The value is optional and if cmd is not given in the definition JSON, it will be 0.
const OnlineStatusRequestCmd uint16 = 500

// OnlineStatusRequest represents the command protocol data structure.
type OnlineStatusRequest struct {
	// Command version of the protocol
	Ver uint8
	// Command ID of the protocol
	Cmd  uint16
	UIDs []string
}

// NewOnlineStatusRequest creates a new instance of OnlineStatusRequest struct.
func NewOnlineStatusRequest() *OnlineStatusRequest {
	return &OnlineStatusRequest{Ver: 2, Cmd: 500, UIDs: make([]string, 0)}
}

// Pack encodes OnlineStatusRequest struct to a byte array to be delivered over the command.
func (proto *OnlineStatusRequest) Pack() []byte {
	bytes := make([]byte, 0)

	/* []string */
	uIDsSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(uIDsSizeBytes, uint16(len(proto.UIDs)))
	bytes = append(bytes, uIDsSizeBytes...)
	for i := 0; i < len(proto.UIDs); i++ {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(len(proto.UIDs[i])))
		bytes = append(bytes, b...)
		bytes = append(bytes, []byte(proto.UIDs[i])...)
	}

	// done
	return bytes
}

// Unpack decodes the command payload byte array to OnlineStatusRequest struct.
func (proto *OnlineStatusRequest) Unpack(bytes []byte) error {
	if len(bytes) < 2 {
		return errors.New("OnlineStatusRequestUnpackError")
	}

	offset := 0

	/* []string */
	uIDsSize := int(binary.BigEndian.Uint16(bytes[offset : offset+2]))
	if uIDsSize+offset > len(bytes) {
		return errors.New("UnpackError")
	}
	offset += 2
	for i := 0; i < uIDsSize; i++ {
		size := int(binary.BigEndian.Uint16((bytes[offset : offset+2])))
		offset += 2
		proto.UIDs = append(proto.UIDs, string(bytes[offset:offset+size]))
		offset += size
	}

	return nil
}

func (proto *OnlineStatusRequest) String() string {
	list := make([]string, 0)
	list = append(list, fmt.Sprint("UIDs = ", proto.UIDs))
	return strings.Join(list, " | ")
}

func (proto *OnlineStatusRequest) GetVer() uint8 {
	return proto.Ver
}
func (proto *OnlineStatusRequest) GetCmd() uint16 {
	return proto.Cmd
}
