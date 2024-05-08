package tcp

import (
	"net"
)

const Type = "TCP"
const Ok = 1
const Bad = 4
const Err = 5

/*
State Client state data structure
*/
type State struct {
	UserID           string
	SID              string
	Conn             *net.TCPConn
	ClientAddr       string
	TTLDuration      int64
	TTL              int64
	Connected        bool
	Version          uint8
	CommandID        uint16
	Payload          []byte
	EncryptionKey    []byte
	EncryptionIV     []byte
	EncryptionMacKey []byte
	ClientKey        string
}

/*
Init [INTERNAL USE ONLY] Flag the TCP server state to be initialized
to prevent registration of commands and hooks BEFORE initialization by user package
*/
func Init() {
}

/*
Setup [INTERNAL USE ONLY] Loads configuration file into memory - pass an empty string to load nothing
*/
func Setup(confpath string) {
}

/*
GetConnectionTTL returns connection TTL
*/
func GetConnectionTTL() int64 {
	return 0
}

/*
IsEncryptionEnabled returns true if encryption/decryption is enabled
*/
func IsEncryptionEnabled() bool {
	return false
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
IsCommandDuplicate returns true if ver and command ID are used elsewhere when this function is called
*/
func IsCommandDuplicate(ver uint8, commandID uint16) bool {
	return false
}

/*
Command Registers a packet handler as a command
*/
func Command(ver uint8, cmd uint16, handlerName string, handler func(*State, func(error))) {
}

/*
OnHeartbeat Registers a handler to be executed on heartbeat
*/
func OnHeartbeat(handler func(*State, func(error))) {
}

/*
OnNewConnection Registers a callback on new TCP connection w/ a client
*/
func OnNewConnection(callback func(*State)) {
}

/*
OnDisconnect Registers a callback on connection disconnect
*/
func OnDisconnect(callback func(string, string, string, error), last bool) {
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
SetPublicEndPoint sets public end point to be used by the client to connect
*/
func SetPublicEndPoint(addr string) {
}

/*
ShutdownTCP Stops TCP server
*/
func ShutdownTCP(next func(error)) {
}

/*
GetEndPoint Returns TCP server end point that is actually bound
*/
func GetEndPoint() string {
	return ""
}

/*
GetStateByAddr Returns a client state struct by client address
*/
func GetStateByAddr(addr string) *State {
	return nil
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

/*
Disconnect Disconnects the socket connection to the client
*/
func (state *State) Disconnect() {
}

/*
Kill Forcefully terminates the socket connection to the client
*/
func (state *State) Kill(err error) {
}
