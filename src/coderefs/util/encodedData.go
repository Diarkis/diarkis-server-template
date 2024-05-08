package util

/*
StringArray represents an array of strings that is goroutine safe without the use of mutex lock
*/
type StringArray struct{}

/*
StringMap represents a map with string key and string value
that is goroutine safe without the use of mutex lock.
map[string]string
*/
type StringMap struct{}

/*
NewStringArray creates a new StringArray instance.
*/
func NewStringArray() *StringArray {
	return nil
}

/*
NewStringArrayFromExportedData creates a new StringArray instance from exported StringArray raw data.
*/
func NewStringArrayFromExportedData(data string) *StringArray {
	return nil
}

/*
NewStringMap creates a new StringMap instance.
*/
func NewStringMap() *StringMap {
	return nil
}

/*
NewStringMapFromExportedData creates a new StringMap instance from exported StringMap raw data.
*/
func NewStringMapFromExportedData(data string) *StringMap {
	return nil
}

/*
Length returns how may keys the map has.
*/
func (sa *StringArray) Length() int {
	return 0
}

/*
GetAt returns the element from the array at the given index.
Returns an empty string if the given index is invalid.
*/
func (sa *StringArray) GetAt(index int) string {
	return ""
}

/*
GetIndex returns the index of the given string value in the array.
Returns -1 if the given string value does not exist in the array.
*/
func (sa *StringArray) GetIndex(value string) int {
	return 0
}

/*
Push pushes a given string value at the end of the array.
*/
func (sa *StringArray) Push(value string) {
}

/*
Pop returns the first element of the array and removes it from the array.
Returns an empty string if there is nothing to pop from the array.
*/
func (sa *StringArray) Pop() string {
	return ""
}

/*
Delete removes the key and its value from the map.
Returns false if the key does not exist in the map.
*/
func (sa *StringArray) Delete(value string) bool {
	return false
}

/*
Clear resets and clears all keys and their values.
*/
func (sa *StringArray) Clear() {
}

/*
Export returns the raw data as a string to be used by NewStringMapFromExportedData.
*/
func (sa *StringArray) Export() string {
	return ""
}

/*
Get returns the value associated with the given key.
Returns an empty string if the key does not exist in the map
*/
func (sm *StringMap) Get(key string) string {
	return ""
}
func (sm *StringMap) Length() int {
	return 0
}

/*
Set assigns the given value to the given key.
*/
func (sm *StringMap) Set(key string, value string) {
}
func (sm *StringMap) Delete(key string) bool {
	return false
}
func (sm *StringMap) Clear() {
}
func (sm *StringMap) Export() string {
	return ""
}
