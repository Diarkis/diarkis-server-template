package dive

/*
SetOnPushByStorageName assigns a callback to be invoked on every LPush and RPush for the specified storage by its name.

	[IMPORTANT] This callback must be assigned on the target node server.
	            Default target node server is HTTP, but it can be configured to something else.

	[IMPORTANT] This function is not goroutine safe. It means that it should be used only during the setup of the server.

	[IMPORTANT] The callback is invoked while the mutex lock is held. You must avoid using mutex lock in the callback to avoid mutex deadlock.

	[IMPORTANT] The callback receives an array of interface{} NOT *Result.
	            In order to handle each interface{} element, you must use util module's To...() function.

	            Example:

	            assigned := SetOnPushByStorageName(storageName, func(list []interface{}) []interface{} {
	              for i := 0; i < len(list); i++ {
	                value, ok := util.ToFloat32(list[i])
	              }
	            })

Parameters

	name - Storage name to assign the callback.
	cb   - Callback function to be executed on every LPush and RPush.
	       The callback will be passed the entire list data and the returned list data will be stored on the storage.
*/
func SetOnPushByStorageName(name string, cb func(key string, list []interface{}) []interface{}) bool {
	return false
}

/*
SetOnPush assigns a callback to be invoked on every LPush and RPush.

	[IMPORTANT] This callback must be assigned on the target node server.
	            Default target node server is HTTP, but it can be configured to something else.

	[IMPORTANT] This function is not goroutine safe. It means that it should be used only during the setup of the server.

	[IMPORTANT] The callback is invoked while the mutex lock is held. You must avoid using mutex lock in the callback to avoid mutex deadlock.

	[IMPORTANT] The callback receives an array of interface{} NOT *Result.
	            In order to handle each interface{} element, you must use util module's To...() function.

	            Example:

	            assigned := storage.SetOnPush(storageName, func(list []interface{}) []interface{} {
	              for i := 0; i < len(list); i++ {
	                value, ok := util.ToFloat32(list[i])
	              }
	            })

Parameters

	cb   - Callback function to be executed on every LPush and RPush.
	       The callback will be passed the entire list data and the returned list data will be stored on the storage.
*/
func (s *Storage) SetOnPush(cb func(key string, list []interface{}) []interface{}) bool {
	return false
}

/*
LPush stores a value as an element of a list associated to the given key on a node derived from the given key
and assigns a TTL for the key and value to expire.

The value as an element is pushed to the front of the array and it updates internal timestamp and TTL.

	[CRITICALLY IMPORTANT] The value must NOT be struct or contain struct.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value data type must be either primitive, string, or []byte.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The key WILL expire in 24 hours.
	            Default TTL can be configured to be other than 24 hours.
	[IMPORTANT] Every LPush updates TTL of the key.

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

	dive.IsPushError(err error)                // Dive Push failure

Parameters

	key   - The key of the given value to be stored.
	value - Value of the given key to be stored.
	        You may pass multiple values to push multiple elements to the list set.
	         Examples:
	                  LPush(key, 1, 2, 3)
	                  LPush(key, []interface{}{1, 2, 3}...)
*/
func (s *Storage) LPush(key string, values ...interface{}) error {
	return nil
}

/*
LPushEx stores a value as an element of a list associated to the given key on a node derived from the given key
and assigns a TTL for the key and value to expire.

The value as an element is pushed to the front of the array and it updates internal timestamp and TTL.

	[CRITICALLY IMPORTANT] The value must NOT be struct or contain struct.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value data type must be either primitive, string, or []byte.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] Every LPushEX updates TTL of the key.

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
	│ TTL exceeds 86400 ms           │ TTL must not exceed 86400 ms (24 hours).                              │
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

	dive.IsPushError(err error)                // Dive Push failure

Parameters

	key   - The key of the given value to be stored.
	ttl   - TTL in seconds for the key and value pair to expire.
	        If greater than 0 is given, TTL of the key will be extended.
	        TTL must not exceed 86400. TTL is in seconds.
	value - Value of the given key to be stored.
	        You may pass multiple values to push multiple elements to the list set.
	         Examples:
	                  LPushEx(key, 1, 2, 3)
	                  LPushEx(key, []interface{}{1, 2, 3}...)
*/
func (s *Storage) LPushEx(key string, ttl uint32, values ...interface{}) error {
	return nil
}

/*
RPush stores a value as an element of a list associated to the given key on a node derived
from the given key and assigns a TTL for the key and value to expire.

The value as an element is pushed to the end of the array and it updates internal timestamp and TTL.

	[CRITICALLY IMPORTANT] The value must NOT be struct or contain struct.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value data type must be either primitive, string, or []byte.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The key WILL expire in 24 hours.
	            Default TTL can be configured to be other than 24 hours.
	[IMPORTANT] Every RPush updates TTL of the key.

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

	dive.IsPushError(err error)                // Dive Push failure

Parameters

	key   - The key of the given value to be stored.
	value - Value of the given key to be stored.
	        You may pass multiple values to push multiple elements to the list set.
	         Examples:
	                  RPush(key, 1, 2, 3)
	                  RPush(key, []interface{}{1, 2, 3}...)
*/
func (s *Storage) RPush(key string, values ...interface{}) error {
	return nil
}

/*
RPushEx stores a value as an element of a list associated to the given key on a node derived from the given key
and assigns a TTL for the key and value to expire.

The value as an element is pushed to the end of the array and it updates internal timestamp and TTL.

	[CRITICALLY IMPORTANT] The value must NOT be struct or contain struct.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] Value data type must be either primitive, string, or []byte.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] Every RPushEX updates TTL of the key.

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
	│ TTL exceeds 86400 ms           │ TTL must not exceed 86400 ms (24 hours).                              │
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

	dive.IsPushError(err error)                // Dive Push failure

Parameters

	key   - The key of the given value to be stored.
	ttl   - TTL in seconds for the key and value pair to expire.
	        If greater than 0 is given, TTL of the key will be extended.
	        TTL must not exceed 86400. TTL is in seconds.
	value - Value of the given key to be stored.
	        You may pass multiple values to push multiple elements to the list set.
	         Examples:
	                  RPushEx(key, 1, 2, 3)
	                  RPushEx(key, []interface{}{1, 2, 3}...)
*/
func (s *Storage) RPushEx(key string, ttl uint32, values ...interface{}) error {
	return nil
}
