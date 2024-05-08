package dive

/*
LPop returns an array of elements from the list set created by LPush and/or RPush from index 0 up to the number given by howmany parameter.

The returned elements will be atomically removed from the list set.

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

	dive.IsPopError(err error)                 // Dive Pop failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is an array of strings
	results, err := storage.LPop(key, 10)

	if err != nil {
	  // handle error...
	}

	for i := 0; i < len(results); i++ {
	  value, err := results[i].ToString()
	}

Parameters

	key     - The key of the value to retrieve.
	howmany - Number of elements from the list set to pop.
	          If the number is greater than the length of the list set, it will pop all elements.
	          If 0 or below given, the entire list will be popped.
*/
func (s *Storage) LPop(key string, howmany int) ([]*Result, error) {
	return nil, nil
}

/*
LPopEx returns an array of elements from the list set created by LPush and/or RPush from index 0 up to the number given by howmany parameter.

The returned elements will be atomically removed from the list set.

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
	│ TTL exceeds 86400 s            │ TTL must not exceed 86400 s (24 hours).                               │
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

	dive.IsPopError(err error)                 // Dive Pop failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is an array of strings
	results, err := storage.LPopEx(key, 10, ttl)

	if err != nil {
	  // handle error...
	}

	for i := 0; i < len(results); i++ {
	  value, err := results[i].ToString()
	}

Parameters

	key     - The key of the value to retrieve.
	howmany - Number of elements from the list set to pop.
	          If the number is greater than the length of the list set, it will pop all elements.
	          If 0 or below given, the entire list will be popped.
	ttl     - If greater than 0, TTL of the key will be extended.
	          TTL must not exceed 86400. TTL is in seconds.
*/
func (s *Storage) LPopEx(key string, howmany int, ttl uint32) ([]*Result, error) {
	return nil, nil
}

/*
RPop returns an array of elements from the list set created by LPush and/or RPush
from the last element up to the number given by howmany in reverse order.

The returned elements will be atomically removed from the list set.

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

	dive.IsPopError(err error)                 // Dive Pop failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is an array of strings
	results, err := storage.RPop(key, 10)

	if err != nil {
	  // handle error...
	}

	for i := 0; i < len(results); i++ {
	  value, err := results[i].ToString()
	}

Parameters

	key     - The key of the value to retrieve.
	howmany - Number of elements from the list set to pop.
	          If the number is greater than the length of the list set, it will pop all elements.
	          If 0 or below given, the entire list will be popped.
*/
func (s *Storage) RPop(key string, howmany int) ([]*Result, error) {
	return nil, nil
}

/*
RPopEx returns an array of elements from the list set created by LPush and/or RPush
from the last element up to the number given by howmany in reverse order.

The returned elements will be atomically removed from the list set.

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
	│ TTL exceeds 86400 s            │ TTL must not exceed 86400 s (24 hours).                               │
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

	dive.IsPopError(err error)                 // Dive Pop failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is an array of strings
	results, err := storage.RPopEx(key, 10, ttl)

	if err != nil {
	  // handle error...
	}

	for i := 0; i < len(results); i++ {
	  value, err := results[i].ToString()
	}

Parameters

	key     - The key of the value to retrieve.
	howmany - Number of elements from the list set to pop.
	          If the number is greater than the length of the list set, it will pop all elements.
	          If 0 or below given, the entire list will be popped.
	ttl     - If greater than 0, TTL of the key will be extended.
	          TTL must not exceed 86400. TTL is in seconds.
*/
func (s *Storage) RPopEx(key string, howmany int, ttl uint32) ([]*Result, error) {
	return nil, nil
}
