// © 2019-2024 Diarkis Inc. All rights reserved.

package customcmds

import (
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"

	"github.com/Diarkis/diarkis-server-template/puffer/go/custom"
)

// Custom sample command IDs
const CustomVer = 2

const helloCmdID = 10
const pushCmdID = 11
const resonanceCmdID = 13

// Client error log
const clientErrLog = 12

// MatchMaker command IDs
const matchmakerAdd = 100
const matchmakerRm = 101
const matchmakerSearch = 102
const matchmakerComplete = 103 // sent when room is full

const MatchedMemberLeaveCmdID = 1011 // sent when matched ticket member leaves

// P2P command IDs
const p2pReportAddr = 110
const p2pInit = 111

// Online Status command ID
const getUserStatusListCmdID = 500

const mmAddInterval = 40 // 40 seconds

var logger = log.New("CUSTOM")

func Expose() {
	// defined in main.go
	diarkisexec.SetServerCommandHandler(CustomVer, helloCmdID, helloCmd)
	diarkisexec.SetServerCommandHandler(CustomVer, pushCmdID, pushCmd)
	// puffer version sample
	diarkisexec.SetServerCommandHandler(custom.EchoVer, custom.EchoCmd, echoPufferCmd)
	// defined in matchmaker.go
	diarkisexec.SetServerCommandHandler(CustomVer, matchmakerAdd, addToMatchMaker)
	diarkisexec.SetServerCommandHandler(CustomVer, matchmakerSearch, searchMatchMaker)
	// defined in p2p.go
	diarkisexec.SetServerCommandHandler(CustomVer, p2pReportAddr, reportP2PAddr)
	diarkisexec.SetServerCommandHandler(CustomVer, p2pInit, initP2P)
	// defined in field.go
	diarkisexec.SetServerCommandHandler(custom.GetFieldInfoVer, custom.GetFieldInfoCmd, getFieldInfo)
	// defined in onlinestatus.go
	diarkisexec.SetServerCommandHandler(CustomVer, getUserStatusListCmdID, getUserStatusList)
	// defined in resonance.go
	diarkisexec.SetServerCommandHandler(CustomVer, resonanceCmdID, resonanceCmd)

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

func pushCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Push command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	// if this is executed as UDP, reliable = true means sending the packet as RUDP
	reliable := true
	// we send a push packet to the client that sent the data to this command
	userData.ServerPush(ver, cmd, payload, reliable)
	// move on to the next command handler if there is any
	next(nil)
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
	userData.ServerRespond(echoData.Pack(), ver, cmd, server.Ok, true) // You can get []byte by using Pack. ( echoData.Pack eauals payload in this example.)
	// move on to the next command handler if there is any
	next(nil)
}
