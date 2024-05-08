package diarkis

const (
	Version = "1.0.0"
	Author  = "Diarkis"
)

/*
IsTerminating returns true when the process is shutting down
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
IsStarted returns true after OnStarted callbacks have been executed.

If true, Diarkis server process has finished all preparation operations and ready.
*/
func IsStarted() bool {
	return false
}

/*
GetCPUNumUsed returns the number of CPU the process is using
*/
func GetCPUNumUsed() int {
	return 0
}

/*
OnStart Registers a callback on diarkis process start.

task - Callback to be invoked when Diarkis server is started.
*/
func OnStart(task func()) {
}

/*
OnReady Registers a callback on diarkis ready (This will be invoked before OnStarted while the server is being started)

This event is raised after all OnStart tasks have completed their operations.

	task - Callback to be invoked while Diarkis server process is being started.
	       Must call next function at the end of the callback to make sure the process moves on to the next task.
*/
func OnReady(task func(func(error))) {
}

/*
OnStarted registers a callback to be invoked when diarkis server process completed all starting operations.

This event is raised after all OnReady tasks have completed their operations.

	callback - Callback to be invoked when Diarkis server finishes all pre-start operations and ready.
*/
func OnStarted(callback func()) {
}

/*
OnReadyLast it is used by ONLY mesh internally
*/
func OnReadyLast(lastTask func(func())) {
}

/*
OnTerminate Registers a callback on diarkis process terminate

	task - Callback to be invoked when the server receives SIGTERM.
	       next function must be called at the end of the callback to move on to the next callbacks.
*/
func OnTerminate(task func(func(error))) {
}

/*
OnSIGUSR1 Registers a callback on SIGUSR1 signal.
When the signal is captured, Diarkis will look for a specific file in /tmp/ directory.

The file in /tmp/ must be named as DIARKIS_SIGUSR1 and the file must contain a task name.

The task name is then read by Diarkis and if a task callback with the task name exists,
it is then invoked.

	taskName - Associated with the callback.
	task     - Callback to be invoked when the server receives the signal.
*/
func OnSIGUSR1(taskName string, task func()) bool {
	return false
}

/*
OnSIGUSR1WithParams assigns a callback on SIGUSR1 signal.
When the signal is captured, Diarkis will look for a specific file in /tmp/ directory.

The file in /tmp/ must be named as DIARKIS_SIGUSR1 and the file must contain a task name.

The task name is then read by Diarkis and if a task callback with the task name exists,
it is then invoked.

Task file format:

The first line must be the task name and the following lines will be parameters.

	$(task_name)
	$(parameter1)
	$(parameter2)
	$(...)

Parameters

	taskName - Associated with the callback.
	task     - Callback to be invoked when the server receives the signal.
*/
func OnSIGUSR1WithParams(taskName string, task func(params []string)) bool {
	return false
}

/*
OnSIGUSR2 Registers a callback on SIGUSR1 signal

When the signal is captured, Diarkis will look for a specific file in /tmp/ directory.

The file in /tmp/ must be named as DIARKIS_SIGUSR2 and the file must contain a task name.

The task name is then read by Diarkis and if a task callback with the task name exists,
it is then invoked.

	taskName - Associated with the callback.
	task     - Callback to be invoked when the server receives the signal.
*/
func OnSIGUSR2(taskName string, task func()) bool {
	return false
}

/*
OnSIGUSR2WithParams assigns a callback on SIGUSR2 signal.
When the signal is captured, Diarkis will look for a specific file in /tmp/ directory.

The file in /tmp/ must be named as DIARKIS_SIGUSR2 and the file must contain a task name.

The task name is then read by Diarkis and if a task callback with the task name exists,
it is then invoked.

Task file format:

The first line must be the task name and the following lines will be parameters.

	$(task_name)
	$(parameter1)
	$(parameter2)
	$(...)

Parameters

	taskName - Associated with the callback.
	task     - Callback to be invoked when the server receives the signal.
*/
func OnSIGUSR2WithParams(taskName string, task func(params []string)) bool {
	return false
}

/*
OnSIGHUP Registers a callback on SIGHUP signal

	task - Callback to be invoked when the server receives the signal.
*/
func OnSIGHUP(task func()) {
}

/*
Start Starts diarkis process and stays running until it is instructed to stop.

This function will block until the process is terminated. It will raise OnReady event.
*/
func Start() {
}

/*
Run Starts diarkis process and stops when all operations are done
*/
func Run() {
}

/*
Stop Stops diarkis process exit code is 0
*/
func Stop() {
}

/*
StopWithError Stops diarkis process with an error exit code is 1
*/
func StopWithError(err error) {
}

/*
IsRunning returns true if diarkis process is running
*/
func IsRunning() bool {
	return false
}

/*
GetDiarkisPath Returns the absolute path to diarkis
*/
func GetDiarkisPath() string {
	return ""
}

/*
GetVer Returns the version of diarkis
*/
func GetVer() string {
	return ""
}
