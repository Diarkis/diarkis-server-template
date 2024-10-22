// © 2019-2024 Diarkis Inc. All rights reserved.

// Code generated by Diarkis Puffer module: DO NOT EDIT.
//
// # Auto-generated by Diarkis Version 1.1.0
//
// - Maximum length of a string is 65535 bytes
// - Maximum length of a byte array is 65535 bytes
// - Maximum length of any array is 65535 elements
package custom

import "encoding/binary"
import "errors"
import "fmt"

//???
import "strings"

//???

// DiarkisCharacterFrameDataVer represents the ver of the protocol's command.
//
//	[NOTE] The value is optional and if ver is not given in the definition JSON, it will be 0.
const DiarkisCharacterFrameDataVer uint8 = 0

// DiarkisCharacterFrameDataCmd represents the command ID of the protocol's command ID.
//
//	[NOTE] The value is optional and if cmd is not given in the definition JSON, it will be 0.
const DiarkisCharacterFrameDataCmd uint16 = 0

// DiarkisCharacterFrameData represents the command protocol data structure.
type DiarkisCharacterFrameData struct {
	// Command version of the protocol
	Ver uint8
	// Command ID of the protocol
	Cmd               uint16
	AnimationBlend    uint8
	AnimationID       uint8
	Position          *DiarkisVector3
	RotationAngles    uint16
	TimestampInterval uint16
}

// NewDiarkisCharacterFrameData creates a new instance of DiarkisCharacterFrameData struct.
func NewDiarkisCharacterFrameData() *DiarkisCharacterFrameData {
	return &DiarkisCharacterFrameData{Ver: 0, Cmd: 0, AnimationBlend: 0, AnimationID: 0, Position: NewDiarkisVector3(), RotationAngles: 0, TimestampInterval: 0}
}

// Pack encodes DiarkisCharacterFrameData struct to a byte array to be delivered over the command.
func (proto *DiarkisCharacterFrameData) Pack() []byte {
	bytes := make([]byte, 0)

	/* uint8 */
	bytes = append(bytes, proto.AnimationBlend)

	/* uint8 */
	bytes = append(bytes, proto.AnimationID)

	/* DiarkisVector3 */
	positionSizeBytes := make([]byte, 2)
	positionPacked := proto.Position.Pack()
	binary.BigEndian.PutUint16(positionSizeBytes, uint16(len(positionPacked)))
	bytes = append(bytes, positionSizeBytes...)
	bytes = append(bytes, positionPacked...)

	/* uint16 */
	bytes = binary.BigEndian.AppendUint16(bytes, uint16(proto.RotationAngles))

	/* uint16 */
	bytes = binary.BigEndian.AppendUint16(bytes, uint16(proto.TimestampInterval))

	// done
	return bytes
}

// Unpack decodes the command payload byte array to DiarkisCharacterFrameData struct.
func (proto *DiarkisCharacterFrameData) Unpack(bytes []byte) error {
	if len(bytes) < 8 {
		return errors.New("DiarkisCharacterFrameDataUnpackError")
	}

	offset := 0

	/* uint8 */
	proto.AnimationBlend = uint8(bytes[offset])
	offset++

	/* uint8 */
	proto.AnimationID = uint8(bytes[offset])
	offset++

	/* DiarkisVector3 */
	positionSize := int(binary.BigEndian.Uint16((bytes[offset : offset+2])))
	if positionSize+offset > len(bytes) {
		return errors.New("UnpackError")
	}
	offset += 2
	positionBytes := bytes[offset : offset+positionSize]
	proto.Position = &DiarkisVector3{Ver: 0, Cmd: 0}
	proto.Position.Unpack(positionBytes)
	offset += positionSize

	/* uint16 */
	proto.RotationAngles = binary.BigEndian.Uint16(bytes[offset : offset+2])
	offset += 2

	/* uint16 */
	proto.TimestampInterval = binary.BigEndian.Uint16(bytes[offset : offset+2])
	offset += 2

	return nil
}

func (proto *DiarkisCharacterFrameData) String() string {
	list := make([]string, 0)
	list = append(list, fmt.Sprint("AnimationBlend = ", proto.AnimationBlend))
	list = append(list, fmt.Sprint("AnimationID = ", proto.AnimationID))
	list = append(list, fmt.Sprint("Position = ", proto.Position.String()))
	list = append(list, fmt.Sprint("RotationAngles = ", proto.RotationAngles))
	list = append(list, fmt.Sprint("TimestampInterval = ", proto.TimestampInterval))
	return strings.Join(list, " | ")
}

func (proto *DiarkisCharacterFrameData) GetVer() uint8 {
	return proto.Ver
}
func (proto *DiarkisCharacterFrameData) GetCmd() uint16 {
	return proto.Cmd
}
