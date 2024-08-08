// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"time"
)

func handleOnConnect(bot *bot) {
	botCounter++
	go func() {
		searchAndJoin(bot)
		for {
			time.Sleep(time.Millisecond * time.Duration(packetInterval))
			if bot.isJoined() {
				broadcast(bot)
			}
		}
	}()
}

func handleOnResponse(bot *bot, ver uint8, cmd uint16, status uint8, payload []byte) {
}

func handleOnPush(bot *bot, ver uint8, cmd uint16, payload []byte) {

}

func handleOnDisconnect() {

}
