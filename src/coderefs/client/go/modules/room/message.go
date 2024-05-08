package room

/*
Message sends a message to a room
- roomID: the ID of the room to send the message to
- memberIDs: the IDs of the members to send the message to
- message: the message to send
- reliable: whether the message should be sent by RUDP or UDP
*/
func (room *Room) Message(roomID string, memberIDs []string, message []byte, reliable bool) {
}

/*
OnMemberMessage is a callback setter that is called when a message is received from a member
The callback is called with the following arguments:
- message: the message that was received
*/
func (room *Room) OnMemberMessage(callback func([]byte)) {
}
