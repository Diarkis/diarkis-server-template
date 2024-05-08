package session

/*
OnSessionGetMessageIDs assigns a callback on session get member IDs event
*/
func (session *Session) OnSessionGetMessageIDs(cb func(uint8, []string)) bool {
	return false
}

/*
GetMemberIDs gets member IDs of the session
*/
func (session *Session) GetMemberIDs(sessionType uint8) {
}
