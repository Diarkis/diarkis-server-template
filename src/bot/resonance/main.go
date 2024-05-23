package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/util"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// http timeout seconds
const authTimeout = 10
const cmdVer uint8 = 2
const cmdResonance uint16 = 13
const minIdleMilliSecond = 500
const maxIdleMilliSecond = 5000
const sleepTime = 10

var interval = 0
var botNum = 0

var botCounter int32
var resonanceCnt = 0
var respondCnt = 0
var pktSize = 0
var host = "127.0.0.1:7000"
var proto = "udp"

type botData struct {
	uid   int
	state int
	udp   *udp.Client
	tcp   *tcp.Client
}

func parseArgs() {
	if len(os.Args) < 4 {
		fmt.Printf("Bot requires 4 parameters: {http host:port} {bot num} {packet interval} {packet size}")
		os.Exit(1)
		return
	}

	var err error
	host = os.Args[1]
	botNum, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("bot num paramter given is invalid #{err}\n")
		os.Exit(1)
	}

	interval, err = strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("bot num paramter given is invalid #{err}\n")
		os.Exit(1)
	}

	pktSize, err = strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("packet size paramter given is invalid #{err}\n")
		os.Exit(1)
	}
}

func main() {
	parseArgs()
	spawnBots()
	for {
		time.Sleep(time.Second * time.Duration(sleepTime))
		timeStamp := util.ZuluTimeFormat(time.Now())
		fmt.Printf("{ \"Time\":\"%v\", \"Bots\":%v, \"ResonanceCnt\":%v, \"RespondCnt\":%v }\n", timeStamp, botCounter, resonanceCnt, respondCnt)
		resonanceCnt = 0
		respondCnt = 0
	}
	fmt.Printf("Exiting...\n")
	os.Exit(0)

}

func spawnBots() {
	for i := 0; i < int(botNum); i++ {
		go spawnUDPBot(i)
	}
}

func spawnUDPBot(id int) {
	time.Sleep(time.Millisecond * time.Duration(int64(util.RandomInt(minIdleMilliSecond, maxIdleMilliSecond))))
	addr, sid, key, iv, mkey, err := auth(id)
	if err != nil {
		fmt.Printf("auth error ID:%v - %v\n", id, err)
		return
	}
	udp.LogLevel(9)
	cli := udp.New(1400, 100)
	bot := new(botData)

	bot.uid = id
	bot.state = 0
	bot.udp = cli
	cli.SetEncryptionKeys(sid, key, iv, mkey)
	cli.OnResponse(func(ver uint8, cmd uint16, status uint8, payload []byte) {
		if ver == cmdVer && cmd == cmdResonance {
			respondCnt++
		}
	})
	cli.OnConnect(func() {
		go startBot(bot)
	})
	cli.OnDisconnect(func() {
		fmt.Printf("BotOnDisconnect. id:%v", id)
		atomic.AddInt32(&botCounter, -1)

		if int(botCounter) >= botNum {
			return
		}
		spawnUDPBot(bot.uid) // recursive
	})
	cli.Connect(addr)
}

func startBot(bot *botData) {
	atomic.AddInt32(&botCounter, 1)
	//waitCounter := int64(0)
	for {
		time.Sleep(time.Millisecond * time.Duration(interval))
		resonance(bot)
	}

}

func resonance(bot *botData) {
	if bot.udp == nil {
		return
	}
	bot.udp.Send(cmdVer, cmdResonance, make([]byte, pktSize))
	resonanceCnt++
}

func auth(uid int) (string, []byte, []byte, []byte, []byte, error) {
	url := fmt.Sprintf("http://%s/auth/%v", host, uid)
	client := &http.Client{
		Timeout: time.Second * authTimeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, nil, nil, nil, err
	}

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
	if _, ok := data[strings.ToUpper(proto)]; ok {
		addr = data[strings.ToUpper(proto)].(string)
	}
	return addr, sid, encKey, encIV, encMacKey, nil
}
