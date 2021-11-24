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
	"sync"
	"time"
)

const cmdVer uint8 = 2
const cmdAdd uint16 = 100
const cmdSearch uint16 = 102
const pushRoomFull uint16 = 103
// waitingTime is in seconds
const waitingTime int64 = 60

var botCounter = 0
var host = "127.0.0.1:7000"
var proto = "udp" // udp or tcp
var howmany = 10
// sleepTime is in seconds
var sleepTime int64 = 60
var searchProps = make(map[string]int)
var addProps = make(map[string]int)
var states = make(map[int]int)
var mutex = new(sync.RWMutex)

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
		if botCounter == 0 {
			break
		}
	}
	fmt.Printf("All bots have finished their works - Exiting the process - Bye!\n")
	os.Exit(0)
}

func spawnBots() {
	for i := 0; i < howmany; i++ {
		go spawnUDPBot(i);
	}
}

func spawnUDPBot(id int) {
	waitMin := 0
	waitMax := 5000
	time.Sleep(time.Millisecond * time.Duration(int64(util.RandomInt(waitMin, waitMax))))
	addr, sid, key, iv, mkey, err := auth(id)
	if err != nil {
		fmt.Printf("Auth error ID:%v - %v\n", id, err)
		return
	}
	rcvByteSize := 1400
	udpSendInterval := int64(100)
	udp.LogLevel(9)
	cli := udp.New(rcvByteSize, udpSendInterval)
	cli.SetEncryptionKeys(sid, key, iv, mkey)
	cli.OnResponse(func(ver uint8, cmd uint16, status uint8, payload []byte) {
		handleOnResponse(id, ver, cmd, status, payload)
	})
	cli.OnPush(func(ver uint8, cmd uint16, payload []byte) {
		handleOnPush(id, ver, cmd, payload)
	})
	cli.OnConnect(func() {
		go startBot(id, cli, nil)
	})
	fmt.Printf("Connecting to %v\n", addr)
	cli.Connect(addr)
}

func startBot(uid int, udpCli *udp.Client, tcpCli *tcp.Client) {
	botCounter++
	fmt.Printf("%v bot started ID:%v - Total bots :%v\n", proto, uid, botCounter)
	// 0 ~ 5 = search
	// 6     = add
	// 10    = room full and disconnect
	currentState := -1
	mutex.Lock()
	states[uid] = 0
	mutex.Unlock()
	for {
		mutex.RLock()
		if currentState == states[uid] {
			mutex.RUnlock()
			continue
		}
		fmt.Printf("Bot ID:%v - state is %v\n", uid, states[uid])
		switch states[uid] {
		case -1, 0, 1, 2, 3, 4, 5:
			search(uid, udpCli, tcpCli)
		case 6:
			add(uid, udpCli, tcpCli)
		case 7:
			// We are waiting
			go func() {
				fmt.Printf("Bot ID:%v is now waiting and will finish in 1 minute\n", uid)
				time.Sleep(time.Second * time.Duration(waitingTime))
				mutex.Lock()
				states[uid] = 10
				mutex.Unlock()
			}()
		case 10:
			// Bot disconnects
			disconnect(uid, udpCli, tcpCli)
		default:
			fmt.Printf("Error corrupt state %v - Bot ID:%v does nothing...\n", states[uid], uid)
			mutex.Lock()
			states[uid] = 10
			mutex.Unlock()
		}
		currentState = states[uid]
		fmt.Printf("Bot ID:%v state updated to %v\n", uid, currentState)
		mutex.RUnlock()
		time.Sleep(time.Millisecond * 100)
	}
}

func search(uid int, udpCli *udp.Client, tcpCli *tcp.Client) {
	fmt.Printf("MatchMaker search client ID:%v\n", uid)
	pkt := dpayload.PackMMSearch([]string{"RankMatch"}, searchProps)
	switch proto {
	case "udp":
		udpCli.RSend(cmdVer, cmdSearch, pkt)
	case "tcp":
		tcpCli.Send(cmdVer, cmdSearch, pkt)
	}
}

func add(uid int, udpCli *udp.Client, tcpCli *tcp.Client) {
	fmt.Printf("MatchMaker add client ID:%v\n", uid)
	pkt := dpayload.PackMMAdd("RankMatch", fmt.Sprintf("%v", uid), addProps,  []byte(fmt.Sprintf("My ID is %v", uid)), uint64(60))
	switch proto {
	case "udp":
		udpCli.RSend(cmdVer, cmdAdd, pkt)
	case "tcp":
		tcpCli.Send(cmdVer, cmdAdd, pkt)
	}
}

func disconnect(uid int, udpCli *udp.Client, tcpCli *tcp.Client) {
	botCounter--
	fmt.Printf("Bot ID:%v finished its work and disconnects - Total bots :%v\n", uid, botCounter)
	switch proto {
	case "udp":
		udpCli.Disconnect()
	case "tcp":
		tcpCli.Disconnect()
	}
	// re-spawn bot
	spawnUDPBot(uid)
}

func handleOnResponse(uid int, ver uint8, cmd uint16, status uint8, payload []byte) {
	if ver != cmdVer {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	switch cmd {
	case cmdAdd:
		if string(payload) == "OK" {
			states[uid] = 7
		}
	case cmdSearch:
		if string(payload) == "Matching not found" {
			states[uid] += 1
		} else {
			states[uid] = 7
		}
	}
}

func handleOnPush(uid int, ver uint8, cmd uint16, payload []byte) {
	if ver != cmdVer {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	switch cmd {
	case pushRoomFull:
		fmt.Printf("Bot ID:%v received room full notification\n", uid)
		// The joined room is full
		states[uid] = 10
	}
}

// returns addr, sid, key, iv, mackey, error
func auth(uid int) (string, []byte, []byte, []byte, []byte, error) {
	url := fmt.Sprintf("http://%s/auth/%v", host, uid)
	fmt.Printf("Connecting to HTTP: %s\n", url)
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
		resp.Body.Close()
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
