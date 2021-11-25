package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Diarkis/diarkis/util"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/client/go/tcp"
	dpayload "{0}/lib/payload"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const cmdVer uint8 = 2
const cmdAdd uint16 = 100
const cmdSearch uint16 = 102
const pushRoomFull uint16 = 103
// waitingTime is in seconds
const waitingTime int64 = 30

var botCounter = 0
var host = "127.0.0.1:7000"
var proto = "udp" // udp or tcp
var howmany = 10
// sleepTime is in seconds
var sleepTime int64 = 60
var searchProps = make(map[string]int)
var addProps = make(map[string]int)
var states = make(map[int]int)
var newSpawns = make([]int, 0)

type botData struct {
	uid   int
	state int
	udp   *udp.Client
	tcp   *tcp.Client
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Bot requires 2 parameters: {http host:port} {how many bots}")
		os.Exit(1)
		return
	}
	host = os.Args[1]
	howmany, _ = strconv.Atoi(os.Args[2])
	fmt.Printf("Starting MatchMaker Bot %v - %v bots\n", proto, howmany)
	searchProps["rank"] = 5
	addProps["rank"] = 5
	spawnBots()
	for {
		time.Sleep(time.Second * time.Duration(sleepTime))
		if botCounter <= 0 {
			//break
		}
	}
	fmt.Printf("All bots have finished their works - Exiting the process - Bye!\n")
	os.Exit(0)
}

func spawnBots() {
	for i := 0; i < howmany; i++ {
		go spawnUDPBot(i, true);
	}
}

func spawnUDPBot(id int, needToWait bool) {
	if needToWait {
		waitMin := 0
		waitMax := 5000
		time.Sleep(time.Millisecond * time.Duration(int64(util.RandomInt(waitMin, waitMax))))
	}
	addr, sid, key, iv, mkey, err := auth(id)
	if err != nil {
		fmt.Printf("Auth error ID:%v - %v\n", id, err)
		return
	}
	rcvByteSize := 1400
	udpSendInterval := int64(100)
	udp.LogLevel(9)
	cli := udp.New(rcvByteSize, udpSendInterval)

	bot := new(botData)
	bot.uid = id
	bot.state = 0
	bot.udp = cli

	cli.SetEncryptionKeys(sid, key, iv, mkey)
	cli.OnResponse(func(ver uint8, cmd uint16, status uint8, payload []byte) {
		handleOnResponse(bot, ver, cmd, status, payload)
	})
	cli.OnPush(func(ver uint8, cmd uint16, payload []byte) {
		handleOnPush(bot, ver, cmd, payload)
	})
	cli.OnConnect(func() {
		go startBot(bot)
	})
	fmt.Printf("Bot ID: %v - Connecting to %v\n", id, addr)
	cli.Connect(addr)
}

func startBot(bot *botData) {
	botCounter++
	// 0 ~ 20 = search
	// 21     = add
	// 22     = wait
	// 23     = room full and disconnect
	currentState := -1
	if util.RandomInt(0, 99) < 30 {
		bot.state = 21
	}
	fmt.Printf("%v bot started ID:%v (state:%v) - Total bots :%v\n", proto, bot.uid, bot.state, botCounter)
	waitCounter := int64(0)
	interval := int64(200)
	for {
		if currentState < 22 && currentState == bot.state {
			time.Sleep(time.Millisecond * time.Duration(interval))
			continue
		}
		fmt.Printf("Bot ID:%v - state is %v\n", bot.uid, bot.state)
		switch bot.state {
		case -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 18, 19, 20:
			fmt.Printf("Bot ID %v (state:%v) - search\n", bot.uid, bot.state)
			search(bot)
		case 21:
			fmt.Printf("Bot ID %v (state:%v) - add\n", bot.uid, bot.state)
			add(bot)
		case 22:
			// We are waiting
			fmt.Printf("Bot ID:%v is now waiting and will finish in %v seconds... %v seconds to go\n",
				bot.uid, waitingTime, (waitingTime * 1000) - waitCounter)
			waitCounter += interval
			if waitCounter >= waitingTime * 1000 {
				bot.state = 23
			}
		case 23:
			// Bot disconnects
			fmt.Printf("Bot ID %v (state:%v) - disconnect\n", bot.uid, bot.state)
			disconnect(bot)
		default:
			fmt.Printf("Error corrupt state %v - Bot ID:%v does nothing...\n", bot.state, bot.uid)
			bot.state = 23
		}
		currentState = bot.state
		fmt.Printf("Bot ID:%v state updated to %v\n", bot.uid, currentState)
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}

func search(bot *botData) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		return
	}
	fmt.Printf("MatchMaker search client ID:%v\n", bot.uid)
	pkt := dpayload.PackMMSearch([]string{"RankMatch"}, searchProps)
	switch proto {
	case "udp":
		bot.udp.RSend(cmdVer, cmdSearch, pkt)
	case "tcp":
		bot.tcp.Send(cmdVer, cmdSearch, pkt)
	}
}

func add(bot *botData) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		return
	}
	fmt.Printf("MatchMaker add client ID:%v\n", bot.uid)
	pkt := dpayload.PackMMAdd("RankMatch", fmt.Sprintf("%v", bot.uid), addProps,  []byte(fmt.Sprintf("My ID is %v", bot.uid)), uint64(60))
	switch proto {
	case "udp":
		bot.udp.RSend(cmdVer, cmdAdd, pkt)
	case "tcp":
		bot.tcp.Send(cmdVer, cmdAdd, pkt)
	}
}

func disconnect(bot *botData) {
	if bot.state == 0 && bot.udp == nil && bot.tcp == nil {
		return
	}
	botCounter--
	fmt.Printf("Bot ID:%v finished its work and disconnects - Total bots :%v\n", bot.uid, botCounter)
	switch proto {
	case "udp":
		bot.udp.Disconnect()
	case "tcp":
		bot.tcp.Disconnect()
	}
	go spawnUDPBot(bot.uid, false);
	bot.udp = nil
	bot.tcp = nil
	bot.uid = -1
	bot.state = 0
}

func handleOnResponse(bot *botData, ver uint8, cmd uint16, status uint8, payload []byte) {
	if ver != cmdVer {
		return
	}
	switch cmd {
	case cmdAdd:
		if string(payload) == "OK" {
			bot.state = 22
		} else {
			bot.state = 23
		}
		fmt.Printf("Bot ID %v added - state %v\n", bot.uid, bot.state)
	case cmdSearch:
		if string(payload) == "Matching not found" {
			bot.state += 1
		} else {
			bot.state = 22
		}
		fmt.Printf("Bot ID %v search - state %v\n", bot.uid, bot.state)
	}
}

func handleOnPush(bot *botData, ver uint8, cmd uint16, payload []byte) {
	if ver != cmdVer {
		return
	}
	switch cmd {
	case pushRoomFull:
		fmt.Printf("Bot ID:%v received room full notification\n", bot.uid)
		// The joined room is full
		bot.state = 23
	}
}

// returns addr, sid, key, iv, mackey, error
func auth(uid int) (string, []byte, []byte, []byte, []byte, error) {
	url := fmt.Sprintf("http://%s/auth/%v", host, uid)
	fmt.Printf("Bot ID: %v - Connecting to HTTP: %s\n", uid, url)
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
	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
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
	if _, ok := data[strings.ToUpper(proto)]; ok {
		addr = data[strings.ToUpper(proto)].(string)
	}
	return addr, sid, encKey, encIV, encMacKey, nil
}
