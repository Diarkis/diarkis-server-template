package session

/*
OnMemberLeave assigns a callback on member leave event
*/
func (session *Session) OnMemberLeave(cb func(uint8, string)) bool {
	return false
}

/*
LeaveSession leaves from a session that you have joined
*/
func (session *Session) LeaveSession(sessionType uint8) {
}
