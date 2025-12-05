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

// sendMessage represents a message to be sent to a user
type sendMessage struct {
	receiverUserID string
	ver            uint8
	cmd            uint16
	payload        []byte
}

// Manager manages LOD (Level of Detail) synchronization for users
type Manager struct {
	userEntities          map[string]*UserEntity
	started               atomic.Bool
	ver                   uint8
	cmd                   uint16
	syncIntervalForNearby int32
	syncIntervalForFar    int32
	maxDistanceForNearby  int32
	maxDistanceForFar     int32
	managerMapMutex       sync.RWMutex
	sendBuffer            chan sendMessage
	senderWg              sync.WaitGroup
}

func (m *Manager) String() string {
	m.managerMapMutex.RLock()
	defer m.managerMapMutex.RUnlock()
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
		started:               atomic.Bool{},
		sendBuffer:            make(chan sendMessage, 10000), // Buffer size: 10000
	}

	lm.started.Store(true)

	// Start send worker goroutine
	lm.senderWg.Add(1)
	go lm.sendWorker()

	// Start LOD loop
	go lm.invokeLodLoop()

	logger.Debugf("NewManager", "syncIntervalForNearby", syncIntervalForNearby, "syncIntervalForFar", syncIntervalForFar, "maxDistanceForNearby", maxDistanceForNearby, "maxDistanceForFar", maxDistanceForFar)
	return lm
}

// AddUserEntity adds or updates a user entity with position and payload data
func (m *Manager) AddUserEntity(userID string, x int32, y int32, payload []byte) {
	m.managerMapMutex.Lock()
	defer m.managerMapMutex.Unlock()
	if m.userEntities[userID] != nil {
		userEntity := m.userEntities[userID]
		userEntity.Payload = payload
		userEntity.X = x
		userEntity.Y = y
		userEntity.ChangedAfterSend = true
		return
	}
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
		userEntity.RememberMutex.Lock()
		delete(userEntity.Remember, userID)
		userEntity.RememberMutex.Unlock()
	}
}

// sendWorker processes messages from sendBuffer and sends them to users
func (m *Manager) sendWorker() {
	defer m.senderWg.Done()

	for msg := range m.sendBuffer {
		receiverUser := user.GetUserBySID(msg.receiverUserID)
		if receiverUser != nil {
			receiverUser.PushToClient(msg.ver, msg.cmd, msg.payload, packet.Unreliable)
		}
	}
	logger.Debugf("sendWorker stopped")
}

// Stop stops the LOD manager and waits for send worker to finish
func (m *Manager) Stop() {
	m.started.Store(false)
	m.senderWg.Wait()
	logger.Debugf("Manager stopped")
}

// shouldSendUpdate determines if an update should be sent to a receiver
// Logs the reason when returning true
func (m *Manager) shouldSendUpdate(
	senderID string,
	senderEntity *UserEntity,
	receiverID string,
	receiverEntity *UserEntity,
	distance int32,
) bool {
	// farther than maxDistanceForFar -> don't send
	if distance > m.maxDistanceForFar {
		return false
	}

	// nearer than maxDistanceForNearby -> send based on nearby rules
	if distance <= m.maxDistanceForNearby {
		if senderEntity.ChangedAfterSend {
			logger.Verbosef("shouldSendUpdate", "from", senderID, "to", receiverID, "distance", distance, "mode", "nearby-changed", "interval", m.syncIntervalForNearby)
			return true
		}

		if time.Since(getLastSendAt(senderEntity, receiverID)) > time.Duration(m.syncIntervalForFar)*time.Millisecond {
			logger.Verbosef("shouldSendUpdate", "from", senderID, "to", receiverID, "distance", distance, "mode", "nearby-unchanged", "interval", m.syncIntervalForFar)
			return true
		}
		return false
	}

	// between maxDistanceForNearby and maxDistanceForFar -> send based on far rules
	if distance <= m.maxDistanceForFar {
		// syncIntervalForFar and maxDistanceForFar is larger than
		// syncIntervalForNearby and maxDistanceForNearby always
		// cf. func loadLodConfigs()
		if senderEntity.ChangedAfterSend {
			interval := (m.syncIntervalForFar - m.syncIntervalForNearby) * (distance - m.maxDistanceForNearby) / (m.maxDistanceForFar - m.maxDistanceForNearby)
			if time.Since(getLastSendAt(senderEntity, receiverID)) > time.Duration(interval)*time.Millisecond {
				logger.Verbosef("shouldSendUpdate", "from", senderID, "to", receiverID, "distance", distance, "mode", "far-changed", "interval", interval, "since", time.Since(getLastSendAt(senderEntity, receiverID)))
				return true
			}
		} else {
			if time.Since(getLastSendAt(senderEntity, receiverID)) > time.Duration(m.syncIntervalForFar)*time.Millisecond {
				logger.Verbosef("shouldSendUpdate", "from", senderID, "to", receiverID, "distance", distance, "mode", "far-unchanged", "interval", m.syncIntervalForFar, "since", time.Since(getLastSendAt(senderEntity, receiverID)))
				return true
			}
		}
	}

	return false
}

// processSenderReceiverPair processes a single sender-receiver pair
// Returns the receiver ID if a message should be sent, empty string otherwise
func (m *Manager) processSenderReceiverPair(
	senderID string,
	senderEntity *UserEntity,
	receiverID string,
	receiverEntity *UserEntity,
) string {
	if senderID == receiverID {
		return ""
	}

	distance := calculateDistance(senderEntity.X, senderEntity.Y, receiverEntity.X, receiverEntity.Y)
	shouldSend := m.shouldSendUpdate(senderID, senderEntity, receiverID, receiverEntity, distance)

	if shouldSend {
		return receiverID
	}

	return ""
}

// processSingleSender processes all receivers for one sender
func (m *Manager) processSingleSender(
	senderUserID string,
	senderUserEntity *UserEntity,
	userEntities map[string]*UserEntity,
) {
	nearbyUserIDs := make([]string, 0, len(userEntities))

	for receiverUserID, receiverUserEntity := range userEntities {
		receiverID := m.processSenderReceiverPair(senderUserID, senderUserEntity, receiverUserID, receiverUserEntity)
		if receiverID != "" {
			nearbyUserIDs = addToSendList(nearbyUserIDs, senderUserEntity, receiverID)
		}
	}
	// send messages to packet 	sender worker
	if len(nearbyUserIDs) > 0 {
		for _, receiverUserID := range nearbyUserIDs {
			select {
			case m.sendBuffer <- sendMessage{
				receiverUserID: receiverUserID,
				ver:            m.ver,
				cmd:            m.cmd,
				payload:        senderUserEntity.Payload,
			}:
			default:
				logger.Warnf("sendBuffer full, dropping message", "receiverUserID", receiverUserID, "senderUserID", senderUserID)
			}
		}
	}
}

// processAllUsers processes all users in the manager
func (m *Manager) processAllUsers() {
	//start := time.Now()
	m.managerMapMutex.RLock()
	userEntities := maps.Clone(m.userEntities)
	m.managerMapMutex.RUnlock()

	for senderUserID, senderUserEntity := range userEntities {
		m.processSingleSender(senderUserID, senderUserEntity, userEntities)
	}

	//logger.Verbosef("invokeLodLoop", "time", time.Since(start))
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
		m.processAllUsers()
	}
	close(m.sendBuffer)
}

// calculateDistance computes Manhattan distance between two entities
func calculateDistance(x1, y1, x2, y2 int32) int32 {
	return abs(x1-x2) + abs(y1-y2)
}

func addToSendList(nearbyUserIDs []string, senderUserEntity *UserEntity, receiverUserID string) []string {
	nearbyUserIDs = append(nearbyUserIDs, receiverUserID)
	senderUserEntity.Remember[receiverUserID] = time.Now()
	senderUserEntity.ChangedAfterSend = false
	return nearbyUserIDs
}

func getLastSendAt(senderUserEntity *UserEntity, receiverUserID string) time.Time {
	senderUserEntity.RememberMutex.RLock()
	defer senderUserEntity.RememberMutex.RUnlock()
	return senderUserEntity.Remember[receiverUserID]
}
