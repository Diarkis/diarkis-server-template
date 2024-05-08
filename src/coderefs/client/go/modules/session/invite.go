package session

/*
OnSessionInvite assigns a callback on session invitation event
*/
func (session *Session) OnSessionInvite(cb func(uint8, string, []byte)) bool {
	return false
}

/*
InviteToSession sends an invitation to a user
*/
func (session *Session) InviteToSession(sessionType uint8, targetUIDs []string, message []byte) {
}
