// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Diarkis/diarkis-server-template/bot/field/custom"
	"github.com/Diarkis/diarkis-server-template/bot/field/fieldlib"
	"github.com/Diarkis/diarkis-server-template/bot/utils"
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/smap"
	"github.com/Diarkis/diarkis/util"
	uuid "github.com/Diarkis/diarkis/uuid/v4"
)

const UDP_STRING string = "udp"
const TCP_STRING string = "tcp"

const (
	STATUS_BEFORE_START = iota
	STATUS_AFTER_START
	STATUS_SYNC
)

// Config represents the structure of the configuration file
type Config struct {
	HostURL                              string `json:"Host"`
	BotCnt                               int    `json:"BotCnt"`
	NewPayloadFormat                     bool   `json:"NewPayloadFormat"`
	MovementIntervalMs                   int    `json:"MoveIntervalMs"`
	AreaWidth                            int    `json:"AreaWidth"`
	MovementRange                        int    `json:"MovementRange"`
	SyncCountPerMovement                 int    `json:"SyncCountPerMovement"`
	MoveDataCountPerSync                 int    `json:"MoveDataCountPerSync"`
	MoveProbabilityPercentagePerInterval int    `json:"MoveProbabilityPercenntagePerInterval"`
	ProtocolSource                       string `json:"Protocol"`
	MovementDuration                     int `json:"MoveDurationMs"`
}

// DefaultConfig holds the default values for the configuration
var DefaultConfig = Config{
	HostURL:                              "127.0.0.1:7000",
	BotCnt:                               250,
	NewPayloadFormat:                     true,
	MovementIntervalMs:                   2000,
	AreaWidth:                            10000,
	MovementRange:                        500,
	MoveProbabilityPercentagePerInterval: 50,
	MoveDataCountPerSync:                 5,
	SyncCountPerMovement:                 3,
	ProtocolSource:                       "udp",
	MovementDuration:		      1000,
}

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
var moveRatio = 5
var mapSize = 4500
var halfMapSize = 2250
var movementRange = 1200
var nbSyncPerMovement = 3
var nbMoveFrame = 16
var movementDuration = 1000
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

var useNewPayloadFormat = false

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
	isMoving   bool
}

// BotTask represents a task for a bot
type BotTask struct {
	ID int
	// Include other necessary bot attributes here
	bot *botData
}

func main() {
	parseFieldArgs()

	if mapSize <= movementRange {
		fmt.Println("Map size must be greater than range")
		os.Exit(1)
	}

	fmt.Printf("Starting Broadcast Bot. protocol: %v, bots num: %v, protocol: %v, broadcast interval: %v map: %v range: %v\n",
		proto, bots, proto, interval, mapSize, movementRange)

	tasks := make(chan BotTask, bots)
	numWorkers := 10 // Define how many workers you want to run concurrently

	go startBotWorkerPool(numWorkers, tasks)

	// Generate bot tasks
	for i := 1; i <= bots; i++ {
		uuid, _ := uuid.New()
		bot := new(botData)
		bot.uid = uuid.String
		for _, idx := range []int{8, 13, 18, 23} {
			bot.uid = bot.uid[:idx] + "-" + bot.uid[idx:]
		}
		bot.state = STATUS_BEFORE_START
		bot.inSightCnt = 0
		bot.userMap = smap.New()
		bot.x = util.RandomInt(-halfMapSize, halfMapSize)
		bot.y = util.RandomInt(-halfMapSize, halfMapSize)

		// Create a new task and send it to the worker pool
		tasks <- BotTask{ID: i, bot: bot}
	}

	close(tasks) // Close the task channel after all tasks are sent

	for {
		time.Sleep(time.Second * time.Duration(sleepTime))
		printBotStatus()
		clearMetricsCounter()
	}
}

func startBotWorkerPool(numWorkers int, tasks <-chan BotTask) {
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go botWorker(i, tasks, &wg)
	}

	wg.Wait()
}

func botWorker(workerID int, tasks <-chan BotTask, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing bot %d\n", workerID, task.ID)
		runBot(task.bot)
		fmt.Printf("Worker %d finished bot %d\n", workerID, task.ID)
	}
}

func runBot(bot *botData) {
	eResp, err := utils.Endpoint(host, bot.uid, proto)
	addr := eResp.ServerHost + ":" + fmt.Sprintf("%v", eResp.ServerPort)
	sid, _ := hex.DecodeString(eResp.Sid)
	key, _ := hex.DecodeString(eResp.EncryptionKey)
	iv, _ := hex.DecodeString(eResp.EncryptionIV)
	mkey, _ := hex.DecodeString(eResp.EncryptionMacKey)

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
		bot.udp.SetEncryptionKeys(sid, key, iv, mkey)
		bot.field = fieldlib.NewFieldAsUDP(bot.udp)
		bot.udp.OnConnect(func() {
			botsMap.Store(bot.uid, bot)
			go startBot(bot)
		})
		bot.udp.OnDisconnect(func() {
			fmt.Printf("Disconnected.")
			botsMap.Delete(bot.uid)
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
	for {
		switch bot.state {
		case STATUS_BEFORE_START:
			bot.field.Join(int64(bot.x), int64(bot.y), 0, 300, 0, nil, false, bot.uid)
			bot.state = STATUS_AFTER_START
		case STATUS_AFTER_START:
			randomSync(bot)
		case STATUS_SYNC:
			randomSync(bot)
		default:
			fmt.Printf("This is unexpected status!!! status is %v\n", bot.state)
			break
		}
		if !bot.isMoving {
			time.Sleep(time.Millisecond * time.Duration(interval))
		}
	}
}

// The rest of your functions (e.g., randomSync, addFieldListener, createMovementPayload, etc.) remain the same
// ...

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

func randomSpawnBot() {
	if botCounter >= bots {
		return
	}
	uuid, _ := uuid.New()
	bot := new(botData)
	bot.uid = uuid.String
	bot.isMoving = false
	for _, idx := range []int{8, 13, 18, 23} {
		bot.uid = bot.uid[:idx] + "-" + bot.uid[idx:]
	}
	bot.state = STATUS_BEFORE_START
	bot.inSightCnt = 0
	bot.userMap = smap.New()
	bot.x = util.RandomInt(-halfMapSize, halfMapSize)
	bot.y = util.RandomInt(-halfMapSize, halfMapSize)
	time.Sleep(time.Millisecond * time.Duration(int64(util.RandomInt(MIN_WAIT_MS, MIN_WAIT_MS))))

	eResp, err := utils.Endpoint(host, bot.uid, proto)
	addr := eResp.ServerHost + ":" + fmt.Sprintf("%v", eResp.ServerPort)
	sid, _ := hex.DecodeString(eResp.Sid)
	key, _ := hex.DecodeString(eResp.EncryptionKey)
	iv, _ := hex.DecodeString(eResp.EncryptionIV)
	mkey, _ := hex.DecodeString(eResp.EncryptionMacKey)

	if err != nil {
		fmt.Printf("Auth error ID:%v - %v\n", bot.uid, err)
		return
	}
	atomic.AddInt64(&httpCnt, 1)

	rcvByteSize := 1400

	switch proto {
	case UDP_STRING:
		udpSendInterval := int64(interval)
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

func float32ToByte(f float32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, f)
	return buf.Bytes()
}

func createNewMovementPayload(direction float32, prevX, prevY, x, y, nbMoveData int, timeStamp int64, frameInterval int, isLast bool) []byte {
	newPayload := custom.NewDiarkisCharacterSyncPayload()
	prevXUnity := float32(prevX) / 100.
	prevYUnity := float32(prevY) / 100.
	xUnity := float32(x) / 100.
	yUnity := float32(y) / 100.

	distanceX := xUnity - prevXUnity
	distanceY := yUnity - prevYUnity
	frameDistanceX := float32(distanceX) / float32(nbMoveData)
	frameDistanceY := float32(distanceY) / float32(nbMoveData)

	currentX := prevXUnity
	currentY := prevYUnity
	newRot := direction
	newPayload.Timestamp = timeStamp
	for i := 0; i < nbMoveData; i++ {
		frameData := custom.NewDiarkisCharacterFrameData()
		frameData.RotationAngles = uint16((((newRot + 180) / 360) * 65535) + 0.5)
		frameData.Position.X = float32(currentX) * 100
		frameData.Position.Z = float32(currentY) * 100
		frameData.Position.Y = 0
		frameData.TimestampInterval = uint16(frameInterval * (i + 1))
		frameData.AnimationBlend = 5.
		frameData.AnimationID = 1
		currentX += frameDistanceX
		currentY += frameDistanceY
		if isLast && i >= nbMoveData-1 {
			frameData.AnimationID = 0
			frameData.AnimationBlend = 0
		}
		newPayload.Frames = append(newPayload.Frames, frameData)
	}

	newPayload.Engine = 0
	payloadBytes := newPayload.Pack()
	return payloadBytes
}

func eulerToQuaternion(angle float64) *custom.DiarkisQuaternion {
	// Convert angle to radians
	angleRad := math.Pi * float64(angle) / 180.0

	// Calculate half angle
	halfAngle := angleRad * 0.5

	// Calculate quaternion components
	w := math.Cos(halfAngle)
	x := 0.0
	y := math.Sin(halfAngle)
	z := 0.0

	// Normalize the quaternion
	length := math.Sqrt(float64(w*w + x*x + y*y + z*z))
	w /= length
	x /= length
	y /= length
	z /= length

	quat := custom.NewDiarkisQuaternion()
	quat.W = float32(w)
	quat.X = float32(x)
	quat.Y = float32(y)
	quat.Z = float32(z)
	return quat
}

func createMovementPayload(direction float32, prevX, prevY, x, y, nbMoveData int, timeStamp int64, frameInterval int, useNewPayload, isLast bool) []byte {
	if useNewPayload {
		return createNewMovementPayload(direction, prevX, prevY, x, y, nbMoveData, timeStamp, frameInterval, isLast)
	}

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
	payload = append(payload, float32ToByte(float32(frameInterval))...)
	payload = append(payload, []byte{byte(nbMoveData)}...)

	for i := 0; i < nbMoveData; i++ {
		payload = append(payload, float32ToByte(frameDistanceX)...)
		payload = append(payload, float32ToByte(frameDistanceY)...)
		payload = append(payload, float32ToByte(float32(0))...)
		payload = append(payload, []byte{0}...)
	}

	return append(payloadSizeBytes, payload...)
}

func randomSync(bot *botData) {
	prevX := bot.x
	prevY := bot.y
	if bot.isMoving == false && rand.Intn(100) < moveRatio {
		bot.isMoving = true
		nextMoveIsInArea := false
		tryLimit := 20
		tryCnt := 0
		for !nextMoveIsInArea && tryCnt < tryLimit {
			r := float64(movementRange)
			angle := float64(util.RandomInt(1, 360))
			theta := (angle / 360.0) * 2.0 * math.Pi
			newX := int(float64(prevX) + r*float64(math.Cos(theta)))
			newY := int(float64(prevY) + r*float64(math.Sin(theta)))
			if newX >= -halfMapSize && newX <= halfMapSize && newY >= -halfMapSize && newY <= halfMapSize {
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
		frameInterval := movementDuration / (nbSyncPerMovement * nbMoveFrame)
		timeStamp := time.Now().UTC().UnixMilli()
		for i := 0; i < nbSyncPerMovement; i++ {
			nextX := currentX + stepX
			nextY := currentY + stepY
			isLast := i >= nbSyncPerMovement - 1
			message := createMovementPayload(bot.angle, currentX, currentY, nextX, nextY, nbMoveFrame, timeStamp, frameInterval, useNewPayloadFormat, isLast)
			bot.field.Sync(int64(nextX), int64(nextY), 0, 0, 0, message, false, bot.uid)
			currentX = nextX
			currentY = nextY
			time.Sleep(time.Millisecond * time.Duration(frameInterval*(nbMoveFrame-1)))
			timeStamp += int64(nbMoveFrame * frameInterval)
		}

		bot.state = STATUS_SYNC
		syncCnt += nbSyncPerMovement + 1
		bot.isMoving = false
	}
}

func parseFieldArgs() {
	configFilePath := "config.json"
	if len(os.Args) > 1 {
		configFilePath = os.Args[1]
	}

	// Check if the config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Config File doesn't exist, create it with default values
		file, err := os.Create(configFilePath)
		if err != nil {
			fmt.Println("Error creating config file:", err)
			return
		}
		defer file.Close()

		// Encode the default config to JSON and write it to the file with indentation
		var buffer bytes.Buffer
		encoder := json.NewEncoder(&buffer)
		encoder.SetIndent("", "  ") // Set indentation for JSON
		if err := encoder.Encode(DefaultConfig); err != nil {
			fmt.Println("Error encoding default config to JSON:", err)
			return
		}

		_, err = file.Write(buffer.Bytes())
		if err != nil {
			fmt.Println("Error writing default config to file:", err)
			return
		}

		fmt.Println("Config file created with default values.")
	} else {
		// File exists, read it
		file, err := os.Open(configFilePath)
		if err != nil {
			fmt.Println("Error opening config file:", err)
			return
		}
		defer file.Close()

		// Decode the JSON data into the Config struct
		var config Config
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			fmt.Println("Error decoding config file:", err)
			return
		}

		fmt.Println("Config file loaded:", config)
		host = config.HostURL
		bots = config.BotCnt
		proto = config.ProtocolSource
		if proto != "udp" && proto != "tcp" {
			fmt.Printf("Protocol value is only udp or tcp %v\n", err)
			os.Exit(1)
			return
		}
		useNewPayloadFormat = config.NewPayloadFormat
		interval = int64(config.MovementIntervalMs)
		mapSize = config.AreaWidth
		movementRange = config.MovementRange
		movementDuration = config.MovementDuration
		moveRatio = config.MoveProbabilityPercentagePerInterval
		nbSyncPerMovement = config.SyncCountPerMovement
		nbMoveFrame = config.MoveDataCountPerSync
		halfMapSize = mapSize / 2
	}

}
