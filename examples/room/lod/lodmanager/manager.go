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
	mu                    sync.RWMutex
	sendBuffer            chan sendMessage
	senderWg              sync.WaitGroup
}

func (m *Manager) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("Manager{userEntities: %v, started: %v, ver: %v, cmd: %v, syncIntervalForNearby: %v, syncIntervalForFar: %v,	 maxDistanceForNearby: %v, maxDistanceForFar: %v}",
		m.userEntities, m.started.Load(), m.ver, m.cmd, m.syncIntervalForNearby, m.syncIntervalForFar, m.maxDistanceForNearby, m.maxDistanceForFar)
}

// NewManager creates a new LOD manager with the specified configuration
func NewManager(ver uint8, cmd uint16, syncIntervalForNearby int32, syncIntervalForFar int32, maxDistanceForNearby int32, maxDistanceForFar int32) *Manager {
	// Initialize metrics on first manager creation

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

	// Increment instance count
	instancesGauge.Inc()

	// Start send worker goroutine
	lm.senderWg.Add(1)
	go lm.sendWorker()

	// Start LOD loop
	go lm.invokeLodLoop()

	return lm
}

// AddUserEntity adds or updates a user entity with position and payload data
func (m *Manager) AddUserEntity(userID string, x int32, y int32, payload []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if userEntity, ok := m.userEntities[userID]; ok { // update
		userEntity.payload = payload
		userEntity.x = x
		userEntity.y = y
		userEntity.changedAfterSend = true
		return
	}
	m.userEntities[userID] = NewUserEntity(x, y, payload)
	// Update active users gauge
	activeUsersGauge.Set(int64(len(m.userEntities)))
}

// RemoveUserEntity removes a user entity from the manager
// and remove remember data of other user entities
func (m *Manager) RemoveUserEntity(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.userEntities[userID]; !ok {
		return
	}
	delete(m.userEntities, userID)
	// delete from others remember data
	for _, userEntity := range m.userEntities {
		userEntity.mu.Lock()
		delete(userEntity.m, userID)
		userEntity.mu.Unlock()
	}
	// Update active users gauge
	activeUsersGauge.Set(int64(len(m.userEntities)))
}

// sendWorker processes messages from sendBuffer and sends them to users
func (m *Manager) sendWorker() {
	defer m.senderWg.Done()

	for msg := range m.sendBuffer {
		receiverUser := user.GetUserBySID(msg.receiverUserID)
		if receiverUser != nil {
			receiverUser.PushToClient(msg.ver, msg.cmd, msg.payload, packet.Unreliable)
			// Track successfully sent messages
			messagesSentCounter.Inc()
		} else {
			// Track errors
			errorsCounter.WithLabelValues("user_not_found").Inc()
		}
	}
}

// Stop stops the LOD manager and waits for send worker to finish
func (m *Manager) Stop() {
	m.started.Store(false)
	m.senderWg.Wait()
	instancesGauge.Dec()
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
		updatesSkippedCounter.WithLabelValues("too_far").Inc()
		return false
	}

	// nearer than maxDistanceForNearby -> send based on nearby rules
	if distance <= m.maxDistanceForNearby {
		if senderEntity.changedAfterSend {
			updatesSentCounter.WithLabelValues("nearby", "changed").Inc()
			return true
		}

		if time.Since(getLastSendAt(senderEntity, receiverID)) > time.Duration(m.syncIntervalForFar)*time.Millisecond {
			updatesSentCounter.WithLabelValues("nearby", "interval").Inc()
			return true
		}
		updatesSkippedCounter.WithLabelValues("interval_not_met").Inc()
		return false
	}

	// between maxDistanceForNearby and maxDistanceForFar -> send based on far rules
	if distance <= m.maxDistanceForFar {
		// syncIntervalForFar and maxDistanceForFar is larger than
		// syncIntervalForNearby and maxDistanceForNearby always
		// cf. func loadLodConfigs()
		if senderEntity.changedAfterSend {
			interval := (m.syncIntervalForFar - m.syncIntervalForNearby) * (distance - m.maxDistanceForNearby) / (m.maxDistanceForFar - m.maxDistanceForNearby)
			if time.Since(getLastSendAt(senderEntity, receiverID)) > time.Duration(interval)*time.Millisecond {
				updatesSentCounter.WithLabelValues("far", "changed").Inc()
				return true
			}
		} else {
			if time.Since(getLastSendAt(senderEntity, receiverID)) > time.Duration(m.syncIntervalForFar)*time.Millisecond {
				updatesSentCounter.WithLabelValues("far", "interval").Inc()
				return true
			}
		}
		updatesSkippedCounter.WithLabelValues("interval_not_met").Inc()
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
		updatesSkippedCounter.WithLabelValues("same_user").Inc()
		return ""
	}

	distance := calculateDistance(senderEntity.x, senderEntity.y, receiverEntity.x, receiverEntity.y)
	userDistancesHistogram.Observe(float64(distance))

	shouldSend := m.shouldSendUpdate(senderID, senderEntity, receiverID, receiverEntity, distance)

	if shouldSend {
		return receiverID
	}

	return ""
}

// processSingleReceiver processes all senders for one receiver
func (m *Manager) processSingleReceiver(
	receiverUserID string,
	receiverUserEntity *UserEntity,
	userEntities map[string]*UserEntity,
) {
	for senderUserID, senderUserEntity := range userEntities {
		targetID := m.processSenderReceiverPair(senderUserID, senderUserEntity, receiverUserID, receiverUserEntity)
		if targetID != "" {
			senderUserEntity.mu.Lock()
			senderUserEntity.m[receiverUserID] = time.Now()
			senderUserEntity.mu.Unlock()

			select {
			case m.sendBuffer <- sendMessage{
				receiverUserID: receiverUserID,
				ver:            m.ver,
				cmd:            m.cmd,
				payload:        senderUserEntity.payload,
			}:
			default:
				logger.Warnf("sendBuffer full, dropping message", "receiverUserID", receiverUserID, "senderUserID", senderUserID)
				bufferDropsCounter.Inc()
			}
		}
	}
}

// processAllUsers processes all users in the manager
func (m *Manager) processAllUsers() {
	start := time.Now()
	m.mu.RLock()
	userEntities := maps.Clone(m.userEntities)
	m.mu.RUnlock()

	for receiverID, receiverEntity := range userEntities {
		m.processSingleReceiver(receiverID, receiverEntity, userEntities)
	}

	for _, userEntity := range userEntities {
		userEntity.changedAfterSend = false
	}

	duration := time.Since(start)

	loopDurationHistogram.Observe(duration.Seconds())
	bufferSizeGauge.Set(int64(len(m.sendBuffer)))
}

// send packets to nearby users
// 1. nearer than maxDistanceForNearby -> send in every syncIntervalForNearby
// 2. farther than maxDistanceForFar -> don't send
// 3. between maxDistanceForNearby and maxDistanceForFar -> send in every syncIntervalForFar
func (m *Manager) invokeLodLoop() {
	tick := time.NewTicker(time.Duration(m.syncIntervalForNearby) * time.Millisecond)
	defer tick.Stop()
	for {
		if !m.started.Load() {
			break
		}
		<-tick.C
		m.processAllUsers()
	}
	close(m.sendBuffer)
}

// calculateDistance computes Manhattan distance between two entities
func calculateDistance(x1, y1, x2, y2 int32) int32 {
	return abs(x1-x2) + abs(y1-y2)
}

func getLastSendAt(senderUserEntity *UserEntity, receiverUserID string) time.Time {
	senderUserEntity.mu.RLock()
	defer senderUserEntity.mu.RUnlock()
	return senderUserEntity.m[receiverUserID]
}
