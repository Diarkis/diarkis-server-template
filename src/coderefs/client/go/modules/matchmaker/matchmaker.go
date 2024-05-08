package matchmaker

import (
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

/*
MatchMaker represents Diarkis MatchMaker client
*/
type MatchMaker struct{}

/*
NewMatchMakerAsTCP creates MatchMaker client as TCP client
*/
func NewMatchMakerAsTCP(tcp *tcp.Client) *MatchMaker {
	return nil
}

/*
NewMatchMakerAsUDP creates MatchMaker client as UDP client
*/
func NewMatchMakerAsUDP(udp *udp.Client) *MatchMaker {
	return nil
}

/*
OnTicketComplete assigns a callback to ticket on complete event
*/
func (mm *MatchMaker) OnTicketComplete(cb func(bool, []byte)) bool {
	return false
}

/*
OnTicketCancel assigns a callback on ticket cancel event
*/
func (mm *MatchMaker) OnTicketCancel(cb func(bool, []byte)) bool {
	return false
}

/*
IssueTicket creates a new ticket
*/
func (mm *MatchMaker) IssueTicket(ticketType uint8) {
}

/*
CancelTicket cancels an existing ticket
*/
func (mm *MatchMaker) CancelTicket(ticketType uint8) {
}

/*
GetRoomID returns the room ID used by MatchMaker

	[IMPORTANT] This is NOT MatchMaker ticket room ID
*/
func (mm *MatchMaker) GetRoomID() string {
	return ""
}
