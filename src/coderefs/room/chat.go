package room

import (
	"github.com/Diarkis/diarkis/user"
)

/*
SetOnChatMessage assigns a callback to be invoked before synchronizing a chat message.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This function does NOT work if the room is not on the same server.

	[NOTE] Uses mutex lock internally.

Return false will reject the incoming chat message and rejected chat message will NOT be sent to other members.
*/
func SetOnChatMessage(roomID string, callback func(*user.User, string) bool) bool {
	return false
}

/*
SyncChatMessage records the chat message data and synchronize the chat message with room members

	[IMPORTANT] This function does NOT work if the room is not on the same server.

	[NOTE]      Uses mutex lock internally.
*/
func SyncChatMessage(roomID string, userData *user.User, ver uint8, cmd uint16, message string) bool {
	return false
}

/*
SetChatHistoryLimit sets the chat history maximum length

	[IMPORTANT] This function does NOT work if the room is not on the same server.

	[NOTE] Uses mutex lock internally.
*/
func SetChatHistoryLimit(roomID string, limit int) bool {
	return false
}

/*
GetChatHistory returns chat history data as an array

	[IMPORTANT] This function does NOT work if the room is not on the same server.

	[NOTE] Uses mutex lock internally.
*/
func GetChatHistory(roomID string) [][]string {
	return nil
}
