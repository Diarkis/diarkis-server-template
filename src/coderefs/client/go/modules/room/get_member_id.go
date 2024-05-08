package room

/*
GetMemberIDs gets the member IDs of a room
*/
func (room *Room) GetMemberIDs() {
}

/*
OnGetMemberIDs assigns a callback on get member IDs event
*/
func (room *Room) OnGetMemberIDs(callback func(bool, []byte, [][]byte)) {
}
