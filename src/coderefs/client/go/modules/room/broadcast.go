package room

/*
OnMemberBroadcast assigns a callback on member broadcast event
*/
func (room *Room) OnMemberBroadcast(callback func([]byte)) {
}

/*
BroadcastTo broadcasts a message to all members in a room
*/
func (room *Room) BroadcastTo(roomID string, message []byte, reliable bool) {
}
