// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Diarkis/diarkis-server-template/testcli/resonance"
	"github.com/Diarkis/diarkis/client/go/test/cli"
)

var (
	tcpResonance *resonance.Resonance
	udpResonance *resonance.Resonance
)

func main() {
	cli.SetupBuiltInCommands()

	// You can add custom commands to the CLI.
	cli.RegisterCommands("test", []cli.Command{
		{CmdName: "resonate", Desc: "Resonate your message", CmdFunc: resonate},
	})

	cli.Connect()

	// setup custom module
	if cli.TCPClient != nil {
		tcpResonance = resonance.SetupAsTCP(cli.TCPClient)
	}
	if cli.UDPClient != nil {
		udpResonance = resonance.SetupAsUDP(cli.UDPClient)
	}

	cli.Run()
}

func resonate() {
	// This is a sample command to add test commands to the CLI.
	fmt.Printf("Which client to join a room? [tcp/udp]")
	client, _ := readLine()

	fmt.Println("Enter the message you want to resonate.")
	message, _ := readLine()

	switch client {
	case "tcp":
		if tcpResonance == nil {
			return
		}
		tcpResonance.Resonate(message)
	case "udp":
		if udpResonance == nil {
			return
		}
		udpResonance.Resonate(message)
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
