package udp

import (
	"net"
)

const Type = "UDP"
const Ok = 1
const Bad = 4
const Err = 5

/*
State Data structure of a client
*/
type State struct {
	UserID                   string
	SID                      string
	ClientAddr               string
	ClientNetAddr            net.Addr
	Version                  uint8
	CommandID                uint16
	Payload                  []byte
	EncryptionKey            []byte
	EncryptionIV             []byte
	EncryptionMacKey         []byte
	ClientKey                string
	P2PEnabled               bool
	ClientLocalAddrListBytes []byte
}

/*
Init [INTERNAL USE ONLY] Flag the UDP server state to be initialized
to prevent registration of commands and hooks BEFORE initialization by user package
*/
func Init() {
}

/*
Setup [INTERNAL USE ONLY] Loads a configuration file into memory - pass an empty string to load nothing
*/
func Setup(confpath string) {
}

/*
IsP2PEnabled returns true if P2P is enabled by the configuration
*/
func IsP2PEnabled() bool {
	return false
}

/*
IsEncryptionEnabled returns true if encryption/decryption is enabled
*/
func IsEncryptionEnabled() bool {
	return false
}

/*
Online flag the server to be online to take new users
*/
func Online() {
}

/*
Taken flag the server to be taken and will not take new users
*/
func Taken() {
}

/*
Offline flag the serve to be offline to be deleted and will not take new users
*/
func Offline() {
}

/*
GetEndPoint Returns its end point that is actually bound
*/
func GetEndPoint() string {
	return ""
}

/*
SetPublicEndPoint sets public end point to be used by the client to connect
*/
func SetPublicEndPoint(addr string) {
}

/*
HookAllCommands Registers a packet handler as a hook to all commands
*/
func HookAllCommands(handlerName string, handler func(*State, func(error))) {
}

/*
HookAllWrites Registers a handler function to all TCP socket writes:
Receives the payload and returns a payload
*/
func HookAllWrites(handler func(*State, []byte) ([]byte, error)) {
}

/*
OnEcho Registers a handler to be executed on echo
*/
func OnEcho(handler func(*State, func(error))) {
}

/*
IsCommandDuplicate returns true if ver and command ID are used elsewhere when this function is called
*/
func IsCommandDuplicate(ver uint8, commandID uint16) bool {
	return false
}

/*
Command Registers a UDP packet handler as a command
*/
func Command(ver uint8, cmd uint16, handlerName string, handler func(*State, func(error))) {
}

/*
ShutdownUDP Stops UDP server
*/
func ShutdownUDP() {
}

/*
GetSeqForDup returns sequence for packet duplicate detection.
*/
func (state *State) GetSeqForDup() uint32 {
	return 0
}

/*
GetClientAddr returns client public address. If it is UDP server and enableP2P is false, it returns 0.0.0.0:0 instead.
*/
func (state *State) GetClientAddr() string {
	return ""
}

/*
Send Sends a response packet to the client
*/
func (state *State) Send(ver uint8, cmd uint16, payload []byte, status uint8) (int, error) {
	return 0, nil
}

/*
Push Sends a push packet to the client
*/
func (state *State) Push(ver uint8, cmd uint16, payload []byte) (int, error) {
	return 0, nil
}
