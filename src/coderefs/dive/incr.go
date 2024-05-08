package dive

/*
IncrBy increments the stored value by the given increment amount.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value to increment must be of the data type int32
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The key and value pair WILL expire in 24 hours.
	            Default TTL can be configured to be other than 24 hours.

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
	│ Key does not exist             │ The key and its value must exist.                                     │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Invalid data type to increment │ The value must be int32.                                              │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Network errors                 │ Internal network error such as reliable communication time out etc.   │
	└────────────────────────────────┴───────────────────────────────────────────────────────────────────────┘

Parameters

	key       - The key of the given value to be stored.
	increment - The increment amount of the value stored.
*/
func (s *Storage) IncrBy(key string, increment int32) *Result {
	return nil
}

/*
IncrByEx increments the stored value by the given increment amount.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value to increment must be of the data type int32
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The key and value pair WILL expire in 24 hours.
	            Default TTL can be configured to be other than 24 hours.

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
	│ Key does not exist             │ The key and its value must exist.                                     │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Invalid data type to increment │ The value must be int32.                                              │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Network errors                 │ Internal network error such as reliable communication time out etc.   │
	└────────────────────────────────┴───────────────────────────────────────────────────────────────────────┘

Parameters

	key       - The key of the given value to be stored.
	increment - The increment amount of the value stored.
	ttl       - If greater than 0 is given, TTL of the key will be extended.
	            TTL must not exceed 86400. TTL is in seconds.
*/
func (s *Storage) IncrByEx(key string, increment int32, ttl int64) *Result {
	return nil
}
