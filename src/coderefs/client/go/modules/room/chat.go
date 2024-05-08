package room

/*
OnChatSync assigns a callback on chat sync event
*/
func (room *Room) OnChatSync(callback func(bool, []byte)) {
}

/*
OnChatPush assigns a callback on chat push event
*/
func (room *Room) OnChatPush(callback func(data ChatData)) {
}

/*
OnChatLog assigns a callback on chat log event
*/
func (room *Room) OnChatLog(callback func(bool, []ChatData, []byte)) {
}

/*
Chat sends a message to a room
*/
func (room *Room) Chat(msg string) {
}

/*
GetChatLog gets the chat log of a room
*/
func (room *Room) GetChatLog() {
}
