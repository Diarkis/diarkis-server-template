// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/config"
	"github.com/Diarkis/diarkis/derror"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"

	"github.com/Diarkis/diarkis-server-template/examples/room/lod/lodmanager"
	proom "github.com/Diarkis/diarkis-server-template/examples/room/lod/puffer/go/room"
)

var ( // lod.json settings
	SyncIntervalForNearby int32
	SyncIntervalForFar    int32
	MaxDistanceForNearby  int32
	MaxDistanceForFar     int32
)

var logger = log.New("LOD")

const (
	configPath               = "configs/shared/lod.json"
	syncIntervalForNearbyKey = "SyncIntervalForNearby"
	syncIntervalForFarKey    = "SyncIntervalForFar"
	maxDistanceForNearbyKey  = "MaxDistanceForNearby"
	maxDistanceForFarKey     = "MaxDistanceForFar"
)

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

	setupLod()
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

	roomID := room.GetRoomID(userData)
	if roomID == "" {
		userData.ServerRespond(derror.ErrData("Not in the room", derror.NotAllowed(0)), ver, cmd, server.Bad, true)
		return
	}

	manager := lodmanager.GetRoomManager(roomID)
	// when room is not setup for lod
	// This if clause is executed only once when the lod broadcast command is received
	if manager == nil {
		manager = lodmanager.NewManager(ver, cmd, SyncIntervalForNearby, SyncIntervalForFar, MaxDistanceForNearby, MaxDistanceForFar)
		lodmanager.SetRoomManager(roomID, manager)
		room.SetOnRoomDiscardByID(roomID, func(roomID string) {
			lodmanager.RemoveRoomManager(roomID)
		})
		room.SetOnLeaveByID(roomID, func(roomID string, userData *user.User) {
			manager.RemoveUserEntity(userData.SID)
		})
	}

	manager.AddUserEntity(userData.SID, proto.X, proto.Y, proto.Payload)
	logger.Sysf("handleRoomBroadcastLoD", "userID", userData.SID, "x", proto.X, "y", proto.Y)
	userData.ServerRespond(nil, ver, cmd, server.Ok, true)
	next(nil)
}

func handleRoomGetLoDInfo(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	// TODO: Get the LoD information from configuration
	lodInfo := proom.NewGetLoDInfoResponse()
	lodInfo.SyncIntervalForNearby = int32(SyncIntervalForNearby)
	lodInfo.SyncIntervalForFar = int32(SyncIntervalForFar)
	lodInfo.MaxDistanceForNearby = int32(MaxDistanceForNearby)
	lodInfo.MaxDistanceForFar = int32(MaxDistanceForFar)

	userData.ServerRespond(lodInfo.Pack(), ver, cmd, server.Ok, true)
	next(nil)
}

func setupLod() {
	loadLodConfigs(configPath)
}

func loadLodConfigs(confPath string) {
	config.Load("Lod", confPath)

	SyncIntervalForNearby = config.GetAsInt32("Lod", syncIntervalForNearbyKey, SyncIntervalForNearby)
	SyncIntervalForFar = config.GetAsInt32("Lod", syncIntervalForFarKey, SyncIntervalForFar)
	if SyncIntervalForFar < SyncIntervalForNearby {
		SyncIntervalForFar = SyncIntervalForNearby
	}
	MaxDistanceForNearby = config.GetAsInt32("Lod", maxDistanceForNearbyKey, MaxDistanceForNearby)
	MaxDistanceForFar = config.GetAsInt32("Lod", maxDistanceForFarKey, MaxDistanceForFar)
	if MaxDistanceForFar < MaxDistanceForNearby {
		MaxDistanceForFar = MaxDistanceForNearby
	}
}
