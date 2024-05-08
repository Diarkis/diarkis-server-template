package group

import (
	"sync"

	"github.com/Diarkis/diarkis/user"
)

/*
Group group data structure
*/
type Group struct {
	ID                  string
	AllowEmpty          bool
	RemoteMemberNodes   map[string]bool
	MemberMeshEndPoints map[string]string
	MemberSIDs          map[string]string
	IsStatic            bool
	TTL                 int64
	sync.RWMutex
}

/*
Member member of a group data structure
*/
type Member struct {
	UID          string
	SID          string
	MeshEndPoint string
}

/*
SyncGroupData represents internally used data
*/
type SyncGroupData struct {
	GroupID      string `json:"groupID"`
	MeshEndpoint string `json:"meshEndPoint"`
	Method       string `json:"method"`
}

/*
SendData represents internally used data
*/
type SendData struct {
	GroupID  string `json:"groupID"`
	Ver      uint8  `json:"ver"`
	Cmd      uint16 `json:"cmd"`
	Message  []byte `json:"message"`
	Reliable bool   `json:"reliable"`
}

/*
GetGroupByID returns a group instance by its ID.

	[IMPORTANT] If the group does not exist or the given ID.
*/
func GetGroupByID(groupID string) *Group {
	return nil
}

/*
SetOnTickByID assigns a callback to a given group ID to be executed at every given interval in seconds.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[CRITICALLY IMPORTANT] If the callback is blocked, it will lead to a goroutine leak and the tick will be blocked entirely.

	[IMPORTANT] The interval evaluation is executed every 2 seconds.
*/
func SetOnTickByID(groupID string, cb func(groupID string), interval int64) error {
	return nil
}

/*
SetOnRejoinCustomMessage assigns a callback on rejoin to create a custom message

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[CRITICALLY IMPORTANT] If the callback is blocked, it will lead to a goroutine leak and the leave operation will be blocked.

	custom func(groupID string, userID string, userSID string) []byte - Custom function to create custom message byte array.
*/
func SetOnRejoinCustomMessage(custom func(groupID string, userID string, userSID string) (message []byte)) {
}

/*
SetOnDiscardCustomMessage sets a custom function to create a custom message to be sent to members on user discard

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[CRITICALLY IMPORTANT] If the callback is blocked, it will lead to a goroutine leak and the join operation will be blocked.

	custom func(groupID string, userID string, userSID string) []byte - Custom function to create custom message byte array.
*/
func SetOnDiscardCustomMessage(custom func(groupID string, userID string, userSID string) (message []byte)) {
}

/*
SetupAsTCPServer Sets up group package as TCP server
*/
func SetupAsTCPServer() {
}

/*
SetupAsUDPServer Sets up group package as UDP/RUDP server
*/
func SetupAsUDPServer() {
}

/*
Setup optionally loads configurations

	configpath string - Absolute path for the configuration file to read.
*/
func Setup(confpath string) {
}

/*
GetGroupIDList returns an array of group IDs

	[NOTE] Uses mutex lock internally.

Parameters

	userData *user.User - Target user to retrieve the list of group IDs from.
*/
func GetGroupIDList(userData *user.User) []string {
	return nil
}

/*
GetLatestGroupID returns the newest group ID the user joined or created

	[NOTE] Uses mutex lock internally.

Parameters

	userData *user.User - Target user to retrieve the last group ID joined.
*/
func GetLatestGroupID(userData *user.User) string {
	return ""
}

/*
NewGroup Creates a new group and returns the group ID as a string to the callback function.

	[NOTE] Uses mutex lock internally.

Error Cases

	┌────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────┐
	│ Error                      │ Reason                                                                                                 │
	├────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────┤
	│ Invalid TTL                │ TTL must be greater than 10 seconds                                                                    │
	├────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────┤
	│ Node is offline            │ A new group cannot be created on a server node that is in offline state (marked and ready to shutdown) │
	├────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────┤
	│ Allow empty and join false │ Allow empty false must be followed by join true, otherwise the new group will be deleted immediately.  │
	└────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────┘

Parameters

	allowEmpty - If true, the group will not be deleted until TTL expires even if all members leave
	join       - If true, the creator user will join the group automatically
	ttl        - TTL of the group in seconds. The value must be greater than 10
	interval   - Broadcast interval in milliseconds. Interval below 100ms will be forced to 0. If it's 0, no packet merge
*/
func NewGroup(userData *user.User, allowEmpty bool, join bool, ttl int64, interval int64) (string, error) {
	return "", nil
}

/*
GetMemberMeshNodes returns an array of member mesh nodes that are participating the group

	[NOTE] Uses mutex lock internally.

Parameters

	groupID string - Target group ID to retrieve the list of server internal addresses that group is stored.
*/
func GetMemberMeshNodes(groupID string) []string {
	return nil
}

/*
Join Joins a group and notify the other members of the group on joining the group
If message is nil or empty, broadcast will NOT be sent

	[NOTE] Uses mutex lock internally.
	[NOTE] This function is asynchronous.

Error Cases

	+------------------------------+-------------------------------------------+
	| Error                        | Reason                                    |
	+------------------------------+-------------------------------------------+
	| User is already in the group | The user is already a member of the group |
	| Failed to add a member       | The group trying to join does not exist   |
	| Reliable Response timed out  | Internal network failed and timed out     |
	+------------------------------+-------------------------------------------+

Parameters

	groupID  - Target group ID to join.
	userData - User to join the group.
	ver      - Command version used for the message sent when join is successful.
	cmd      - Command ID used for the message sent when join is successful.
	message  - Message byte array to be sent when join is successful.
	           If message is either nil or empty, the message will not be sent.

	callback func(error) - Callback invoked when join operation is completed.
*/
func Join(groupID string, userData *user.User, ver uint8, cmd uint16, message []byte, callback func(error)) {
}

/*
Leave Leaves from a group and notify the other members on leaving
If message is nil or empty, broadcast will NOT be sent

	[NOTE] Uses mutex lock internally.
	[NOTE] This function is asynchronous.

Error Cases

	+--------------------------+-----------------------------------------+
	| Error                    | Reason                                  |
	+--------------------------+-----------------------------------------+
	| User is not in the group | The user is not a member of the group   |
	| Group not found          | The group to leave from does not exists |
	+--------------------------+-----------------------------------------+

Parameters

	groupID  - Target group ID to leave from.
	userData - User to leave the group.
	ver      - Command version used for the message sent when leave is successful.
	cmd      - Command ID used for the message sent when leave is successful.
	message  - Message byte array to be sent when leave is successful.
	           If message is either nil or empty, the message will not be sent.

	callback func(error) - Callback invoked when join operation is completed.
*/
func Leave(groupID string, userData *user.User, ver uint8, cmd uint16, message []byte, callback func(error)) {
}

/*
Broadcast Sends a broadcast message to the other members in the group

	[NOTE] Uses mutex lock internally.
	[NOTE] This function is asynchronous.
	[NOTE] This function may fail silently.

Parameters

	groupID    - Target group ID.
	senderUSer - User that sends the broadcast.
	ver        - Command version used for the broadcast message.
	cmd        - Command ID used for the broadcast message.
	message    - Message byte array to be sent as the broadcast.
	reliable   - If true, UDP will become RUDP.
*/
func Broadcast(groupID string, senderUser *user.User, ver uint8, cmd uint16, message []byte, reliable bool) {
}

/*
Announce sends a broadcast message to the other members in the group without having the "sender":
MUST not be called by client directly for security reason

	[NOTE] Uses mutex lock internally.
	[NOTE] This function is asynchronous.

Parameters

	groupID    - Target group ID.
	senderUSer - User that sends the broadcast.
	ver        - Command version used for the broadcast message.
	cmd        - Command ID used for the broadcast message.
	message    - Message byte array to be sent as the broadcast.
	reliable   - If true, UDP will become RUDP.
*/
func Announce(groupID string, ver uint8, cmd uint16, message []byte, reliable bool) {
}
