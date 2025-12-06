// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Diarkis/diarkis-server-template/examples/room/lod/bot/utils"
	proom "github.com/Diarkis/diarkis-server-template/examples/room/lod/puffer/go/room"
	"github.com/Diarkis/diarkis/client/go/modules/room"
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/uuid/v4"
)

const UDP_STRING string = "udp"
const TCP_STRING string = "tcp"

const (
	STATUS_BROADCAST = iota
)

// udp client settings
const RcvByteSize = 1400
const DiarkisClientLogLevel = 70

var logger *slog.Logger

// metrics
// metrics
var botCounter = 0
var joinedCnt atomic.Int64
var broadcastSendCnt atomic.Int64
var broadcastReceiveCnt atomic.Int64
var createdRoomMap sync.Map

// sleepTime is in seconds
var sleepTime int64 = 1

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
}

func (b *bot) isJoined() bool {
	return b.room.ID != ""
}

var bm botManager

func Run() {
	// Setup logger
	programLevel := new(slog.LevelVar)
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))

	// Load configuration
	loadBotConfig()
	loadBotLodConfig()

	// Set log level
	switch logLevel {
	case "debug":
		programLevel.Set(slog.LevelDebug)
	case "info":
		programLevel.Set(slog.LevelInfo)
	case "warn":
		programLevel.Set(slog.LevelWarn)
	case "error":
		programLevel.Set(slog.LevelError)
	default:
		programLevel.Set(slog.LevelInfo)
	}

	logger.Info("bot args",
		"host", host,
		"bots", bots,
		"authInterval", authInterval,
		"roomSize", roomSize,
		"averageRoomMember", averageRoomMember,
		"packetInterval", packetInterval,
		"packetSize", packetSize,
		"logLevel", logLevel,
		"protocol", protocol,
	)
	spawnBots()
	for {
		time.Sleep(time.Second * time.Duration(sleepTime))
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

func spawnBots() {
	// Calculate how many rooms are needed
	// If averageRoomMember is set (>0), use it to calculate numRooms.
	// Otherwise, default to filling rooms to capacity (roomSize).
	memberPerRoom := roomSize
	if averageRoomMember > 0 {
		memberPerRoom = averageRoomMember
	}
	numRooms := (bots + memberPerRoom - 1) / memberPerRoom

	logger.Info("Room creation plan",
		"totalBots", bots,
		"roomSize", roomSize,
		"averageRoomMember", averageRoomMember,
		"targetMemberPerRoom", memberPerRoom,
		"numRoomsToCreate", numRooms)

	for i := 0; i < bots; i++ {
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

		go spawnBot(botUuid.String, isCreator, targetRoomIdx)
		time.Sleep(time.Millisecond * time.Duration(authInterval))
	}
}

func newBot(id string, isCreator bool, targetRoomIndex int) *bot {
	eResp, err := utils.Endpoint(host, id, protocol)
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

	rcvByteSize := RcvByteSize
	udpSendInterval := int64(udpClientSendInterval)
	udp.LogLevel(DiarkisClientLogLevel)
	cli := udp.New(rcvByteSize, udpSendInterval)
	bot := new(bot)
	bot.uid = id
	bot.state = 0
	bot.udp = cli
	bot.isCreator = isCreator
	bot.targetRoomIndex = targetRoomIndex
	// Initialize bot position randomly within map boundaries
	bot.x = utils.RandomInt32(MapMinX, MapMaxX)
	bot.y = utils.RandomInt32(MapMinY, MapMaxY)
	cli.SetEncryptionKeys(sid, key, iv, macKey)
	cli.OnResponse(func(ver uint8, cmd uint16, status uint8, payload []byte) {
		handleOnResponse(bot, ver, cmd, status, payload)
	})
	cli.OnPush(func(ver uint8, cmd uint16, payload []byte) {
		handleOnPush(bot, ver, cmd, payload)
	})
	cli.OnConnect(func() {
		handleOnConnect(bot)
	})
	cli.OnDisconnect(func() {
		handleOnDisconnect()
	})
	addr := eResp.ServerHost + ":" + fmt.Sprintf("%v", eResp.ServerPort)
	cli.Connect(addr)
	bot.room = new(room.Room)
	bot.room.SetupAsUDP(bot.udp)

	return bot
}

func spawnBot(id string, isCreator bool, targetRoomIndex int) {
	bot := newBot(id, isCreator, targetRoomIndex)
	bm.bots = append(bm.bots, bot)
}

func broadcast(bot *bot) {
	message := make([]byte, packetSize)
	// Use LOD broadcast instead of regular room broadcast
	proto := proom.NewBroadcastLoD()
	proto.X = bot.x
	proto.Y = bot.y
	proto.Payload = message
	bot.udp.RSend(proto.Ver, proto.Cmd, proto.Pack())
	bot.broadcastSendCnt.Add(1)
	broadcastSendCnt.Add(1)
}

func createRoom(bot *bot) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		logger.Error("bot is not connected to any server")
		return
	}
	roomCli := new(room.Room)
	switch protocol {
	case UDP_STRING:
		roomCli.SetupAsUDP(bot.udp)
	case TCP_STRING:
		roomCli.SetupAsTCP(bot.tcp)
	}

	// Use Create to make a room.
	roomCli.Create(uint16(roomSize), false, true, 60, 0)
	bot.room = roomCli

	roomCli.OnCreate(func(success bool, roomID string, createdTime uint) {
		if success {
			joinedCnt.Add(1)
			bot.state = STATUS_BROADCAST
			// Store roomID in map with index
			createdRoomMap.Store(bot.targetRoomIndex, roomID)

			logger.Info("Room created",
				"bot.uid", bot.uid,
				"roomID", roomID,
				"roomIndex", bot.targetRoomIndex,
				"createdTime", createdTime,
				"maxMembers", roomSize)
		} else {
			logger.Error("OnCreate failed",
				"bot.uid", bot.uid)
		}
	})

	roomCli.OnJoin(func(success bool, createdTime uint) {
		if success {
			joinedCnt.Add(1)
			bot.state = STATUS_BROADCAST
			logger.Info("Creator joined existing room",
				"bot.uid", bot.uid,
				"roomID", bot.room.ID,
				"createdTime", createdTime)
		}
	})

	roomCli.OnMemberLeave(func(message []byte) {
		logger.Debug("OnMemberLeave",
			"bot.uid", bot.uid,
			"message", message)
	})
	roomCli.OnMemberBroadcast(func(bytes []byte) {
		broadcastReceiveCnt.Add(1)
		bot.broadcastRcvCnt.Add(1)
	})
}

func joinRoom(bot *bot) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		logger.Error("bot is not connected to any server")
		return
	}
	roomCli := new(room.Room)
	switch protocol {
	case UDP_STRING:
		roomCli.SetupAsUDP(bot.udp)
	case TCP_STRING:
		roomCli.SetupAsTCP(bot.tcp)
	}

	var targetRoomID string
	// Loop until the target room is created
	for {
		if val, ok := createdRoomMap.Load(bot.targetRoomIndex); ok {
			targetRoomID = val.(string)
			break
		}
		time.Sleep(time.Second * 1)
		logger.Info("Waiting for room to be created",
			"bot.uid", bot.uid,
			"targetRoomIndex", bot.targetRoomIndex)
	}

	logger.Info("Joining room",
		"bot.uid", bot.uid,
		"roomID", targetRoomID,
		"targetRoomIndex", bot.targetRoomIndex)
	roomCli.Join(targetRoomID, []byte(""))
	bot.room = roomCli

	roomCli.OnJoin(func(success bool, createdTime uint) {
		if success {
			joinedCnt.Add(1)
			bot.state = STATUS_BROADCAST
			logger.Info("Joiner joined room",
				"bot.uid", bot.uid,
				"roomID", bot.room.ID,
				"createdTime", createdTime)
		} else {
			logger.Warn("OnJoin failed, retrying...",
				"bot.uid", bot.uid)
			// Retry after delay
			time.Sleep(time.Millisecond * 500)
			joinRoom(bot)
		}
	})

	roomCli.OnCreate(func(success bool, name string, createdTime uint) {
		// Joiners shouldn't create rooms, but if they do, log it
		if success {
			joinedCnt.Add(1)
			bot.state = STATUS_BROADCAST
			logger.Warn("Joiner unexpectedly created room",
				"bot.uid", bot.uid,
				"roomID", name)
		}
	})

	roomCli.OnMemberLeave(func(message []byte) {
		logger.Debug("OnMemberLeave",
			"bot.uid", bot.uid,
			"message", message)
	})
	roomCli.OnMemberBroadcast(func(bytes []byte) {
		broadcastReceiveCnt.Add(1)
		bot.broadcastRcvCnt.Add(1)
	})
}
