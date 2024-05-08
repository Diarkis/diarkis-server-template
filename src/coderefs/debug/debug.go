package debug

import (
	_ "net/http/pprof"
)

/*
Enable Turns on debug functions
*/
func Enable() {
}

/*
Disable Turns off debug functions
*/
func Disable() {
}

/*
IsEnabled returns true, if debug.Enable() has been called
*/
func IsEnabled() bool {
	return false
}

/*
RunMemoryWatch Reads memory usage statistics and outputs them to stdout stream
*/
func RunMemoryWatch() {
}
