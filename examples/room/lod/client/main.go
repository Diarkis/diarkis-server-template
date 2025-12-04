// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	proom "github.com/Diarkis/diarkis-server-template/examples/room/lod/puffer/go/room"
	"github.com/Diarkis/diarkis/client/go/test/cli"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/log"
)

var logger = log.New("CLI")

var udpLod *lod

func main() {
	cli.SetupBuiltInRoomCommands()
	cli.RegisterCommands("lod", []cli.Command{
		{CmdName: "b", Desc: "Broadcast the LoD(alias of broadcast)", CmdFunc: broadcastLoD},
		{CmdName: "broadcast", Desc: "Broadcast the LoD", CmdFunc: broadcastLoD},
		{CmdName: "getinfo", Desc: "Get the LoD configuration", CmdFunc: getLoDInfo},
	})
	cli.Connect()
	udpLod = setupLoD(cli.UDPClient)
	cli.Run()
}

func broadcastLoD() {
	reader := bufio.NewReader(os.Stdin)
	xInt := int64(0)
	yInt := int64(0)
	var err error
	for {
		fmt.Println("Enter X (int32):")
		x, _ := reader.ReadString('\n')
		x = strings.Trim(x, "\r\n")
		xInt, err = strconv.ParseInt(x, 10, 32)
		if err != nil {
			fmt.Printf("Failed to parse X: %v\n", err)
			continue
		}
		break
	}
	for {
		fmt.Println("Enter Y (int32):")
		y, _ := reader.ReadString('\n')
		y = strings.Trim(y, "\r\n")
		yInt, err = strconv.ParseInt(y, 10, 32)
		if err != nil {
			fmt.Printf("Failed to parse Y: %v\n", err)
			continue
		}
		break
	}

	fmt.Println("Enter Payload (string):")
	payload, _ := reader.ReadString('\n')
	payload = strings.Trim(payload, "\r\n")
	payloadBytes := []byte(payload)

	udpLod.broadcastLoD(int32(xInt), int32(yInt), payloadBytes)
}

func getLoDInfo() {
	udpLod.getLoDInfo()
}

type lod struct {
	udp *udp.Client
}

func setupLoD(c *udp.Client) *lod {
	l := &lod{udp: c}
	l.udp.OnResponse(l.onResponse)
	l.udp.OnPush(l.onPush)
	return l
}

func (l *lod) onResponse(ver uint8, cmd uint16, status uint8, payload []byte) {
	switch {
	case ver == proom.BroadcastLoDVer && cmd == proom.BroadcastLoDCmd:
		if status != uint8(1) {
			fmt.Printf("Broadcast LoD failed: %v\n", string(payload))
			return
		}
		fmt.Printf("Broadcast LoD successful: %v\n", string(payload))

	case ver == proom.GetLoDInfoVer && cmd == proom.GetLoDInfoCmd:
		if status != uint8(1) {
			fmt.Printf("Get LoD Info failed: %v\n", string(payload))
			return
		}
		res := proom.NewGetLoDInfoResponse()
		err := res.Unpack(payload)
		if err != nil {
			fmt.Printf("Failed to unpack get lod info response: %v\n", err)
			return
		}
		fmt.Printf("Get LoD Info successful: %v\n", res)
	}
}

func (l *lod) onPush(ver uint8, cmd uint16, payload []byte) {
	switch {
	case ver == proom.BroadcastLoDPushVer && cmd == proom.BroadcastLoDPushCmd:
		fmt.Printf("Broadcast LoD Push successful: %v\n", string(payload))
	}
}

func (l *lod) broadcastLoD(x int32, y int32, payload []byte) {
	proto := proom.NewBroadcastLoD()
	proto.X = x
	proto.Y = y
	proto.Payload = payload
	l.udp.RSend(proto.Ver, proto.Cmd, proto.Pack())
}

func (l *lod) getLoDInfo() {
	proto := proom.NewGetLoDInfo()
	l.udp.RSend(proto.Ver, proto.Cmd, proto.Pack())
}
