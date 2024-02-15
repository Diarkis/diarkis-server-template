package customcmds

import (
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

// Custom sample command IDs
const customVer = 2
const helloCmdID = 10
const pushCmdID = 11

// Client error log
const clientErrLog = 12

// MatchMaker command IDs
const matchmakerAdd = 100
const matchmakerRm = 101
const matchmakerSearch = 102
const matchmakerComplete = 103 // sent when room is full
// P2P command IDs
const p2pReportAddr = 110
const p2pInit = 111

const mmAddInterval = 40 // 40 seconds

var logger = log.New("CUSTOM")

func Expose() {
	// defined in main.go
	server.HandleCommand(customVer, helloCmdID, helloCmd)
	server.HandleCommand(customVer, helloCmdID, afterHelloCmd)
	server.HandleCommand(customVer, pushCmdID, pushCmd)
	server.HandleCommand(customVer, pushCmdID, outputClientErrLog)
	// defined in matchmaker.go
	server.HandleCommand(customVer, matchmakerAdd, addToMatchMaker)
	server.HandleCommand(customVer, matchmakerSearch, searchMatchMaker)
	// defined in p2p.go
	server.HandleCommand(customVer, p2pReportAddr, reportP2PAddr)
	server.HandleCommand(customVer, p2pInit, initP2P)
	server.HandleCommand(GetFieldInfoVer, GetFieldInfoCmd, getFieldInfo)
}

func helloCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Hello command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	// if this is executed as UDP, reliable = true means sending the packet as RUDP
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
	// if this is executed as UDP, reliable = true means sending the packet as RUDP
	reliable := true
	// we send a push packet to the client that sent the data to this command
	userData.ServerPush(ver, cmd, payload, reliable)
	// move on to the next command handler if there is any
	next(nil)
}

func outputClientErrLog(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Error("Client error log - client sid:%v uid:%v: %v", userData.SID, userData.ID, string(payload))
	next(nil)
}
