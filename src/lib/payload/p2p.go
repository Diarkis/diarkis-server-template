// © 2019-2024 Diarkis Inc. All rights reserved.

package payload

import (
	"encoding/binary"
	"github.com/Diarkis/diarkis/packet"
)

const dataType uint16 = 1000
const ver uint16 = 1
const messageType uint8 = 0
const headerSize = 5

func PackP2PReport(addr string) []byte {
	header := packHeader()
	return append(header, []byte(addr)...)
}

func UnpackP2PReport(bytes []byte) string {
	if len(bytes) <= headerSize {
		return ""
	}
	return string(bytes[headerSize:])
}

func PackP2PInit(list []string) []byte {
	header := packHeader()
	bytes := packet.StringListToBytes(list)
	return append(header, bytes...)
}

func UnpackP2PInit(bytes []byte) []string {
	if len(bytes) <= headerSize {
		return make([]string, 0)
	}
	return packet.BytesToStringList(bytes[headerSize:])
}

func packHeader() []byte {
	twoBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(twoBytes, dataType)
	bytes := make([]byte, 0)
	bytes = append(bytes, twoBytes...)
	binary.BigEndian.PutUint16(twoBytes, ver)
	bytes = append(bytes, twoBytes...)
	typeBytes := make([]byte, 1)
	typeBytes[0] = messageType
	bytes = append(bytes, typeBytes...)
	return bytes
}
