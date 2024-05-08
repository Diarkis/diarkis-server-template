package user

/*
GetAsBool returns user data by its key as a bool.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.

	[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsBool(k string) (bool, bool) {
	return false, false
}

/*
GetAsString returns user data by its key as a string.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.

	[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsString(k string) (string, bool) {
	return "", false
}

/*
GetAsBytes returns user data by its key as a string.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the value held is a base64 encoded string, it will be returned as a byte array as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsBytes(k string) ([]byte, bool) {
	return nil, false
}

/*
GetAsUint8 returns user data by its key as a uint8.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsUint8(k string) (uint8, bool) {
	return 0, false
}

/*
GetAsUint16 returns user data by its key as a uint16.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsUint16(k string) (uint16, bool) {
	return 0, false
}

/*
GetAsUint returns user data by its key as a uint.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsUint(k string) (uint, bool) {
	return 0, false
}

/*
GetAsUint32 returns user data by its key as a uint32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsUint32(k string) (uint32, bool) {
	return 0, false
}

/*
GetAsUint64 returns user data by its key as a uint32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsUint64(k string) (uint64, bool) {
	return 0, false
}

/*
GetAsInt8 returns user data by its key as a int8.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsInt8(k string) (int8, bool) {
	return 0, false
}

/*
GetAsInt16 returns user data by its key as a int16.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsInt16(k string) (int16, bool) {
	return 0, false
}

/*
GetAsInt returns user data by its key as a int.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsInt(k string) (int, bool) {
	return 0, false
}

/*
GetAsInt32 returns user data by its key as a int32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsInt32(k string) (int32, bool) {
	return 0, false
}

/*
GetAsInt64 returns user data by its key as a int64.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsInt64(k string) (int64, bool) {
	return 0, false
}

/*
GetAsFloat32 returns user data by its key as a float32.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsFloat32(k string) (float32, bool) {
	return 0, false
}

/*
GetAsFloat64 returns user data by its key as a float64.

	[IMPORTANT] The second returned value indicates if the value exists or not.
	            The second value will be true, if the value is invalid or missing.
	[IMPORTANT] If the held value is float64, it will be converted to the correct type and returned as a valid value.

[NOTE] Uses mutex lock internally.
*/
func (u *User) GetAsFloat64(k string) (float64, bool) {
	return 0, false
}
