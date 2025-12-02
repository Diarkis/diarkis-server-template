package lodmanager

import (
	"fmt"
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
	cache                 *cache // TODO: is this needed?
}

// UserEntity represents a user's position and payload data
type UserEntity struct {
	X                int32
	Y                int32
	Payload          []byte
	LastSendAt       time.Time // TODO: this needs to store time for each user, so actually, Map is correct.
	ChangedAfterSend bool
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
	lm.started.Store(true)
	logger.Debugf("NewManager", "syncIntervalForNearby", syncIntervalForNearby, "syncIntervalForFar", syncIntervalForFar, "maxDistanceForNearby", maxDistanceForNearby, "maxDistanceForFar", maxDistanceForFar)
	go lm.invokeLodLoop()
	return lm
}

// AddUserEntity adds or updates a user entity with position and payload data
func (m *Manager) AddUserEntity(userID string, x int32, y int32, payload []byte) {
	m.userEntities[userID] = &UserEntity{X: x, Y: y, Payload: payload, LastSendAt: time.Now(), ChangedAfterSend: true}
}

// RemoveUserEntity removes a user entity from the manager
func (m *Manager) RemoveUserEntity(userID string) {
	delete(m.userEntities, userID)
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

		for senderUserID, senderUserEntity := range m.userEntities {
			if senderUserEntity == nil {
				m.RemoveUserEntity(senderUserID)
				continue
			}
			nearbyUserIDs := make([]string, 0, len(m.userEntities))
			for receiverUserID, receiverUserEntity := range m.userEntities {
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
						logger.Debugf("invokeLodLoop", "from", senderUserID, "to", receiverUserID, "distance", distance, "interval", m.syncIntervalForNearby)
						continue
					} else {
						if time.Since(senderUserEntity.LastSendAt) > time.Duration(m.syncIntervalForFar)*time.Millisecond {
							addToSendList(&nearbyUserIDs, senderUserEntity, receiverUserID)
							logger.Debugf("invokeLodLoop", "from", senderUserID, "to", receiverUserID, "distance", distance, "interval", m.syncIntervalForFar)
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
					interval := (m.syncIntervalForFar - m.syncIntervalForNearby) * (distance - m.maxDistanceForNearby) / (m.maxDistanceForFar - m.maxDistanceForNearby)
					if time.Since(senderUserEntity.LastSendAt) > time.Duration(interval)*time.Millisecond {
						addToSendList(&nearbyUserIDs, senderUserEntity, receiverUserID)
						logger.Debugf("invokeLodLoop", "from", senderUserID, "to", receiverUserID, "distance", distance, "interval", interval)
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
	}
}

func addToSendList(nearbyUserIDs *[]string, senderUserEntity *UserEntity, receiverUserID string) {
	*nearbyUserIDs = append(*nearbyUserIDs, receiverUserID)
	senderUserEntity.LastSendAt = time.Now()
	senderUserEntity.ChangedAfterSend = false
}
