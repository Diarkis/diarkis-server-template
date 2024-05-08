package synchronizer

import (
	"sync"
)

/*
Properties represents synchronizer based on properties.

Property synchronizer allows its members to updated shared properties and synchronize at fixed internal in milliseconds.

Each synchronized property value will invoke the callback and the application may decide how it wants to synchronize the property.

	[IMPORTANT] If a property is updated multiple times within the time window of synchronize interval, only the last updated value of the property is synchronized.
*/
type Properties struct{ sync.RWMutex }

/*
EncodeData is used encode Properties
*/
type EncodeData struct {
	UpdateVer    uint8  `json:"updateVer"`
	RemoveVer    uint8  `json:"removeVer"`
	UpdateCmd    uint16 `json:"updateCmd"`
	RemoveCmd    uint16 `json:"removeCmd"`
	Interval     int64  `json:"interval"`
	SyncAllDelay int64  `json:"syncAllDelay"`
	TTL          int64  `json:"ttl"`
	CombineLimit int    `json:"combineLimit"`
	Encoded      []byte `json:"encoded"`
}

/*
NewProperties creates a new Properties instance.

Parameters

	updateVer    - Server push command version for property update.
	removeVer    - Server push command version for property removal.
	updateCmd    - Server push command ID for property update.
	removeCmd    - Server push command ID for property removal.
	interval     - Synchronization interval in milliseconds.
	               Minimum value for interval is 17 milliseconds.
	syncAllDelay - Delays the invocation of synchronization of all properties in milliseconds.
	ttl          - TTL of property. If a property is not updated for more than TTL, it will be removed internally
	               TTL is in milliseconds and the minimum allowed value is 500.
	               By passing TTL lower than the minimum value of 500, you will disable TTL function for properties.
	combineLimit - Maximum number of outbound synchronization packets to be combined per member.
	               Default is 30.
*/
func NewProperties(updateVer, removeVer uint8, updateCmd, removeCmd uint16, interval, syncAllDelay, ttl int64, combineLimit uint8) *Properties {
	return nil
}

/*
Encode encodes its data into byte array for transport and storage
*/
func Encode(p *Properties) ([]byte, error) {
	return nil, nil
}

/*
Decode decodes encoded data byte array and overwrites its members with the decoded data
*/
func Decode(encoded []byte) (*Properties, error) {
	return nil, nil
}

/*
GetInterval returns the synchronization interval in milliseconds.
*/
func (p *Properties) GetInterval() int64 {
	return 0
}

/*
AddMember adds a user as a member to join by SID
*/
func (p *Properties) AddMember(sid string) error {
	return nil
}

/*
SyncAllProperties sends all existing properties to the target user by the given SID.
*/
func (p *Properties) SyncAllProperties(sid string) {
}

/*
RemoveMember removes a user by SID

How to remove user member when the user client disconnects:

	// "github.com/Diarkis/diarkis/user"

	// This must be written to be executed when the server process starts
	user.OnDiscard(func(userData *user.User) {
		// this is the user that is disconnected
		properties.RemoveMember(userData.SID)
	})
*/
func (p *Properties) RemoveMember(sid string) bool {
	return false
}

/*
Reset clears all properties and members
*/
func (p *Properties) Reset() {
}

/*
Count returns the number of existing properties
*/
func (p *Properties) Count() int {
	return 0
}

/*
IsLinkedToUser returns true if the given key is linked the given user SID.
*/
func (p *Properties) IsLinkedToUser(key string, sid string) bool {
	return false
}

/*
UpdateProperty updates a property by key to be synchronize.

Synchronization payload format:

Fragment [1]

	+---------------------+--------------+-----------------------+----------------+
	| size header for key | property key | size header for value | property value |
	+---------------------+--------------+-----------------------+----------------+
	|       1 byte        |   variable   |        1 byte         |    variable    |
	+---------------------+--------------+-----------------------+----------------+

Fragment [1] will repeat for as long as there are updated properties on every synchronization tick.
*/
func (p *Properties) UpdateProperty(key string, value []byte, reliable bool) error {
	return nil
}

/*
UpsertProperty updates a property by key to be synchronize using existing value.

It will simply set a new value if there is no existing value for the key given.

Synchronization payload format:

Fragment [1]

	+---------------------+--------------+-----------------------+----------------+
	| size header for key | property key | size header for value | property value |
	+---------------------+--------------+-----------------------+----------------+
	|       1 byte        |   variable   |        1 byte         |    variable    |
	+---------------------+--------------+-----------------------+----------------+

Fragment [1] will repeat for as long as there are updated properties on every synchronization tick.
*/
func (p *Properties) UpsertProperty(key string, reliable bool, cb func(exists bool, current []byte) []byte) error {
	return nil
}

/*
UpdatePropertyAsMemberProperty updates a property by key to be synchronize.

The property will be treated as the user's property and it will be removed when the user is removed as a member.

Synchronization payload format:

Fragment [1]

	+---------------------+--------------+-----------------------+----------------+
	| size header for key | property key | size header for value | property value |
	+---------------------+--------------+-----------------------+----------------+
	|       1 byte        |   variable   |        1 byte         |    variable    |
	+---------------------+--------------+-----------------------+----------------+

Fragment [1] will repeat for as long as there are updated properties on every synchronization tick.
*/
func (p *Properties) UpdatePropertyAsMemberProperty(sid string, key string, value []byte, reliable bool) error {
	return nil
}

/*
RemoveProperty removes the property of the given key.

Synchronization payload format:

Fragment [1]

	+---------------------+--------------+
	| size header for key | property key |
	+---------------------+--------------+
	|       4 bytes       |   variable   |
	+---------------------+--------------+

Fragment [1] will repeat for as long as there are updated properties on every synchronization tick.
*/
func (p *Properties) RemoveProperty(key string) error {
	return nil
}

/*
Stop stops the synchronization loop
*/
func (p *Properties) Stop() bool {
	return false
}

/*
Start starts the synchronization loop
*/
func (p *Properties) Start() bool {
	return false
}
