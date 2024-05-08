package session

import (
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

/*
Session represents Diarkis Session client
*/
type Session struct{}

/*
NewSessionAsTCP creates a Diarkis Session client as TCP
*/
func NewSessionAsTCP(tcp *tcp.Client) *Session {
	return nil
}

/*
NewSessionAsUDP creates a Diarkis Session client as UDP
*/
func NewSessionAsUDP(udp *udp.Client) *Session {
	return nil
}
