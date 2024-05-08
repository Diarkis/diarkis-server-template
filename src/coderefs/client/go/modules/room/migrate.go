package room

/*
OnMigrate assigns a callback on migrate event
*/
func (room *Room) OnMigrate(callback func(bool, string)) {
}

/*
Migrate migrates a room that the user is currently in.
This operation can only be performed by the owner of the room.
When successful, the owner user raises OnConnect event with reconnect flag true.
All members of the room will automatically join the new migrated room without sending join broadcast.
*/
func (room *Room) Migrate() {
}

/*
Move moves a room that the user is currently in.
*/
func (room *Room) Move(newRoomID string, leaveMessage []byte, joinMessage []byte, callback func(bool, uint)) {
}
