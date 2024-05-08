package packet

import (
	"sync"
)

const UDPProtoHeaderSize int = 4
const HeaderSize int = 10
const StatusOk uint8 = 0x00
const StatusBad uint8 = 0x01
const StatusPush uint8 = 0xff

var UDPProto uint8 = 0x1
var RUDPProtoSyn uint8 = 0x2
var RUDPProtoDat uint8 = 0x3
var RUDPProtoAck uint8 = 0x4
var RUDPProtoRst uint8 = 0x5
var RUDPProtoEack uint8 = 0x6
var RUDPProtoFin uint8 = 0x7

/*
SplitPacket a packet that are split into smaller chunks
*/
type SplitPacket struct {
	ID uint16
	sync.RWMutex
}

/*
RequestHeader Request header data structure
*/
type RequestHeader struct {
	Version     uint8
	CommandID   uint16
	PayloadSize uint32
}

/*
ResponseHeader Response header data structure
*/
type ResponseHeader struct {
	Version     uint8
	CommandID   uint16
	PayloadSize uint32
	Status      uint8
}

/*
RequestPacket Request packet data structure
*/
type RequestPacket struct {
	Header    *RequestHeader
	Payload   []byte
	SeqForDup uint32
}

/*
ResponsePacket Response packet data structure
*/
type ResponsePacket struct {
	Header  *ResponseHeader
	Payload []byte
	Push    bool
}

/*
BoolToBytes converts bool to bytes
*/
func BoolToBytes(val bool) []byte {
	return nil
}

/*
BytesToBool converts bytes to bool
*/
func BytesToBool(bytes []byte) bool {
	return false
}

/*
NewSplitPacket creates a split packet receiver
*/
func NewSplitPacket(bytes []byte) *SplitPacket {
	return nil
}

/*
IsSplitPacket returns true if the evaluated bytes is a split packet chunk
*/
func IsSplitPacket(bytes []byte) bool {
	return false
}

/*
Add adds split packet chunk to split packet receiver
*/
func (sp *SplitPacket) Add(bytes []byte) error {
	return nil
}

/*
ConsumeBytes returns the reconstructed split bytes
*/
func (sp *SplitPacket) ConsumeBytes() ([]byte, bool) {
	return nil, false
}

/*
GetSplitPacketID returns the ID of split packet
*/
func GetSplitPacketID(bytes []byte) uint16 {
	return 0
}

/*
CreateSplitPacket creates an array of split packets
*/
func CreateSplitPacket(id uint16, bytes []byte, splitSize int) [][]byte {
	return nil
}

/*
IsInvalidPacketErr Returns true if the given error is an invalid packet error
*/
func IsInvalidPacketErr(err error) bool {
	return false
}

/*
IsPushPacket Returns true if the given packet status is a push packet
*/
func IsPushPacket(status uint8) bool {
	return false
}

/*
StringListToBytesWithSize Converts an array of strings to byte array with a specified header size

It takes two parameters:

 1. list, slice of string that you want to convert to bytes; and
 2. headerSize, how many bytes it uses for its header;

And it returns two values:

 1. converted byte list;
 3. error if it fails to convert.
*/
func StringListToBytesWithSize(list []string, headerSize int) ([]byte, error) {
	return nil, nil
}

/*
BytesToStringListWithSize Converts a byte array to an array of strings with a specified header size
bytes need to be exact size for the parsed list;

It takes two parameters:

 1. bytes, byte list that includes header byte for each which indicates its size; and
 2. headerSize, how many bytes it uses for its header.

And it returns two values:

 1. parsed string list; and
 2. error if it fails to convert.
*/
func BytesToStringListWithSize(bytes []byte, headerSize int) ([]string, error) {
	return nil, nil
}

/*
BytesToStringListWithLength Converts a byte array to an array of strings with a specified header size and array length
bytes do not need to be exact size, it is allowed to have extra bytes after the list

It takes three parameters:

 1. bytes, byte list that includes header byte for each which indicates its size;
 2. headerSize, how many bytes it uses for its header; and
 3. length, recurring time.

And it returns three values:

 1. parsed string list;
 2. total size after parsing ; and
 3. error if it fails to convert.
*/
func BytesToStringListWithLength(bytes []byte, headerSize int, length int) ([]string, int, error) {
	return nil, 0, nil
}

/*
StringListToBytes Converts an array of strings to byte array
*/
func StringListToBytes(list []string) []byte {
	return nil
}

/*
BytesToStringList Converts a byte array to an array of strings
*/
func BytesToStringList(bytes []byte) []string {
	return nil
}

/*
BytesListToBytes converts an array of byte array to byte array
*/
func BytesListToBytes(list [][]byte) []byte {
	return nil
}

/*
BytesToBytesList converts a byte array to an array of byte array
*/
func BytesToBytesList(bytes []byte) [][]byte {
	return nil
}

/*
BytesToBytesListMin converts a byte array to an array of byte array(used by state sync)
*/
func BytesToBytesListMin(bytes []byte) [][]byte {
	return nil
}

/*
BytesToFloat64Map converts a byte array to a map[string]float64
*/
func BytesToFloat64Map(bytes []byte) map[string]float64 {
	return nil
}

/*
Float64MapToBytes converts a map[string]float64 to a byte array
*/
func Float64MapToBytes(m map[string]float64) []byte {
	return nil
}

/*
BytesToFloat64 converts a byte array to a float64
*/
func BytesToFloat64(bytes []byte) float64 {
	return 0
}

/*
Float64ToBytes converts a float64 to a byte array
*/
func Float64ToBytes(v float64) []byte {
	return nil
}

/*
ParseReconnectPayload converts the given payload to reconnect address
*/
func ParseReconnectPayload(payload []byte) string {
	return ""
}

/*
CreateReconnectPayload Creates a payload to instruct the client to reconnect
*/
func CreateReconnectPayload(addr string) []byte {
	return nil
}

/*
CreateUDPPacket Creates a UDP packet from the byte array created by

	CreateRequestPacket(), CreateResponsePacket(), and CreatePushPacket()
*/
func CreateUDPPacket(flag uint8, seq uint32, packet []byte) ([]byte, error) {
	return nil, nil
}

/*
CreateSecureRequestPayload Creates an encrypted request packet payload
*/
func CreateSecureRequestPayload(sid []byte, key []byte, iv []byte, mackey []byte, payload []byte) ([]byte, error) {
	return nil, nil
}

/*
CreateSecureResponsePayload Creates an encrypted response packet payload
*/
func CreateSecureResponsePayload(key []byte, iv []byte, mackey []byte, payload []byte) ([]byte, error) {
	return nil, nil
}

/*
CreateRequestPacket Creates a request packet
*/
func CreateRequestPacket(version uint8, commandID uint16, payload []byte) []byte {
	return nil
}

/*
CreateResponsePacket Creates a response packet - by giving StatusPush as status, it creates a push packet
*/
func CreateResponsePacket(version uint8, commandID uint16, status uint8, payload []byte) []byte {
	return nil
}

/*
CreatePushPacket Creates a push packet
*/
func CreatePushPacket(version uint8, commandID uint16, payload []byte) []byte {
	return nil
}

/*
GetSidFromPayload Returns sid (session ID) and encrypted payload from the given payload
*/
func GetSidFromPayload(payload []byte) ([]byte, []byte) {
	return nil, nil
}

/*
ParseUDPPacket Parses a packet created by CreateUDPPacket
*/
func ParseUDPPacket(packet []byte) (uint8, uint32, []byte, error) {
	return 0, 0, nil, nil
}

/*
ParseRequestPacket Parses a request packet
*/
func ParseRequestPacket(packet []byte) (*RequestPacket, int, error) {
	return nil, 0, nil
}

/*
ParseResponsePacket Parses a response packet
*/
func ParseResponsePacket(packet []byte) (*ResponsePacket, int, error) {
	return nil, 0, nil
}
