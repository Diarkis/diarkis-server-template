// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/rand"
	"sync/atomic"
	"time"

	proom "github.com/Diarkis/diarkis-server-template/examples/room/lod/puffer/go/room"
	"github.com/Diarkis/diarkis/client/go/modules/room"
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/uuid/v4"
)

const (
	UDP_STRING               string = "udp"
	TCP_STRING               string = "tcp"
	STATUS_BROADCAST                = iota
	RCV_BYTE_SIZE                   = 1400
	DIARKIS_CLIENT_LOG_LEVEL        = 70
)

var (
	joinedCnt           atomic.Int64
	broadcastSendCnt    atomic.Int64
	broadcastReceiveCnt atomic.Int64
	logger              *slog.Logger
)

type bot struct {
	uid              string
	state            int
	udp              *udp.Client
	tcp              *tcp.Client
	room             *room.Room
	broadcastSendCnt atomic.Int64
	broadcastRcvCnt  atomic.Int64
	x                int32 // position x for LOD
	y                int32 // position y for LOD
	isCreator        bool  // whether this bot creates a room
	targetRoomIndex  int   // index of the room this bot should join/create

	settings        *Settings
	roomCoordinator *RoomCoordinator
}

func (b *bot) isJoined() bool {
	return b.room.ID != ""
}

var bm botManager

func Run() {
	// Load configuration
	settings := LoadSettings()
	logger = setupLogger(settings)

	roomCoordinator := NewRoomCoordinator()

	spawnBots(settings, roomCoordinator)
	for { // every second print metrics
		time.Sleep(time.Second)
		printMetrics()
	}
}

func printMetrics() {
	for _, bot := range bm.bots {
		logger.Debug("bot status",
			"bot.uid", bot.uid,
			"bot.state", bot.state,
			"bot.room", bot.room,
			"bot.broadcastSendCnt", bot.broadcastSendCnt.Load(),
			"bot.broadcastRcvCnt", bot.broadcastRcvCnt.Load())
	}
	logger.Info("-----------------------------------------------------")
	logger.Info("bot metrics",
		"joinedBotCnt", bm.getJoinedCnt(),
		"rcvCntTotal", bm.getBroadcastReceiveCntTotal(),
		"rcvAvg", bm.getBroadcastReceiveAvg(),
		"sendCntTotal", bm.getBroadcastSendCntTotal(),
		"sendAvg", bm.getBroadcastSendAvg(),
		"joinedRoom", bm.getJoinedRooms(),
	)
	logger.Info("-----------------------------------------------------")
	bm.resetCnt()
}

func spawnBots(settings *Settings, rc *RoomCoordinator) {
	// Calculate how many rooms are needed
	// If averageRoomMember is set (>0), use it to calculate numRooms.
	// Otherwise, default to filling rooms to capacity (roomSize).
	memberPerRoom := settings.RoomSize
	if settings.AverageRoomMember > 0 {
		memberPerRoom = settings.AverageRoomMember
	}
	numRooms := (settings.BotsCount + memberPerRoom - 1) / memberPerRoom

	logger.Info("Room creation plan",
		"totalBots", settings.BotsCount,
		"roomSize", settings.RoomSize,
		"averageRoomMember", settings.AverageRoomMember,
		"targetMemberPerRoom", memberPerRoom,
		"numRoomsToCreate", numRooms)

	for i := 0; i < settings.BotsCount; i++ {
		botUuid, _ := uuid.New()
		// Determine which room this bot belongs to
		targetRoomIdx := i % numRooms
		// The first bot assigned to a room index becomes the creator for that room
		// This simple logic works because we iterate i=0..bots.
		// i < numRooms logic was: 0..numRooms-1 are creators.
		// New logic: We need exactly ONE creator per targetRoomIdx.
		// We can assign creators where i < numRooms.
		// e.g. bots=10, numRooms=2.
		// i=0 -> idx=0. Creator? Yes (0 < 2)
		// i=1 -> idx=1. Creator? Yes (1 < 2)
		// i=2 -> idx=0. Creator? No.
		// ...
		isCreator := i < numRooms

		go spawnBot(botUuid.String, isCreator, targetRoomIdx, settings, rc)
		time.Sleep(time.Millisecond * time.Duration(settings.AuthInterval))
	}
}

func newBot(id string, isCreator bool, targetRoomIndex int, settings *Settings, rc *RoomCoordinator) *bot {
	eResp, err := endpoint(settings.Host, id, settings.Protocol)
	if err != nil {
		logger.Error("Auth error",
			"bot.uid", id,
			"err", err)
		return nil
	}
	logger.Debug("eResponse",
		"eResp.Sid", eResp.Sid,
		"eResp.EncryptionKey", eResp.EncryptionKey,
		"eResp.EncryptionIV", eResp.EncryptionIV,
		"eResp.EncryptionMacKey", eResp.EncryptionMacKey)
	sid, _ := hex.DecodeString(eResp.Sid)
	key, _ := hex.DecodeString(eResp.EncryptionKey)
	iv, _ := hex.DecodeString(eResp.EncryptionIV)
	macKey, _ := hex.DecodeString(eResp.EncryptionMacKey)

	rcvByteSize := RCV_BYTE_SIZE
	udpSendInterval := int64(settings.UDPClientSendInterval)
	udp.LogLevel(DIARKIS_CLIENT_LOG_LEVEL)
	cli := udp.New(rcvByteSize, udpSendInterval)
	bot := &bot{
		uid:             id,
		udp:             cli,
		isCreator:       isCreator,
		targetRoomIndex: targetRoomIndex,
		settings:        settings,
		roomCoordinator: rc,
		x:               randomInt32(settings.MapMinX, settings.MapMaxX),
		y:               randomInt32(settings.MapMinY, settings.MapMaxY),
	}
	cli.SetEncryptionKeys(sid, key, iv, macKey)
	cli.OnResponse(func(ver uint8, cmd uint16, status uint8, payload []byte) {
		bot.handleOnResponse(ver, cmd, status, payload)
	})
	cli.OnPush(func(ver uint8, cmd uint16, payload []byte) {
		bot.handleOnPush(ver, cmd, payload)
	})
	cli.OnConnect(func() {
		bot.handleOnConnect()
	})
	cli.OnDisconnect(func() {

	})
	addr := eResp.ServerHost + ":" + fmt.Sprintf("%v", eResp.ServerPort)
	cli.Connect(addr)
	bot.room = new(room.Room)
	bot.room.SetupAsUDP(bot.udp)

	return bot
}

func spawnBot(id string, isCreator bool, targetRoomIndex int, settings *Settings, rc *RoomCoordinator) {
	bot := newBot(id, isCreator, targetRoomIndex, settings, rc)
	if bot != nil {
		bm.bots = append(bm.bots, bot)
	}
}

func (b *bot) LodBroadcast() {
	message := make([]byte, b.settings.PacketSize)
	// Use LOD broadcast instead of regular room broadcast
	proto := proom.NewBroadcastLoD()
	proto.X = b.x
	proto.Y = b.y
	proto.Payload = message
	b.udp.RSend(proto.Ver(), proto.Cmd(), proto.Pack())
	b.broadcastSendCnt.Add(1)
	broadcastSendCnt.Add(1)
}

func (b *bot) CreateRoom() {
	if b.state == 0 && b.udp == nil && b.tcp == nil {
		logger.Error("bot is not connected to any server")
		return
	}

	switch b.settings.Protocol {
	case UDP_STRING:
		b.room.SetupAsUDP(b.udp)
	case TCP_STRING:
		b.room.SetupAsTCP(b.tcp)
	}

	// Use Create to make a room.
	b.room.Create(uint16(b.settings.RoomSize), false, true, 60, 0)

	b.room.OnCreate(func(success bool, roomID string, createdTime uint) {
		if success {
			joinedCnt.Add(1)
			b.state = STATUS_BROADCAST
			// Store roomID in map with index
			b.roomCoordinator.AddRoom(b.targetRoomIndex, roomID)

			logger.Info("Room created",
				"bot.uid", b.uid,
				"roomID", roomID,
				"roomIndex", b.targetRoomIndex,
				"createdTime", createdTime,
				"maxMembers", b.settings.RoomSize)
		} else {
			logger.Error("OnCreate failed",
				"bot.uid", b.uid)
		}
	})

	b.room.OnJoin(func(success bool, createdTime uint) {
		if success {
			joinedCnt.Add(1)
			b.state = STATUS_BROADCAST
			logger.Info("Creator joined existing room",
				"bot.uid", b.uid,
				"roomID", b.room.ID,
				"createdTime", createdTime)
		}
	})

	b.setupRoomCallbacks()
}

func (b *bot) JoinRoom() {
	if b.state == 0 && b.udp == nil && b.tcp == nil {
		logger.Error("bot is not connected to any server")
		return
	}

	switch b.settings.Protocol {
	case UDP_STRING:
		b.room.SetupAsUDP(b.udp)
	case TCP_STRING:
		b.room.SetupAsTCP(b.tcp)
	}

	var targetRoomID string
	// Loop until the target room is created
	for {
		if val, ok := b.roomCoordinator.GetRoom(b.targetRoomIndex); ok {
			targetRoomID = val
			break
		}
		time.Sleep(time.Second * 1)
		logger.Info("Waiting for room to be created",
			"bot.uid", b.uid,
			"targetRoomIndex", b.targetRoomIndex)
	}

	logger.Info("Joining room",
		"bot.uid", b.uid,
		"roomID", targetRoomID,
		"targetRoomIndex", b.targetRoomIndex)
	b.room.Join(targetRoomID, []byte(""))

	b.room.OnJoin(func(success bool, createdTime uint) {
		if success {
			joinedCnt.Add(1)
			b.state = STATUS_BROADCAST
			logger.Info("Joiner joined room",
				"bot.uid", b.uid,
				"roomID", b.room.ID,
				"createdTime", createdTime)
		} else {
			logger.Warn("OnJoin failed, retrying...",
				"bot.uid", b.uid)
			// Retry after delay
			time.Sleep(time.Millisecond * 500)
			b.JoinRoom()
		}
	})

	b.room.OnCreate(func(success bool, name string, createdTime uint) {
		// Joiners shouldn't create rooms, but if they do, log it
		if success {
			joinedCnt.Add(1)
			b.state = STATUS_BROADCAST
			logger.Warn("Joiner unexpectedly created room",
				"bot.uid", b.uid,
				"roomID", name)
		}
	})

	b.setupRoomCallbacks()
}

func (b *bot) setupRoomCallbacks() {
	b.room.OnMemberLeave(func(message []byte) {
		logger.Debug("OnMemberLeave",
			"bot.uid", b.uid,
			"message", message)
	})
	b.room.OnMemberBroadcast(func(bytes []byte) {
		broadcastReceiveCnt.Add(1)
		b.broadcastRcvCnt.Add(1)
	})
}

// Handlers moved from callbacks.go

func (b *bot) handleOnConnect() {
	// Start movement loop
	go func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(b.settings.MovementInterval))
			if b.isJoined() {
				b.move()
			}
		}
	}()
	// Start broadcast loop
	go func() {
		if b.isCreator {
			logger.Info("Bot creating room", "bot.uid", b.uid)
			b.CreateRoom()
		} else {
			time.Sleep(time.Second * 2)
			logger.Info("Bot joining room", "bot.uid", b.uid)
			b.JoinRoom()
		}

		for {
			time.Sleep(time.Millisecond * time.Duration(b.settings.PacketInterval))
			if b.isJoined() {
				b.LodBroadcast()
			}
		}
	}()
}

func (b *bot) handleOnResponse(ver uint8, cmd uint16, status uint8, payload []byte) {
}

func (b *bot) handleOnPush(ver uint8, cmd uint16, payload []byte) {
	if ver == proom.BroadcastLoDPushVer && cmd == proom.BroadcastLoDPushCmd {
		broadcastReceiveCnt.Add(1)
		b.broadcastRcvCnt.Add(1)
	}
}

// from movement.go
func (b *bot) move() {
	// Random movement in 8 directions (N, NE, E, SE, S, SW, W, NW)
	dx := int32(rand.Intn(3) - 1) // -1, 0, or 1
	dy := int32(rand.Intn(3) - 1) // -1, 0, or 1

	// Apply movement speed
	newX := b.x + (dx * b.settings.MovementSpeed)
	newY := b.y + (dy * b.settings.MovementSpeed)

	// Clamp to map boundaries
	if newX < b.settings.MapMinX {
		newX = b.settings.MapMinX
	} else if newX > b.settings.MapMaxX {
		newX = b.settings.MapMaxX
	}

	if newY < b.settings.MapMinY {
		newY = b.settings.MapMinY
	} else if newY > b.settings.MapMaxY {
		newY = b.settings.MapMaxY
	}

	// Update bot position
	b.x = newX
	b.y = newY
}
