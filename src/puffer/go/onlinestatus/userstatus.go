// © 2019-2024 Diarkis Inc. All rights reserved.

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

// UserStatusVer represents the ver of the protocol's command.
//
//	[NOTE] The value is optional and if ver is not given in the definition JSON, it will be 0.
const UserStatusVer uint8 = 0

// UserStatusCmd represents the command ID of the protocol's command ID.
//
//	[NOTE] The value is optional and if cmd is not given in the definition JSON, it will be 0.
const UserStatusCmd uint16 = 0

// UserStatus represents the command protocol data structure.
type UserStatus struct {
	// Command version of the protocol
	Ver uint8
	// Command ID of the protocol
	Cmd         uint16
	InRoom      bool
	SessionData []*UserSessionData
	UID         string
}

// NewUserStatus creates a new instance of UserStatus struct.
func NewUserStatus() *UserStatus {
	return &UserStatus{Ver: 0, Cmd: 0, InRoom: false, SessionData: make([]*UserSessionData, 0), UID: ""}
}

// Pack encodes UserStatus struct to a byte array to be delivered over the command.
func (proto *UserStatus) Pack() []byte {
	bytes := make([]byte, 0)

	/* bool */
	inRoomBytes := make([]byte, 1)
	if proto.InRoom {
		inRoomBytes[0] = uint8(1) /* true */
	} else {
		inRoomBytes[0] = uint8(2) /* false */
	}
	bytes = append(bytes, inRoomBytes...)

	/* []UserSessionData */
	sessionDataLengthBytes := make([]byte, 2)
	sessionDataLength := len(proto.SessionData)
	binary.BigEndian.PutUint16(sessionDataLengthBytes, uint16(sessionDataLength))
	bytes = append(bytes, sessionDataLengthBytes...)
	for i := 0; i < sessionDataLength; i++ {
		sessionDataSizeBytes := make([]byte, 2)
		sessionDataPacked := proto.SessionData[i].Pack()
		binary.BigEndian.PutUint16(sessionDataSizeBytes, uint16(len(sessionDataPacked)))
		bytes = append(bytes, sessionDataSizeBytes...)
		bytes = append(bytes, sessionDataPacked...)
	}

	/* string */
	uIDSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(uIDSizeBytes, uint16(len(proto.UID)))
	bytes = append(bytes, uIDSizeBytes...)
	bytes = append(bytes, []byte(proto.UID)...)

	// done
	return bytes
}

// Unpack decodes the command payload byte array to UserStatus struct.
func (proto *UserStatus) Unpack(bytes []byte) error {
	if len(bytes) < 5 {
		return errors.New("UserStatusUnpackError")
	}

	offset := 0

	/* bool */
	inRoomByte := bytes[offset]
	if inRoomByte == 0x01 {
		proto.InRoom = true
	} else {
		proto.InRoom = false
	}
	offset++

	/* []UserSessionData */
	sessionDataLength := int(binary.BigEndian.Uint16((bytes[offset : offset+2])))
	offset += 2
	proto.SessionData = make([]*UserSessionData, sessionDataLength)
	for i := 0; i < sessionDataLength; i++ {
		sessionDataSize := int(binary.BigEndian.Uint16((bytes[offset : offset+2])))
		if sessionDataSize+offset > len(bytes) {
			return errors.New("UnpackError")
		}
		offset += 2
		sessionDataBytes := bytes[offset : offset+sessionDataSize]
		item := &UserSessionData{Ver: 0, Cmd: 0}
		item.Unpack(sessionDataBytes)
		proto.SessionData[i] = item
		offset += sessionDataSize
	}

	/* string */
	uIDSize := int(binary.BigEndian.Uint16(bytes[offset : offset+2]))
	if uIDSize+offset > len(bytes) {
		return errors.New("UnpackError")
	}
	offset += 2
	proto.UID = string(bytes[offset : offset+uIDSize])
	offset += uIDSize

	return nil
}

func (proto *UserStatus) String() string {
	list := make([]string, 0)
	list = append(list, fmt.Sprint("InRoom = ", proto.InRoom))
	for i, item := range proto.SessionData {
		list = append(list, fmt.Sprint("SessionData[", i, "] = ", "[", item.String(), "]"))
	}
	list = append(list, fmt.Sprint("UID = ", proto.UID))
	return strings.Join(list, " | ")
}

func (proto *UserStatus) GetVer() uint8 {
	return proto.Ver
}
func (proto *UserStatus) GetCmd() uint16 {
	return proto.Cmd
}
