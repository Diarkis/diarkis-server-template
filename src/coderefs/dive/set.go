package dive

/*
Set stores the key value pair on a node derived from the given key.

	[CRITICALLY IMPORTANT] The value must NOT be struct or contain struct.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value data type must be either primitive, string, or []byte.
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
	│ Invalid data type              │ The value to increment must be of type int32.                         │
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

	dive.IsInvalidValueTypeError(err error)    // Provided value type is invalid

	dive.IsSetError(err error)                 // Dive Set failure

Parameters

	key   - The key of the given value to be stored.
	value - Value of the given key to be stored.
*/
func (s *Storage) Set(key string, value interface{}) error {
	return nil
}

/*
SetEx stores the key value pair on a node derived from the given key and assigns a TTL in seconds.

	[CRITICALLY IMPORTANT] The value must NOT be struct or contain struct.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value data type must be either primitive, string, or []byte.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.

Error Cases

	┌────────────────────────────────┬───────────────────────────────────────────────────────────────────────┐
	│ Error                          │ Reason                                                                │
	╞════════════════════════════════╪═══════════════════════════════════════════════════════════════════════╡
	│ Setup must be invoked          │ In order to use Dive module,                                          │
	│                                │ dive.Setup() must be called before calling diarkis.Start()            │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key must not be empty          │ Input given key is an empty string.                                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ TTL must not be 0              │ TTL must be greater than 0 second.                                    │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage node address not found │ There is no node that stores the given key.                           │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage does not exist         │ The storage cannot be found on the designated node.                   │
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

	dive.IsInvalidValueTypeError(err error)    // Provided value type is invalid

	dive.IsInvalidTTLError(err error)          // Provided TTL is invalid

	dive.IsSetError(err error)                 // Dive Set failure

Parameters

	key   - The key of the given value to be stored.
	value - Value of the given key to be stored.
	ttl   - TTL in seconds for the key and value pair to expire.
*/
func (s *Storage) SetEx(key string, value interface{}, ttl uint32) error {
	return nil
}

/*
SetIfNotExists stores the key value pair on a node derived from the given key.

	[CRITICALLY IMPORTANT] The value must NOT be struct or contain struct.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value data type must be either primitive, string, or []byte.
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
	│ Key already exists             │ The given key already exists.                                         │
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

	dive.IsInvalidValueTypeError(err error)    // Provided value type is invalid

	dive.IsSetIfNotExistError(err error)       // Dive SetIfNotExist failure

Parameters

	key   - The key of the given value to be stored.
	value - Value of the given key to be stored.
*/
func (s *Storage) SetIfNotExists(key string, value interface{}) error {
	return nil
}

/*
SetIfNotExistsEx stores the key value pair on a node derived from the given key and assigns a TTL for the key and value to expire.

	[CRITICALLY IMPORTANT] The value must NOT be struct or contain struct.

	[IMPORTANT]            The function uses mutex lock internally.
	[IMPORTANT]            Value data type must be either primitive, string, or []byte.
	[IMPORTANT]            The function is asynchronous internally and accesses a remote node in the Diarkis cluster.

Error Cases

	┌────────────────────────────────┬───────────────────────────────────────────────────────────────────────┐
	│ Error                          │ Reason                                                                │
	╞════════════════════════════════╪═══════════════════════════════════════════════════════════════════════╡
	│ Setup must be invoked          │ In order to use Dive module,                                          │
	│                                │ dive.Setup() must be called before calling diarkis.Start()            │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key must not be empty          │ Input given key is an empty string.                                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ TTL must not be 0              │ TTL must be greater than 0 second.                                    │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage node address not found │ There is no node that stores the given key.                           │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Storage does not exist         │ The storage cannot be found on the designated node.                   │
	├────────────────────────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Key already exists             │ The given key already exists.                                         │
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

	dive.IsInvalidValueTypeError(err error)    // Provided value type is invalid

	dive.IsInvalidTTLError(err error)          // Provided TTL is invalid

	dive.IsSetIfNotExistError(err error)       // Dive SetIfNotExist failure

Parameters

	key   - The key of the given value to be stored.
	value - Value of the given key to be stored.
	ttl   - TTL in seconds for the key and value pair to expire.
*/
func (s *Storage) SetIfNotExistsEx(key string, value interface{}, ttl uint32) error {
	return nil
}
