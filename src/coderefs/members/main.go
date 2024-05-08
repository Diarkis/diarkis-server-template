package members

import (
	"github.com/Diarkis/diarkis/user"
)

const NotLimited = -1
const Unexpected = -2

/*
Members a container to store member data.
*/
type Members struct{}

/*
Client is a representation of the user client with its IP addresses and IDs.
The struct is derived from *user.User
*/
type Client struct {
	SID              string
	ID               string
	MeshAddr         string
	PublicAddr       string
	PrivateAddrBytes []byte
	UserData         map[string]interface{}
}

/*
UnlockExclusive forcefully unlocks members exclusive
*/
func UnlockExclusive(userData *user.User, key string) {
}

/*
IsUserInRoom returns true if the user with the matching name is in a members container.
*/
func IsUserInRoom(name string, userData *user.User) bool {
	return false
}

/*
ResetUser resets the user key data with the matching name.
*/
func ResetUser(name string, userData *user.User) {
}

/*
New creates a members container without a limit to the number of members.
*/
func New(name string) *Members {
	return nil
}

/*
NewLimited creates a members container with a limit to the number of members.
*/
func NewLimited(name string, limit uint) *Members {
	return nil
}

/*
NewExclusive creates an exclusive members container
that does NOT allow the same user to be in another members container
without a limit to the number of members.
*/
func NewExclusive(name string) *Members {
	return nil
}

/*
NewExclusiveLimited creates an exclusive members container
that does NOT allow the same user to be in another members container
with a limit to the number of members.
*/
func NewExclusiveLimited(name string, limit uint) *Members {
	return nil
}

/*
EncodeClient converts the client struct instance into a string.
Use DecodeClient to revert it back to the client struct instance.
*/
func EncodeClient(client *Client) string {
	return ""
}

/*
DecodeClient reverts back the encoded client struct instance.
It returns a nil if the decoding fails.
*/
func DecodeClient(str string) *Client {
	return nil
}

/*
GetName returns the name of members container.
*/
func (m *Members) GetName() string {
	return ""
}

/*
AddReserve makes a reservation for the given user to make sure that the user may join the members container later for sure.
*/
func (m *Members) AddReserve(userData *user.User) bool {
	return false
}

/*
AddReserveByUserID makes a reservation for the given user ID to make sure that the user may join the members container later for sure.
*/
func (m *Members) AddReserveByUserID(userID string) bool {
	return false
}

/*
RemoveReserve cancels the reservation for the given user.
*/
func (m *Members) RemoveReserve(userData *user.User) bool {
	return false
}

/*
RemoveReserveByUserID cancels the reservation for the given user ID.
*/
func (m *Members) RemoveReserveByUserID(userID string) bool {
	return false
}

/*
Add adds a userData.ID as key and userData.SID as value to the container.
*/
func (m *Members) Add(userData *user.User, meshAddr string) error {
	return nil
}

/*
AddIf allows the user given to be added to the members if the given callback returns true.
*/
func (m *Members) AddIf(userData *user.User, meshAddr string, check func() bool) error {
	return nil
}

/*
IsAddAllowed returns true if the members contain still has room to add more users.
*/
func (m *Members) IsAddAllowed(clientsToBeAdded int) bool {
	return false
}

/*
GetExclusiveLockKey returns the key of user exclusive lock
*/
func (m *Members) GetExclusiveLockKey() string {
	return ""
}

/*
UpdateClientByID updates internal client data
*/
func (m *Members) UpdateClientByID(uid, sid, pubAddr string, privateAddrs []byte, meshAddr string) bool {
	return false
}

/*
Remove removes a member data from the container.
*/
func (m *Members) Remove(userData *user.User) bool {
	return false
}

/*
RemoveAll removes all member users form the members container.
*/
func (m *Members) RemoveAll() {
}

/*
Exists returns true if the given key exists in the members container.
*/
func (m *Members) Exists(userID string) bool {
	return false
}

/*
Count returns the number of members in the container.
*/
func (m *Members) Count() int {
	return 0
}

/*
GetOne returns the member user of the given user ID.
If the key does not exist, it returns a nil instead.
*/
func (m *Members) GetOne(userID string) *Client {
	return nil
}

/*
GetAll returns all member SIDs as a map.
The key of the map is user ID and the value of the map is user.
*/
func (m *Members) GetAll() map[string]*Client {
	return nil
}

/*
Keys returns all member keys as an array.
The keys are an array of user IDs.
*/
func (m *Members) Keys() []string {
	return nil
}

/*
GetLimit returns the maximum allowed number of members.

	[NOTE] If there is no limit, it returns NotLimited -1.
	[NOTE] If *Member is somehow nil, it returns Unexpected -2.
*/
func (m *Members) GetLimit() int {
	return 0
}
