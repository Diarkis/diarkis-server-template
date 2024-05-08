package log

/*
ChangeLogLevelByName is automatically called
when the server receives SIGUSR1 signal and when there is /tmp/DIARKIS_SIGUSR1.

This function acts as a toggle.

# DIARKIS_SIGUSR1 File Format

The first line must be ChangeLogLevelByName and the rest will be parameters.

	[NOTE] line breaks will be the delimiter.

# How To Disable Runtime Log Level Change

	params[0] = "disable"

# How To Enable Runtime Log For Specific Log By Name

	params[0] = "enable"
	params[1] = "$(log_name)"  // e.i. MESH, SERVER etc.
	params[2] = "$(log_level)" // e.i. verbose, network etc.
*/
func ChangeLogLevelByName(params []string) {
}
