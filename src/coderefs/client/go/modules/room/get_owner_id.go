package room

/*
GetOwnerID gets the owner ID of a room
- roomID: the ID of the room to get the owner ID of
*/
func (room *Room) GetOwnerID(roomID string) {
}

/*
OnGetOwnerID assigns a callback on get owner ID event
The callback is called with the following arguments:
- success: whether the room was joined successfully
- ownerID: the ID of the owner of the room
*/
func (room *Room) OnGetOwnerID(callback func(bool, []byte)) {
}
