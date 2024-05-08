package payload

/*
MMAdd represents MatchMaker add data
*/
type MMAdd struct {
	TTL      int64
	ID       string
	UID      string
	Props    map[string]int
	Metadata []byte
}

/*
MMRemove represents MatchMaker remove data
*/
type MMRemove struct {
	ID   string
	UIDs []string
}

/*
MMSearch represents MatchMaker search data
*/
type MMSearch struct {
	IDs   []string
	Props map[string]int
}

/*
PackMMAdd encodes the given parameters into byte array
*/
func PackMMAdd(mmID, uniqueID string, props map[string]int, metadata []byte, ttl uint64) []byte {
	return nil
}

/*
UnpackMMAdd decodes the given byte array into an instance of MMAdd
*/
func UnpackMMAdd(bytes []byte) *MMAdd {
	return nil
}

/*
PackMMRemove encodes the given parameters into a byte array
*/
func PackMMRemove(mmID string, uniqueIDs []string) []byte {
	return nil
}

/*
UnpackMMRemove decodes the given byte array to an instance of MMRemove
*/
func UnpackMMRemove(bytes []byte) *MMRemove {
	return nil
}

/*
PackMMSearch encodes the given parameters into a byte array
*/
func PackMMSearch(mmIDs []string, props map[string]int) []byte {
	return nil
}

/*
UnpackMMSearch decodes the given byte array into an instance of MMSearch
*/
func UnpackMMSearch(bytes []byte) *MMSearch {
	return nil
}
