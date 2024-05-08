package field

import (
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

/*
Field represents Diarkis Field client
*/
type Field struct{}

/*
NewFieldAsTCP creates a Diarkis Field client as TCP
*/
func NewFieldAsTCP(tcpClient *tcp.Client) *Field {
	return nil
}

/*
NewFieldAsUDP creates a Diarkis Field client as UDP
*/
func NewFieldAsUDP(udpClient *udp.Client) *Field {
	return nil
}

/*
SetupAsTCP sets up Diarkis Field client as TCP
*/
func (f *Field) SetupAsTCP(tcpClient *tcp.Client) bool {
	return false
}

/*
SetupAsUDP sets up Diarkis Field client as UDP
*/
func (f *Field) SetupAsUDP(udpClient *udp.Client) bool {
	return false
}

/*
Sync sends sync data to the server.
*/
func (f *Field) Sync(x, y, z int64, syncLimit uint16, customFilterID uint8, msg []byte, reliable bool) {
}

/*
Disappear sends out disappear to users in view
*/
func (f *Field) Disappear() {
}

/*
OnSync assigns a callback to sync event
*/
func (f *Field) OnSync(cb func([]byte)) {
}

/*
OnReconnect assigns a callback to re-connect event
*/
func (f *Field) OnReconnect(cb func()) {
}

/*
OnDisappear assigns a callback to disappear event
*/
func (f *Field) OnDisappear(cb func(string)) {
}

/*
OnServerSync deprecated
*/
func (f *Field) OnServerSync(cb func(bool, []byte)) {
}
