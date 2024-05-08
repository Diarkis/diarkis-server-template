package room

/*
Register a room (the client is already a member of) with a type (used by FindRoomsByType).
OnRegister is called when the room is registered.
Register registers a room with the given parameters
- roomType: the type of the room
- roomName: the name of the room
- roomMetadata: Extra string data for the room registered.
*/
func (room *Room) Register(roomType int, roomName string, roomMetadata string) {
}

/*
OnRegister is a callback setter that is called when a room is created
The callback is called with the following arguments:
- success: whether the room was created successfully
- msg: Error message byte array if success was false.
*/
func (room *Room) OnRegister(callback func(success bool, msg []byte)) {
}
