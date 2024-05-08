package diarkisexec

import (
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/user"
)

/*
Options represents module configurations.
*/
type Options struct {
	ConfigPath     string
	ExposeCommands bool
}

/*
Modules represents module declarations.
*/
type Modules struct {
	Dive       *Options
	DM         *Options
	Field      *Options
	Group      *Options
	Room       *Options
	MatchMaker *Options
	Session    *Options
	Metrics    *Options
	Notifier   *Options
	P2P        *Options
}

/*
CommandHandler callback for a UDP/TCP server command.

	[IMPORTANT] At the end of callback operations, you MUST call next func(error).
	            Not calling next func(error) will cause the server not to be able to handle the next incoming commands from the client.
*/
type CommandHandler func(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error))

/*
MeshCommandHandler callback for a mesh command.

	[IMPORTANT] In order to return a response value for SendMeshRequest, the handler function must return the response value as a byte array.
*/
type MeshCommandHandler func(req map[string]interface{}) ([]byte, error)

/*
SetupDiarkis initializes core and optional modules.

	[NOTE] The configuration paths must be relative to the executable server binary file location.
*/
func SetupDiarkis(logConfigPath, meshConfigPath string, m *Modules) {
}

/*
SetupDiarkisHTTPServer declares the server to be HTTP server.

	[IMPORTANT] A server cannot have multiple network protocols (HTTP, TCP, and UDP).

	[NOTE] The configuration paths must be relative to the executable server binary file location.
*/
func SetupDiarkisHTTPServer(config string) {
}

/*
SetupDiarkisTCPServer declares the server to be TCP server.

	[IMPORTANT] A server cannot have multiple network protocols (HTTP, TCP, and UDP).

	[NOTE] The configuration paths must be relative to the executable server binary file location.
*/
func SetupDiarkisTCPServer(config string) {
}

/*
SetupDiarkisUDPServer declares the server to be UDP server.

	[IMPORTANT] A server cannot have multiple network protocols (HTTP, TCP, and UDP).

	[NOTE] The configuration paths must be relative to the executable server binary file location.
*/
func SetupDiarkisUDPServer(config string) {
}

/*
SetServerCommandHandler assigns a callback to the given command ver and cmd.

	[IMPORTANT] This function must be invoked BEFORE calling diarkisexec.StartDiarkis()
	[IMPORTANT] Diarkis' built-in commands use ver ranging from 0 to 1.
	            You may NOT assign callbacks with those ver values.
	[IMPORTANT] You may NOT assign a callback to the same ver and cmd combination.
	[IMPORTANT] Attempting to assign multiple callbacks to the same ver and cmd combination
	            as the built-in commands of Diarkis will result in panic.
	            Even if you do not use any built-in commands by setting ExposeCommands to false,
	            Diarkis start does NOT allow you to use built-in commands' ver and cmd.
*/
func SetServerCommandHandler(ver uint8, cmd uint16, handler CommandHandler) {
}

/*
NewLogger creates a new instance of log.Logger.

	[IMPORTANT] This function must be invoked BEFORE calling diarkisexec.StartDiarkis()
*/
func NewLogger(name string) *log.Logger {
	return nil
}

/*
StartDiarkis starts the Diarkis server.

	[IMPORTANT] This function blocks because it starts the process as a server,
	            which means that no operations after calling of this function will be executed at all.
	[IMPORTANT] SetupDiarkis must be called before calling StartDiarkis to properly setup Diarkis server.
*/
func StartDiarkis() {
}
