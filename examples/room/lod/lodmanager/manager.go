package lodmanager

import (
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
}

// UserEntity represents a user's position and payload data
type UserEntity struct {
	X          int32
	Y          int32
	Payload    []byte
	LastSyncAt time.Time
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
	go lm.invokeLodLoop()
	return lm
}

// AddUserEntity adds or updates a user entity with position and payload data
func (m *Manager) AddUserEntity(userID string, x int32, y int32, payload []byte) {
	m.userEntities[userID] = &UserEntity{X: x, Y: y, Payload: payload, LastSyncAt: time.Now()}
}

// RemoveUserEntity removes a user entity from the manager
func (m *Manager) RemoveUserEntity(userID string) {
	delete(m.userEntities, userID)
}

func (m *Manager) invokeLodLoop() {
	for {
		if !m.started.Load() {
			break
		}
		time.Sleep(time.Duration(m.syncIntervalForNearby) * time.Millisecond)
		for userID, userEntity := range m.userEntities {
			if userEntity == nil {
				m.RemoveUserEntity(userID)
				continue
			}
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
