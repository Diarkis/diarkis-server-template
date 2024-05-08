package config

/*
Load Load configurations from a file into memory
*/
func Load(name string, path string) map[string]interface{} {
	return nil
}

/*
Set Sets a configuration value by name and its property name
*/
func Set(name string, property string, value interface{}) {
}

/*
Read retrieves configurations by name and return as a map
*/
func Read(name string) map[string]interface{} {
	return nil
}

/*
ReadProperty returns a property of a configuration object map
*/
func ReadProperty(name string, property string) interface{} {
	return nil
}

/*
GetAsString retrieves a configuration value as a string
*/
func GetAsString(name string, property string, _default string) string {
	return ""
}

/*
GetAsInt retrieves a configuration value as an int
*/
func GetAsInt(name string, property string, _default int) int {
	return 0
}

/*
GetAsInt8 retrieves a configuration value as an int8
*/
func GetAsInt8(name string, property string, _default int8) int8 {
	return 0
}

/*
GetAsUint8 retrieves a configuration value as a uint8
*/
func GetAsUint8(name string, property string, _default uint8) uint8 {
	return 0
}

/*
GetAsInt16 retrieves a configuration value as an int16
*/
func GetAsInt16(name string, property string, _default int16) int16 {
	return 0
}

/*
GetAsUint16 retrieves a configuration value as a uint16
*/
func GetAsUint16(name string, property string, _default uint16) uint16 {
	return 0
}

/*
GetAsInt32 retrieves a configuration value as an int32
*/
func GetAsInt32(name string, property string, _default int32) int32 {
	return 0
}

/*
GetAsUint32 retrieves a configuration value as a uint32
*/
func GetAsUint32(name string, property string, _default uint32) uint32 {
	return 0
}

/*
GetAsInt64 retrieves a configuration value as an int64
*/
func GetAsInt64(name string, property string, _default int64) int64 {
	return 0
}

/*
GetAsUint64 retrieves a configuration value as a uint64
*/
func GetAsUint64(name string, property string, _default uint64) uint64 {
	return 0
}

/*
GetAsFloat64 retrieves a configuration value as a float64
*/
func GetAsFloat64(name string, property string, _default float64) float64 {
	return 0
}

/*
GetAsBool Retrieves a configuration value as a boolean
*/
func GetAsBool(name string, property string, _default bool) bool {
	return false
}

/*
GetAsStruct Retrieves a configuration value and parse it into the passed struct
_struct should be a pointer for a struct instance
*/
func GetAsStruct(name string, _struct any) {
}
