package process

/*
IsTerminating returns true when the process is in termination mode after invoking Stop
*/
func IsTerminating() bool {
	return false
}

/*
GetShutdownTimeout returns the server process shutdown timeout in seconds.

	[IMPORTANT] Default is 10 seconds.
	[IMPORTANT] If DIARKIS_SHUTDOWN_TIMEOUT env is given, shutdown timeout changes.
*/
func GetShutdownTimeout() int64 {
	return 0
}

/*
OnStart Registers a callback on process start to be executed BEFORE the process is ready
*/
func OnStart(task func()) {
}

/*
OnStop registers a callback on process stop to be executed BEFORE termination tasks are executed
*/
func OnStop(cb func()) {
}

/*
OnSIGHUP Registers a callback function on SIGHUP signal
*/
func OnSIGHUP(callback func()) {
}

/*
OnSIGUSR1 Registers a callback function on SIGUSR1 signal
*/
func OnSIGUSR1(callback func()) {
}

/*
OnSIGUSR2 Registers a callback function on SIGUSR2 signal
*/
func OnSIGUSR2(callback func()) {
}

/*
OnTerminate Registers a callback on process stop to be executed BEFORE the process termination
*/
func OnTerminate(task func(func(error))) {
}

/*
Start Starts the process
*/
func Start() {
}

/*
Stop Stops the process with exit code 0
*/
func Stop() {
}

/*
StopWithError Stops the process with an error and exit code 1
*/
func StopWithError(err error) {
}

/*
IsRunning returns true if the process is running
*/
func IsRunning() bool {
	return false
}
