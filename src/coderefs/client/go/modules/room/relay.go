package room

/*
Relay sends a message to a room
*/
func (room *Room) Relay(message []byte, reliable bool) {
}

/*
RelayTo sends a message to a room memebers
*/
func (room *Room) RelayTo(memberIDs []string, message []byte, reliable bool) {
}

/*
RelayProfile sends a message to a room
*/
func (room *Room) RelayProfile(message []byte, reliable bool) {
}

/*
RelayToProfile sends a message to a room memebers
*/
func (room *Room) RelayToProfile(memberIDs []string, message []byte, reliable bool) {
}

/*
OnRelayResponse assigns a callback on relay response event
*/
func (room *Room) OnRelayResponse(callback func(success bool, msg []byte)) {
}

/*
OnRelayPush assigns a callback on relay push event
*/
func (room *Room) OnRelayPush(callback func(msg []byte)) {
}

/*
OnRelayToResponse assigns a callback on relay to response event
*/
func (room *Room) OnRelayToResponse(callback func(success bool, msg []byte)) {
}

/*
OnRelayToPush assigns a callback on relay to push event
*/
func (room *Room) OnRelayToPush(callback func(msg []byte)) {
}

/*
OnRelayToProfileResponse assigns a callback on relay to profile response event
*/
func (room *Room) OnRelayToProfileResponse(callback func(success bool, msg []byte)) {
}

/*
OnRelayToProfilePush assigns a callback on relay to profile push event
*/
func (room *Room) OnRelayToProfilePush(callback func(msg []byte)) {
}

/*
OnRelayProfileResponse assigns a callback on relay profile response event
*/
func (room *Room) OnRelayProfileResponse(callback func(success bool, msg []byte)) {
}

/*
OnRelayProfilePush assigns a callback on relay profile push event
*/
func (room *Room) OnRelayProfilePush(callback func(msg []byte)) {
}
