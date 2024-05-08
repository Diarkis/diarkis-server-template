package matching

/*
SearchData represents internally used matchmaking data
*/
type SearchData struct {
	MatchingID string           `json:"matchingID"`
	Tag        string           `json:"tag"`
	Props      map[string][]int `json:"props"`
	Limit      int              `json:"limit"`
}

/*
SearchReturnData represents internally used data
*/
type SearchReturnData struct {
	Results []*SearchReturnItemData `json:"results"`
}

/*
SearchReturnItemData represents internally used data
*/
type SearchReturnItemData struct {
	ID    string      `json:"id"`
	TTL   int64       `json:"ttl"`
	Value interface{} `json:"value"`
}

/*
SearchWithRange searches for matched data based on the given values of props in range.

Each props will have an array of elements with property values that should range from minimum value to maximum value.
The function will use those values per property to performance the search.

	[IMPORTANT] SearchWithRange performs search operations on remote servers
	            and the number of servers to perform the search is calculated based on the number of the server and distributionRate configuration.
	            This means that there is a chance that SearchWithRange function may miss the server(s) with the desired matchmaking data
	            resulting in not finding the intended matchmaking data.

	[IMPORTANT] The number of allowed range properties is limited to two.

Example:

	searchProps["level"]     = []int{ 1, 2, 3, 4, 5 } // range property
	searchProps["rank"]      = []int{ 1, 2, 3 } // range property
	searchProps["matchType"] = []int{ 1 } // regular property
	searchProps["league"]    = []int{ 10 } // regular property

Error Cases

	╒═════════════════════╤══════════════════════════════════════════════════════════════╕
	│ Error               │ Reason                                                       │
	╞═════════════════════╪══════════════════════════════════════════════════════════════╡
	│ Invalid profile IDs │ Either missing or invalid matchmaking profile IDs give.      │
	├─────────────────────┼──────────────────────────────────────────────────────────────┤
	│ Invalid properties  │ Either missing or invalid properties given.                  │
	├─────────────────────┼──────────────────────────────────────────────────────────────┤
	│ Reliable timeout    │ Communication to another server process timed out or failed. │
	╘═════════════════════╧══════════════════════════════════════════════════════════════╛

Parameters

	profileIDList - A list of matchmaking profile IDs to search by. The order of the list is the order of search attempts.
	tag           - A tag is used to isolate and group search meaning the search will not include tags that do not match.
	searchProps   - A map with search condition values.
	                You are allowed to have maximum two properties with range (more than 1 element in the int array).
	                If you have more than two elements with range, it will give you an error.
	limit         - A hard limit to the number of results.
	callback      - A callback with the search results or an error.

Diagram below shows how it works:

	     ┏━━━┓
	     ┃ 8 ┃ 8 with ±3 would fall into 5 ~ 11 and that means it matches with items in the bucket of 0 to 10 and 11 to 20
	     ┗━┳━┛
	     ┏━┻━━━━┓
	│    ▼    │ ▼       │         │
	│  0 ~ 10 │ 11 ~ 20 │ 21 ~ 30 │
	└─────────┴─────────┴─────────┘

To make ranged search more precise, consider using SetOnTicketAllowMatchIf callback for MatchMaker Ticket.
*/
func SearchWithRange(profileIDList []string, tag string, searchProps map[string][]int, limit int, callback func(err error, results []interface{})) {
}

/*
Search searches for matched data based on the given props' values

	[IMPORTANT] Search performs search operations on remote servers
	            and the number of servers to perform the search is calculated based on the number of the server and distributionRate configuration.
	            This means that there is a chance that Search function may miss the server(s) with the desired matchmaking data
	            resulting in not finding the intended matchmaking data.

	[NOTE] Uses mutex lock internally.

Error Cases

	╒═════════════════════╤══════════════════════════════════════════════════════════════╕
	│ Error               │ Reason                                                       │
	╞═════════════════════╪══════════════════════════════════════════════════════════════╡
	│ Invalid profile IDs │ Either missing or invalid matchmaking profile IDs give.      │
	├─────────────────────┼──────────────────────────────────────────────────────────────┤
	│ Invalid properties  │ Either missing or invalid properties given.                  │
	├─────────────────────┼──────────────────────────────────────────────────────────────┤
	│ Reliable timeout    │ Communication to another server process timed out or failed. │
	╘═════════════════════╧══════════════════════════════════════════════════════════════╛

Parameters

	profileIDList - A list of matchmaking profile IDs to search by. The order of the list is the order of search attempts.
	tag           - A tag is used to isolate and group search meaning the search will not include tags that do not match.
	props         - A property map to act as match making conditions.
	limit         - A hard limit to the number of results.
	callback      - A callback with the search results or an error.
*/
func Search(profileIDList []string, tag string, props map[string]int, limit int, callback func(err error, results []interface{})) {
}
