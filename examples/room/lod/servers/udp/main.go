// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/derror"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"

	proom "github.com/Diarkis/diarkis-server-template/examples/room/lod/puffer/go/room"
)

var logger = log.New("UDP")

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := ""

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Room: &diarkisexec.Options{ConfigPath: "configs/shared/room.json", ExposeCommands: true},
		Dive: &diarkisexec.Options{ConfigPath: "configs/shared/dive.json"},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")
	diarkisexec.SetServerCommandHandler(proom.BroadcastLoDVer, proom.BroadcastLoDCmd, handleRoomBroadcastLoD)
	diarkisexec.SetServerCommandHandler(proom.GetLoDInfoVer, proom.GetLoDInfoCmd, handleRoomGetLoDInfo)

	diarkisexec.StartDiarkis()
}

func handleRoomBroadcastLoD(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {

	proto := proom.NewBroadcastLoD()
	err := proto.Unpack(payload)
	if err != nil {
		userData.ServerRespond(derror.ErrData(err.Error(), derror.InvalidParameter(0)), ver, cmd, server.Bad, true)
		next(err)
		return
	}

	// TODO: Broadcast the LoD data
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		userData.ServerRespond(derror.ErrData("Not in the room", derror.NotAllowed(0)), ver, cmd, server.Bad, true)
		return
	}

	room.Broadcast(roomID, userData, ver, cmd, proto.Payload, true)
	userData.ServerRespond(nil, ver, cmd, server.Ok, true)
	next(nil)
}

func handleRoomGetLoDInfo(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {

	// TODO: Get the LoD information from configuration
	lodInfo := proom.NewGetLoDInfoResponse()
	lodInfo.SyncIntervalForNearby = 16
	lodInfo.SyncIntervalForFar = 2000
	lodInfo.MaxDistanceForNearby = 10000
	lodInfo.MaxDistanceForFar = 40000

	userData.ServerRespond(lodInfo.Pack(), ver, cmd, server.Ok, true)
	next(nil)
}
