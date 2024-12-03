// © 2019-2024 Diarkis Inc. All rights reserved.

package loadtest

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/Diarkis/diarkis-server-template/bot/utils"
	"github.com/Diarkis/diarkis/client/go/modules/dm"
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/util"
)

type Bot struct {
	uid string
	dm  *dm.DirectMessage
}

type Params struct {
	Host     string
	Howmany  int
	Protocol string
	Size     int
	Interval int64
}

type Report struct {
	Sent     int
	Received int
}

var bots = make([]*Bot, 0)
var reports = new(Report)
var mutex = new(sync.RWMutex)

func Spawn(params *Params) {
	bot := &Bot{dm: &dm.DirectMessage{}}

	var tcpCli *tcp.Client
	var udpCli *udp.Client

	switch params.Protocol {
	case "TCP":
		rcvByteSize := 1400
		tcpSendInterval := int64(100)
		tcpHbInterval := int64(5000)
		tcp.LogLevel(50)
		tcpCli = tcp.New(rcvByteSize, tcpSendInterval, tcpHbInterval)
		bot.dm.SetupAsTCP(tcpCli)
		bot.uid = tcpCli.ID
	case "UDP":
		rcvByteSize := 1400
		udpSendInterval := int64(100)
		udp.LogLevel(50)
		udpCli = udp.New(rcvByteSize, udpSendInterval)
		bot.dm.SetupAsUDP(udpCli)
		bot.uid = udpCli.ID
	}

	bot.dm.OnPeerSend(func(uid string, message []byte) {
		go countReceive()
	})

	fmt.Println("Bot spawned -> UID", bot.uid)

	eResp, err := utils.Endpoint(params.Host, bot.uid, params.Protocol)
	addr := eResp.ServerHost + ":" + fmt.Sprintf("%v", eResp.ServerPort)
	sid, _ := hex.DecodeString(eResp.Sid)
	key, _ := hex.DecodeString(eResp.EncryptionKey)
	iv, _ := hex.DecodeString(eResp.EncryptionIV)
	mkey, _ := hex.DecodeString(eResp.EncryptionMacKey)

	if err != nil {
		fmt.Println("Failed to get endpoint for the client", err)
		return
	}

	await := util.Async(1)

	if tcpCli != nil {
		tcpCli.SetEncryptionKeys(sid, key, iv, mkey)
		tcpCli.OnConnect(func() {
			await.Done()
		})
		tcpCli.OnDisconnect(func() {
			fmt.Println("Bot disconnected -> UID", bot.uid)
		})
		fmt.Println("Bot connecting -> UID", bot.uid, addr)
		tcpCli.Connect(addr)
	}
	if udpCli != nil {
		udpCli.SetEncryptionKeys(sid, key, iv, mkey)
		udpCli.OnConnect(func() {
			await.Done()
		})
		udpCli.OnDisconnect(func() {
			fmt.Println("Bot disconnected -> UID", bot.uid)
		})
		fmt.Println("Bot connecting -> UID", bot.uid, addr)
		udpCli.Connect(addr)
	}

	await.Wait()

	bots = append(bots, bot)
}

func StartLoadTest(params *Params) {
	for i := 0; i < params.Howmany; i++ {
		go spam(bots[i], params)
	}
}

func GetReport() (sent int, received int) {
	mutex.Lock()
	defer mutex.Unlock()
	sent = reports.Sent
	received = reports.Received
	reports.Sent = 0
	reports.Received = 0
	return sent, received
}

func spam(bot *Bot, params *Params) {
	fmt.Println("Bot load test started -> UID", bot.uid)
	targetUID := getTargetUID(bot.uid)
	message := make([]byte, params.Size)

	for {
		bot.dm.Send(targetUID, message)
		go countSend()
		time.Sleep(time.Millisecond * time.Duration(params.Interval))
	}
}

func countSend() {
	mutex.Lock()
	defer mutex.Unlock()
	reports.Sent++
}

func countReceive() {
	mutex.Lock()
	defer mutex.Unlock()
	reports.Received++
}

func getTargetUID(myUID string) string {
	for {
		uid := bots[util.RandomInt(0, len(bots)-1)].uid

		if uid != myUID {
			return uid
		}
	}
}
