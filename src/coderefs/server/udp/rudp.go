package udp

import (
	"sync"
)

/*
Connection RUDP connection struct
*/
type Connection struct{ sync.RWMutex }

/*
GetConnectionTTL returns connection TTL
*/
func GetConnectionTTL() int64 {
	return 0
}

/*
OnDisconnect Registers a callback function on RUDP disconnect
*/
func OnDisconnect(callback func(string, string, string), last bool) {
}

/*
OnNewConnection [INTERNAL USE ONLY]
*/
func OnNewConnection(callback func(string) bool) {
}

/*
 && outs == 0

*/

/*
GetConnection returns RUDP connection struct
*/
func (state *State) GetConnection() *Connection {
	return nil
}

/*
Disconnect Disconnects UDP client (RUDP) from the server
*/
func (state *State) Disconnect() {
}

/*
RSend Reliable UDP send with status
We turn outs bytes into one RUDP packet when we actually send it
it means that all combined packets are considered "ONE" RUDP packet
*/
func (state *State) RSend(conn *Connection, ver uint8, cmd uint16, payload []byte, status uint8) {
}

/*
RPush Reliable UDP push from the server
Same as RSend, we buffer the push packets and combine them into one packet at certain interval
*/
func (state *State) RPush(conn *Connection, ver uint8, cmd uint16, payload []byte) {
}
