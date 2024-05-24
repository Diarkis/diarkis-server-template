package app

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

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
var botCounter = 0
var joinedCnt atomic.Int64
var broadcastSendCnt atomic.Int64
var broadcastReceiveCnt atomic.Int64

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
}

func (b *bot) isJoined() bool {
	return b.room.ID != ""
}

var bm botManager

func Run() {
	parseArgs()
	logger.Info("bot args",
		"host", host,
		"bots", bots,
		"authInterval", authInterval,
		"roomSize", roomSize,
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
	for i := 0; i < bots; i++ {
		botUuid, _ := uuid.New()
		go spawnBot(botUuid.String)
		time.Sleep(time.Millisecond * time.Duration(authInterval))
	}
}

func newBot(id string) *bot {
	eResp, err := endpoint(host, id, protocol)
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

func spawnBot(id string) {
	bot := newBot(id)
	bm.bots = append(bm.bots, bot)
}

func broadcast(bot *bot) {
	message := make([]byte, packetSize, packetSize)
	bot.room.BroadcastTo(bot.room.ID, message, false)
	bot.broadcastSendCnt.Add(1)
	broadcastSendCnt.Add(1)
}

func searchAndJoin(bot *bot) {
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
	joinMessage := []byte("")
	roomCli.JoinRandom(roomSize, 60, joinMessage, 0)
	bot.room = roomCli

	roomCli.OnJoin(func(success bool, createdTime uint) {
		joinedCnt.Add(1)
		logger.Debug("OnJoin",
			"bot.uid", bot.uid,
			"success", success,
			"createdTime", createdTime)
		if success {
			bot.state = STATUS_BROADCAST
		}
	})
	roomCli.OnCreate(func(success bool, name string, createdTime uint) {
		logger.Info("OnCreate",
			"bot.uid", bot.uid,
			"success", success,
			"createdTime", createdTime)
		if success {
			bot.state = STATUS_BROADCAST
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
