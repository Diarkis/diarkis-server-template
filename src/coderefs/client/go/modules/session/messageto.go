package session

/*
OnSessionMessageTo assigns a callback on session message event
*/
func (session *Session) OnSessionMessageTo(cb func(uint8, string)) bool {
	return false
}

/*
MessageTo sends a message to selected members of the session
*/
func (session *Session) MessageTo(sessionType uint8, recipientUIDs []string, message []byte) {
}
