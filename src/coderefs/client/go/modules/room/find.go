package room

/*
FindRoomsByType finds rooms by type
- roomType: the type of the room
- limit: the maximum number of rooms to return
*/
func (room *Room) FindRoomsByType(roomType int, limit int) {
}

/*
OnFindRoomsByType is a callback setter that is called when a room is created
The callback is called with the following arguments:
- success: whether the room was created successfully
- roomListItem: the list of roomListItem
*/
func (room *Room) OnFindRoomsByType(callback func(success bool, room []ListItem)) {
}
