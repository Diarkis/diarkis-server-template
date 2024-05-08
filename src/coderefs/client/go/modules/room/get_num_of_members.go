package room

/*
GetNumOfMembers gets the number of members in a room
*/
func (room *Room) GetNumOfMembers(roomID string) {
}

/*
OnGetNumOfMembers assigns a callback on get number of members event
*/
func (room *Room) OnGetNumOfMembers(callback func(bool, []byte, int, int)) {
}
