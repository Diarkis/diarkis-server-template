package dive

/*
SetOnLRangeByStorageName assigns a callback to be invoked on every LRange call by storage name.

	[IMPORTANT] This callback must be assigned on the target node server.
	            Default target node server is HTTP, but it can be configured to something else.

	[IMPORTANT] This function is not goroutine safe. It means that it should be used only during the setup of the server.

	[IMPORTANT] The callback receives an array of interface{} NOT *Result.
	            In order to handle each interface{} element, you must use util module's To...() function.

	            Example:

	            assigned := SetOnLRangeByStorageName(storageName, func(list []interface{}) []interface{} {
	              for i := 0; i < len(list); i++ {
	                value, ok := util.ToFloat32(list[i])
	              }
	            })

Parameters

	name - Storage name to assign the callback.
	cb   - Callback function to be executed on every LRange.
	       The callback will be passed the entire list data and the returned list data will be the returned value of LRange.
*/
func SetOnLRangeByStorageName(name string, cb func(list []interface{}) []interface{}) bool {
	return false
}

/*
SetOnLRange assigns a callback to be invoked on every LRange.

	[IMPORTANT] This callback must be assigned on the target node server.
	            Default target node server is HTTP, but it can be configured to something else.

	[IMPORTANT] This function is not goroutine safe. It means that it should be used only during the setup of the server.

	[IMPORTANT] The callback receives an array of interface{} NOT *Result.
	            In order to handle each interface{} element, you must use util module's To...() function.

	            Example:

	            assigned := storage.SetOnLRange(func(list []interface{}) []interface{} {
	              for i := 0; i < len(list); i++ {
	                value, ok := util.ToFloat32(list[i])
	              }
	            })

Parameters

	cb   - Callback function to be executed on every LRange.
	       The callback will be passed the entire list data and the returned list data will be the returned value of LRange.
*/
func (s *Storage) SetOnLRange(cb func(list []interface{}) []interface{}) bool {
	return false
}

/*
Get retrieves the value of the given key.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The function may fail during key migration.

Error Cases

	┌────────────────────────────────┬───────────────────────────────────────────────────────────────────────┐
	│ Error                          │ Reason                                                                │
	╞════════════════════════════════╪═══════════════════════════════════════════════════════════════════════╡
	│ Setup must be invoked          │ In order to use Dive module,                                          │
	│                                │ dive.Setup() must be called before calling diarkis.Start()            │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key must not be empty          │ Input given key is an empty string.                                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage node address not found │ There is no node that stores the given key.                           │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage does not exist         │ The storage cannot be found on the designated node.                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key has expired                │ The given key has expired, but not yet deleted from the memory.       │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key does not exists            │ The given key does not exist in the storage.                          │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Network errors                 │ Internal network error such as reliable communication time out etc.   │
	└────────────────────────────────┴───────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

Besides the following errors, there are Mesh module errors as well.

	dive.IsSetupError(err error)               // Dive has not been setup

	dive.IsInvalidKeyError(err error)          // Invalid or empty key provided

	dive.IsNodeNotFoundError(err error)        // Dive storage server node not found

	dive.IsNodeAddressNotFoundError(err error) // Dive storage server node address not found by the given key

	dive.IsGetError(err error)                 // Dive Get failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is a string
	value, err := storage.Get(key).ToString()

	// If the intended result value is an array of float32
	values, err := storage.Get(anotherKey).ToFloat32Array()

Parameters

	key - The key of the value to retrieve.
*/
func (s *Storage) Get(key string) *Result {
	return nil
}

/*
GetEx retrieves the value of the given key.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The function may fail during key migration.

Error Cases

	┌────────────────────────────────┬───────────────────────────────────────────────────────────────────────┐
	│ Error                          │ Reason                                                                │
	╞════════════════════════════════╪═══════════════════════════════════════════════════════════════════════╡
	│ Setup must be invoked          │ In order to use Dive module,                                          │
	│                                │ dive.Setup() must be called before calling diarkis.Start()            │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key must not be empty          │ Input given key is an empty string.                                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage node address not found │ There is no node that stores the given key.                           │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage does not exist         │ The storage cannot be found on the designated node.                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ TTL exceeds 86400 ms           │ TTL must not exceed 86400 ms (24 hours).                              │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key has expired                │ The given key has expired, but not yet deleted from the memory.       │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key does not exists            │ The given key does not exist in the storage.                          │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Network errors                 │ Internal network error such as reliable communication time out etc.   │
	└────────────────────────────────┴───────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

Besides the following errors, there are Mesh module errors as well.

	dive.IsSetupError(err error)               // Dive has not been setup

	dive.IsInvalidKeyError(err error)          // Invalid or empty key provided

	dive.IsNodeNotFoundError(err error)        // Dive storage server node not found

	dive.IsNodeAddressNotFoundError(err error) // Dive storage server node address not found by the given key

	dive.IsInvalidTTLError(err error)          // Provided TTL is invalid

	dive.IsGetError(err error)                 // Dive Get failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is a string
	value, err := storage.GetEx(key, ttl).ToString()

	// If the intended result value is an array of float32
	values, err := storage.GetEx(anotherKey, ttl).ToFloat32Array()

Parameters

	key - The key of the value to retrieve.
	ttl - If greater than 0 is given, TTL of the key will be extended.
	      TTL must not exceed 86400. TTL is in seconds.
*/
func (s *Storage) GetEx(key string, ttl uint32) *Result {
	return nil
}

/*
LRange retrieves the value of the given key as an array. LRange can only work with either LPush or RPush.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The function may fail during key migration.
	[IMPORTANT] If the total length of the value list is smaller than the given range, LRange will retrieve the entire value list.

Error Cases

	┌────────────────────────────────┬───────────────────────────────────────────────────────────────────────┐
	│ Error                          │ Reason                                                                │
	╞════════════════════════════════╪═══════════════════════════════════════════════════════════════════════╡
	│ Setup must be invoked          │ In order to use Dive module,                                          │
	│                                │ dive.Setup() must be called before calling diarkis.Start()            │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key must not be empty          │ Input given key is an empty string.                                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage node address not found │ There is no node that stores the given key.                           │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage does not exist         │ The storage cannot be found on the designated node.                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Invalid range                  │ Given range is invalid: "from" must be less than "to".                │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Invalid value                  │ The value of the given key is not a list.                             │
	│                                │ It means that the value is not stored by either LPush or RPush.       │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key has expired                │ The given key has expired, but not yet deleted from the memory.       │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key does not exists            │ The given key does not exist in the storage.                          │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Network errors                 │ Internal network error such as reliable communication time out etc.   │
	└────────────────────────────────┴───────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

Besides the following errors, there are Mesh module errors as well.

	dive.IsSetupError(err error)                  // Dive has not been setup

	dive.IsInvalidKeyError(err error)             // Invalid or empty key provided

	dive.IsNodeNotFoundError(err error)           // Dive storage server node not found

	dive.IsNodeAddressNotFoundError(err error)    // Dive storage server node address not found by the given key

	dive.IsInvalidRangeParametersError(err error) // Range from and to parameters are invalid

	dive.IsRangeError(err error)                  // Dive Range failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is an array of strings
	results, err := storage.LRange(key, 0, 10)

	if err != nil {
	  // handle error...
	}

	for i := 0; i < len(results); i++ {
	  value, err := results[i].ToString()
	}

Parameters

	key  - The key of the value to retrieve.
	from - The index offset to read from. Example: from=0 to=10 will retrieve list[0:10]
	to   - The index offset to read to.   Example: from=5 to=15 will retrieve list[5:15]
	       If you give 0 or below, LRange will retrieve the entire length of the list value.
*/
func (s *Storage) LRange(key string, from, to int) ([]*Result, error) {
	return nil, nil
}

/*
LRangeEx retrieves the value of the given key as an array. LRange can only work with either LPush or RPush.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The function may fail during key migration.
	[IMPORTANT] If the total length of the value list is smaller than the given range, LRange will retrieve the entire value list.

Error Cases

	┌────────────────────────────────┬───────────────────────────────────────────────────────────────────────┐
	│ Error                          │ Reason                                                                │
	╞════════════════════════════════╪═══════════════════════════════════════════════════════════════════════╡
	│ Setup must be invoked          │ In order to use Dive module,                                          │
	│                                │ dive.Setup() must be called before calling diarkis.Start()            │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key must not be empty          │ Input given key is an empty string.                                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage node address not found │ There is no node that stores the given key.                           │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage does not exist         │ The storage cannot be found on the designated node.                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Invalid range                  │ Given range is invalid: "from" must be less than "to".                │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Invalid value                  │ The value of the given key is not a list.                             │
	│                                │ It means that the value is not stored by either LPush or RPUsh.       │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key has expired                │ The given key has expired, but not yet deleted from the memory.       │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key does not exists            │ The given key does not exist in the storage.                          │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ TTL exceeds 86400 ms           │ TTL must not exceed 86400 ms (24 hours).                              │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Network errors                 │ Internal network error such as reliable communication time out etc.   │
	└────────────────────────────────┴───────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

Besides the following errors, there are Mesh module errors as well.

	dive.IsSetupError(err error)                  // Dive has not been setup

	dive.IsInvalidKeyError(err error)             // Invalid or empty key provided

	dive.IsNodeNotFoundError(err error)           // Dive storage server node not found

	dive.IsNodeAddressNotFoundError(err error)    // Dive storage server node address not found by the given key

	dive.IsInvalidRangeParametersError(err error) // Range from and to parameters are invalid

	dive.IsInvalidTTLError(err error)             // Provided TTL is invalid

	dive.IsRangeError(err error)                  // Dive Range failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is an array of strings
	results, err := storage.LRangeEx(key, 0, 10, ttl)

	if err != nil {
	  // handle error...
	}

	for i := 0; i < len(results); i++ {
	  value, err := results[i].ToString()
	}

Parameters

	key  - The key of the value to retrieve.
	from - The index offset to read from. Example: from=0 to=10 will retrieve list[0:10]
	to   - The index offset to read to.   Example: from=5 to=15 will retrieve list[5:15]
	       If you give 0 or below, LRange will retrieve the entire length of the list value.
	ttl  - If greater than 0 is given, TTL of the key will be extended.
	       TTL must not exceed 86400. TTL is in seconds.
*/
func (s *Storage) LRangeEx(key string, from, to int, ttl uint32) ([]*Result, error) {
	return nil, nil
}
