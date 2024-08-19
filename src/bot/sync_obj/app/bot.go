// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"encoding/hex"
	"fmt"
	"github.com/Diarkis/diarkis/client/go/modules/room"
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/uuid/v4"
	"log/slog"
	"math/rand/v2"
	"sync/atomic"
	"time"
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
var joinedCnt atomic.Int64
var syncSendCnt atomic.Int64
var syncRcvCnt atomic.Int64

// sleepTime is in seconds
var sleepTime int64 = 1

type bot struct {
	uid         string
	state       int
	udp         *udp.Client
	tcp         *tcp.Client
	room        *room.Room
	syncSendCnt atomic.Int64
	syncRcvCnt  atomic.Int64
}

func (b *bot) isJoined() bool {
	return b.room.ID != ""
}

var bm botManager

func Run() {
	parseArgs()
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
			"bot.syncSendCnt", bot.syncSendCnt.Load(),
			"bot.syncRcvCnt", bot.syncRcvCnt.Load())
	}
	logger.Info("-----------------------------------------------------")
	logger.Info("bot metrics",
		"joinedBotCnt", bm.getJoinedCnt(),
		"rcvCntTotal", bm.getSyncReceiveCntTotal(),
		"rcvAvg", bm.getSyncReceiveAvg(),
		"sendCntTotal", bm.getSyncReceiveCntTotal(),
		"sendAvg", bm.getSyncSendAvg(),
		"joinedRoom", bm.getJoinedRooms(),
	)
	logger.Info("-----------------------------------------------------")
	bm.resetCnt()
}

func spawnBots() {
	for i := 0; i < *bots; i++ {
		botUuid, _ := uuid.New()
		go spawnBot(botUuid.String)
		time.Sleep(time.Millisecond * time.Duration(*authInterval))
	}
}

func newBot(id string) *bot {
	eResp, err := endpoint(*host, id, UDP_STRING)
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
	udpSendInterval := int64(100)
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
	bot.room.OnObjectUpdateResponse(func(msg []byte) {
		handleOnObjectUpdateResponse(bot, msg)
	})

	bot.room.OnObjectUpdatePush(func(status uint8, msg string, obj map[string]float64) {
		handleOnObjectUpdatePush(bot, status, msg, obj)
	})
	return bot
}

func spawnBot(id string) {
	bot := newBot(id)
	bm.bots = append(bm.bots, bot)
}

func updateObj(bot *bot) {
	obj := make(map[string]float64)
	if rand.IntN(100) < *changePercent {
		obj["test"] = 10
	} else {
		obj["test"] = 11
	}
	for i := 0; i < 10; i++ {
		obj[fmt.Sprintf("test%d", i)] = float64(i)
	}
	if *mode == "incr" {
		bot.room.UpdateObject(room.UpdateObjectIncrMode, "test", obj, false)
	} else if *mode == "set" {
		bot.room.UpdateObject(room.UpdateObjectSetMode, "test", obj, false)
	}
	//bot.room.UpdateObject(room.UpdateObjectSetMode, "test", obj, false)
	//bot.room.UpdateObject(room.UpdateObjectIncrMode, "test", obj, false)
}

func searchAndJoin(bot *bot) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		logger.Error("bot is not connected to any server")
		return
	}
	roomCli := new(room.Room)
	roomCli.SetupAsUDP(bot.udp)
	joinMessage := []byte("")
	roomCli.JoinRandom(*roomSize, 60, joinMessage, 0)
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
		bot.syncRcvCnt.Add(1)
	})
}
