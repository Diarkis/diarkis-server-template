// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Diarkis/diarkis/client/go/udp"

	//	"github.com/Diarkis/diarkis/client/go/tcp"

	"github.com/Diarkis/diarkis-server-template/bot/utils"
	"github.com/Diarkis/diarkis/client/go/modules/matchmaker"
	"github.com/Diarkis/diarkis/util"
	v4 "github.com/Diarkis/diarkis/uuid/v4"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var running = true

// number of tickets issued in 60 seconds
var ticketCnt = 0

// number of tickets successlly completed in 60 seconds
var ticketSuccessCnt = 0

// HTTP endpoint
var host = "127.0.0.1:7000"

// UDP or TCP
var serverType string = ""

// total number of bots
var totalBots = 0
var botCnt = 0
var timedoutCnt = 0
var clientTimeout int64 = 90 // 90s for a bot to timeout and die no matter what
// interval in seconds to add bots.
// Bots will be added upto totalBots per interval.
var interval int64 = 100

// log map to create a JSON data when process is terminated
var logmap = make([]map[string]int, 0)

// total count of tickets issued
var ticketTotalCnt = 0

// total count of successful tickets
var successTicketTotalCnt = 0

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Bot requires 4 parameters: {http host and port xxxx:0000} {server type UDP or TCP} {how many bots} {interval in seconds}\n")
		os.Exit(1)
		return
	}
	setupSignalHandler()
	host = os.Args[1]
	serverType = os.Args[2]
	totalBotsSource, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Number of bots must be a valid integer %v\n", err)
		os.Exit(1)
		return
	}
	totalBots = totalBotsSource
	intervalSource, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Interval must be a valid integer %v\n", err)
		os.Exit(1)
		return
	}
	interval = int64(intervalSource)
	go spawnBots()
	start := util.NowSeconds()
	cnt := 0
	for running {
		time.Sleep(time.Second * 5)
		cnt++
		if cnt < 6 {
			continue
		}
		cnt = 0
		elapsed := util.NowSeconds() - start
		ticketTotalCnt += ticketCnt
		successTicketTotalCnt += ticketSuccessCnt
		fmt.Printf("=============\n")
		fmt.Printf("Time elapsed       %v seconds\n", elapsed)
		fmt.Printf("Number of timeouts %v\n", timedoutCnt)
		fmt.Printf("Number of tickets  %v - Total:%v\n", ticketCnt, ticketTotalCnt)
		fmt.Printf("Successful tickets %v - Total:%v\n", ticketSuccessCnt, successTicketTotalCnt)
		fmt.Printf("Success Rate       %v percent - Total:%v percent\n",
			int(float64(ticketSuccessCnt)/float64(ticketCnt)*float64(100)),
			int(float64(successTicketTotalCnt)/float64(ticketTotalCnt)*float64(100)))
		fmt.Printf("=============\n")
		data := make(map[string]int)
		data["Time"] = int(elapsed)
		data["Tickets"] = ticketCnt
		data["Success"] = ticketSuccessCnt
		data["Rate"] = int(float64(ticketSuccessCnt) / float64(ticketCnt) * float64(100))
		logmap = append(logmap, data)
		timedoutCnt = 0
		ticketCnt = 0
		ticketSuccessCnt = 0
	}
	fmt.Printf("Bot terminated\n")
	encoded, _ := json.Marshal(logmap)
	fmt.Printf("===== JSON Report =====\n")
	fmt.Printf("%s\n", string(encoded))
	fmt.Printf("=======================\n")
	fmt.Printf("====== CSV Report =====\n")
	fmt.Printf("%s\n", convertToCSVFromLogData())
	fmt.Printf("=======================\n")
}

func setupSignalHandler() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGQUIT)
	go handleSignal(ch)
}

func handleSignal(ch chan os.Signal) {
	sig := <-ch
	fmt.Printf("Signal captured %v => terminate bot\n", sig)
	running = false
}

func spawnBots() {
	for true {
		fmt.Printf("interval: %v, totalbots: %v, botCnt: %v\n", interval, totalBots, botCnt)
		if totalBots <= botCnt {
			time.Sleep(time.Second * time.Duration(interval))
			continue
		}
		fmt.Printf("Bots to spawn %v\n", totalBots-botCnt)
		for i := 0; i < totalBots-botCnt; i++ {
			uuid, _ := v4.New()
			spawnBot(uuid.String)
			time.Sleep(time.Millisecond * 100)
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func spawnBot(id string) {
	eResp, err := utils.Endpoint(host, id, serverType)
	sid, _ := hex.DecodeString(eResp.Sid)
	key, _ := hex.DecodeString(eResp.EncryptionKey)
	iv, _ := hex.DecodeString(eResp.EncryptionIV)
	mkey, _ := hex.DecodeString(eResp.EncryptionMacKey)

	addr := eResp.ServerHost + ":" + fmt.Sprintf("%v", eResp.ServerPort)

	if err != nil {
		fmt.Printf("Endpoint error => %v\n", err)
		return
	}
	dead := false
	rcvByteSize := 1400
	udpSendInterval := int64(100)
	udp.LogLevel(50)
	cli := udp.New(rcvByteSize, udpSendInterval)
	cli.SetEncryptionKeys(sid, key, iv, mkey)
	cli.OnDisconnect(func() {
		botCnt--
		dead = true
	})
	cli.OnConnect(func() {
		botCnt++
		mm := matchmaker.NewMatchMakerAsUDP(cli)
		mm.OnTicketComplete(func(success bool, data []byte) {
			if success {
				ticketSuccessCnt++
			}
			cli.Disconnect()
		})
		ticketCnt++
		mm.IssueTicket(0)
	})
	go func() {
		time.Sleep(time.Second * time.Duration(clientTimeout))
		if dead {
			return
		}
		botCnt--
		timedoutCnt++
		cli.Disconnect()
	}()
	cli.Connect(addr)
}

func convertToCSVFromLogData() string {
	if len(logmap) == 0 {
		return ""
	}
	csv := ""
	names := make([]string, 0)
	list := make([]string, 0)
	for k := range logmap[0] {
		names = append(names, k)
		list = append(list, fmt.Sprintf("\"%s\"", k))
	}
	csv = strings.Join(list, ",") + "\n"
	for i := 0; i < len(logmap); i++ {
		row := logmap[i]
		list = make([]string, 0)
		for _, name := range names {
			list = append(list, fmt.Sprintf("%v", row[name]))
		}
		csv += strings.Join(list, ",")
		if i < len(logmap)-1 {
			csv += "\n"
		}
	}
	return csv
}
