package p2p

import (
	"github.com/Diarkis/diarkis/user"
)

/*
StartP2PFromRoom starts P2P connections using Room. Returns an error and error bytes if unsuccessful.
*/
func StartP2PFromRoom(ver uint8, cmd uint16, payload []byte, userData *user.User) ([]byte, error) {
	return nil, nil
}

/*
CreateLinkGroupsFromRoom creates an array of link list of users to be used for P2P Links message relay

The users will be sorted by their server-to-client latency and will be grouped accordingly
*/
func CreateLinkGroupsFromRoom(roomID string, maxLinkMembers int) ([][]*user.User, [][]string) {
	return nil, nil
}

/*
RelayLinkMessage sends a message via room to selected users to relay among their link groups
*/
func RelayLinkMessage(roomID string, ver uint8, cmd uint16, message []byte, userData *user.User, reliable bool) bool {
	return false
}

/*
CreateP2PSyncBytes creates a byte array with encryption keys and a list of user client IDs and addresses.

*
* === Encryption Key Part ===
* +----------------+---------------+--------------------+
* | Encryption Key | Encryption IV | Encryption Mac Key |
* +----------------+---------------+--------------------+
* |     16 bytes   |    16 bytes   |       16 bytes     |
* +----------------+---------------+--------------------+
*
* === User client ID and address Part ===
* User ID and User client address may be repeated as a set.
* Big Endian
* +--------------+---------+---------------------------+-----------------------+
* | User ID Size | User ID | Base64 Encode String Size | Base64 Encoded String |
* +--------------+---------+---------------------------+-----------------------+
* |    4 bytes   |         |           4 bytes         |                       |
* +--------------+---------+---------------------------+-----------------------+
*
* === Base64 Encoded String ===
* address list as a set.
* The first address is always the client public address and the rest is local address.
* Big Endian
* +--------------+---------+
* | Address Size | Address |
* +--------------+---------+
* |    4 bytes   |         |
* +--------------+---------+
*
* === Outcome byte array ===
* Currently used
* +---------------------+---------------------------------+
* | Encryption Key Part | User client ID and address Part |
* +---------------------+---------------------------------+
* |        48 bytes     |          Variable byte size     |
* +---------------------+---------------------------------+
*
* Currently NOT used
* +--------------------------------------+---------------------------------+---------------------+
* | User client ID and address Part size | User client ID and address Part | Encryption Key Part |
* +--------------------------------------+---------------------------------+---------------------+
* |           2 bytes big endian         |        Variable byte size       |      48 bytes       |
* +--------------------------------------+---------------------------------+---------------------+
*/
func CreateP2PSyncBytes(roomID string, users []*user.User) ([]byte, error) {
	return nil, nil
}
