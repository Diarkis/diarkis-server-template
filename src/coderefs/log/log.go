package log

/*
Logger Logger data structure
*/
type Logger struct{ Name string }

/*
IsVerbose returns true if the log level is verbose
*/
func IsVerbose() bool {
	return false
}

/*
IsNetwork returns true if the log level is network
*/
func IsNetwork() bool {
	return false
}

/*
IsSys returns true if the log level is sys
*/
func IsSys() bool {
	return false
}

/*
IsDebug returns true if the log level is debug
*/
func IsDebug() bool {
	return false
}

/*
IsInfo returns true if the log level is info
*/
func IsInfo() bool {
	return false
}

/*
IsNotice returns true if the log level is notice
*/
func IsNotice() bool {
	return false
}

/*
IsWarn returns true if the log level is warn
*/
func IsWarn() bool {
	return false
}

/*
IsError returns true if the log level is error
*/
func IsError() bool {
	return false
}

/*
IsFatal returns true if the log level is fatal
*/
func IsFatal() bool {
	return false
}

/*
Setup Load a configuration file into memory - pass an empty string to load nothing
*/
func Setup(confpath string) {
}

/*
UpdateConfigs Apply configurations
It can be executed in runtime to update the configurations also
*/
func UpdateConfigs(configs map[string]interface{}) {
}

/*
UseCustomOutput Enables custom output for all log data.
*/
func UseCustomOutput() {
}

/*
SetCustomOutput Registers a custom function to create a log output string
*/
func SetCustomOutput(custom func(formatted bool, name string, prefix string, level string, vals []interface{}) string) {
}

/*
UseJSONFormat formats the log with JSON format.
*/
func UseJSONFormat() {
}

/*
UseStackdriverLogging is an alias of UseJSONFormat
*/
func UseStackdriverLogging() {
}

/*
AddFormat assigns a formatting function to the given name.

	[CRITICALLY IMPORTANT] You must NOT read from a map directly in the callback because of possible concurrent map access panic.

When formatted logging is used (Verbosef, Networkf, Sysf, Debugf, Infof, Warnf, Errorf, Fatalf), with the given name,
it will execute the assigned callback.

# Example

The example below will output a specific user property safely (UserFlag is just an example, of course).

	log.AddFormat("UserFlag", func(v interface{}) interface{} {

	  // If the value is not *User, we ignore and do nothing
	  if _, ok := v.(*User) {
	    return v
	  }

	  // We make sure to lock the user to print
	  ud := v.(*User)

	  ud.RLock()
	  defer ud.RUnlock()

	  // Format the output as a string
	  formatted := log.Fmt("Properties:", ud.Data["UserFlag"])
	  return formatted

	})
*/
func AddFormat(name string, cb func(v interface{}) interface{}) bool {
	return false
}

/*
SetPrefix Sets a prefix for every log
*/
func SetPrefix(prefix string) {
}

/*
CloseFileOutput [INTERNAL USE ONLY] Closes the file descriptor of log data if used
*/
func CloseFileOutput(next func(error)) {
}

/*
SetVerboseLevel set log level to verbose
*/
func SetVerboseLevel() {
}

/*
SetNetworkLevel set log level to sys
*/
func SetNetworkLevel() {
}

/*
SetSysLevel set log level to sys
*/
func SetSysLevel() {
}

/*
SetDebugLevel set log level to debug
*/
func SetDebugLevel() {
}

/*
SetInfoLevel set log level to info
*/
func SetInfoLevel() {
}

/*
SetNoticeLevel set log level to notice
*/
func SetNoticeLevel() {
}

/*
SetWarnLevel set log level to warn
*/
func SetWarnLevel() {
}

/*
SetErrorLevel set log level to error
*/
func SetErrorLevel() {
}

/*
SetFatalLevel set log level to fatal
*/
func SetFatalLevel() {
}

/*
Level Sets logging level
*/
func Level(_level int) {
}

/*
FmtGrey Colorize the given string
*/
func FmtGrey(msg string) string {
	return ""
}

/*
FmtBlue Colorize the given string
*/
func FmtBlue(msg string) string {
	return ""
}

/*
FmtLBlue Colorize the given string
*/
func FmtLBlue(msg string) string {
	return ""
}

/*
FmtDBlue Colorize the given string
*/
func FmtDBlue(msg string) string {
	return ""
}

/*
FmtGreen Colorize the given string
*/
func FmtGreen(msg string) string {
	return ""
}

/*
FmtYellow Colorize the given string
*/
func FmtYellow(msg string) string {
	return ""
}

/*
FmtPurple Colorize the given string
*/
func FmtPurple(msg string) string {
	return ""
}

/*
FmtRed Colorize the given string
*/
func FmtRed(msg string) string {
	return ""
}

/*
FmtRedBg Colorize the given string
*/
func FmtRedBg(msg string) string {
	return ""
}

/*
New Creates a new logger
*/
func New(name string) *Logger {
	return nil
}

/*
EnableCustom allows logger.Write()
*/
func (logger *Logger) EnableCustom() {
}

/*
DisableCustom blocks logger.Write()
*/
func (logger *Logger) DisableCustom() {
}

/*
Write outputs a log to stdout stream
*/
func (logger *Logger) Write(vals ...interface{}) {
}

/*
Verbose Outputs a log to stdout stream
*/
func (logger *Logger) Verbose(vals ...interface{}) {
}

/*
Network Outputs a log to stdout stream
*/
func (logger *Logger) Network(vals ...interface{}) {
}

/*
Sys Outputs a log to stdout stream
*/
func (logger *Logger) Sys(vals ...interface{}) {
}

/*
Trace Outputs a log to stdout stream as Sys level log with stack trace. Useful for debugging
*/
func (logger *Logger) Trace(vals ...interface{}) {
}

/*
Debug Outputs a log to stdout stream
*/
func (logger *Logger) Debug(vals ...interface{}) {
}

/*
Info Outputs a log to stdout stream
*/
func (logger *Logger) Info(vals ...interface{}) {
}

/*
Notice Outputs a log to stdout stream
*/
func (logger *Logger) Notice(vals ...interface{}) {
}

/*
Warn Outputs a log to stdout stream
*/
func (logger *Logger) Warn(vals ...interface{}) {
}

/*
Error Outputs a log to stdout stream
*/
func (logger *Logger) Error(vals ...interface{}) {
}

/*
Fatal Outputs a log to stdout stream
*/
func (logger *Logger) Fatal(vals ...interface{}) {
}

/*
Verbosef outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Verbosef("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Verbosef(vals ...interface{}) {
}

/*
Networkf outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Networkf("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Networkf(vals ...interface{}) {
}

/*
Sysf outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Sysf("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Sysf(vals ...interface{}) {
}

/*
Debugf outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Debugf("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Debugf(vals ...interface{}) {
}

/*
Infof outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Infof("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Infof(vals ...interface{}) {
}

/*
Noticef outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Noticef("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Noticef(vals ...interface{}) {
}

/*
Warnf outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Warnf("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Warnf(vals ...interface{}) {
}

/*
Errorf outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Errorf("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Errorf(vals ...interface{}) {
}

/*
Fatalf outputs formatted logging.

This is useful when you need to format the logging message.

	[IMPORTANT] FormattedText supports JSON format log output as well.

Example:

	logger.Fatalf("This is a log", "intVar", 123, "stringVar", "aaa", "object", obj)
*/
func (logger *Logger) Fatalf(vals ...interface{}) {
}

/*
FmtVerbose Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtVerbose(vals ...interface{}) string {
	return ""
}

/*
FmtNetwork Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtNetwork(vals ...interface{}) string {
	return ""
}

/*
FmtSys Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtSys(vals ...interface{}) string {
	return ""
}

/*
FmtDebug Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtDebug(vals ...interface{}) string {
	return ""
}

/*
FmtInfo Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtInfo(vals ...interface{}) string {
	return ""
}

/*
FmtNotice Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtNotice(vals ...interface{}) string {
	return ""
}

/*
FmtWarn Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtWarn(vals ...interface{}) string {
	return ""
}

/*
FmtError Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtError(vals ...interface{}) string {
	return ""
}

/*
FmtFatal Formats given strings and returns it for logging purpose
*/
func (logger *Logger) FmtFatal(vals ...interface{}) string {
	return ""
}

/*
Fmt alias of fmt.Sprint
*/
func Fmt(vals ...interface{}) string {
	return ""
}
