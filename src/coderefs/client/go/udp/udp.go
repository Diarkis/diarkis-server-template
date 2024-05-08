package udp

import (
	"time"
)

/*
Client UDP client
*/
type Client struct {
	ID               string
	IDBytes          []byte
	ServerEndPoint   string
	ClientKey        string
	SID              []byte
	SIDString        string
	EncryptionKey    []byte
	EncryptionIV     []byte
	EncryptionMacKey []byte
	IsPufferEnabled  bool
}

/*
LogLevel sets log level
*/
func LogLevel(lvl int) {
}

/*
New Creates a new UDP Client
*/
func New(rcvMaxSize int, interval int64) *Client {
	return nil
}

/*
EnablePuffer enables or disables Puffer
*/
func (cli *Client) EnablePuffer(isPufferEnabled bool) {
}

/*
SetID sets ID
*/
func (cli *Client) SetID(id string) {
}

/*
OnConnect registers a callback on Connect
*/
func (cli *Client) OnConnect(callback func()) {
}

/*
OnEcho registers a callback on Echo
Returns round trip time between the server and inside global IP address
*/
func (cli *Client) OnEcho(callback func(rtt time.Duration, address string)) {
}

/*
OnReconnect registers a callack on Reconnect
*/
func (cli *Client) OnReconnect(callback func()) {
}

/*
OnDisconnect registers a callback on Disconnect
*/
func (cli *Client) OnDisconnect(callback func()) {
}

/*
OnOffline registers a callback called when server gets Offline state
*/
func (cli *Client) OnOffline(callback func()) {
}

/*
IsOffline returns if the connected server is offline state
*/
func (cli *Client) IsOffline() bool {
	return false
}

/*
SetClientKey sets up client key
*/
func (cli *Client) SetClientKey(clientKey string) {
}

/*
SetEncryptionKeys Sets up encryption keys
*/
func (cli *Client) SetEncryptionKeys(sid []byte, key []byte, iv []byte, macKey []byte) {
}

/*
Connect connects UDP Client
*/
func (cli *Client) Connect(addr string) {
}

/*
Reconnect restarts UDP Client
*/
func (cli *Client) Reconnect(addr string) {
}

/*
Disconnect disconnects the UDP Client
*/
func (cli *Client) Disconnect() {
}

/*
Die stops server communication ungracefully
*/
func (cli *Client) Die() {
}

/*
Ping sends a ping packet over UDP
*/
func (cli *Client) Ping() {
}

/*
Send Sends a UDP packet to server
*/
func (cli *Client) Send(ver uint8, cmd uint16, payload []byte) {
}

/*
RSend Sends a Reliable UDP packet to server
*/
func (cli *Client) RSend(ver uint8, cmd uint16, payload []byte) {
}

/*
OnResponse Registers a callback on response receive
*/
func (cli *Client) OnResponse(callback func(uint8, uint16, uint8, []byte)) {
}

/*
OnPush Registers a callback on push receive
*/
func (cli *Client) OnPush(callback func(uint8, uint16, []byte)) {
}

/*
CatchOnReconnect registers a callback on reconnect or response w/ ver and cmd of your choice ONCE
*/
func (cli *Client) CatchOnReconnect(ver uint8, cmd uint16, callback func(bool)) {
}

/*
RemoveOnReconnect removes a callback function on reconnect event
*/
func (cli *Client) RemoveOnReconnect(callbackToRemove func()) {
}

/*
RemoveOnResponse removes a callback function on response event
*/
func (cli *Client) RemoveOnResponse(callbackToRemove func(uint8, uint16, uint8, []byte)) {
}

/*
RemoveAllOnResponse Removes all callbacks of OnResponse
*/
func (cli *Client) RemoveAllOnResponse() {
}

/*
RemoveOnPush removes a callback function on push event
*/
func (cli *Client) RemoveOnPush(callbackToRemove func(uint8, uint16, []byte)) {
}

/*
RemoveAllOnPush Removes all callbacks of OnPush
*/
func (cli *Client) RemoveAllOnPush() {
}

/*
RemoveAll Removes all callbacks
*/
func (cli *Client) RemoveAll() {
}
