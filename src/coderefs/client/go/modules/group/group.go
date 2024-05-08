package group

import (
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

/*
Group represents Diarkis Group client
*/
type Group struct{ ID string }

/*
SetupAsTCP sets up the Group client as TCP client
*/
func (group *Group) SetupAsTCP(tcpClient *tcp.Client) bool {
	return false
}

/*
SetupAsUDP sets up the Group client as UDP client
*/
func (group *Group) SetupAsUDP(udpClient *udp.Client) bool {
	return false
}

/*
OnCreate assigns a callback on create event
*/
func (group *Group) OnCreate(callback func(bool, string)) {
}

/*
OnJoin assigns a callback on join event
*/
func (group *Group) OnJoin(callback func(bool, string)) {
}

/*
OnLeave assigns a callback on leave event
*/
func (group *Group) OnLeave(callback func(bool)) {
}

/*
OnMemberJoin assigns a callback on member join event
*/
func (group *Group) OnMemberJoin(callback func([]byte)) {
}

/*
OnMemberLeave assigns a callback on member leave event
*/
func (group *Group) OnMemberLeave(callback func([]byte)) {
}

/*
OnMemberBroadcast assigns a callback on member broadcast event
*/
func (group *Group) OnMemberBroadcast(callback func([]byte)) {
}

/*
Create creates a new group
*/
func (group *Group) Create(allowEmpty bool, join bool, ttl uint16) {
}

/*
Join joins a group
*/
func (group *Group) Join(groupID string, message []byte) {
}

/*
JoinRandom joins a random group or creates a new group if no group is found
*/
func (group *Group) JoinRandom(ttl int, msg []byte, allowEmpty bool, interval int) {
}

/*
Leave leaves from a group that you have joined
*/
func (group *Group) Leave(groupID string, message []byte) {
}

/*
BroadcastTo sends a message to all group members
*/
func (group *Group) BroadcastTo(groupID string, message []byte, reliable bool) {
}
