package util

/*
GetAsBytesArray returns a value of the given map by its key as an array of bytes.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsBytesArray(data map[string]interface{}, k string) ([][]byte, bool) {
	return nil, false
}

/*
GetAsBoolArray returns a value of the given map by its key as an array of bool.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsBoolArray(data map[string]interface{}, k string) ([]bool, bool) {
	return nil, false
}

/*
GetAsStringArray returns a value of the given map by its key as an array of string.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsStringArray(data map[string]interface{}, k string) ([]string, bool) {
	return nil, false
}

/*
GetAsIntArray returns a value of the given map by its key as an array of int.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsIntArray(data map[string]interface{}, k string) ([]int, bool) {
	return nil, false
}

/*
GetAsUintArray returns a value of the given map by its key as an array of uint.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsUintArray(data map[string]interface{}, k string) ([]uint, bool) {
	return nil, false
}

/*
GetAsUint8Array returns a value of the given map by its key as an array of uint8.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsUint8Array(data map[string]interface{}, k string) ([]uint8, bool) {
	return nil, false
}

/*
GetAsUint16Array returns a value of the given map by its key as an array of uint16.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsUint16Array(data map[string]interface{}, k string) ([]uint16, bool) {
	return nil, false
}

/*
GetAsUint32Array returns a value of the given map by its key as an array of uint32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsUint32Array(data map[string]interface{}, k string) ([]uint32, bool) {
	return nil, false
}

/*
GetAsUint64Array returns a value of the given map by its key as an array of uint64.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsUint64Array(data map[string]interface{}, k string) ([]uint64, bool) {
	return nil, false
}

/*
GetAsInt8Array returns a value of the given map by its key as an array of int8.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsInt8Array(data map[string]interface{}, k string) ([]int8, bool) {
	return nil, false
}

/*
GetAsInt16Array returns a value of the given map by its key as an array of int16.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsInt16Array(data map[string]interface{}, k string) ([]int16, bool) {
	return nil, false
}

/*
GetAsInt32Array returns a value of the given map by its key as an array of int32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsInt32Array(data map[string]interface{}, k string) ([]int32, bool) {
	return nil, false
}

/*
GetAsInt64Array returns a value of the given map by its key as an array of int64.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsInt64Array(data map[string]interface{}, k string) ([]int64, bool) {
	return nil, false
}

/*
GetAsFloat32Array returns a value of the given map by its key as an array of float32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsFloat32Array(data map[string]interface{}, k string) ([]float32, bool) {
	return nil, false
}

/*
GetAsFloat64Array returns a value of the given map by its key as an array of float64.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.

	[NOTE]      Uses mutex lock internally.
*/
func GetAsFloat64Array(data map[string]interface{}, k string) ([]float64, bool) {
	return nil, false
}

/*
ToBytesArray returns given interface{} as an array of bytes.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToBytesArray(data interface{}) ([][]byte, bool) {
	return nil, false
}

/*
ToBoolArray returns given interface{} as an array of bool.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToBoolArray(data interface{}) ([]bool, bool) {
	return nil, false
}

/*
ToStringArray returns given interface{} as an array of string.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToStringArray(data interface{}) ([]string, bool) {
	return nil, false
}

/*
ToIntArray returns given interface{} as an array of string.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToIntArray(data interface{}) ([]int, bool) {
	return nil, false
}

/*
ToUintArray returns given interface{} as an array of uint.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToUintArray(data interface{}) ([]uint, bool) {
	return nil, false
}

/*
ToUint8Array returns given interface{} as an array of uint8.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToUint8Array(data interface{}) ([]uint8, bool) {
	return nil, false
}

/*
ToUint16Array returns given interface{} as an array of uint16.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToUint16Array(data interface{}) ([]uint16, bool) {
	return nil, false
}

/*
ToUint32Array returns given interface{} as an array of uint16.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToUint32Array(data interface{}) ([]uint32, bool) {
	return nil, false
}

/*
ToUint64Array returns given interface{} as an array of uint16.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToUint64Array(data interface{}) ([]uint64, bool) {
	return nil, false
}

/*
ToInt8Array returns given interface{} as an array of int8.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToInt8Array(data interface{}) ([]int8, bool) {
	return nil, false
}

/*
ToInt16Array returns given interface{} as an array of int16.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToInt16Array(data interface{}) ([]int16, bool) {
	return nil, false
}

/*
ToInt32Array returns given interface{} as an array of int32.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToInt32Array(data interface{}) ([]int32, bool) {
	return nil, false
}

/*
ToInt64Array returns given interface{} as an array of int64.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToInt64Array(data interface{}) ([]int64, bool) {
	return nil, false
}

/*
ToFloat32Array returns given interface{} as an array of float32.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToFloat32Array(data interface{}) ([]float32, bool) {
	return nil, false
}

/*
ToFloat64Array returns given interface{} as an array of float64.

	[IMPORTANT] The second returned value indicates if the value is valid or not.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
*/
func ToFloat64Array(data interface{}) ([]float64, bool) {
	return nil, false
}

/*
IsArray returns a true if the given data interface{} is an array (slice).
*/
func IsArray(data interface{}) bool {
	return false
}

/*
ArrayEqual returns a true if the given interface{} share the same values as arrays (slices).

Valid data types for the comparison:

	[]uint8
	[]uint16
	[]uint32
	[]uint64
	[]int
	[]int8
	[]int16
	[]int32
	[]int64
	[]int
	[]float32
	[]float64
	[][]byte
*/
func ArrayEqual(d1 interface{}, d2 interface{}) bool {
	return false
}

/*
Float64SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Float64SliceEqual(v1 []float64, v2 []float64) bool {
	return false
}

/*
Float32SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Float32SliceEqual(v1 []float32, v2 []float32) bool {
	return false
}

/*
Uint8SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Uint8SliceEqual(v1 []uint8, v2 []uint8) bool {
	return false
}

/*
Uint16SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Uint16SliceEqual(v1 []uint16, v2 []uint16) bool {
	return false
}

/*
Uint32SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Uint32SliceEqual(v1 []uint32, v2 []uint32) bool {
	return false
}

/*
Uint64SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Uint64SliceEqual(v1 []uint64, v2 []uint64) bool {
	return false
}

/*
UintSliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func UintSliceEqual(v1 []uint, v2 []uint) bool {
	return false
}

/*
Int8SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Int8SliceEqual(v1 []int8, v2 []int8) bool {
	return false
}

/*
Int16SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Int16SliceEqual(v1 []int16, v2 []int16) bool {
	return false
}

/*
Int32SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Int32SliceEqual(v1 []int32, v2 []int32) bool {
	return false
}

/*
Int64SliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func Int64SliceEqual(v1 []int64, v2 []int64) bool {
	return false
}

/*
IntSliceEqual compares two slices and returns true if all the elements in the arrays match.
*/
func IntSliceEqual(v1 []int, v2 []int) bool {
	return false
}

/*
BytesSliceEqual compares two arrays of byte array.

Returns true if the two arrays' elements contain the same values.
*/
func BytesSliceEqual(v1 [][]byte, v2 [][]byte) bool {
	return false
}

/*
InterfaceSliceEqual compares two arrays of interface{} and returns true,
if all the elements are equal values.
*/
func InterfaceSliceEqual(v1 []interface{}, v2 []interface{}) bool {
	return false
}

/*
InterfaceEqual compares two interface{} variables and returns true,
if the two variables have the same value.
*/
func InterfaceEqual(d1 interface{}, d2 interface{}) bool {
	return false
}
