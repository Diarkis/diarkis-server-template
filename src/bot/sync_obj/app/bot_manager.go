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

func (b *botManager) getSyncSendCntTotal() int64 {
	var sum int64
	for _, bot := range b.bots {
		sum += bot.syncSendCnt.Load()
	}
	return sum
}

func (b *botManager) getSyncReceiveCntTotal() int64 {
	var sum int64
	for _, bot := range b.bots {
		sum += bot.syncRcvCnt.Load()
	}
	return sum
}

func (b *botManager) getSyncSendAvg() int64 {
	if len(b.bots) == 0 {
		return 0
	}
	return b.getSyncSendCntTotal() / int64(len(b.bots))
}

func (b *botManager) getSyncReceiveAvg() int64 {
	if len(b.bots) == 0 {
		return 0
	}
	return b.getSyncReceiveCntTotal() / int64(len(b.bots))
}

func (b *botManager) resetCnt() {
	for _, bot := range b.bots {
		bot.syncRcvCnt.Store(0)
		bot.syncSendCnt.Store(0)
	}
}
