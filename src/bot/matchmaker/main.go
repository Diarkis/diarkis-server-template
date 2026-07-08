// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Diarkis/diarkis-server-template/bot/utils"
	"github.com/Diarkis/diarkis/client/go/modules/matchmaker"
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	pmm "github.com/Diarkis/diarkis/proto/go/matching"
	"github.com/Diarkis/diarkis/util"
	"github.com/Diarkis/diarkis/uuid/v4"
)

// waitingTime is in seconds
const waitingTime int64 = 30

var botCounter = 0
var searchCnt = 0
var matchedCnt = 0
var completedCnt = 0
var authErrCnt = 0
var host = "127.0.0.1:7000"
var proto = "udp" // udp or tcp
var howmany = 10
var maxmembers uint16 = 10

// sleepTime is in seconds
var sleepTime int64 = 10
var searchProps = make(map[string]int)
var addProps = make(map[string]int)
var hosts int
var interval int64
var profileID = "RankMatch"

type botData struct {
	uid   string
	state int
	udp   *udp.Client
	tcp   *tcp.Client
	mm    *matchmaker.MatchMaker
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Bot requires 4 parameters: {http host:port} {how many bots} {number of host} {search interval in milliseconds}\n")
		os.Exit(1)
		return
	}
	host = os.Args[1]
	howmanySource, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("How many bot parameter given is invalid %v\n", err)
		os.Exit(1)
		return
	}
	howmany = howmanySource
	hosts, err = strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Ratio of hosts parameter given is invalid %v\n", err)
		os.Exit(1)
		return
	}
	intervalSource, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Interval of searches parameter given is invalid %v\n", err)
		os.Exit(1)
		return
	}
	interval = int64(intervalSource)
	fmt.Printf("Starting MatchMaker Bot %v - %v bots [hosts:%v per cent] - hosts:%v guests:%v search interval %vms\n",
		proto, howmany, hosts, float64(howmany)*(float64(hosts)/float64(100)),
		float64(howmany)-(float64(howmany)*(float64(hosts)/float64(100))), interval)
	searchProps["rank"] = 5
	addProps["rank"] = 5
	spawnBots()
	for {
		time.Sleep(time.Second * time.Duration(sleepTime))
		timestamp := util.ZuluTimeFormat(time.Now())
		fmt.Printf("{ \"Time\":\"%v\", \"Bots\":%v, \"Searches\":%v, \"Matches\":%v, \"Completed\":%v, \"AuthError\":%v }\n",
			timestamp, botCounter, searchCnt, matchedCnt, uint16(completedCnt)/maxmembers, authErrCnt)
		searchCnt = 0
		matchedCnt = 0
		completedCnt = 0
		authErrCnt = 0
	}
}

func spawnBots() {
	for i := 0; i < howmany; i++ {
		id, _ := uuid.New()
		go spawnUDPBot(id.String, true)
	}
}

func spawnUDPBot(id string, needToWait bool) {
	if needToWait {
		waitMin := 0
		waitMax := 5000
		time.Sleep(time.Millisecond * time.Duration(int64(util.RandomInt(waitMin, waitMax))))
	}
	eResp, err := utils.Endpoint(host, id, "UDP")
	if err != nil {
		authErrCnt++
		return
	}
	sid, _ := hex.DecodeString(eResp.Sid)
	key, _ := hex.DecodeString(eResp.EncryptionKey)
	iv, _ := hex.DecodeString(eResp.EncryptionIV)
	mkey, _ := hex.DecodeString(eResp.EncryptionMacKey)
	addr := eResp.ServerHost + ":" + fmt.Sprintf("%v", eResp.ServerPort)
	rcvByteSize := 1400
	udpSendInterval := int64(100)
	udp.LogLevel(50)
	cli := udp.New(rcvByteSize, udpSendInterval)

	bot := new(botData)
	bot.uid = id
	bot.state = 0
	bot.udp = cli

	cli.SetEncryptionKeys(sid, key, iv, mkey)

	cli.OnConnect(func() {
		bot.mm = matchmaker.NewMatchMakerAsUDP(cli)
		bot.mm.OnHostMatchmakingResponse(func(status bool, roomID string) {
			if status {
				bot.state = 22
			} else {
				bot.state = 23
			}
		})
		bot.mm.OnResultResponse(func(status bool, results []*pmm.SearchResult) {
			searchCnt++
			if !status {
				bot.state += 1
			} else {
				matchedCnt++
				bot.state = 22
			}
		})
		// A proxy for completion
		bot.mm.OnClaimPush(func(*pmm.ClaimMatchmakingPush) {
			bot.state = 23
			completedCnt++
		})
		go startBot(bot)
	})
	cli.OnDisconnect(func() {
		botCounter--
		if botCounter >= howmany {
			return
		}
		spawnUDPBot(bot.uid, true)

	})

	cli.Connect(addr)
}

func startBot(bot *botData) {
	botCounter++
	// 0 ~ 20 = search
	// 21     = add
	// 22     = wait
	// 23     = room full and disconnect
	currentState := -1
	if util.RandomInt(0, 99) < hosts {
		bot.state = 21
	}
	waitCounter := int64(0)
	for {
		if currentState < 22 && currentState == bot.state {
			time.Sleep(time.Millisecond * time.Duration(interval))
			continue
		}
		switch bot.state {
		case -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 18, 19, 20:
			search(bot)
		case 21:
			add(bot)
		case 22:
			// We are waiting
			waitCounter += interval
			if waitCounter >= waitingTime*1000 {
				bot.state = 23
			}
		case 23:
			disconnect(bot)
			return
		default:
			bot.state = 23
			disconnect(bot)
			return
		}
		currentState = bot.state
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}

func search(bot *botData) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		return
	}
	if bot.mm != nil {
		bot.mm.SearchList([]string{profileID}, nil, searchProps, 10, false, false)
	}
}

func add(bot *botData) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		return
	}
	if bot.mm != nil {
		bot.mm.HostMatchmaking(maxmembers, 60, profileID, nil, addProps, []byte("metadata"), false, "")
	}
}

func disconnect(bot *botData) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		return
	}
	switch proto {
	case "udp":
		if bot.udp != nil {
			bot.udp.Disconnect()
		}
	case "tcp":
		if bot.tcp != nil {
			bot.tcp.Disconnect()
		}
	}
}
