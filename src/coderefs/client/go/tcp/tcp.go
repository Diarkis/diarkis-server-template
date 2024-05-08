package tcp

/*
Client TCP Client
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
New Creates a new TCP Client
*/
func New(rcvMaxSize int, interval int64, hbInterval int64) *Client {
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
OnConnect registers a callback on connect
*/
func (cli *Client) OnConnect(callback func()) {
}

/*
OnHeartbeat registers a callback on heartbeat
*/
func (cli *Client) OnHeartbeat(callback func()) {
}

/*
OnDisconnect registers a callback on disconnect
*/
func (cli *Client) OnDisconnect(callback func()) {
}

/*
OnReconnect registers a callback on reconnect
*/
func (cli *Client) OnReconnect(callback func()) {
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
Connect Starts TCP Client
*/
func (cli *Client) Connect(addr string) {
}

/*
Reconnect restarts TCP Client
*/
func (cli *Client) Reconnect(addr string) {
}

/*
Die stops server communication ungracefully
*/
func (cli *Client) Die() {
}

/*
Disconnect Stops TCP Client
*/
func (cli *Client) Disconnect() {
}

/*
Send Sends a TCP packet to server
*/
func (cli *Client) Send(ver uint8, cmd uint16, payload []byte) {
}

/*
PauseHeartbeat Stops sending heartbeat command to the server
*/
func (cli *Client) PauseHeartbeat() {
}

/*
ResumeHeartbeat Resumes sending heartbeat command to the server
*/
func (cli *Client) ResumeHeartbeat() {
}

/*
OnResponse Registers a callback on response reveive
*/
func (cli *Client) OnResponse(callback func(uint8, uint16, uint8, []byte)) {
}

/*
OnPush Registers a callback on response reveive
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
RemoveAllOnResponse Removes all callbacks on OnResponse
*/
func (cli *Client) RemoveAllOnResponse() {
}

/*
RemoveAllOnPush Removes all callbacks on OnPush
*/
func (cli *Client) RemoveAllOnPush() {
}

/*
RemoveAll Removes all callbacks on OnResponse and OnPush
*/
func (cli *Client) RemoveAll() {
}
