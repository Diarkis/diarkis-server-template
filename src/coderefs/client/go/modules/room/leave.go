package room

/*
OnLeave assigns a callback on leave event
*/
func (room *Room) OnLeave(callback func(bool)) {
}

/*
Leave leaves a room
roomID: the room ID to leave
message: the message to send to the room
*/
func (room *Room) Leave(roomID string, message []byte) {
}

/*
OnMemberLeave assigns a callback on member leave event
*/
func (room *Room) OnMemberLeave(callback func([]byte)) {
}
