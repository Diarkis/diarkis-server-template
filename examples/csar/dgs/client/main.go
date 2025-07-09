// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Diarkis/diarkis-server-template/examples/csar/dgs/client/dgs"
	"github.com/Diarkis/diarkis/client/go/test/cli"
)

var (
	tcpDGSClient *dgs.DGS
	udpDGSClient *dgs.DGS
)

func main() {
	cli.SetupBuiltInCommands()

	// You can add custom commands to the CLI.
	cli.RegisterCommands("test", []cli.Command{
		{CmdName: "dgs create", Desc: "DGS create from room", CmdFunc: dgsCreate},
		{CmdName: "dgs backfill", Desc: "DGS backfill from room", CmdFunc: dgsBackfill},
	})

	cli.Connect()

	// setup custom module
	if cli.TCPClient != nil {
		tcpDGSClient = dgs.SetupAsTCP(cli.TCPClient)
	}
	if cli.UDPClient != nil {
		udpDGSClient = dgs.SetupAsUDP(cli.UDPClient)
	}

	cli.Run()
}

func dgsCreate() {
	// This is a sample command to add test commands to the CLI.
	fmt.Printf("Which client to create DGS? [tcp/udp]")
	client, _ := readLine()

	switch client {
	case "tcp":
		if tcpDGSClient == nil {
			return
		}
		tcpDGSClient.CreateSessionFromRoom()
	case "udp":
		if udpDGSClient == nil {
			return
		}
		udpDGSClient.CreateSessionFromRoom()
	default:
		fmt.Println("Invalid input. Please provide tcp or udp.")
	}
}

func dgsBackfill() {
	fmt.Printf("Which client to backfill DGS? [tcp/udp]")
	client, _ := readLine()

	switch client {
	case "tcp":
		if tcpDGSClient == nil {
			return
		}
		tcpDGSClient.BackfillSessionFromRoom()
	case "udp":
		if udpDGSClient == nil {
			return
		}
		udpDGSClient.BackfillSessionFromRoom()
	default:
		fmt.Println("Invalid input. Please provide tcp or udp.")
	}
}

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	str = strings.TrimSuffix(str, "\n")
	// On Windows you need to strip the \r control character too.
	// See issue #2851
	str = strings.TrimSuffix(str, "\r")

	return str, nil
}
