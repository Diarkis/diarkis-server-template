package user

/*
GetAsBoolArray returns user data by its key as an array of bool.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsBoolArray(k string) ([]bool, bool) {
	return nil, false
}

/*
GetAsStringArray returns user data by its key as an array of string.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsStringArray(k string) ([]string, bool) {
	return nil, false
}

/*
GetAsIntArray returns user data by its key as an array of int.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsIntArray(k string) ([]int, bool) {
	return nil, false
}

/*
GetAsUint8Array returns user data by its key as an array of uint8.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsUint8Array(k string) ([]uint8, bool) {
	return nil, false
}

/*
GetAsUint16Array returns user data by its key as an array of uint16.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsUint16Array(k string) ([]uint16, bool) {
	return nil, false
}

/*
GetAsUint32Array returns user data by its key as an array of uint32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsUint32Array(k string) ([]uint32, bool) {
	return nil, false
}

/*
GetAsUint64Array returns user data by its key as an array of uint64.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsUint64Array(k string) ([]uint64, bool) {
	return nil, false
}

/*
GetAsInt8Array returns user data by its key as an array of int8.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsInt8Array(k string) ([]int8, bool) {
	return nil, false
}

/*
GetAsInt16Array returns user data by its key as an array of int16.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsInt16Array(k string) ([]int16, bool) {
	return nil, false
}

/*
GetAsInt32Array returns user data by its key as an array of int32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsInt32Array(k string) ([]int32, bool) {
	return nil, false
}

/*
GetAsInt64Array returns user data by its key as an array of int64.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.
	[IMPORTANT] If one or more elements in the array are to be float64,
	            it is an exception that will be auto-converted to the specified numerical data type
	            and the second return value will be true.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsInt64Array(k string) ([]int64, bool) {
	return nil, false
}

/*
GetAsFloat64Array returns user data by its key as an array of float64.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be false, if the value data type is not invalid or missing.
	            Valid means that the value is of the correct data type.
	[IMPORTANT] If one or more elements in the array are not the specified data type,
	            the second return value will be false.

	[NOTE]      Uses mutex lock internally.
*/
func (u *User) GetAsFloat64Array(k string) ([]float64, bool) {
	return nil, false
}
