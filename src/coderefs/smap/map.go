package smap

/*
SyncMap does NOT allow another SyncMap within
*/
type SyncMap []*shard

/*
KeysIfOperation filters keys to be returned by KeysIf.
*/
type KeysIfOperation func(key string) bool

/*
KeysIfValueOperation filters keys to be returned by KeysIfValue.
*/
type KeysIfValueOperation func(key string, value interface{}) bool

/*
UpsertOperation is a callback to be invoked when executing Upsert.
The callback is called while the lock is held
and because of that you MUST NOT try to access the internal map with the same key.
Returned value will be stored as an updated value.
*/
type UpsertOperation func(exists bool, storedValue interface{}, updateValue interface{}) (updatedValue interface{})

/*
RemoveOperation is a callback to be invoked when executing RemoveIf.
The callback is called while the lock is held
and because of that you MUST NOT try to access the internal map with the same key.
If the callback returns true, the key and its value will be removed from the map.
*/
type RemoveOperation func(exists bool, storedValue interface{}) bool

/*
New creates and returns a new SyncMap instance
*/
func New() SyncMap {
	return nil
}

/*
Encode encodes the stored keys and values as a snapshot.
Use Decode to convert it back to *SyncMap
*/
func Encode(sm SyncMap) ([]byte, error) {
	return nil, nil
}

/*
Decode converts the encoded *SyncMap back to *SyncMap.
If invalid data is given, the function may return an error.
*/
func Decode(encoded []byte) (SyncMap, error) {
	return nil, nil
}

/*
AllowAllTypes allows the value type of be anything including *SyncMap data type to be stored.

	[IMPORTANT] Encode and Decode do NOT support *SyncMap and other structs.
*/
func (sm SyncMap) AllowAllTypes(key string) {
}

/*
AllowNonPrimitive allows the given key's value to be non-primitive.
This operation is per internal shared map.
To retrieve the stored value that is not primitive, use Get.
*/
func (sm SyncMap) AllowNonPrimitive(key string) {
}

/*
Set sets the given value with the key.

	[IMPORTANT] Allowed value is either a primitive, string, or byte array.

Returns an error if the given value data type is not primitive or byte array.
*/
func (sm SyncMap) Set(key string, value interface{}) error {
	return nil
}

/*
MSet sets multiple keys and values.

	[IMPORTANT] Allowed value is either a primitive, string, or byte array.

Returns an error if the given data map contains an invalid data type.

When it returns an error no key and value will be stored at all.
*/
func (sm SyncMap) MSet(data map[string]interface{}, allowNonPrimitive bool) error {
	return nil
}

/*
SetIf sets the given key along with the value if the given check function returns true
and the key does not exists.

	[IMPORTANT] check function is executed while SetIf holds a mutex lock.
	            Do not use mutex lock in check function to avoid deadlock.

If either the check function returns false or the key already exists, it returns an error.
*/
func (sm SyncMap) SetIf(key string, value interface{}, check func() bool) error {
	return nil
}

/*
SetIfNotExists sets the given key along with the value if the key does not exist.

	[IMPORTANT] Allowed value is either a primitive, string, or byte array.

If the key already exists, it returns an error.
*/
func (sm SyncMap) SetIfNotExists(key string, value interface{}) error {
	return nil
}

/*
Upsert will updates the existing value using the return value of UpsertOperation callback if the value of the key exists.

If the key and value does not exist, it will simply add the new key along with the value.

Returns the updated or inserted value as an interface{}.

Use To...() function such as ToInt, ToString etc. To convert it to appropriate data type.

	[IMPORTANT] Allowed value is either a primitive, string, or byte array.

	[IMPORTANT] Callback is invoked while the lock is held.
	            You must NOT use functions that uses lock inside the callback.
*/
func (sm SyncMap) Upsert(key string, value interface{}, cb UpsertOperation) (interface{}, error) {
	return nil, nil
}

/*
IncrAsUint8 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsUint8(key string, incr uint8) (uint8, error) {
	return 0, nil
}

/*
IncrAsUint16 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsUint16(key string, incr uint16) (uint16, error) {
	return 0, nil
}

/*
IncrAsUint32 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsUint32(key string, incr uint32) (uint32, error) {
	return 0, nil
}

/*
IncrAsUint64 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsUint64(key string, incr uint64) (uint64, error) {
	return 0, nil
}

/*
IncrAsUint increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsUint(key string, incr uint) (uint, error) {
	return 0, nil
}

/*
IncrAsInt8 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsInt8(key string, incr int8) (int8, error) {
	return 0, nil
}

/*
IncrAsInt16 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsInt16(key string, incr int16) (int16, error) {
	return 0, nil
}

/*
IncrAsInt32 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsInt32(key string, incr int32) (int32, error) {
	return 0, nil
}

/*
IncrAsInt64 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsInt64(key string, incr int64) (int64, error) {
	return 0, nil
}

/*
IncrAsInt increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsInt(key string, incr int) (int, error) {
	return 0, nil
}

/*
IncrAsFloat32 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsFloat32(key string, incr float32) (float32, error) {
	return 0, nil
}

/*
IncrAsFloat64 increments the value of the given key by the given incr value.
*/
func (sm SyncMap) IncrAsFloat64(key string, incr float64) (float64, error) {
	return 0, nil
}

/*
Get returns the value of the given key as an interface{}.
*/
func (sm SyncMap) Get(key string) (interface{}, bool) {
	return nil, false
}

/*
GetAsUint8 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint8(key string) (uint8, bool) {
	return 0, false
}

/*
ToUint8 converts interface{}.
*/
func (sm SyncMap) ToUint8(v interface{}) (uint8, bool) {
	return 0, false
}

/*
GetAsUint16 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint16(key string) (uint16, bool) {
	return 0, false
}

/*
ToUint16 converts interface{}.
*/
func (sm SyncMap) ToUint16(v interface{}) (uint16, bool) {
	return 0, false
}

/*
GetAsUint32 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint32(key string) (uint32, bool) {
	return 0, false
}

/*
ToUint32 converts interface{}
*/
func (sm SyncMap) ToUint32(v interface{}) (uint32, bool) {
	return 0, false
}

/*
GetAsUint64 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint64(key string) (uint64, bool) {
	return 0, false
}

/*
ToUint64 converts interface{}
*/
func (sm SyncMap) ToUint64(v interface{}) (uint64, bool) {
	return 0, false
}

/*
GetAsInt8 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt8(key string) (int8, bool) {
	return 0, false
}

/*
ToInt8 converts interface{}.
*/
func (sm SyncMap) ToInt8(v interface{}) (int8, bool) {
	return 0, false
}

/*
GetAsInt16 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt16(key string) (int16, bool) {
	return 0, false
}

/*
ToInt16 converts interface{}.
*/
func (sm SyncMap) ToInt16(v interface{}) (int16, bool) {
	return 0, false
}

/*
GetAsInt32 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt32(key string) (int32, bool) {
	return 0, false
}

/*
ToInt32 converts interface{}
*/
func (sm SyncMap) ToInt32(v interface{}) (int32, bool) {
	return 0, false
}

/*
GetAsInt64 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt64(key string) (int64, bool) {
	return 0, false
}

/*
ToInt64 converts interface{}
*/
func (sm SyncMap) ToInt64(v interface{}) (int64, bool) {
	return 0, false
}

/*
GetAsUint returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint(key string) (uint, bool) {
	return 0, false
}

/*
ToUint converts interface{}
*/
func (sm SyncMap) ToUint(v interface{}) (uint, bool) {
	return 0, false
}

/*
GetAsInt returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt(key string) (int, bool) {
	return 0, false
}

/*
ToInt converts interface{}
*/
func (sm SyncMap) ToInt(v interface{}) (int, bool) {
	return 0, false
}

/*
GetAsFloat32 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsFloat32(key string) (float32, bool) {
	return 0, false
}

/*
ToFloat32 converts interface{}.
*/
func (sm SyncMap) ToFloat32(v interface{}) (float32, bool) {
	return 0, false
}

/*
GetAsFloat64 returns the stored value of the given key.
*/
func (sm SyncMap) GetAsFloat64(key string) (float64, bool) {
	return 0, false
}

/*
ToFloat64 converts interface{}
*/
func (sm SyncMap) ToFloat64(v interface{}) (float64, bool) {
	return 0, false
}

/*
GetAsBool returns the stored value of the given key.
*/
func (sm SyncMap) GetAsBool(key string) (bool, bool) {
	return false, false
}

/*
ToBool converts interface{}
*/
func (sm SyncMap) ToBool(v interface{}) (bool, bool) {
	return false, false
}

/*
GetAsString returns the stored value of the given key.
*/
func (sm SyncMap) GetAsString(key string) (string, bool) {
	return "", false
}

/*
ToString converts interface{}
*/
func (sm SyncMap) ToString(v interface{}) (string, bool) {
	return "", false
}

/*
GetAsBytes returns the stored value of the given key.
*/
func (sm SyncMap) GetAsBytes(key string) ([]byte, bool) {
	return nil, false
}

/*
ToBytes converts interface{}
*/
func (sm SyncMap) ToBytes(v interface{}) ([]byte, bool) {
	return nil, false
}

/*
GetAsUint8Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint8Array(key string) ([]uint8, bool) {
	return nil, false
}

/*
ToUint8Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToUint8Array(v interface{}) ([]uint8, bool) {
	return nil, false
}

/*
GetAsUint16Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint16Array(key string) ([]uint16, bool) {
	return nil, false
}

/*
ToUint16Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToUint16Array(v interface{}) ([]uint16, bool) {
	return nil, false
}

/*
GetAsUint32Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint32Array(key string) ([]uint32, bool) {
	return nil, false
}

/*
ToUint32Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToUint32Array(v interface{}) ([]uint32, bool) {
	return nil, false
}

/*
GetAsUint64Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUint64Array(key string) ([]uint64, bool) {
	return nil, false
}

/*
ToUint64Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToUint64Array(v interface{}) ([]uint64, bool) {
	return nil, false
}

/*
GetAsUintArray returns the stored value of the given key.
*/
func (sm SyncMap) GetAsUintArray(key string) ([]uint, bool) {
	return nil, false
}

/*
ToUintArray converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToUintArray(v interface{}) ([]uint, bool) {
	return nil, false
}

/*
GetAsInt8Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt8Array(key string) ([]int8, bool) {
	return nil, false
}

/*
ToInt8Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToInt8Array(v interface{}) ([]int8, bool) {
	return nil, false
}

/*
GetAsInt16Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt16Array(key string) ([]int16, bool) {
	return nil, false
}

/*
ToInt16Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToInt16Array(v interface{}) ([]int16, bool) {
	return nil, false
}

/*
GetAsInt32Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt32Array(key string) ([]int32, bool) {
	return nil, false
}

/*
ToInt32Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToInt32Array(v interface{}) ([]int32, bool) {
	return nil, false
}

/*
GetAsInt64Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsInt64Array(key string) ([]int64, bool) {
	return nil, false
}

/*
ToInt64Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToInt64Array(v interface{}) ([]int64, bool) {
	return nil, false
}

/*
GetAsIntArray returns the stored value of the given key.
*/
func (sm SyncMap) GetAsIntArray(key string) ([]int, bool) {
	return nil, false
}

/*
ToIntArray converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToIntArray(v interface{}) ([]int, bool) {
	return nil, false
}

/*
GetAsFloat32Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsFloat32Array(key string) ([]float32, bool) {
	return nil, false
}

/*
ToFloat32Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToFloat32Array(v interface{}) ([]float32, bool) {
	return nil, false
}

/*
GetAsFloat64Array returns the stored value of the given key.
*/
func (sm SyncMap) GetAsFloat64Array(key string) ([]float64, bool) {
	return nil, false
}

/*
ToFloat64Array converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToFloat64Array(v interface{}) ([]float64, bool) {
	return nil, false
}

/*
GetAsStringArray returns the stored value of the given key.
*/
func (sm SyncMap) GetAsStringArray(key string) ([]string, bool) {
	return nil, false
}

/*
ToStringArray converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToStringArray(v interface{}) ([]string, bool) {
	return nil, false
}

/*
GetAsBoolArray returns the stored value of the given key.
*/
func (sm SyncMap) GetAsBoolArray(key string) ([]bool, bool) {
	return nil, false
}

/*
ToBoolArray converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToBoolArray(v interface{}) ([]bool, bool) {
	return nil, false
}

/*
GetAsBytesArray returns the stored value of the given key.
*/
func (sm SyncMap) GetAsBytesArray(key string) ([][]byte, bool) {
	return nil, false
}

/*
ToBytesArray converts interface{}. If other data type found in the array,
it returns an empty array with false.
*/
func (sm SyncMap) ToBytesArray(v interface{}) ([][]byte, bool) {
	return nil, false
}

/*
Count returns the number of elements stored.
*/
func (sm SyncMap) Count() int {
	return 0
}

/*
Range iterates over all keys and invokes the callback on every key.
The lock is held while the callback is invoked
and because of that you MUST NOT try to access the internal map with the same key.
*/
func (sm SyncMap) Range(cb func(key string, value interface{})) {
}

/*
Keys returns all keys as an array of string.
*/
func (sm SyncMap) Keys() []string {
	return nil
}

/*
GetKeysByRange returns keys as an array of the length specified

	[IMPORTANT] If the total number of keys is smaller than the given length,
	            the returned array will be smaller than the given length.
*/
func (sm SyncMap) GetKeysByRange(howmany int) []string {
	return nil
}

/*
KeysIf returns all keys that the callback KeysIfOperation returns true.

	[IMPORTANT] Callback is invoked while the lock is held and because of that
	            you must NOT use functions that uses lock inside the callback.
*/
func (sm SyncMap) KeysIf(cb KeysIfOperation) []string {
	return nil
}

/*
KeysIfValue returns all keys that the callback KeysIfOperation returns true.

	[IMPORTANT] Callback is invoked while the lock is held and because of that
	            you must NOT use functions that uses lock inside the callback.
*/
func (sm SyncMap) KeysIfValue(cb KeysIfValueOperation) []string {
	return nil
}

/*
Exists returns true if the given key exists.
*/
func (sm SyncMap) Exists(key string) bool {
	return false
}

/*
Remove removes an element of the given key.
Returns true if the element has been removed.
*/
func (sm SyncMap) Remove(key string) bool {
	return false
}

/*
RemoveIf removes an element of the given key if the callback returns true.

The return value of RemoveIf indicates if the key is deleted or not.

	[IMPORTANT] Callback is invoked while the lock is held and because of that
	            you must NOT use functions that uses lock inside the callback.
*/
func (sm SyncMap) RemoveIf(key string, cb RemoveOperation) bool {
	return false
}

/*
Clear removes all keys and values.
*/
func (sm SyncMap) Clear() {
}
