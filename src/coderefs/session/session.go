package session

import (
	"github.com/Diarkis/diarkis/user"
)

/*
GetID returns the session ID.
*/
func (s *Session) GetID() string {
	return ""
}

/*
GetType returns the type of the session.
*/
func (s *Session) GetType() uint8 {
	return 0
}

/*
GetMemberIDs returns an array of member IDs.
*/
func (s *Session) GetMemberIDs() []string {
	return nil
}

/*
GetMemberSIDs returns an array of member SIDs.
*/
func (s *Session) GetMemberSIDs() []string {
	return nil
}

/*
IsMemberByUID returns true if the given user ID is the UID of the session member user.
*/
func (s *Session) IsMemberByUID(uid string) bool {
	return false
}

/*
IsJoinAllowed returns true if the session has room to add another user.
*/
func (s *Session) IsJoinAllowed(memberCnt int) bool {
	return false
}

/*
GetMemberSIDByUID returns the member's SID.

It returns and empty string if the member is not found or invalid.
*/
func (s *Session) GetMemberSIDByUID(uid string) string {
	return ""
}

/*
GetMemberUsers returns an array of member user copies.

	[IMPORTANT] The return array contains copies of member users.
*/
func (s *Session) GetMemberUsers() []*user.User {
	return nil
}

/*
GetMemberMeshAddrList returns an array of internal server address of each user.
*/
func (s *Session) GetMemberMeshAddrList() []string {
	return nil
}

/*
GetMemberMeshAddrByUID returns a mesh address of a member.

Returns an empty string if the member is not found or invalid.
*/
func (s *Session) GetMemberMeshAddrByUID(uid string) string {
	return ""
}

/*
GetOwnerUser returns the ticket owner user.

	[IMPORTANT] Returned owner user is NOT a copy.
*/
func (s *Session) GetOwnerUser() (*user.User, bool) {
	return nil, false
}

/*
SetOnDelete assigns a callback on session deletion that is invoked before the session is deleted.

	[IMPORTANT] Only one callback can be assigned to a session.
*/
func (s *Session) SetOnDelete(cb func(id string)) bool {
	return false
}

/*
SetOnDeleted assigns a callback on session deletion that is invoked after the session is deleted.

	[IMPORTANT] Only one callback can be assigned to a session.
*/
func (s *Session) SetOnDeleted(cb func(id string)) bool {
	return false
}

/*
SetOnJoin assigns a callback on session to be invoked when a new member is attempting to join the session.

The callback returns a bool and if you return false, the user will be rejected and will not join the session.

	[IMPORTANT] Only one callback can be assigned to a session.
	[IMPORTANT] userData passed to the callback is a copy of the actual user is attempting to join.
*/
func (s *Session) SetOnJoin(cb func(id string, userData *user.User) bool) bool {
	return false
}

/*
SetOnJoined assigns a callback on session to be invoked when a new member joins.

	[IMPORTANT] Only one callback can be assigned to a session.
	[IMPORTANT] If you need to use Broadcast in the callback,
	            you must pass the owner user to Broadcast because the newly joined user does not have the session ID yet.
	[IMPORTANT] userData passed to the callback is a copy of the actual user that joined.
*/
func (s *Session) SetOnJoined(cb func(id string, userData *user.User)) bool {
	return false
}

/*
SetOnLeft assigns a callback on session to be invoked when a member of user leaves the session.

	[IMPORTANT] Only one callback can be assigned to a session.
	[IMPORTANT] If you need to use Broadcast in the callback,
	            you must pass the owner user to Broadcast because the user left no longer has the session ID.
	[IMPORTANT] userData passed to the callback is a copy of the actual user that left.
*/
func (s *Session) SetOnLeft(cb func(id string, userData *user.User)) bool {
	return false
}

/*
SetOnTickStop assigns a callback to be invoked when a tick of the session stops.

	[IMPORTANT] Only one callback can be assigned to a session.
*/
func (s *Session) SetOnTickStop(interval uint16, cb func(id string)) bool {
	return false
}

/*
SetOnTick assigns a callback to be invoked at every given interval.

	[IMPORTANT] A single callback may be assigned per tick interval.
	            You may not assign multiple callbacks a tick with the same interval.

Parameters

	interval - Tick interval in seconds.
	cb       - Callback to be invoked at every tick.
*/
func (s *Session) SetOnTick(interval uint16, cb func(id string)) bool {
	return false
}

/*
StopAllTicks stops all tick loops.
*/
func (s *Session) StopAllTicks() bool {
	return false
}

/*
SetPropertyIfNotExists stores a key and value data to session if the same key does not exist.

Properties are only primitive values and does not support reference type data such as array and map.
*/
func (s *Session) SetPropertyIfNotExists(key string, value interface{}) bool {
	return false
}

/*
SetProperty stores a key and value data to session.

If the same key exists, it overwrites the existing value of the same key.

Properties are only primitive values and does not support reference type data such as array and map.
*/
func (s *Session) SetProperty(key string, value interface{}) {
}

/*
SetProperties stores a collection of keys and their values to session.

If the same key exists, it overwrites the existing value of the same key.

Properties are only primitive values and does not support reference type data such as array and map.
*/
func (s *Session) SetProperties(data map[string]interface{}) error {
	return nil
}

/*
UpdateProperty changes the existing property value of session.

The callback is invoked while the internal lock is still held, locking inside the callback may cause mutex deadlock.

Properties are only primitive values and does not support reference type data such as array and map.

	key   - A key of the property to be updated.
	value - A value of the property to be updated with.
	cb    - Callback to be invoked on every key and value pair to handle the update.
	        func(exists bool, storedValue interface{}, updateValue interface{}) (updatedValue interface{})
	         - exists      - Indicates if the same key already exists or not
	         - storedValue - Existing value that is stored as a property. If the key does not exist it is a nil.
	         - updateValue - The value to be used to update/replace or set.
*/
func (s *Session) UpdateProperty(key string, value interface{}, cb func(bool, interface{}, interface{}) interface{}) {
}

/*
UpdateProperties changes the existing property values of session.

The callback is invoked while the internal lock is still held, locking inside the callback may cause mutex deadlock.

Properties are only primitive values and does not support reference type data such as array and map.

	data - A map of key and value pair to be stored as properties.
	cb   - Callback to be invoked on every key and value pair to handle the update.
	       func(exists bool, storedValue interface{}, updateValue interface{}) (updatedValue interface{})
	         - exists      - Indicates if the same key already exists or not
	         - storedValue - Existing value that is stored as a property. If the key does not exist it is a nil.
	         - updateValue - The value to be used to update/replace or set.
*/
func (s *Session) UpdateProperties(data map[string]interface{}, cb func(bool, interface{}, interface{}) interface{}) {
}

/*
GetProperty returns the value of the given key and if the key does not exist, the second return value will be a false.

Properties are only primitive values and does not support reference type data such as array and map.

The returned property value is an interface{}, in order to type assert safely, please use Diarkis' util package functions.

Example:

	v, ok := r.GetProperty("someKey")

	if !ok {
	  // handle error here
	}

	// If the value data type is an uint8, of course ;)
	v, ok := util.ToUint8(v)
*/
func (s *Session) GetProperty(key string) (interface{}, bool) {
	return nil, false
}

/*
GetProperties returns key and value pairs as a map.

Properties are only primitive values and does not support reference type data such as array and map.

If a value of a given key does not exist, the returned map will have a nil as a value of the key.

The returned property value is an interface{}, in order to type assert safely, please use Diarkis' util package functions.

Example:

	values, ok := r.GetProperties([]string{ "someKey" })

	if !ok {
	  // handle error here
	}

	for key, v := range values {
	  // If the value data type is an uint8, of course ;)
	  value, ok := util.ToUint8(v)
	}
*/
func (s *Session) GetProperties(keys []string) map[string]interface{} {
	return nil
}
