package msg

const LvlCustom = " CUSTOM"
const LvlVerbose = "VERBOSE"
const LvlSys = " SYSTEM"
const LvlNetwork = " NETWORK"
const LvlTrace = "  TRACE"
const LvlDebug = "  DEBUG"
const LvlInfo = "   INFO"
const LvlNotice = " NOTICE"
const LvlWarn = "WARNING"
const LvlErr = "  ERROR"
const LvlFatal = "  FATAL"
const OutputPrefix = "[%!s(MISSING)]%!s(MISSING)%!s(MISSING) %!s(MISSING)"
const OutputSuffix = "\n"
const TimeZoneUTC = 0
const TimeZoneLocal = 1

/*
SetCustomOutput Registers a custom function to create a log output string
*/
func SetCustomOutput(custom func(formatted bool, name string, prefix string, level string, vals []interface{}) string) {
}

/*
UseCustomOutput Enables custom output
*/
func UseCustomOutput() {
}

/*
SetTimeZone Sets time zone of the log
*/
func SetTimeZone(_timeZone string) {
}

/*
SetPrefix Sets a prefix for every log
*/
func SetPrefix(pf string) {
}

/*
EnableFlat Turns on flatten format of log output
*/
func EnableFlat() {
}

/*
DisableFlat Turns off flatten format of log output
*/
func DisableFlat() {
}

/*
Write output custom level log
*/
func Write(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Verbose Outputs log
*/
func Verbose(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Network Outputs log
*/
func Network(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Sys Outputs log
*/
func Sys(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Trace Outputs log
*/
func Trace(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Debug Outputs log
*/
func Debug(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Info Outputs log
*/
func Info(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Notice Outputs log
*/
func Notice(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Warn Outputs log
*/
func Warn(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Error Outputs log
*/
func Error(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
Fatal Outputs log
*/
func Fatal(formatted bool, name string, vals ...interface{}) string {
	return ""
}

/*
GetStackTrace returns stack trace [INTERNAL USE ONLY]
*/
func GetStackTrace(delimiter string, depth int) string {
	return ""
}
