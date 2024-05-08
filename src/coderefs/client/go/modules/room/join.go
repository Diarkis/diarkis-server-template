package room

/*
Join joins a room with the given roomID
*/
func (room *Room) Join(roomID string, message []byte) {
}

/*
OnMemberJoin assigns a callback on member join event
*/
func (room *Room) OnMemberJoin(callback func([]byte)) {
}

/*
OnJoin is a callback setter that is called when a room is joined
The callback is called with the following arguments:
- success: whether the room was joined successfully
- roomCreatedTime: the time at which the room was created
*/
func (room *Room) OnJoin(callback func(bool, uint)) {
}
