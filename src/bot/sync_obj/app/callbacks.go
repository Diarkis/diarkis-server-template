// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"time"
)

func handleOnConnect(bot *bot) {
	go func() {
		searchAndJoin(bot)
		for {
			time.Sleep(time.Millisecond * time.Duration(*packetInterval))
			if bot.isJoined() {
				updateObj(bot)
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

func handleOnObjectUpdateResponse(bot *bot, payload []byte) {
	bot.syncSendCnt.Add(1)
	logger.Debug("onObjectUpdateResponse", "bot.uid", bot.uid, "payload", payload)
}

func handleOnObjectUpdatePush(bot *bot, status uint8, payload string, obj map[string]float64) {
	bot.syncRcvCnt.Add(1)
	logger.Debug("onObjectUpdatePush", "bot.uid", bot.uid, "status", status, "payload", payload, "obj", obj)
}
