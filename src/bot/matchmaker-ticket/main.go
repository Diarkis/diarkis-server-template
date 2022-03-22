package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Diarkis/diarkis/client/go/udp"
//	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/modules/matchmaker"
	"github.com/Diarkis/diarkis/util"
	"io/ioutil"
	"net/http"
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
// interval in seconds to add bots.
// Bots will be added upto totalBots per interval.
var interval int64 = 100
// log map to create a JSON data when process is terminated
var logmap = make([]map[string]int, 0)

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
		if cnt < 12 {
			continue
		}
		cnt = 0
		elapsed := util.NowSeconds() - start
		fmt.Printf("=============\n")
		fmt.Printf("Time elapsed       %v seconds\n", elapsed)
		fmt.Printf("Number of tickets  %v\n", ticketCnt)
		fmt.Printf("Successful tickets %v\n", ticketSuccessCnt)
		fmt.Printf("Success Rate       %v percent\n", int(float64(ticketSuccessCnt)/float64(ticketCnt)*float64(100)))
		fmt.Printf("=============\n")
		data := make(map[string]int)
		data["Time"] = int(elapsed)
		data["Tickets"] = ticketCnt
		data["Success"] = ticketSuccessCnt
		data["Rate"] = int(float64(ticketSuccessCnt)/float64(ticketCnt)*float64(100))
		logmap = append(logmap, data)
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
	id := 0
	for true {
		for i := 0; i < totalBots; i++ {
			id++
			go spawnBot(id)
			time.Sleep(time.Millisecond * 100)
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func spawnBot(id int) {
	addr, sid, key, iv, mkey, err := endpoint(id)
	if err != nil {
		fmt.Printf("Endpoint error => %v\n", err)
		return
	}
	rcvByteSize := 1400
	udpSendInterval := int64(100)
	udp.LogLevel(9)
	cli := udp.New(rcvByteSize, udpSendInterval)
	cli.SetEncryptionKeys(sid, key, iv, mkey)
	cli.OnDisconnect(func() {
		botCnt--
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
		mm.IssueTicket()
	})
	cli.Connect(addr)
}

// returns addr, sid, key, iv, mackey, error
func endpoint(uid int) (string, []byte, []byte, []byte, []byte, error) {
	url := fmt.Sprintf("http://%s/endpoint/type/%s/user/%v", host, serverType, uid)
	// fmt.Printf("Bot ID: %v - Connecting to HTTP: %s\n", uid, url)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	//req.Header.Add("ClientKey", clientKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	if resp.StatusCode > 300 {
		err := errors.New(fmt.Sprintf("Error response status %v - body:%v", resp.StatusCode, string(body)))
		return "", nil, nil, nil, nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	if _, ok := data["sid"]; !ok {
		return "", nil, nil, nil, nil, err
	}
	sid, err := hex.DecodeString(data["sid"].(string))
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	encKey, err := hex.DecodeString(data["encryptionKey"].(string))
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	encIV, err := hex.DecodeString(data["encryptionIV"].(string))
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	encMacKey, err := hex.DecodeString(data["encryptionMacKey"].(string))
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	addr := ""
	if _, ok := data[serverType]; ok {
		addr = data[serverType].(string)
	}
	return addr, sid, encKey, encIV, encMacKey, nil
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
		if i < len(logmap) - 1 {
			csv += "\n"
		}
	}
	return csv
}
