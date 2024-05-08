package packet

/*
MMAdd data structure of MatchMaker's matchmaking add and wait command

	BigEndian
	+-------------+--------------+--------+-------------+-----+-------------+-----+-----+
	| Max Members | Reserve Flag |   TTL  | Size of (1) | (1) | Size of (2) | (2) | (3) |
	+-------------+--------------+--------+-------------+-----+-------------+-----+-----+
	|   2 bytes   |    1 byte    | 2 byte |    2 byte   | (1) |    2 byte   | (2) | (3) |
	+-------------+--------------+--------+-------------+-----+-------------+-----+-----+

	(1) - Matchmaking profile ID, UID, and Tag
	+------------------------+------------------------+-------------+-------------+-------------+-------------+
	| Size of matchmaking ID |      Matchmaking ID    | Size of UID |     UID     | Size of Tag |     Tag     |
	+------------------------+------------------------+-------------+-------------+-------------+-------------+
	|         4 byte         | Size of matchmaking ID |    4 byte   | Size of UID |    4 byte   | Size of Tag |
	+------------------------+------------------------+-------------+-------------+-------------+-------------+

	(2) - Properties
	For multiple propery values, value, size of name, and name repeats as a data set
	+---------------+-----------------------+-----------------------+
	| Propety Value | Size of property name | Propery name          | ...
	+---------------+-----------------------+-----------------------+
	|     4 byte    |          2 byte       | Size of property name | ...
	+---------------+-----------------------+-----------------------+

	(3) - Metadata
	+--------------------------------------------------------------+
	|                         Metadata                             |
	+--------------------------------------------------------------+
	| From the end of properties to the end of the data byte array |
	+--------------------------------------------------------------+
*/
type MMAdd struct {
	TTL         int64
	ID          string
	UID         string
	Tag         string
	Props       map[string]int
	Metadata    []byte
	MaxMembers  int
	ReserveOnly bool
}

/*
MMRemove data structure of MatchMaker's remove and abort of matchmaking command

	BigEndian

	+-----------+------------------------+----------------+-------------+-----+-----------------------------------------------+
	| Halt Flag | Size of matchmaking ID | Matchmaking ID | Size of UID | UID |                    Message                    |
	+-----------+------------------------+----------------+-------------+-----+-----------------------------------------------+
	|   1 byte  |         4 byte         |                |    4 byte   |     | From end of UID to the end of data byte array |
	+-----------+------------------------+----------------+-------------+-----+-----------------------------------------------+
*/
type MMRemove struct {
	HaltFlag bool
	ID       string
	UIDs     []string
	Message  []byte
}

/*
MMSearch data structure of MatchMaker's search command

	BigEndian

	+------------------+-----------+-------------------------+-----------------+-------------+-----+--------------------+-----+-----+
	| How many results | Join Flag | Size of matchmaking IDs | Matchmaking IDs | Size of Tag |  Tag | Size of properties | (1) | (2) |
	+------------------+-----------+-------------------------+-----------------+-------------+------+--------------------+-----------+
	|      2 byte      |   1 byte  |          2 byte         |                 |    2 byte   |      |         2 byte     | (1) | (2) |
	+------------------+-----------+-------------------------+-----------------+-------------+------+--------------------+-----+-----+

	(1) - Properties
	For multiple propery values, value, size of name, and name repeats as a data set
	+---------------+-----------------------+-----------------------+
	| Propety Value | Size of property name | Propery name          | ...
	+---------------+-----------------------+-----------------------+
	|     4 byte    |          2 byte       | Size of property name | ...
	+---------------+-----------------------+-----------------------+

	(2) Message
	+------------------------------------------------------+
	|                     Message                          |
	+------------------------------------------------------+
	| From end of properties to the end of data byte array |
	+------------------------------------------------------+
*/
type MMSearch struct {
	HowMany int
	Join    bool
	IDs     []string
	Tag     string
	Props   map[string]int
	Message []byte
}

/*
MMClaim MatchMaker's reservation claim command data structure

	BigEndian

	+-----------------+---------+----------------------------------------------+
	| Size of room ID | Room ID |                     Message                  |
	+-----------------+---------+----------------------------------------------+
	|      2 byte     |         | end of room ID to the end of data byte array |
	+-----------------+---------+----------------------------------------------+
*/
type MMClaim struct {
	RoomID  string
	Message []byte
}

/*
MMIssueTicket data structure of MatchMaker's issue ticket command

	BigEndian

	+-------------------------+-----------------+--------------------+-----+-----+
	| Size of matchmaking IDs | Matchmaking IDs | Size of properties | (1) | (2) |
	+-------------------------+-----------------+--------------------+-----------+
	|          2 byte         |                 |         2 byte     | (1) | (2) |
	+-------------------------+-----------------+--------------------+-----+-----+

	(1) - Properties
	For multiple property values, value, size of name, and name repeats as a data set
	+----------------+-----------------------+-----------------------+
	| Property Value | Size of property name | Property name         | ...
	+----------------+-----------------------+-----------------------+
	|     4 byte     |          2 byte       | Size of property name | ...
	+----------------+-----------------------+-----------------------+
*/
type MMIssueTicket struct {
	IDs   []string
	Props map[string]int
}

/*
PackMMAdd packs command data into data byte array
*/
func PackMMAdd(mmID, uniqueID string, tag string, maxMembers uint16, reserveOnly bool, props map[string]int, metadata []byte, ttl uint16) []byte {
	return nil
}

/*
UnpackMMAdd unpacks the data byte array to data structure
*/
func UnpackMMAdd(bytes []byte) *MMAdd {
	return nil
}

/*
PackMMRemove packs the command data into data byte array
*/
func PackMMRemove(mmID string, haltFlag bool, uniqueIDs []string, msg []byte) []byte {
	return nil
}

/*
UnpackMMRemove unpacks the data byte array to data struct
*/
func UnpackMMRemove(bytes []byte) *MMRemove {
	return nil
}

/*
PackMMSearch packs the command data into data byte array
*/
func PackMMSearch(howmany uint16, joinFlag bool, mmIDs []string, tag string, props map[string]int, msg []byte) []byte {
	return nil
}

/*
UnpackMMSearch unpacks data byte array to data structure
*/
func UnpackMMSearch(bytes []byte) *MMSearch {
	return nil
}

/*
PackMMClaim packs the command data into data byte array
*/
func PackMMClaim(roomID string, message []byte) []byte {
	return nil
}

/*
UnpackMMClaim unpacks command data byte array
*/
func UnpackMMClaim(bytes []byte) *MMClaim {
	return nil
}

/*
PackMMIssueTicket packs the command data into data byte array
*/
func PackMMIssueTicket(mmIDs []string, props map[string]int) []byte {
	return nil
}

/*
UnpackMMIssueTicket unpacks the command data array
*/
func UnpackMMIssueTicket(bytes []byte) *MMIssueTicket {
	return nil
}
