package room

import (
	"github.com/Diarkis/diarkis/user"
)

const MemberAdded = 1
const MemberRemoved = 2
const PropertyChanged = 3

/*
Room Room data structure
*/
type Room struct {
	ID         string
	MaxMembers int
	AllowEmpty bool
	Members    []string
	MemberIDs  map[string]string
}

/*
Setup Sets up room package but without any server association
*/
func Setup() {
}

/*
SetupAsTCPServer Sets up room package
*/
func SetupAsTCPServer() {
}

/*
SetupAsUDPServer Sets up room package
*/
func SetupAsUDPServer() {
}

/*
SetJoinCondition registers a condition function to control room join. Returning an error will cause room join to fail

Join condition is NOT protected by mutex lock.

	SetJoinCondition(condition func(roomID string, userData *user.User))

	cond func(roomID string, userData *user.User) error - Custom function to be invoked before room join to evaluate join conditions.
*/
func SetJoinCondition(cond func(roomID string, userData *user.User) error) {
}

/*
SetRollbackOnJoinFailureByID registers a callback to perform rollback operations on join failure

	[NOTE] Uses mutex lock internally.
*/
func SetRollbackOnJoinFailureByID(roomID string, callback func(roomID string, userData *user.User)) bool {
	return false
}

/*
SetOnLeaveByID registers a callback on leave

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.
	[CRITICALLY IMPORTANT] The callbacks are invoked in room.Leave function. If a callback blocks, room.Leave will block also.

	[NOTE] Uses mutex lock internally.
*/
func SetOnLeaveByID(roomID string, callback func(roomID string, userData *user.User)) bool {
	return false
}

/*
SetOnJoinCompleteByID registers a callback on join operation complete.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[NOTE] Uses mutex lock internally.

# Joining a room that requires server re-connect

When joining a room that is on a different server,
the client is required to re-connect to the server where the room is.
This is handle by Diarkis internally.

# Where SetOnJoinCompleteByID callback is executed

The diagram below explains where the callback of SetOnJoinCompleteByID is executed
when the client re-connects:

The number in < > indicates the order of execution of callbacks and operations.

	┌──────────────────┐            ┌────────────────────────┐
	│ UDP/TCP server A │  <2> Join  │    UDP/TCP server B    │
	│                  │───────────▶︎│ [ Room A exists here ] │
	└──────────────────┘            └────────────────────────┘
	        ▲                                     ▲ <4> Callback of SetOnJoinCompleteByID is on server B
	        │                                     │
	        │ <1> Join                            │ <3> Re-connect
	        │                                     │
	    ╭────────╮                                │
	    │ Client │────────────────────────────────┘
	    ╰────────╯ Client re-connects to where the room is

	SetOnJoinCompleteByID waits until the client completes re-connection to the new server and finishes the join operation.
*/
func SetOnJoinCompleteByID(roomID string, callback func(roomID string, userData *user.User)) bool {
	return false
}

/*
SetOnJoinByID registers a callback to be executed just before the join operation to evaluate,
if the user should join or not.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.
	[CRITICALLY IMPORTANT] If the callback blocks, the room will NOT execute Join that follows.

	[NOTE] Uses mutex lock internally.

Returning value false will cause the join operation to fail.

When joining a room that is on a different server,
the client is required to re-connect to the server where the room is.
This is handle by Diarkis internally.

Diagram below shows where the callback of SetOnJoinByID is executed.

	┌──────────────────┐            ┌────────────────────────┐
	│ UDP/TCP server A │  <2> Join  │    UDP/TCP server B    │
	│                  │───────────▶︎│ [ Room A exists here ] │
	└──────────────────┘            └────────────────────────┘
	        ▲                 <3> Callback of SetOnJoinByID is executed on server B
	        │
	        │ <1> Join
	        │
	    ╭────────╮
	    │ Client │
	    ╰────────╯
*/
func SetOnJoinByID(roomID string, callback func(roomID string, userData *user.User) bool) bool {
	return false
}

/*
SetOnCustomUserDiscardByID registers a callback to be executed when a room member disconnects unexpected.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[CRITICALLY IMPORTANT] If the callback blocks, it will NOT allow the room to execute Leave when the user is disconnected.

	[NOTE] Uses mutex lock internally.

The purpose of this callback is to use custom command ver, cmd, and message data to be sent to other members.
*/
func SetOnCustomUserDiscardByID(roomID string, callback func(roomID string, userData *user.User) (uint8, uint16, []byte)) bool {
	return false
}

/*
SetOnClock assigns a custom callback to be invoked at a given interval until it is stopped by ClockStop callback

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID   - Target room ID of the room to have a clock callback loop
	name     - Unique name of the clock callback. Duplicates will be ignored
	interval - Clock interval in seconds
	callback - Callback function to be invoked
*/
func SetOnClock(roomID string, name string, interval int64, callback func(roomID string)) bool {
	return false
}

/*
SetOnClockStop assigns a callback to be invoked when a clock callback is stopped by StopClock

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID   - Target room ID of the room to stop the clock callback
	name     - Unique name of the clock to stop
	callback - Callback function to be invoked when clock is stopped by StopClock
*/
func SetOnClockStop(roomID string, name string, callback func(string)) bool {
	return false
}

/*
StopClock stops a clock callback loop

	[NOTE] Uses mutex lock internally.

Parameters

	roomID - Target room ID of the room to stop the clock callback
	name   - Unique name of the clock to stop
*/
func StopClock(roomID string, name string) bool {
	return false
}

/*
SetOnTick registers a callback on every 5 seconds tick on a room

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] If you access room properties in the callback, it updates the TTL of the room internally.
	            It means that as long as the callback keeps accessing the properties, it will keep the room
	            from being deleted beyond its TTL.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID   string              - Target room ID to set a tick callback.
	callback func(roomID string) - Callback to be invoked at every tick.

Usage Example: Use this to call matching.Add so that you may keep the room available for matchmaking as long as you with.
*/
func SetOnTick(roomID string, callback func(string)) bool {
	return false
}

/*
SetOnTickStop assigns a callback to be invoked when room tick stops

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID   string       - Target room ID
	callback func(string) - Callback to be invoked on tick stop. The callback passes roomID
*/
func SetOnTickStop(roomID string, callback func(string)) bool {
	return false
}

/*
StopTick stops room's tick event loop

	[NOTE] Uses mutex lock internally.
*/
func StopTick(roomID string) {
}

/*
SetOnDiscardCustomMessage sets a custom function to create a custom message on user discard to be sent to members
The callback receives roomID string, user ID string, user SID string

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	custom func(roomID string, userID string, userSID string) []byte
	- Custom function to create a custom message byte array to be sent to other members. When the member client disconnects and leave.
*/
func SetOnDiscardCustomMessage(custom func(string, string, string) []byte) bool {
	return false
}

/*
SetOnRoomDiscard registers a callback function on room deletions

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	cb func(roomID string) - Function to be invoked on room's deletion.
*/
func SetOnRoomDiscard(cb func(roomID string)) bool {
	return false
}

/*
SetOnRoomDiscardByID registers a callback on room deletion by room ID

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This function does NOT work if the room is not on the same server.

	[NOTE] Uses mutex lock internally.
*/
func SetOnRoomDiscardByID(roomID string, callback func(roomID string)) bool {
	return false
}

/*
SetOnRoomChange registers a callback function on room change such as create, join, leave, property updates

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This function does NOT work if the room is not on the same server.

Parameters

	cb - Callback
	     cb func(roomID string, changeEvent int, memberUserIDs []string, roomProperties map[string]interface{})
	     changeEvent is an enum value that tells us what change took place: room.MemberAdded, room.MemberRemoved, room.PropertyChanged
*/
func SetOnRoomChange(cb func(string, int, []string, map[string]interface{})) bool {
	return false
}

/*
SetOnRoomChangeByID registers a callback function on room change such as create, join, leave, property updates

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This function does NOT work if the room is not on the same server.

	changeEvent is an enum value that tells us what change took place: room.MemberAdded, room.MemberRemoved, room.PropertyChanged
	cb func(roomID string, changeEvent int, memberUserIDs []string, roomProperties map[string]interface{})
*/
func SetOnRoomChangeByID(roomID string, cb func(roomID string, changeEvent int, memberIDs []string, properties map[string]interface{})) bool {
	return false
}

/*
SetupBackup Enables room data backup to different nodes
- maximum backup nodes you can set is 2

	backupNum int - Number of backup servers.
*/
func SetupBackup(backupNum int) {
}

/*
SetOnRoomPropertyUpdate sets a callback for room property updates.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[NOTE] Uses mutex lock internally.

Parameters

	cb func(roomID string, properties map[string]interface{}) - Callback to be invoked when room property changes.
*/
func SetOnRoomPropertyUpdate(cb func(roomID string, props map[string]interface{})) {
}

/*
SetOnRoomOwnerChange sets a callback for room owner change

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] The owner ID may be an empty string.
	[NOTE] Uses mutex lock internally.

Parameters

	cb func(roomID string, ownerID string) - Callback to be invoked when the room owner changes.
*/
func SetOnRoomOwnerChange(cb func(roomID string, ownerID string)) {
}

/*
SetOnAnnounce assigns a listener callback on Announce of room

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[NOTE] Uses mutex lock internally.
	[NOTE] ver and cmd passed to the callback are the ver and cmd used to send Broadcast or Message.

It captures Broadcast and Message as well

	roomID string - Target room ID.
	cb            - func(roomID, ver uint8, cmd uint16, message []byte) - Callback to be invoked on every Broadcast and Message.

	Usage Example: Useful when capturing and sending the message data to external database or service.
*/
func SetOnAnnounce(roomID string, cb func(string, uint8, uint16, []byte)) bool {
	return false
}

/*
DebugDataDump returns an array of all rooms on the server.

	[IMPORTANT] This is for debug purpose only.
*/
func DebugDataDump() []map[string]interface{} {
	return nil
}

/*
GetAllRoomIDsOnNode returns all valid room IDs on the current node

	[NOTE] Uses mutex lock internally.
*/
func GetAllRoomIDsOnNode() []string {
	return nil
}

/*
GetRoomNodeAddressList returns a list of internal node addresses that the room is stored

	roomID - Target room ID to extract the server internal node addresses from.

Returns a list of server internal node addresses of the room and error if it fails.
*/
func GetRoomNodeAddressList(roomID string) ([]string, error) {
	return nil, nil
}

/*
IsUserInRoom returns true if the user is in the room

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.
	[NOTE] Uses mutex lock internally.

Parameters

	roomID   - Target room ID to check.
	userData - Target user to see if the user is in the room.
*/
func IsUserInRoom(roomID string, userData *user.User) bool {
	return false
}

/*
IsRoomFull returns true if the room is full.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID - Target room ID to check.
*/
func IsRoomFull(roomID string) bool {
	return false
}

/*
GetRoomID returns roomID that the user has joined

	[IMPORTANT] This function can be used only on the server that the room exists.
	[NOTE] Uses mutex lock internally.

Parameters

	userData - Target user to get the room ID that the user is currently member of.
*/
func GetRoomID(userData *user.User) string {
	return ""
}

/*
GetRoomByID returns a room

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.
	[NOTE] This will not update TTL of the room.

Parameters

	roomID - Target room ID to get the room object.
*/
func GetRoomByID(roomID string) *Room {
	return nil
}

/*
GetCreatedTime returns the created time (in seconds) of the room

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.
*/
func GetCreatedTime(roomID string) int64 {
	return 0
}

/*
GetRoomOwnerID returns the user ID of the owner of the room

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID - Target room ID to get the owner from.
*/
func GetRoomOwnerID(roomID string) string {
	return ""
}

/*
Exists return true if the given room ID's room still exists

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID - Target room ID to check if the room exists or not.
*/
func Exists(roomID string) bool {
	return false
}

/*
ChangeTTL changes TTL of the targeted room

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[NOTE] Uses mutex lock internally.

Parameters

	roomID - Target room ID to change its TTL.
	ttl    - New TTL of the room.
*/
func ChangeTTL(roomID string, ttl int64) bool {
	return false
}

/*
GetMemberIDsAndOwner returns an array of all member IDs and owner ID and ownerSID

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally for both room and user.

Parameters

	roomID - Target room ID to get the list of member user IDs and its owner user ID from.

	Returns member user IDs as an array, owner user ID, and owner user SID.
*/
func GetMemberIDsAndOwner(roomID string) ([]string, string, string) {
	return nil, "", ""
}

/*
GetMemberIDs returns the list of member IDs (not SIDs)

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally for both room and user.

Parameters

	roomID - Target room ID to get the list of member user IDs.
*/
func GetMemberIDs(roomID string) []string {
	return nil
}

/*
GetMemberSIDsAndOwner returns an array of all member IDs and owner ID and owner SID

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally for both room and user.

Parameters

	roomID - Target room ID to get the list of member SIDs, owner user ID, and owner SID.

Returns an array of member user SIDs, owner user ID, and owner SID.
*/
func GetMemberSIDsAndOwner(roomID string) ([]string, string, string) {
	return nil, "", ""
}

/*
GetMemberSIDs returns the list of member SIDs

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID - Target room ID to get a list of member SIDs.
*/
func GetMemberSIDs(roomID string) []string {
	return nil
}

/*
GetMemberUsers returns the list of member user struct instances

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.
	[NOTE] Returned *user.User may be nil if the user has disconnected.

Parameters

	roomID - Target room ID to get the list of member userData.
*/
func GetMemberUsers(roomID string) []*user.User {
	return nil
}

/*
GetMemberUsersAndOwnerUser returns an array of all member IDs and owner ID

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.
	[NOTE] Returned *user.User may be nil if the user has disconnected.

Parameters

	roomID - Target room ID to get member *user.User and owner *user.User.

Returns member *user.User as an array and owner *user.User.
*/
func GetMemberUsersAndOwnerUser(roomID string) ([]*user.User, *user.User) {
	return nil, nil
}

/*
GetMemberByID returns member user by room ID and user ID

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[NOTE] Uses mutex lock internally.

Parameters

	roomID - Target room ID to get the member *user.User from.
	id     - Target user ID to get *user.User of.
*/
func GetMemberByID(roomID string, id string) *user.User {
	return nil
}

/*
GetNumberOfRoomMembers retrieves number of room members.

This can be executed by anyone and they do not have to be part of the room.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.
	[IMPORTANT] This function works for all users including the users that are not members of the target room.

Parameters

	roomID   - Target room ID to get the number of members from.
	callback - Callback to be invoked when the number of current members and max members retrieved.
	           func(err error, currentMembers int, maxMembers int)
*/
func GetNumberOfRoomMembers(roomID string, callback func(err error, currentMembers int, maxMembers int)) {
}
