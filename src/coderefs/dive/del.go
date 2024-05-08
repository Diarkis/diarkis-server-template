package dive

/*
Del deletes the key and value pair after retrieves the value of the given key.

	[IMPORTANT] The function uses mutex lock internally.
	[IMPORTANT] The function is asynchronous internally and accesses a remote node in the Diarkis cluster.
	[IMPORTANT] The function may fail during key migration.

The returned value is the value of the given key that has been deleted.

Error Cases

	┌────────────────────────────────┬─────────────────────────────────────────────────────────────────────────┐
	│ Error                          │ Reason                                                                  │
	╞════════════════════════════════╪═════════════════════════════════════════════════════════════════════════╡
	│ Setup must be invoked          │ In order to use Dive module,                                            │
	│                                │ dive.Setup() must be called before calling diarkis.Start()              │
	├────────────────────────────────┼─────────────────────────────────────────────────────────────────────────┤
	│ Key must not be empty          │ Input given key is an empty string.                                     │
	├────────────────────────────────┼─────────────────────────────────────────────────────────────────────────┤
	│ Storage node address not found │ There is no node that stores the given key.                             │
	├────────────────────────────────┼─────────────────────────────────────────────────────────────────────────┤
	│ Storage does not exist         │ The storage cannot be found on the designated node.                     │
	├────────────────────────────────┼─────────────────────────────────────────────────────────────────────────┤
	│ Key does not exists            │ The given key does not exist in the storage.                            │
	├────────────────────────────────┼─────────────────────────────────────────────────────────────────────────┤
	│ Network errors                 │ Internal network error such as reliable communication time out etc.     │
	└────────────────────────────────┴─────────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

Besides the following errors, there are Mesh module errors as well.

	dive.IsSetupError(err error)               // Dive has not been setup

	dive.IsInvalidKeyError(err error)          // Invalid or empty key provided

	dive.IsNodeNotFoundError(err error)        // Dive storage server node not found

	dive.IsNodeAddressNotFoundError(err error) // Dive storage server node address not found by the given key

	dive.IsDelError(err error)                 // Dive Del failure

▶︎ Return Data

The function returns an instance of *Result. To retrieve the data according to the correct data type, *Result has To...() functions.

Example:

	// If the intended result value is a string
	value, err := storage.Del(key).ToString()

	// If the intended result value is an array of float32
	values, err := storage.Del(anotherKey).ToFloat32Array()

Parameters

	key - The key of the value to retrieve and delete.
*/
func (s *Storage) Del(key string) *Result {
	return nil
}
