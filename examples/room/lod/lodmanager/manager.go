// © 2019-2025 Diarkis Inc. All rights reserved.

package lodmanager

import (
	"fmt"
	"maps"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/user"
)

// Manager manages LOD (Level of Detail) synchronization for users
type Manager struct {
	userEntities          map[string]*UserEntity
	started               *atomic.Bool
	ver                   uint8
	cmd                   uint16
	syncIntervalForNearby int32
	syncIntervalForFar    int32
	maxDistanceForNearby  int32
	maxDistanceForFar     int32
	managerMapMutex       sync.Mutex
}

func (m *Manager) String() string {
	return fmt.Sprintf("Manager{userEntities: %v, started: %v, ver: %v, cmd: %v, syncIntervalForNearby: %v, syncIntervalForFar: %v, maxDistanceForNearby: %v, maxDistanceForFar: %v}", m.userEntities, m.started, m.ver, m.cmd, m.syncIntervalForNearby, m.syncIntervalForFar, m.maxDistanceForNearby, m.maxDistanceForFar)
}

// NewManager creates a new LOD manager with the specified configuration
func NewManager(ver uint8, cmd uint16, syncIntervalForNearby int32, syncIntervalForFar int32, maxDistanceForNearby int32, maxDistanceForFar int32) *Manager {
	lm := &Manager{
		ver:                   ver,
		cmd:                   cmd,
		syncIntervalForNearby: syncIntervalForNearby,
		syncIntervalForFar:    syncIntervalForFar,
		maxDistanceForNearby:  maxDistanceForNearby,
		maxDistanceForFar:     maxDistanceForFar,
		userEntities:          make(map[string]*UserEntity),
		started:               &atomic.Bool{},
	}

	if lm.started.CompareAndSwap(false, true) {
		go lm.invokeLodLoop()
		logger.Debugf("NewManager", "syncIntervalForNearby", syncIntervalForNearby, "syncIntervalForFar", syncIntervalForFar, "maxDistanceForNearby", maxDistanceForNearby, "maxDistanceForFar", maxDistanceForFar)
		return lm
	}
	return nil
}

// AddUserEntity adds or updates a user entity with position and payload data
func (m *Manager) AddUserEntity(userID string, x int32, y int32, payload []byte) {
	m.managerMapMutex.Lock()
	defer m.managerMapMutex.Unlock()
	m.userEntities[userID] = NewUserEntity(x, y, payload)
}

// RemoveUserEntity removes a user entity from the manager
// and remove remember data of other user entities
func (m *Manager) RemoveUserEntity(userID string) {
	m.managerMapMutex.Lock()
	defer m.managerMapMutex.Unlock()
	if _, ok := m.userEntities[userID]; !ok {
		return
	}
	delete(m.userEntities, userID)
	// delete from others remember data
	for _, userEntity := range m.userEntities {
		delete(userEntity.Remember, userID)
	}
}

// send packets to nearby users
// 1. nearer than maxDistanceForNearby -> send in every syncIntervalForNearby
// 2. farther than maxDistanceForFar -> don't send
// 3. between maxDistanceForNearby and maxDistanceForFar -> send in every syncIntervalForFar
func (m *Manager) invokeLodLoop() {
	for {
		if !m.started.Load() {
			break
		}
		time.Sleep(time.Duration(m.syncIntervalForNearby) * time.Millisecond)
		start := time.Now()
		m.managerMapMutex.Lock()
		userEntities := maps.Clone(m.userEntities)
		m.managerMapMutex.Unlock()
		for senderUserID, senderUserEntity := range userEntities {
			if senderUserEntity == nil {
				m.RemoveUserEntity(senderUserID)
				continue
			}
			nearbyUserIDs := make([]string, 0, len(userEntities))
			for receiverUserID, receiverUserEntity := range userEntities {
				if senderUserID == receiverUserID {
					continue
				}
				distance := abs(senderUserEntity.X-receiverUserEntity.X) + abs(senderUserEntity.Y-receiverUserEntity.Y)
				// farther than maxDistanceForFar -> don't send
				if distance > m.maxDistanceForFar {
					continue
				}
				// nearer than maxDistanceForNearby -> send in every syncIntervalForNearby
				if distance <= m.maxDistanceForNearby {
					if receiverUserEntity.ChangedAfterSend {
						addToSendList(&nearbyUserIDs, senderUserEntity, receiverUserID)
						logger.Verbosef("invokeLodLoop", "from", senderUserID, "to", receiverUserID, "distance", distance, "mode", "nearby-changed", "interval", m.syncIntervalForNearby)
						continue
					} else {
						if time.Since(getLastSendAt(senderUserEntity, receiverUserID)) > time.Duration(m.syncIntervalForFar)*time.Millisecond {
							addToSendList(&nearbyUserIDs, senderUserEntity, receiverUserID)
							logger.Verbosef("invokeLodLoop", "from", senderUserID, "to", receiverUserID, "distance", distance, "mode", "nearby-interval", "interval", m.syncIntervalForFar)
							continue
						}
					}
				}

				// between maxDistanceForNearby and maxDistanceForFar
				//  -> send in every syncIntervalForFar
				if distance <= m.maxDistanceForFar {
					// syncIntervalForFar and maxDistanceForFar is larger than
					// syncIntervalForNearby and maxDistanceForNearby always
					// cf. func loadLodConfigs()
					if receiverUserEntity.ChangedAfterSend {
						interval := (m.syncIntervalForFar - m.syncIntervalForNearby) * (distance - m.maxDistanceForNearby) / (m.maxDistanceForFar - m.maxDistanceForNearby)
						if time.Since(getLastSendAt(senderUserEntity, receiverUserID)) > time.Duration(interval)*time.Millisecond {
							addToSendList(&nearbyUserIDs, senderUserEntity, receiverUserID)
							logger.Verbosef("invokeLodLoop", "from", senderUserID, "to", receiverUserID, "distance", distance, "mode", "far-changed", "interval", interval)
						}
					} else {
						if time.Since(getLastSendAt(senderUserEntity, receiverUserID)) > time.Duration(m.syncIntervalForFar)*time.Millisecond {
							addToSendList(&nearbyUserIDs, senderUserEntity, receiverUserID)
							logger.Verbosef("invokeLodLoop", "from", senderUserID, "to", receiverUserID, "distance", distance, "mode", "far-interval", "interval", m.syncIntervalForFar)
						}
					}
				}
			}
			if len(nearbyUserIDs) > 0 {
				for _, receiverUserID := range nearbyUserIDs {
					receiverUser := user.GetUserBySID(receiverUserID)
					if receiverUser != nil {
						receiverUser.PushToClient(m.ver, m.cmd, senderUserEntity.Payload, packet.Unreliable)
					}
				}
			}

		}
		logger.Debugf("invokeLodLoop", "time", time.Since(start))
	}
}

func addToSendList(nearbyUserIDs *[]string, senderUserEntity *UserEntity, receiverUserID string) {
	*nearbyUserIDs = append(*nearbyUserIDs, receiverUserID)
	senderUserEntity.Remember[receiverUserID] = time.Now()
	senderUserEntity.ChangedAfterSend = false
}

func getLastSendAt(senderUserEntity *UserEntity, receiverUserID string) time.Time {
	return senderUserEntity.Remember[receiverUserID]
}
