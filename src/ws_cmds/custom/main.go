package customcmds

import (
	"encoding/json"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

const customVer = 2
const helloCmdID = 10
const pushCmdID = 11

var logger = log.New("CUSTOM")

func Expose() {
	server.Command(customVer, helloCmdID, helloCmd)
	server.Command(customVer, helloCmdID, afterHelloCmd)
	server.Command(customVer, pushCmdID, pushCmd)
}

func helloCmd(userData *user.User, next func(error)) {
	logger.Debug("Hello command has received %#v from the client SID:%s - UID:%s", userData.Data, userData.SID, userData.ID)
	payload, err := json.Marshal(userData.Data)
	if err != nil {
		userData.Respond([]byte(err.Error()), server.Bad, true)
		next(err)
		return
	}
	// we send a response back to the client with the byte array sent from the client
	userData.Respond(payload, server.Ok, true)
	// move on to the next command handler if there is any
	next(nil)
}

func afterHelloCmd(userData *user.User, next func(error)) {
	logger.Debug("This is executed after Hello command has been handled")
	next(nil)
}

func pushCmd(userData *user.User, next func(error)) {
	logger.Debug("Push command has received %#v from the client SID:%s - UID:%s", userData.Data, userData.SID, userData.ID)
	payload, err := json.Marshal(userData.Data)
	if err != nil {
		logger.Error("Failed to create payload %v", err)
		next(err)
		return
	}
	// we send a push packet to the client that sent the data to this command
	userData.Push(payload, true)
	// move on to the next command handler if there is any
	next(nil)
}
