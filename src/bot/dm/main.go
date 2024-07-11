// © 2019-2024 Diarkis Inc. All rights reserved.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Diarkis/diarkis/util"

	"github.com/Diarkis/diarkis-server-template/bot/dm/loadtest"
	"github.com/Diarkis/diarkis-server-template/bot/dm/parameters"
)

func main() {
	if len(os.Args) < 5 {
		msg := "Bot requires 5 parameters:"
		params := "host=$(host:port) bots=$(how many bots) protocol=$(UDP|TCP) size=$(packet size) interval=$(send message interval in milliseconds)"
		fmt.Println(msg, params)
		os.Exit(1)
		return
	}

	params := parameters.ParseParams(os.Args)

	if params == nil {
		fmt.Println("Invalid parameters...", os.Args)
		os.Exit(1)
	}

	// spawn all bot clients
	for i := 0; i < params.Howmany; i++ {
		loadtest.Spawn(params)
		time.Sleep(time.Millisecond * 100)
	}

	// once we finish spawning all bot clients,
	// we start DM spam for load test
	loadtest.StartLoadTest(params)

	fmt.Println("Bot is working hard...")

	// this loop keeps the bot alive
	for {
		time.Sleep(time.Second * 30)
		sent, received := loadtest.GetReport()
		now := time.Now()
		fmt.Println(util.ZuluTimeFormat(now), "Sent messages", sent, "Received messages", received)
	}
}
