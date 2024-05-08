package session

import (
	"sync"

	"github.com/Diarkis/diarkis/user"
)

/*
Session represents session.

A session allows multiple users join a members up to the pre-determines.

All session members may send and receive messages and set and get shared properties within the session.
*/
type Session struct{ sync.RWMutex }

/*
PropertyKeyValue represents session property key value pair.
*/
type PropertyKeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

/*
MessageData represents internally used broadcast message data
*/
type MessageData struct {
	ID         string   `json:"id"`
	Ver        uint8    `json:"ver"`
	Cmd        uint16   `json:"cmd"`
	Msg        []byte   `json:"msg"`
	MemberSIDs []string `json:"memberSIDs"`
}

/*
UserData represents internally used user data update
*/
type UserData struct {
	SessionID        string `json:"id"`
	UserID           string `json:"ID"`
	SID              string `json:"SID"`
	PublicAddr       string `json:"PublicAddr"`
	PrivateAddrBytes []byte `json:"PrivateAddrBytes"`
	MeshAddr         string `json:"meshAddr"`
}

/*
FindOwnerData represents internally used data
*/
type FindOwnerData struct {
	SID string `json:"sid"`
}

/*
LeaveSessionData represents internally used data
*/
type LeaveSessionData struct {
	Type uint8  `json:"sessionType"`
	ID   string `json:"id"`
	UID  string `json:"uid"`
	SID  string `json:"sid"`
}

/*
LeaveSessionReturnData represents internally used data
*/
type LeaveSessionReturnData struct {
	LockKey string `json:"lockKey"`
}

/*
JoinSessionData represents internally used data
*/
type JoinSessionData struct {
	ID       string                 `json:"id"`
	UID      string                 `json:"uid"`
	SID      string                 `json:"sid"`
	MeshAddr string                 `json:"meshAddr"`
	UserData map[string]interface{} `json:"userData"`
}

/*
IsOwnerData represents internally used data
*/
type IsOwnerData struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
}

/*
IsOwnerReturnData represents internally used data
*/
type IsOwnerReturnData struct {
	IsOwner bool `json:"isOwner"`
}

/*
PropertyData represents internally used data
*/
type PropertyData struct {
	ID    string      `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Flag  bool        `json:"flag"`
}

/*
PropertyReturnData represents internally used data
*/
type PropertyReturnData struct {
	Success bool        `json:"success"`
	Key     string      `json:"key"`
	Value   interface{} `json:"value"`
}

/*
PropertiesData represents internally used data
*/
type PropertiesData struct {
	ID     string              `json:"id"`
	KVList []*PropertyKeyValue `json:"kvList"`
}

/*
GetMemberData represents internally used data
*/
type GetMemberData struct {
	ID   string   `json:"id"`
	List []string `json:"list"`
}

/*
GetNumberOfMembersData represents internally used data
*/
type GetNumberOfMembersData struct {
	ID         string `json:"id"`
	MaxMembers int    `json:"maxMembers"`
	Number     int    `json:"number"`
}

/*
Setup must be invoked before calling diarkis.Start in order to use session package.
*/
func Setup() {
}

/*
GetSessionIDByUser returns the session ID of the session that the user is member of.

[IMPORTANT] this function works only if the session exists on the same server.

	[IMPORTANT] If the second returned value is false, it means that the user is not a member of any session by the given type.
*/
func GetSessionIDByUser(sessionType uint8, userData *user.User) (string, bool) {
	return "", false
}

/*
GetSessionByID corresponding session to passed id.

[IMPORTANT] this function works only if the session exists on the same server.
*/
func GetSessionByID(id string) (*Session, bool) {
	return nil, false
}

/*
NewSession creates a new session instance with the given session type.

	[IMPORTANT] The user given becomes a member of the session automatically.
	[IMPORTANT] You may not create a new session if the server is in offline state.

Error Cases

	+-----------------------------+------------------------------------------------------------------------------------------------------+
	| Error                       | Reason                                                                                               |
	+-----------------------------+------------------------------------------------------------------------------------------------------+
	| Server is offline           | No session can be created on a server that is in offline state due to receiving SIGTERM.             |
	| Max member of session       | Max member of session must be greater than 1.                                                        |
	| Failed to create session ID | When the server fails to generate UUID v4 ID string and/or invalid internal server address is found. |
	+-----------------------------+------------------------------------------------------------------------------------------------------+

	When it fails to generate session ID, it most likely is caused by the incorrect setup of Diarkis server cluster or a server.

Parameters

	userData    - A user that creates a new session and becomes the first member.
	sessionType - A session type that the new session will be identified by.
	              A user cannot be a member of multiple sessions with the same session type.
	maxMembers  - Maximum allowed number of session members.
	ttl         - TTL of the session instance when it becomes empty (no members) in seconds.
	              TTL count is not exact that the deletion of an empty session may happen earlier or later than TTL.
	              Value smaller than 30s for TTL will automatically be changed to 30s.
*/
func NewSession(userData *user.User, sessionType uint8, maxMembers uint16, ttl uint16) (*Session, error) {
	return nil, nil
}

/*
IsOwner checks if the given user is the owner of the session or not.

	[IMPORTANT] The owner user is used internally,
	            and it may change when the current owner user disconnects or leaves the session.
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+----------------------------------------------+---------------------------------------------------------------------------------------+
	| Error                                        | Reason                                                                                |
	+----------------------------------------------+---------------------------------------------------------------------------------------+
	| Failed to ascertain if the user is the owner | The given user is not a member of the session or session being not available anymore. |
	| Failed to ascertain if the user is the owner | Internal server-to-server communication failed due to either from server load         |
	|                                              | and/or incorrectly setup Diarkis cluster.                                             |
	+----------------------------------------------+---------------------------------------------------------------------------------------+

Parameters

	sessionType - Session type.
	userData    - Session member user to be checked if the user is the internal owner of the session or not.
	cb          - Callback.
	              func(err error, isOwner bool)
	              err     - Error to indicate the failure of the operation.
	              isOwner - If true, the given user is the internal owner of the session.
*/
func IsOwner(sessionType uint8, userData *user.User, cb func(err error, isOwner bool)) {
}

/*
SendInvite sends an invitation to a user to join the session.

The invitation message (sessionType, sessionID, and custom invitation message) is sent to the user client device.

The user may join the session by accepting the invite via Join.

	[IMPORTANT] Invitation does NOT guarantee the invited user to be able to join the session.
	            If the session is full or no longer available when the invited user accepts the session,
	            joining the session will fail.
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+------------------------------+--------------------------------------------------------------------------------------+
	| Error                        | Reason                                                                               |
	+------------------------------+--------------------------------------------------------------------------------------+
	| Must a member                | Only the member of the session may send an invitation.                               |
	| Cannot invite yourself       | You may not invite yourself.                                                         |
	| Failed to deliver invitation | Invitation failed to be delivered to the intended user. e.i. User is not found etc.  |
	|                              | If you send invitations to multiple users and only one fails,                        |
	|                              | the callback will return with an error.                                              |
	+------------------------------+--------------------------------------------------------------------------------------+

Parameters

	sessionType  - The session type.
	id           - The session ID.
	userData     - The member user of the session to send an invite.
	targetUserID - The user ID to be sent the invitation to.
	ver          - Command version of the invitation message to be sent.
	cmd          - Command ID of the invitation message to be sent.
	message      - Message of the invitation to be sent.
	cb           - Callback to be invoked when the invitation is sent
	               (not to the client device but when the server finishes the send operation).
*/
func SendInvite(sessionType uint8, id string, userData *user.User, targetUserIDs []string, ver uint8, cmd uint16, message []byte, cb func(err error)) {
}

/*
JoinSessionByID joins a session as a member by session's type and ID.

	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+----------------+-------------------------------------------------------------------------------+
	| Error          | Reason                                                                        |
	+----------------+-------------------------------------------------------------------------------+
	| Failed to join | Internal server-to-server communication failed due to either from server load |
	|                | and/or incorrectly setup Diarkis cluster.                                     |
	| Failed to join | Session is no longer available.                                               |
	| Failed to join | User is rejected by On join callback check.                                   |
	+----------------+-------------------------------------------------------------------------------+

Parameters

	sessionType       - The session type.
	id                - The session ID.
	userData          - The user that will join the session.
	cb                - Callback to be invoked when join operation completes regardless of success or failure.
*/
func JoinSessionByID(sessionType uint8, id string, userData *user.User, cb func(error)) {
}

/*
KickFromSessionWithUserData kicks out a user from a session.

This gets sessionID out of userData and passes it to KickFromSession.

	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

	[NOTE] See the doc for session.KickFromSession for more Error Cases.

Error Cases

	+-----------------+--------------------------------------+
	| Error           | Reason                               |
	+-----------------+--------------------------------------+
	| Kick out failed | User is not a member of the session. |
	+-----------------+--------------------------------------+

Parameters

	sessionType  - Session type.
	userData     - Session member user.
	targetUserID - User ID to kick out.
	cb           - Callback.
*/
func KickFromSessionWithUserData(sessionType uint8, userData *user.User, targetUserID string, cb func(error)) {
}

/*
KickFromSession kicks out a user from a session.

	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+----------------------+--------------------------------------------------------------------------------+
	| Error                | Reason                                                                         |
	+----------------------+--------------------------------------------------------------------------------+
	| User is not a member | The target user to kick out is not a member of the session.                    |
	| Kick out failed      | Session ID is invalid.                                                         |
	|                      | The most likely cause is internal bug.                                         |
	| Kick out failed      | Internal server-to-server communication failed due to either from server load. |
	|                      | and/or incorrectly setup Diarkis cluster.                                      |
	+----------------------+--------------------------------------------------------------------------------+

Parameters

	sessionType  - Session type.
	sessionID    - Session ID.
	targetUserID - User ID to kick out.
	cb           - Callback.
*/
func KickFromSession(sessionType uint8, sessionID string, targetUserID string, cb func(error)) {
}

/*
LeaveSessionByID leaves from the session by its type and ID.

	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+-----------------------+-------------------------------------------------------------------------+
	| Error                 | Reason                                                                  |
	+-----------------------+-------------------------------------------------------------------------+
	| Session ID is invalid | The given session ID is not a valid session ID.                         |
	| User is not a member  | The given user is not a member of the session.                          |
	| Failed to leave       | Session is no longer available.                                         |
	| Failed to leave       | Server-to-server communication failed.                                  |
	|                       | Possible cause maybe server load, incorrectly setup server cluster etc. |
	+-----------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - The session type.
	id          - The session ID.
	userData    - The user that will leave from the session.
	cb          - Callback to be invoked when leave operation completes regardless of success or failure.
*/
func LeaveSessionByID(sessionType uint8, id string, userData *user.User, cb func(error)) {
}

/*
SetPropertyIfNotExists stores a key and a value as a property to the ticket if the same key does not exist.

	[IMPORTANT] Properties are only primitive values and does not support reference type data such as array and map.
	[IMPORTANT] This function is asynchronous.
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+------------------------+-------------------------------------------------------------------------+
	| Error                  | Reason                                                                  |
	+------------------------+-------------------------------------------------------------------------+
	| Session ID is invalid  | The given session ID is not a valid session ID.                         |
	| User is not a member   | The given user is not a member of the session.                          |
	| Failed to set property | Session is no longer available.                                         |
	| Failed to set property | Server-to-server communication failed.                                  |
	|                        | Possible cause maybe server load, incorrectly setup server cluster etc. |
	| Failed to set property | The property with the same key already exists.                          |
	+------------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - Type of session.
	userData    - Session member user.
	key         - Property key.
	value       - Property value to set.
	cb          - Callback to be invoked when the operation is finished.
	              func(err error, success bool)
	              err     - If not nil, operation error out.
	              success - If true, setting of the key and value was a success.
*/
func SetPropertyIfNotExists(sessionType uint8, userData *user.User, key string, value interface{}, cb func(err error, success bool)) {
}

/*
SetProperty stores a key and a value as a property to the ticket.

	[IMPORTANT] If the same key exists, it overwrites the existing value of the same key.
	[IMPORTANT] Properties are only primitive values and does not support reference type data such as array and map.
	[IMPORTANT] This function is asynchronous.
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+------------------------+-------------------------------------------------------------------------+
	| Error                  | Reason                                                                  |
	+------------------------+-------------------------------------------------------------------------+
	| Session ID is invalid  | The given session ID is not a valid session ID.                         |
	| User is not a member   | The given user is not a member of the session.                          |
	| Failed to set property | Session is no longer available.                                         |
	| Failed to set property | Server-to-server communication failed.                                  |
	|                        | Possible cause maybe server load, incorrectly setup server cluster etc. |
	+------------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - Type of session.
	userData    - Session member user.
	key         - Property key.
	value       - Property value to set.
	cb          - Callback to be invoked when the operation is finished.
	              func(err error, success bool)
	              err     - If not nil, operation error out.
	              success - If true, setting of the key and value was a success.
*/
func SetProperty(sessionType uint8, userData *user.User, key string, value interface{}, cb func(err error, success bool)) {
}

/*
SetProperties stores a collection of keys and their values to ticket as properties.

	[IMPORTANT] If the same key exists, it overwrites the existing value of the same key.
	[IMPORTANT] Properties are only primitive values and does not support reference type data such as array and map.
	[IMPORTANT] If there is an error no properties will be set.
	[IMPORTANT] This function is asynchronous.
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+--------------------------+-------------------------------------------------------------------------+
	| Error                    | Reason                                                                  |
	+--------------------------+-------------------------------------------------------------------------+
	| Session ID is invalid    | The given session ID is not a valid session ID.                         |
	| User is not a member     | The given user is not a member of the session.                          |
	| Failed to set properties | Session is no longer available.                                         |
	| Failed to set properties | Server-to-server communication failed.                                  |
	|                          | Possible cause maybe server load, incorrectly setup server cluster etc. |
	+--------------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - Session type of the session.
	userData    - Session member user.
	kvList      - A list of property key and value pair to be set.
	cb          - Callback to be invoked when the operation is completed.
	              func(err error)
	              err - Error if the operation fails.
*/
func SetProperties(sessionType uint8, userData *user.User, kvList []*PropertyKeyValue, cb func(err error)) {
}

/*
GetProperty returns the value of the given key and if the key does not exist, the second return value will be a false.

	[IMPORTANT] Properties are only primitive values and does not support reference type data such as array and map.
	[IMPORTANT] The returned property value is an interface{}, in order to type assert safely, please use Diarkis' util package functions.
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+------------------------+-------------------------------------------------------------------------+
	| Error                  | Reason                                                                  |
	+------------------------+-------------------------------------------------------------------------+
	| Session ID is invalid  | The given session ID is not a valid session ID.                         |
	| User is not a member   | The given user is not a member of the session.                          |
	| Failed to get property | Session is no longer available.                                         |
	| Failed to get property | Server-to-server communication failed.                                  |
	|                        | Possible cause maybe server load, incorrectly setup server cluster etc. |
	+------------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - Session type of the session.
	userData    - Session member user.
	key         - Property key.
	cb          - Callback to be invoked when fetching the property is complete.
	              func(err error, value interface{}, exists bool)
	              err    - Error when fetching of the property fails.
	              value  - Value of the property key.
	              exists - If the property key is not found, it is false.

Example

	GetProperty(sessionType, userData, "someKey", func(err error, value interface{}, exists bool)) {

		if err != nil {
			// handle error here...
		}

		if !exists {
			// key does not exist...
		}

		// If the value data type is an uint8, of course ;)
		v, ok := util.ToUint8(v)
	})
*/
func GetProperty(sessionType uint8, userData *user.User, key string, cb func(err error, value interface{}, exists bool)) {
}

/*
GetProperties returns key and value pairs as a map.

	[IMPORTANT] Properties are only primitive values and does not support reference type data such as array and map.
	[IMPORTANT] If a value of a given key does not exist, the returned map will not contain the key without the value.
	[IMPORTANT] The returned property value is an interface{}, in order to type assert safely, please use Diarkis' util package functions.
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+--------------------------+-------------------------------------------------------------------------+
	| Error                    | Reason                                                                  |
	+--------------------------+-------------------------------------------------------------------------+
	| Session ID is invalid    | The given session ID is not a valid session ID.                         |
	| User is not a member     | The given user is not a member of the session.                          |
	| Failed to get properties | Session is no longer available.                                         |
	| Failed to get properties | Server-to-server communication failed.                                  |
	|                          | Possible cause maybe server load, incorrectly setup server cluster etc. |
	+--------------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - Session type of the session.
	userData    - Session member user.
	keys        - An array of property keys to fetch.
	cb          - Callback to be invoked when the fetch operation is complete.
	              func(err error, props map[string]interface{})
	              err   - Error to inform if the fetch operation fails.
	              props - A map of property keys and values.

Example

	GetProperties(sessionType, userData, []string{ "someKey" }, func(err error, props map[string]interface{}) {

	  if !ok {
	    // handle error here
	  }

	  for key, v := range props {
	    // If the value data type is an uint8, of course ;)
	    value, ok := util.ToUint8(v)
	  }
	})
*/
func GetProperties(sessionType uint8, userData *user.User, keys []string, cb func(err error, props map[string]interface{})) {
}

/*
GetMemberIDs returns the list of member user IDs with the callback.

	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+--------------------------+-------------------------------------------------------------------------+
	| Error                    | Reason                                                                  |
	+--------------------------+-------------------------------------------------------------------------+
	| Session ID is invalid    | The given session ID is not a valid session ID.                         |
	| User is not a member     | The given user is not a member of the session.                          |
	| Failed to get member IDs | Session is no longer available.                                         |
	| Failed to get member IDs | Server-to-server communication failed.                                  |
	|                          | Possible cause maybe server load, incorrectly setup server cluster etc. |
	+--------------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - Session type.
	userData    - Session member user.
	cb          - Callback to be invoked when the operation is complete.
	              func(err error, memberIDs []string)
	              err       - Error when the operation fails.
	              memberIDs - A list of session member IDs.
*/
func GetMemberIDs(sessionType uint8, userData *user.User, cb func(err error, memberIDs []string)) {
}

/*
GetMemberSIDs returns the list of member user IDs with the callback.

	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+---------------------------+-------------------------------------------------------------------------+
	| Error                     | Reason                                                                  |
	+----------------------------+-------------------------------------------------------------------------+
	| Session ID is invalid     | The given session ID is not a valid session ID.                         |
	| User is not a member      | The given user is not a member of the session.                          |
	| Failed to get member SIDs | Session is no longer available.                                         |
	| Failed to get member SIDs | Server-to-server communication failed.                                  |
	|                           | Possible cause maybe server load, incorrectly setup server cluster etc. |
	+---------------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - Session type.
	userData    - Session member user.
	cb          - Callback to be invoked when the operation is complete.
	              func(err error, memberIDs []string)
	              err       - Error when the operation fails.
	              memberSIDs - A list of session member SIDs.
*/
func GetMemberSIDs(sessionType uint8, userData *user.User, cb func(err error, memberSIDs []string)) {
}

/*
GetNumberOfSessionMembers returns the number of current members and the maximum number of members in the session.

	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+-----------------------+-------------------------------------------------------------------------+
	| Error                 | Reason                                                                  |
	+-----------------------+-------------------------------------------------------------------------+
	| Session ID is invalid | The given session ID is not a valid session ID.                         |
	| User is not a member  | The given user is not a member of the session.                          |
	| Failed to get number  | Session is no longer available.                                         |
	| Failed to get number  | Server-to-server communication failed.                                  |
	|                       | Possible cause maybe server load, incorrectly setup server cluster etc. |
	+-----------------------+-------------------------------------------------------------------------+

Parameters

	sessionType - Session type.
	userData    - Session member user.
	cb          - Callback to be invoked when the operation is complete.
	              func(err error, currentMembers int, maxMembers int)
	              err            - Error when the operation fails.
	              currentMembers - The number of current members in the session.
	              maxMembers     - The maximum number of members in the session.
*/
func GetNumberOfSessionMembers(sessionType uint8, userData *user.User, cb func(err error, currentMembers int, maxMembers int)) {
}

/*
Broadcast sends a reliable message to all member users with the given ver, cmd, and message byte array.

	[IMPORTANT] This function can be executed by any user as long as it is provided with a valid session ID..
	[IMPORTANT] If the session no longer exists, it will not propagate the error, but silently fails internally.
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+-----------------------+-------------------------------------------------+
	| Error                 | Reason                                          |
	+-----------------------+-------------------------------------------------+
	| Session ID is invalid | The given session ID is not a valid session ID. |
	| Session not found     | The session no longer exists.                   |
	+-----------------------+-------------------------------------------------+

Parameters

	sessionType - Type of the session.
	sessionID   - Session ID.
	ver         - Command ver to be used as the broadcast message.
	cmd         - Command ID to be used as the broadcast message.
	msg         - Broadcast message to be sent.
*/
func Broadcast(sessionType uint8, sessionID string, ver uint8, cmd uint16, msg []byte) error {
	return nil
}

/*
MessageTo sends a reliable message to specified member users with the given ver, cmd, recipient UIDs and message byte array.

	[IMPORTANT] This function can be executed by any user as long as it is provided with a valid session ID..
	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

Error Cases

	+-----------------------+-------------------------------------------------+
	| Error                 | Reason                                          |
	+-----------------------+-------------------------------------------------+
	| Session ID is invalid | The given session ID is not a valid session ID. |
	| Session not found     | The session no longer exists.                   |
	| User    not found     | A recipient was not found in the session.       |
	+-----------------------+-------------------------------------------------+k

Parameters

	sessionType   - Type of the session.
	sessionID     - Session ID.
	senderUID     - Sender user ID.
	recipientUIDs - An array of recipient user IDs to send the message to.
	ver           - Command ver to be used as the message command.
	cmd           - Command ID to be used as the message command.
	msg           - Message to be sent.
*/
func MessageTo(sessionType uint8, sessionID string, senderUID string, recipientUIDs []string, ver uint8, cmd uint16, msg []byte) error {
	return nil
}

/*
GetAllSessionMemberData collects and returns an array of UserData from the servers where Session members are.
Since the client data in Members can be outdated, you might want to use this func when you need the latest data.

	[IMPORTANT] The function communicates with another server internally.
	[IMPORTANT] This function works on any server.

	[NOTE] It works only when you are a member of Session specified with passed sessionType.

Error Cases:

	+-----------------------+-------------------------------------------------------------------------------------------+
	| Error                 | Reason                                                                                    |
	+-----------------------+-------------------------------------------------------------------------------------------+
	| Cannot get session ID | The passed sessionType is wrong or you are not a member of the Session anymore.           |
	| Invalid session ID    | Session ID is wrong format. This normally does not happen as it is from userData.         |
	+-----------------------+-------------------------------------------------------------------------------------------+

Parameters

	sessionType - sessionType that you want to get the data from.
	userData    - Sender user data: sender.
	cb          - Callback invoked on get session member data. It returns an array of UserData.
*/
func GetAllSessionMemberData(sessionType uint8, userData *user.User, cb func(err error, sud []*UserData)) {
}

/*
GetSessionID returns a session of the given session type that the given user is a member of.

	[IMPORTANT] It returns an empty string if the user is not a member of any session by the given type.

	[NOTE] If the user is not a member of any session of the given sessionType, it returns an empty string.
*/
func GetSessionID(sessionType uint8, userData *user.User) string {
	return ""
}
