package matching

/*
RemoveData represents internally used matchmaking removal data
*/
type RemoveData struct {
	MatchingID string   `json:"matchingID"`
	UniqueIDs  []string `json:"uniqueIDs"`
}

/*
Remove removes a list of searchable matching candidate items by their unique IDs

	[NOTE] This function is very expensive as it will send a message to all related mesh nodes
	[NOTE] There maybe some delay with removal operation on multiple servers
	       as the instruction to remove must traverse all servers in the cluster.

Parameters

	matchingID   - Target matching profile ID to remove items from
	uniqueIDList - A list of unique IDs of the searchable items to remove
	limit        - Mesh network relay limit
*/
func Remove(matchingID string, uniqueIDList []string, limit int) {
}
