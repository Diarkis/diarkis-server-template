package session

/*
OnSessionBroadcast assigns a callback on session broadcast event
*/
func (session *Session) OnSessionBroadcast(cb func(uint8, string)) bool {
	return false
}

/*
BroadcastSession sends a message to all session members
*/
func (session *Session) BroadcastSession(sessionType uint8, message []byte) {
}
