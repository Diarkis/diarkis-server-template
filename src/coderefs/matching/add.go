package matching

/*
AddData represents internally used matchmaking data
*/
type AddData struct {
	InternalID string                 `json:"internalID"`
	MatchingID string                 `json:"matchingID"`
	Tag        string                 `json:"tag"`
	Props      map[string]int         `json:"props"`
	Value      map[string]interface{} `json:"value"`
	TTL        int64                  `json:"ttl"`
}

/*
Add Adds a searchable matching candidate data.

	[NOTE] In order to have long lasting (TTL longer than 60 seconds) searchable items,
	       the application must "add" searchable items repeatedly with TTL=60.

	[NOTE] Uses mutex lock internally.

Error Cases

	╒═════════════════════╤══════════════════════════════════════════════════════════╕
	│ Error               │ Reason                                                   │
	╞═════════════════════╪══════════════════════════════════════════════════════════╡
	│ Invalid profile IDs │ Either missing or invalid matchmaking profile IDs given. │
	├─────────────────────┼──────────────────────────────────────────────────────────┤
	│ Unique ID missing   │ Unique ID must not be an empty string.                   │
	├─────────────────────┼──────────────────────────────────────────────────────────┤
	│ Properties missing  │ Properties must not be a nil.                            │
	├─────────────────────┼──────────────────────────────────────────────────────────┤
	│ Invalid limit value │ Limit must be greater than 1.                            │
	╘═════════════════════╧══════════════════════════════════════════════════════════╛

Parameters

	mid      - Matching profile ID to add the searchable data to.
	unqiueID - Unique ID of the searchable ID.
	tag      - Tag is used to isolate and group add and search of matchmaking.
	           If an empty string is given, it will be ignored.
	props    - Searchable condition properties.
	value    - Searchable data to be returned along with the search results.
	ttl      - TTL of the searchable item to be added in seconds. Maximum 60 seconds.
	limit    - Number of search node to propagate the searchable item data at a time.
*/
func Add(mid, uqid, tag string, props map[string]int, value map[string]interface{}, ttl int64, limit int) error {
	return nil
}
