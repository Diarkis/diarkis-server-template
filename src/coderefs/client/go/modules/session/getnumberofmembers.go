package session

/*
OnSessionGetNumberOfMembers assigns a callback on session get number of members event
*/
func (session *Session) OnSessionGetNumberOfMembers(cb func(uint8, uint16, uint16)) bool {
	return false
}

/*
GetNumberOfMembers gets the number of session members
*/
func (session *Session) GetNumberOfMembers(sessionType uint8) {
}
