package dm

import (
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

/*
OnResponse is a callback function for on response event
*/
type OnResponse func(success bool, payload []byte)

/*
OnPush is a callback function for on push event
*/
type OnPush func(uid string, message []byte)

/*
DirectMessage represents Diarkis DirectMessage client
*/
type DirectMessage struct{}

/*
SetupAsTCP sets up Diarkis DirectMessage client as TCP
*/
func (dm *DirectMessage) SetupAsTCP(tcp *tcp.Client) bool {
	return false
}

/*
SetupAsUDP sets up Diarkis DirectMessage client as UDP
*/
func (dm *DirectMessage) SetupAsUDP(udp *udp.Client) bool {
	return false
}

/*
OnSend assigns a callback on send response event
*/
func (dm *DirectMessage) OnSend(cb func(success bool, payload []byte)) {
}

/*
OnDisconnect assigns a callback on disconnect response event
*/
func (dm *DirectMessage) OnDisconnect(cb func(success bool, payload []byte)) {
}

/*
OnPeerSend assigns a callback on peer send event
*/
func (dm *DirectMessage) OnPeerSend(cb func(uid string, message []byte)) {
}

/*
OnPeerDisconnect assigns a callback on peer disconnect event
*/
func (dm *DirectMessage) OnPeerDisconnect(cb func(uid string, message []byte)) {
}

/*
Send sends out a message to the peer
*/
func (dm *DirectMessage) Send(uid string, message []byte, reliable bool) {
}

/*
Disconnect sends a disconnect message to the peer
*/
func (dm *DirectMessage) Disconnect(uid string, message []byte) {
}
