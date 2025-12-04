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
	// Handle LOD broadcast push notifications
	if ver == 2 && cmd == 1001 { // BroadcastLoDPush
		broadcastReceiveCnt.Add(1)
		bot.broadcastRcvCnt.Add(1)
		logger.Debug("LOD broadcast received",
			"bot.uid", bot.uid,
			"payloadSize", len(payload))
	}
}

func handleOnDisconnect() {

}
