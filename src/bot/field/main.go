package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"{0}/bot/field/fieldlib"

	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/smap"
	"github.com/Diarkis/diarkis/util"
	v4 "github.com/Diarkis/diarkis/uuid/v4"
)

const UDP_STRING string = "udp"
const TCP_STRING string = "tcp"

const (
	STATUS_BEFORE_START = iota
	STATUS_AFTER_START
	STATUS_SYNC
)

const MIN_WAIT_MS = 1000
const MAX_WAIT_MS = 90000

const PACKET_SIZE = 100

const SERVER_COUNT = 1

// parameters
var proto = "udp" // udp or tcp
var host string
var bots = 10
var packetSize = 10
var interval int64
var moveRatio = 50
var mapSize = 4500
var movementRange = 1200

// metrics counter
var botCounter = 0
var syncCnt int
var onSyncCnt = smap.New()
var onDisappearCnt int
var latencyTotal int64
var httpCnt int64
var latencyMax int64
var latencyCnt int64

// sleep time is in seconds
var sleepTime int64 = 1
var botsMap sync.Map

type botData struct {
	uid        string
	state      int
	udp        *udp.Client
	tcp        *tcp.Client
	field      *fieldlib.Field
	x          int
	y          int
	angle      float32
	userMap    smap.SyncMap
	inSightCnt int
}

func main() {
	if len(os.Args) < 7 {
		fmt.Printf("Bot requires 4 parameters: {http host:port} {how many bots} {protocol} {packet interval in milliseconds} {map size} {range}")
		os.Exit(1)
		return
	}

	rand.Seed(time.Now().UnixNano())

	parseFieldArgs()

	if mapSize <= movementRange {
		fmt.Println("Map size must be greater than range")
		os.Exit(1)
	}

	fmt.Printf("Starting Broadcast Bot. protocol: %v, bots num: %v, protocol: %v, broadcast interval: %v map: %v range: %v\n",
		proto, bots, proto, interval, mapSize, movementRange)

	spawnBots()
	for {
		time.Sleep(time.Second * time.Duration(sleepTime))
		printBotStatus()
		clearMetricsCounter()
	}

	fmt.Printf("All bots have finished their works - Exiting the process - Bye!\n")
	os.Exit(0)
}

func clearMetricsCounter() {
	onSyncCnt.Clear()
	syncCnt = 0
	onDisappearCnt = 0
	latencyCnt = 0
	latencyTotal = 0
	latencyMax = 0
	httpCnt = 0
}

func printBotStatus() {
	timestamp := util.ZuluTimeFormat(time.Now())
	inSightMax := 0
	inSightTotal := 0

	botNum := 0
	inSightMap := make(map[int]int)
	botsMap.Range(func(k, v interface{}) bool {
		bot := v.(*botData)
		if bot.inSightCnt > inSightMax {
			inSightMax = bot.inSightCnt
		}
		_, ok := inSightMap[bot.inSightCnt]
		if ok {
			inSightMap[bot.inSightCnt]++
		} else {
			inSightMap[bot.inSightCnt] = 1
		}
		inSightTotal += bot.inSightCnt
		botNum++
		bot.inSightCnt = 0
		bot.userMap = smap.New()
		return true
	})
	inSightAvg := 0
	if botNum != 0 { // prevent divide by zero
		inSightAvg = inSightTotal / botNum
	}

	fmt.Printf("{ \"Time\":\"%v\", \"Bots\":%v, \"http\":%v, \"Sync\": %v , \"inSightMax\": %v , \"inSightAvg\": %v , \"onDisappear\": %v}\n",
		timestamp, botCounter, httpCnt, syncCnt, inSightMax, inSightAvg, onDisappearCnt)
	inSightMap = make(map[int]int)
}

func spawnBots() {
	for i := 0; i < bots; i++ {
		go randomSpawnBot()
	}
}

func addFieldListener(bot *botData) {
	bot.field.OnDisappear(func(message string) {
		onDisappearCnt++
	})
	bot.field.OnSync(func(message []byte) {
		uid := fmt.Sprintf("%v", message[28])
		bot.userMap.Set(uid, 1)
		bot.inSightCnt = len(bot.userMap.Keys())
	})
}

func endpoint(uid string, host string, proto string) (string, []byte, []byte, []byte, []byte, error) {
	url := fmt.Sprintf("http://%s/endpoint/type/%s/user/%s", host, "UDP", uid)
	client := &http.Client{
		Timeout: time.Second * 10,
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
	if resp.StatusCode > 300 {
		err := errors.New(fmt.Sprintf("Error response status %v - body:%v", resp.StatusCode, string(body)))
		return "", nil, nil, nil, nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)

	fmt.Println(err)
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
	addr := data["serverHost"].(string) + ":" + strconv.Itoa(int(data["serverPort"].(float64)))
	return addr, sid, encKey, encIV, encMacKey, nil
}

func randomSpawnBot() {
	if botCounter >= bots {
		return
	}
	id, _ := v4.New()
	bot := new(botData)
	bot.uid = id.String
	for _, idx := range []int{8, 13, 18, 23} {
		bot.uid = bot.uid[:idx] + "-" + bot.uid[idx:]
	}
	bot.state = STATUS_BEFORE_START
	bot.inSightCnt = 0
	bot.userMap = smap.New()
	bot.x = util.RandomInt(-mapSize/2, mapSize/2)
	bot.y = util.RandomInt(-mapSize/2, mapSize/2)
	time.Sleep(time.Millisecond * time.Duration(int64(util.RandomInt(MIN_WAIT_MS, MIN_WAIT_MS))))
	//addr, sid, key, iv, mkey, err := botUtil.FieldAuth(id.String, host, proto, bot.x, bot.y, 0)
	addr, sid, key, iv, mkey, err := endpoint(bot.uid, host, proto)
	if err != nil {
		fmt.Printf("Auth error ID:%v - %v\n", bot.uid, err)
		return
	}
	atomic.AddInt64(&httpCnt, 1)

	rcvByteSize := 1400

	switch proto {
	case UDP_STRING:
		udpSendInterval := int64(100)
		bot.udp = udp.New(rcvByteSize, udpSendInterval)
		//udp.LogLevel(9)
		bot.udp.SetEncryptionKeys(sid, key, iv, mkey)
		bot.field = fieldlib.NewFieldAsUDP(bot.udp)
		bot.udp.OnConnect(func() {
			botsMap.Store(bot.uid, bot)
			go startBot(bot)
		})
		bot.udp.OnDisconnect(func() {
			fmt.Printf("Disconnected.")
			botsMap.Delete(bot.uid)
			// botCounter--
			if botCounter >= bots {
				return
			}
			time.Sleep(time.Millisecond * time.Duration(int64(util.RandomInt(MIN_WAIT_MS, MAX_WAIT_MS))))
			randomSpawnBot()
		})
	case TCP_STRING:
		tcpSendInterval := int64(100)
		tcpHbInterval := int64(1000)
		bot.tcp = tcp.New(rcvByteSize, tcpSendInterval, tcpHbInterval)
		tcp.LogLevel(9)
		bot.tcp.SetEncryptionKeys(sid, key, iv, mkey)
		bot.field = fieldlib.NewFieldAsTCP(bot.tcp)
		bot.tcp.OnConnect(func() {
			botsMap.Store(bot.uid, bot)
			startBot(bot)
		})
		bot.tcp.OnDisconnect(func() {
			fmt.Printf("Disconnected.")
			botsMap.Delete(bot.uid)
			botCounter--
			if botCounter >= bots {
				return
			}
			time.Sleep(time.Millisecond * time.Duration(int64(util.RandomInt(MIN_WAIT_MS, MAX_WAIT_MS))))
			randomSpawnBot()
		})
	}
	//Bots are not listening to the server for now
	addFieldListener(bot)

	switch proto {
	case UDP_STRING:
		bot.udp.Connect(addr)
	case TCP_STRING:
		bot.tcp.Connect(addr)
	}
}

func startBot(bot *botData) {
	botCounter++
	bot.state = STATUS_AFTER_START

	for {
		switch bot.state {
		case STATUS_BEFORE_START:
			bot.state = STATUS_AFTER_START
		case STATUS_AFTER_START:
			randomSync(bot)
		case STATUS_SYNC:
			randomSync(bot)
		default:
			fmt.Printf("This is unexpected status!!! status is %v\n", bot.state)
			break
		}
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}

func float32ToByte(f float32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, f)
	return buf.Bytes()
}

func createMovementPayload(direction float32, prevX, prevY, x, y, nbMoveData, fps int) []byte {
	payloadSize := 1 + (4 * 8) + (nbMoveData)*(13)
	payload := []byte{}
	payloadSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(payloadSizeBytes, uint16(payloadSize))

	payload = append(payload, float32ToByte(float32(prevX))...)
	payload = append(payload, float32ToByte(float32(prevY))...)
	payload = append(payload, float32ToByte(float32(0))...)
	payload = append(payload, float32ToByte(float32(direction))...)

	distanceX := float32(x) - float32(prevX)
	distanceY := float32(y) - float32(prevY)

	frameDistanceX := float32(distanceX) / float32(nbMoveData)
	frameDistanceY := float32(distanceY) / float32(nbMoveData)

	velocityX := 1.0
	velocityY := 1.0

	if prevX == x {
		velocityX = 0
	}

	if prevY == y {
		velocityY = 0
	}

	payload = append(payload, float32ToByte(float32(velocityX))...)
	payload = append(payload, float32ToByte(float32(velocityY))...)
	payload = append(payload, float32ToByte(float32(0))...)
	payload = append(payload, float32ToByte(float32(fps))...)
	payload = append(payload, []byte{byte(nbMoveData)}...)

	for i := 0; i < nbMoveData; i++ {
		payload = append(payload, float32ToByte(frameDistanceX)...)
		payload = append(payload, float32ToByte(frameDistanceY)...)
		payload = append(payload, float32ToByte(float32(0))...)
		payload = append(payload, []byte{ 0 }...)
	}

	return append(payloadSizeBytes, payload...)
}

func randomSync(bot *botData) {
	prevX := bot.x
	prevY := bot.y
	nbSyncPerMovement := 9
	if rand.Intn(100) < moveRatio {
		nextMoveIsInArea := false
		tryLimit := 20
		tryCnt := 0
		for !nextMoveIsInArea {
			if tryCnt >= tryLimit {
				return
			}
			r := float64(movementRange)
			angle := float64(util.RandomInt(1, 360))
			theta := (angle / 360.0) * 2.0 * math.Pi
			newX := int(float64(prevX) + r*float64(math.Cos(theta)))
			newY := int(float64(prevY) + r*float64(math.Sin(theta)))
			limitXUp := mapSize / 2
			limitXDown := -mapSize / 2
			limitYUp := mapSize / 2
			limitYDown := -mapSize / 2
			if prevX < 0 {
				limitXUp = 0
			} else {
				limitXDown = 0
			}
			if prevY < 0 {
				limitYUp = 0
			} else {
				limitYDown = 0
			}
			if !(newX > limitXUp || newX < limitXDown || newY > limitYUp || newY < limitYDown) {
				nextMoveIsInArea = true
				bot.angle = float32(angle)
				bot.x = newX
				bot.y = newY
			}
			tryCnt++
		}
		stepX := (bot.x - prevX) / nbSyncPerMovement
		stepY := (bot.y - prevY) / nbSyncPerMovement
		currentX := prevX
		currentY := prevY
		for i := 0; i < nbSyncPerMovement; i++ {
			nextX := currentX + stepX
			nextY := currentY + stepY
			message := createMovementPayload(bot.angle, currentX, currentY, nextX, nextY, 6, 60)
			go bot.field.Sync(int64(nextX), int64(nextY), 0, 300, 0, message, false, bot.uid)
			currentX = nextX
			currentY = nextY
			time.Sleep(time.Millisecond * time.Duration(100))
		}
		message := createMovementPayload(bot.angle, bot.x, bot.y, bot.x, bot.y, 6, 60)
		go bot.field.Sync(int64(bot.x), int64(bot.y), 0, 300, 0, message, false, bot.uid)

		bot.state = STATUS_SYNC
		syncCnt += nbSyncPerMovement + 1
	} else {
		currentX := prevX
		currentY := prevY
		for i := 0; i < nbSyncPerMovement; i++ {
			nextX := currentX
			nextY := currentY
			message := createMovementPayload(bot.angle, currentX, currentY, currentX, currentY, 6, 60)
			go bot.field.Sync(int64(nextX), int64(nextY), 0, 300, 0, message, false, bot.uid)
			currentX = nextX
			currentY = nextY
			time.Sleep(time.Millisecond * time.Duration(100))
		}
		message := createMovementPayload(bot.angle, bot.x, bot.y, bot.x, bot.y, 6, 60)
		go bot.field.Sync(int64(bot.x), int64(bot.y), 0, 300, 0, message, false, bot.uid)

		bot.state = STATUS_SYNC
		syncCnt += nbSyncPerMovement + 1
	}
}

func parseFieldArgs() {
	host = os.Args[1]
	botsSource, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("How many bot parameter given is invalid %v\n", err)
		os.Exit(1)
		return
	}
	bots = botsSource

	protocolSource := strings.ToLower(os.Args[3]) // only tcp or udp
	if protocolSource != "udp" && protocolSource != "tcp" {
		fmt.Printf("Protocol value is only udp or tcp %v\n", err)
		os.Exit(1)
		return
	}
	proto = protocolSource

	intervalSource, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Interval of broadcast parameter given is invalid %v\n", err)
		os.Exit(1)
		return
	}
	interval = int64(intervalSource)

	mapSizeSource, err := strconv.Atoi(os.Args[5])
	if err != nil {
		fmt.Printf("Map size parameter given is invalid %v\n", err)
		os.Exit(1)
		return
	}
	mapSize = mapSizeSource

	rangeSource, err := strconv.Atoi(os.Args[6])
	if err != nil {
		fmt.Printf("Movement range parameter given is invalid %v\n", err)
		os.Exit(1)
		return
	}
	movementRange = rangeSource
}
