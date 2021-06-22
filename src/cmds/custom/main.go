package customcmds

import (
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

const customVer = 2
const helloCmdID = 10
const pushCmdID = 11

var logger = log.New("CUSTOM")

func Expose() {
	server.HandleCommand(customVer, helloCmdID, helloCmd)
	server.HandleCommand(customVer, helloCmdID, afterHelloCmd)
	server.HandleCommand(customVer, pushCmdID, pushCmd)
}

func helloCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Hello command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	// if this is executed as UDP, reliable = true means sending the packet as RRUDP
	reliable := true
	// we send a response back to the client with the byte array sent from the client
	userData.ServerRespond(payload, ver, cmd, server.Ok, reliable)
	// move on to the next command handler if there is any
	next(nil)
}

func afterHelloCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("This is executed after Hello command has been handled")
	next(nil)
}

func pushCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Push command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	// if this is executed as UDP, reliable = true means sending the packet as RRUDP
	reliable := true
	// we send a push packet to the client that sent the data to this command
	userData.ServerPush(ver, cmd, payload, reliable)
	// move on to the next command handler if there is any
	next(nil)
}
