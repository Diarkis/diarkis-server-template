// © 2019-2024 Diarkis Inc. All rights reserved.

package payload

import (
	"encoding/binary"
	"github.com/Diarkis/diarkis/packet"
)

type MMAdd struct {
	TTL      int64
	ID       string
	UID      string
	Props    map[string]int
	Metadata []byte
}

type MMRemove struct {
	ID   string
	UIDs []string
}

type MMSearch struct {
	IDs   []string
	Props map[string]int
}

func PackMMAdd(mmID, uniqueID string, props map[string]int, metadata []byte, ttl uint64) []byte {
	res := make([]byte, 0)
	sizeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sizeBytes[0:8], ttl)
	res = append(res, sizeBytes[0:8]...)
	list := make([]string, 2)
	list[0] = mmID
	list[1] = uniqueID
	bytes := packet.StringListToBytes(list)
	binary.BigEndian.PutUint16(sizeBytes[0:2], uint16(len(bytes)))
	res = append(res, sizeBytes[0:2]...)
	res = append(res, bytes...)
	bytes = make([]byte, 0)
	for name, val := range props {
		binary.BigEndian.PutUint32(sizeBytes[0:4], uint32(val))
		bytes = append(bytes, sizeBytes[0:4]...)
		binary.BigEndian.PutUint16(sizeBytes[0:2], uint16(len([]byte(name))))
		bytes = append(bytes, sizeBytes[0:2]...)
		bytes = append(bytes, []byte(name)...)
	}
	binary.BigEndian.PutUint16(sizeBytes[0:2], uint16(len(bytes)))
	res = append(res, sizeBytes[0:2]...)
	res = append(res, bytes...)
	res = append(res, metadata...)
	return res
}

func UnpackMMAdd(bytes []byte) *MMAdd {
	if len(bytes) <= 20 {
		return nil
	}
	offset := 10
	ttl := int64(binary.BigEndian.Uint64(bytes[0:8]))
	size := int(binary.BigEndian.Uint16(bytes[8:10]))
	offset += size
	list := packet.BytesToStringList(bytes[10:offset])
	if len(list) < 2 {
		return nil
	}
	mmID := list[0]
	uniqueID := list[1]
	props := make(map[string]int)
	size = int(binary.BigEndian.Uint16(bytes[offset : offset+2]))
	offset += 2
	propsBytes := bytes[offset : offset+size]
	offset += size
	index := 0
	for index < size {
		val := int(binary.BigEndian.Uint32(propsBytes[index : index+4]))
		index += 4
		nameSize := int(binary.BigEndian.Uint16(propsBytes[index : index+2]))
		index += 2
		propName := string(propsBytes[index : index+nameSize])
		props[propName] = val
		index += nameSize
	}
	mmAdd := new(MMAdd)
	mmAdd.TTL = ttl
	mmAdd.ID = mmID
	mmAdd.UID = uniqueID
	mmAdd.Props = props
	mmAdd.Metadata = bytes[offset:]
	return mmAdd
}

func PackMMRemove(mmID string, uniqueIDs []string) []byte {
	list := make([]string, 1)
	list[0] = mmID
	list = append(list, uniqueIDs...)
	return packet.StringListToBytes(list)
}

func UnpackMMRemove(bytes []byte) *MMRemove {
	if len(bytes) <= 4 {
		return nil
	}
	list := packet.BytesToStringList(bytes)
	// must contain mmID and at least one uniqueID
	if len(list) < 2 {
		return nil
	}
	// first item in the list is mmID
	mmRemove := new(MMRemove)
	mmRemove.ID = list[0]
	mmRemove.UIDs = list[1:]
	return mmRemove
}

func PackMMSearch(mmIDs []string, props map[string]int) []byte {
	res := packet.StringListToBytes(mmIDs)
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint16(sizeBytes[0:2], uint16(len(res)))
	res = append(sizeBytes[0:2], res...)
	bytes := make([]byte, 0)
	for name, val := range props {
		binary.BigEndian.PutUint32(sizeBytes[0:4], uint32(val))
		bytes = append(bytes, sizeBytes[0:4]...)
		binary.BigEndian.PutUint16(sizeBytes[0:2], uint16(len([]byte(name))))
		bytes = append(bytes, sizeBytes[0:2]...)
		bytes = append(bytes, []byte(name)...)
	}
	res = append(res, bytes...)
	return res
}

func UnpackMMSearch(bytes []byte) *MMSearch {
	if len(bytes) <= 4 {
		return nil
	}
	size := int(binary.BigEndian.Uint16(bytes[0:2]))
	matchingIDs := packet.BytesToStringList(bytes[2 : 2+size])
	propsBytes := bytes[2+size:]
	index := 0
	props := make(map[string]int)
	for index < len(propsBytes) {
		val := int(binary.BigEndian.Uint32(propsBytes[index : index+4]))
		index += 4
		nameSize := int(binary.BigEndian.Uint16(propsBytes[index : index+2]))
		index += 2
		propName := string(propsBytes[index : index+nameSize])
		props[propName] = val
		index += nameSize
	}
	mmSearch := new(MMSearch)
	mmSearch.IDs = matchingIDs
	mmSearch.Props = props
	return mmSearch
}
