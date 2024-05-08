package room

/*
OnCreate is a callback setter that is called when a room is created
The callback is called with the following arguments:
- success: whether the room was created successfully
- roomID: the ID of the room that was created
- roomCreateTime: the time at which the room was create
*/
func (room *Room) OnCreate(callback func(bool, string, uint)) {
}

/*
Create creates a room with the given parameters
- maxMembers: the maximum number of members that can join the room
- allowEmpty: whether the room should be destroyed when the last member leaves
- join: whether the client should join the room after creating it
- ttl: the time to live of the room
*/
func (room *Room) Create(maxMembers uint16, allowEmpty bool, join bool, ttl uint16, interval uint16) {
}
