package session

/*
OnSessionCreate assigns a callback on session create event
*/
func (session *Session) OnSessionCreate(cb func(uint8, string)) bool {
	return false
}

/*
CreateSession creates a new session
*/
func (session *Session) CreateSession(sessionType uint8, maxMembers uint8, ttl uint8) {
}
