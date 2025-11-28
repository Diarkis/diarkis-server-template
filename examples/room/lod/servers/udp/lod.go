package main

import (
	"github.com/Diarkis/diarkis/log"

	proom "github.com/Diarkis/diarkis-server-template/examples/room/lod/puffer/go/room"
	"github.com/Diarkis/diarkis/derror"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

var (
	logger                      = log.New("LOD")
	SyncIntervalForNearby int32 = 16
	SyncIntervalForFar    int32 = 2000
	MaxDistanceForNearby  int32 = 10000
	MaxDistanceForFar     int32 = 40000
)

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

	lodManager := getRoomLodManager(roomID)
	if lodManager == nil {
		lodManager = newLodManager(ver, cmd, SyncIntervalForNearby, SyncIntervalForFar, MaxDistanceForNearby, MaxDistanceForFar)
		setRoomLodManager(roomID, lodManager)
	}
	lodManager.addUserEntity(userData.SID, proto.X, proto.Y, proto.Payload)
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
