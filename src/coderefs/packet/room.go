package packet

/*
PackChatMessage creates a chat message data byte array
*/
func PackChatMessage(senderID string, timestamp int64, message string) []byte {
	return nil
}

/*
PackChatHistory creates a byte array of list of chat data packed by PackChatMessage
*/
func PackChatHistory(packedChatHistory [][]byte) []byte {
	return nil
}
