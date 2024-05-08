package util

/*
ToStringKeyIntMap returns a map with string keys with int values.
The second returned value indicates if the data is valid or not.
*/
func ToStringKeyIntMap(d interface{}) (map[string]int, bool) {
	return nil, false
}

/*
ToStringKeyBytesMap converts an interface to a map of byte array with string keys.
*/
func ToStringKeyBytesMap(d interface{}) (map[string][]byte, bool) {
	return nil, false
}

/*
ToStringKeyInt64Map converts an intrerface to a map of int64 with string keys.
*/
func ToStringKeyInt64Map(d interface{}) (map[string]int64, bool) {
	return nil, false
}

/*
GetAsStringKeyInt64Map returns a map with string keys and int64 values.
from a given map[string]interface{}.
The second returned value indicates if the data is valid or not.
*/
func GetAsStringKeyInt64Map(data map[string]interface{}, key string) (map[string]int64, bool) {
	return nil, false
}

/*
GetAsStringKeyBytesMap returns a map with string keys and byte array values.
from a given map[string]interface{}.
The second returned value indicates if the data is valid or not.
*/
func GetAsStringKeyBytesMap(data map[string]interface{}, key string) (map[string][]byte, bool) {
	return nil, false
}

/*
GetAsStringKeyIntMap returns a map with string keys and int values
from a given map[string]interface{}.
The second returned value indicates if the data is valid or not.
*/
func GetAsStringKeyIntMap(data map[string]interface{}, key string) (map[string]int, bool) {
	return nil, false
}

/*
ToStringKeyInterfaceMap returns a map with string keys with interface values.
The second returned value indicates if the data is valid or not.
*/
func ToStringKeyInterfaceMap(d interface{}) (map[string]interface{}, bool) {
	return nil, false
}

/*
GetAsStringKeyInterfaceMap returns a map with string keys and interface values
from a given map[string]interface{}.
The second returned value indicates if the data is valid or not.
*/
func GetAsStringKeyInterfaceMap(data map[string]interface{}, key string) (map[string]interface{}, bool) {
	return nil, false
}
