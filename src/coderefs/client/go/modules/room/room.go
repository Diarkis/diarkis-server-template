package room

import (
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

/*
Room represents Diarkis Room client
*/
type Room struct{ ID string }

/*
ListItem represents a room list item
This is used for room finding
*/
type ListItem struct{}

/*
ChatData represents a chat data
This is used for chat log
*/
type ChatData struct {
	SenderUID string
	Timestamp int64
	Message   string
}

/*
SetupAsTCP sets up Room client as TCP client
*/
func (room *Room) SetupAsTCP(tcpClient *tcp.Client) bool {
	return false
}

/*
SetupAsUDP sets up Room client as UDP client
*/
func (room *Room) SetupAsUDP(udpClient *udp.Client) bool {
	return false
}

/*
SetupOnJoinEvent has been deprecated
*/
func (room *Room) SetupOnJoinEvent(ver uint8, cmd uint16) {
}

/*
SetupOnOffline registers a callback which is triggered when server goes offline state
*/
func (room *Room) SetupOnOffline(callback func()) {
}
