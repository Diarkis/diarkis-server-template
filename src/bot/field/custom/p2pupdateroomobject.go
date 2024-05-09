// Code generated by Diarkis Puffer module: DO NOT EDIT.
//
// Auto-generated by Diarkis Version 1.0.0
//
// - Maximum length of a string is 65535 bytes
// - Maximum length of a byte array is 65535 bytes
// - Maximum length of any array is 65535 elements
package custom

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	util "github.com/Diarkis/diarkis/util"
)
// P2PUpdateRoomObjectVer represents the ver of the protocol's command.
//
//	[NOTE] The value is optional and if ver is not given in the definition JSON, it will be 0.
const P2PUpdateRoomObjectVer uint8 = 2

// P2PUpdateRoomObjectCmd represents the command ID of the protocol's command ID.
//
//	[NOTE] The value is optional and if cmd is not given in the definition JSON, it will be 0.
const P2PUpdateRoomObjectCmd uint16 = 10101

// P2PUpdateRoomObject represents the command protocol data structure.
type P2PUpdateRoomObject struct {
	// Command version of the protocol
	Ver uint8
	// Command ID of the protocol
	Cmd uint16
	ObjectName string
	ObjectPropValues []float64
	ObjetcPropKeys []string
	UpdateMode uint8
}

// NewP2PUpdateRoomObject creates a new instance of P2PUpdateRoomObject struct.
func NewP2PUpdateRoomObject() *P2PUpdateRoomObject {
	return &P2PUpdateRoomObject{ Ver: 2, Cmd: 10101, ObjectName: "", ObjetcPropKeys: make([]string, 0), ObjectPropValues: make([]float64, 0), UpdateMode: 0 }
}

// Pack encodes P2PUpdateRoomObject struct to a byte array to be delivered over the command.
func (proto *P2PUpdateRoomObject) Pack() []byte {
	bytes := make([]byte, 0)

	/* string */
	objectNameSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(objectNameSizeBytes, uint16(len(proto.ObjectName)))
	bytes = append(bytes, objectNameSizeBytes...)
	bytes = append(bytes, []byte(proto.ObjectName)...)

	/* []float64 */
	objectPropValuesSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(objectPropValuesSizeBytes, uint16(len(proto.ObjectPropValues)))
	bytes = append(bytes, objectPropValuesSizeBytes...)
	for i := 0; i < len(proto.ObjectPropValues); i++ {
		b := make([]byte, 8)
		bits := math.Float64bits(proto.ObjectPropValues[i])
		b[0] = byte(bits >> 56)
		b[1] = byte(bits >> 48)
		b[2] = byte(bits >> 40)
		b[3] = byte(bits >> 32)
		b[4] = byte(bits >> 24)
		b[5] = byte(bits >> 16)
		b[6] = byte(bits >> 8)
		b[7] = byte(bits)
		b = util.ReverseBytes(b)
		bytes = append(bytes, b...)
	}

	/* []string */
	objetcPropKeysSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(objetcPropKeysSizeBytes, uint16(len(proto.ObjetcPropKeys)))
	bytes = append(bytes, objetcPropKeysSizeBytes...)
	for i := 0; i < len(proto.ObjetcPropKeys); i++ {
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(len(proto.ObjetcPropKeys[i])))
		bytes = append(bytes, b...)
		bytes = append(bytes, []byte(proto.ObjetcPropKeys[i])...)
	}

	/* uint8 */
	updateModeBytes := make([]byte, 1)
	updateModeBytes[0] = proto.UpdateMode
	bytes = append(bytes, updateModeBytes...)

	// done
	return bytes
}

// Unpack decodes the command payload byte array to P2PUpdateRoomObject struct.
func (proto *P2PUpdateRoomObject) Unpack(bytes []byte) error {
	if len(bytes) < 7 {
		return errors.New("P2PUpdateRoomObjectUnpackError")
	}

	offset := 0

	/* string */
	objectNameSize := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
	if objectNameSize + offset > len(bytes) {
		return errors.New("UnpackError")
	}
	offset += 2
	proto.ObjectName = string(bytes[offset:offset + objectNameSize])
	offset += objectNameSize

	/* []float64 */
	objectPropValuesSize := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
	if objectPropValuesSize + offset > len(bytes) {
		return errors.New("UnpackError")
	}
	offset += 2
	for i := 0; i < objectPropValuesSize; i++ {
		b := util.ReverseBytes(bytes[offset:offset + 8])
		bits := binary.BigEndian.Uint64(b)
		offset += 8
		proto.ObjectPropValues = append(proto.ObjectPropValues, math.Float64frombits(bits))
	}

	/* []string */
	objetcPropKeysSize := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
	if objetcPropKeysSize + offset > len(bytes) {
		return errors.New("UnpackError")
	}
	offset += 2
	for i := 0; i < objetcPropKeysSize; i++ {
		size := int(binary.BigEndian.Uint16((bytes[offset:offset + 2])))
		offset += 2
		proto.ObjetcPropKeys = append(proto.ObjetcPropKeys, string(bytes[offset:offset + size]))
		offset += size
	}

	/* uint8 */
	proto.UpdateMode = uint8(bytes[offset])
	offset++


	return nil
}

func (proto *P2PUpdateRoomObject) String() string {
	list := make([]string, 0)
	list = append(list, fmt.Sprint("ObjectName = ", proto.ObjectName))
	list = append(list, fmt.Sprint("ObjectPropValues = ", proto.ObjectPropValues))
	list = append(list, fmt.Sprint("ObjetcPropKeys = ", proto.ObjetcPropKeys))
	list = append(list, fmt.Sprint("UpdateMode = ", proto.UpdateMode))
	return strings.Join(list, " | ")
}

func (proto *P2PUpdateRoomObject) GetVer() uint8 {
	return proto.Ver
}
func (proto *P2PUpdateRoomObject) GetCmd() uint16 {
	return proto.Cmd
}
