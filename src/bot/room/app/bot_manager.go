// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import "slices"

type botManager struct {
	bots []*bot
}

func (b *botManager) getJoinedCnt() int {
	var sum int
	for _, bot := range b.bots {
		if bot.room != nil && bot.room.ID != "" {
			sum++
		}
	}
	return sum
}

func (b *botManager) getJoinedRooms() []string {
	var rooms []string
	for _, bot := range b.bots {
		if bot.room != nil && bot.room.ID != "" {
			if !slices.Contains(rooms, bot.room.ID) {
				rooms = append(rooms, bot.room.ID)
			}
		}
	}
	return rooms
}

func (b *botManager) getBroadcastSendCntTotal() int64 {
	var sum int64
	for _, bot := range b.bots {
		sum += bot.broadcastSendCnt.Load()
	}
	return sum
}

func (b *botManager) getBroadcastReceiveCntTotal() int64 {
	var sum int64
	for _, bot := range b.bots {
		sum += bot.broadcastRcvCnt.Load()
	}
	return sum
}

func (b *botManager) getBroadcastSendAvg() int64 {
	return b.getBroadcastSendCntTotal() / int64(len(b.bots))
}

func (b *botManager) getBroadcastReceiveAvg() int64 {
	return b.getBroadcastReceiveCntTotal() / int64(len(b.bots))
}

func (b *botManager) resetCnt() {
	for _, bot := range b.bots {
		bot.broadcastSendCnt.Store(0)
		bot.broadcastRcvCnt.Store(0)
	}
}
