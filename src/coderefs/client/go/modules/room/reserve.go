package room

/*
OnReserve assigns a callback on reserve event
*/
func (room *Room) OnReserve(callback func(bool, []byte)) {
}

/*
Reserve reserves a room
*/
func (room *Room) Reserve(userIDs []string) {
}

/*
CancelReservation cancels a room reservation
*/
func (room *Room) CancelReservation(userIDs []string) {
}
