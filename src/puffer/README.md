This is a readme describing puffer, a packet definition generation tool.
Reference: https://docs.diarkis.io/docs/server/v1.0.0-rc1/diarkis/puffer/index.html

## Usage

### command

`make gen` - generate packet definition files
`make clean` - remove generated files

## Code Example

```
func Expose() { // add command using puffer handler
    ...
	// puffer version sample
	diarkisexec.SetServerCommandHandler(custom.EchoVer, custom.EchoCmd, echoPufferCmd) // sample puffer echo command
    ...
}

func echoPufferCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Hello puffer command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	// unpack []byte to struct
	echoData := custom.NewEcho()
	err := echoData.Unpack(payload) // You can unpack []byte to go struct

	if err != nil {
		logger.Error("Failed to unpack echo data: %v", err)
		userData.ServerRespond(nil, ver, cmd, server.Err, true)
		next(nil)
		return
	}

	logger.Debug("Unpacked echo data: %#v", echoData)
    userData.ServerRespond(echoData.Pack(), ver, cmd, server.Ok, true) // You can get []byte by using Pack. ( echoData.Pack equals payload in this example.)
	// move on to the next command handler if there is any
	next(nil)
}

```

## directory architecture

`make gen` will generate the following directories.
go - go code
cs - c# code
cpp - C++ code

The code is generated for each runtime, and can be used on the client for cs and cpp, and can be used on the client and server in common for go.
`json_definitions` contains puffer definitions.
`json_definitions/samples` is a sample definition that is not used in the code right now.
We recommend that you delete it when actually using it.

## Note

The current version of puffer has a problem with overwriting the same command name even if the package name is different.
This will be fixed in the next version.
