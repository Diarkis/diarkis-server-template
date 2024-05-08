package debug

import (
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

/*
Debug represents the data structure of Debug
*/
type Debug struct{}

/*
NewDebugAsTCP creates a new Debug instance for TCP client
*/
func NewDebugAsTCP(tcp *tcp.Client) *Debug {
	return nil
}

/*
NewDebugAsUDP creates a new Debug instance for UDP client
*/
func NewDebugAsUDP(udp *udp.Client) *Debug {
	return nil
}

/*
Online sets the server status to ONLINE
*/
func (d *Debug) Online() {
}

/*
Taken sets the server status to TAKEN
*/
func (d *Debug) Taken() {
}

/*
Offline sets the server status to OFFLINE
*/
func (d *Debug) Offline() {
}

/*
Terminate instructs the server to stop
*/
func (d *Debug) Terminate() {
}
