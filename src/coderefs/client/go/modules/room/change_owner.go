package room

/*
OnRoomOwnerChange is a callback setter that is called when the owner of the room changes.
The callback is called with the following arguments:
- roomID: the new ID of the room
*/
func (room *Room) OnRoomOwnerChange(callback func(message []byte)) {
}
