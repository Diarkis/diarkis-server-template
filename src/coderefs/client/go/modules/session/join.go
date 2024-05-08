package session

/*
OnMemberJoin assigns a callback on member join event
*/
func (session *Session) OnMemberJoin(cb func(uint8, string)) bool {
	return false
}

/*
JoinSession joins a session
*/
func (session *Session) JoinSession(sessionType uint8, sessionID string) {
}
