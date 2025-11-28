package main

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/user"
)

var (
	lodManagerMap      = make(map[string]*lodManager)
	lodManagerMapMutex = sync.RWMutex{}
)

func getRoomLodManager(roomID string) *lodManager {
	lodManagerMapMutex.RLock()
	defer lodManagerMapMutex.RUnlock()
	return lodManagerMap[roomID]
}

func setRoomLodManager(roomID string, lodManager *lodManager) {
	lodManagerMapMutex.Lock()
	defer lodManagerMapMutex.Unlock()
	lodManagerMap[roomID] = lodManager
}

type lodManager struct {
	userEntities          map[string]*userEntity
	started               *atomic.Bool
	ver                   uint8
	cmd                   uint16
	syncIntervalForNearby int32
	syncIntervalForFar    int32
	maxDistanceForNearby  int32
	maxDistanceForFar     int32
}

type userEntity struct {
	X          int32
	Y          int32
	Payload    []byte
	LastSyncAt time.Time
}

func newLodManager(ver uint8, cmd uint16, syncIntervalForNearby int32, syncIntervalForFar int32, maxDistanceForNearby int32, maxDistanceForFar int32) *lodManager {
	lm := &lodManager{ver: ver, cmd: cmd, syncIntervalForNearby: syncIntervalForNearby,
		syncIntervalForFar: syncIntervalForFar, maxDistanceForNearby: maxDistanceForNearby,
		maxDistanceForFar: maxDistanceForFar, userEntities: make(map[string]*userEntity),
	}
	lm.started.Store(false)
	go lm.invokeLodLoop()
	return lm

}

func (m *lodManager) addUserEntity(userID string, x int32, y int32, payload []byte) {
	m.userEntities[userID] = &userEntity{X: x, Y: y, Payload: payload, LastSyncAt: time.Now()}
}

func (m *lodManager) removeUserEntity(userID string) {
	delete(m.userEntities, userID)
}

func (m *lodManager) invokeLodLoop() {
	for {
		time.Sleep(time.Duration(SyncIntervalForNearby) * time.Millisecond)

		for userID, userEntity := range m.userEntities {
			nearbyUserIDs := make([]string, 0, len(m.userEntities))
			for otherUserID, otherUserEntity := range m.userEntities {
				if userID == otherUserID {
					continue
				}
				distance := abs(userEntity.X-otherUserEntity.X) + abs(userEntity.Y-otherUserEntity.Y)
				if distance < m.maxDistanceForNearby && time.Since(userEntity.LastSyncAt) > time.Duration(m.syncIntervalForNearby)*time.Millisecond {
					nearbyUserIDs = append(nearbyUserIDs, otherUserID)
					userEntity.LastSyncAt = time.Now()
				}
			}
			if len(nearbyUserIDs) > 0 {
				user := user.GetUserBySID(userID)
				if user != nil {
					user.PushToClient(m.ver, m.cmd, userEntity.Payload, packet.Unreliable)
				}
			}
		}
	}
}

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
