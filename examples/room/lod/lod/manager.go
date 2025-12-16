// © 2019-2025 Diarkis Inc. All rights reserved.

package lod

import (
	"fmt"
	"maps"
	"sync"
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
	userEntities           map[string]*UserEntity
	roomID                 string
	ver                    uint8
	cmd                    uint16
	syncIntervalForNearby  time.Duration
	syncIntervalForFar     time.Duration
	syncIntervalProportion int64
	maxDistanceForNearby   int32
	maxDistanceForFar      int32
	mu                     sync.RWMutex
	sendBuffer             chan sendMessage
	sendDone               chan struct{} // channel to signal send worker to stop
	lodDone                chan struct{} // channel to signal lod worker to stop
}

func (m *Manager) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("Manager{roomID: %v, ver: %v, cmd: %v, syncIntervalForNearby: %v, syncIntervalForFar: %v, syncIntervalProportion: %v, maxDistanceForNearby: %v, maxDistanceForFar: %v}",
		m.roomID, m.ver, m.cmd, m.syncIntervalForNearby, m.syncIntervalForFar, m.syncIntervalProportion, m.maxDistanceForNearby, m.maxDistanceForFar)
}

// NewManager creates a new LOD manager with the specified configuration
func NewManager(roomID string, ver uint8, cmd uint16, syncIntervalForNearby time.Duration, syncIntervalForFar time.Duration, maxDistanceForNearby int32, maxDistanceForFar int32) *Manager {
	// Initialize metrics on first manager creation
	logger.Info("NewManager", "roomID", roomID, "ver", ver, "cmd", cmd, "syncIntervalForNearby", syncIntervalForNearby, "syncIntervalForFar", syncIntervalForFar, "maxDistanceForNearby", maxDistanceForNearby, "maxDistanceForFar", maxDistanceForFar)
	lm := &Manager{
		roomID:                 roomID,
		ver:                    ver,
		cmd:                    cmd,
		syncIntervalForNearby:  syncIntervalForNearby,
		syncIntervalForFar:     syncIntervalForFar,
		maxDistanceForNearby:   maxDistanceForNearby,
		maxDistanceForFar:      maxDistanceForFar,
		userEntities:           make(map[string]*UserEntity),
		syncIntervalProportion: int64(syncIntervalForFar-syncIntervalForNearby) / int64(maxDistanceForFar-maxDistanceForNearby),
		sendBuffer:             make(chan sendMessage, 10000), // Buffer size: 10000
		sendDone:               make(chan struct{}),
		lodDone:                make(chan struct{}),
	}

	// Increment instance count
	instancesGauge.Inc()

	// Start send worker goroutine
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
	activeUsersGauge.WithLabelValues(m.roomID).Set(int64(len(m.userEntities)))
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
	activeUsersGauge.WithLabelValues(m.roomID).Set(int64(len(m.userEntities)))
}

// sendWorker processes messages from sendBuffer and sends them to users
func (m *Manager) sendWorker() {
	for {
		select {
		case msg := <-m.sendBuffer:
			receiverUser := user.GetUserBySID(msg.receiverUserID)
			if receiverUser != nil {
				receiverUser.PushToClient(msg.ver, msg.cmd, msg.payload, packet.Unreliable)
				// Track successfully sent messages
				messagesSentCounter.Inc()
			} else {
				// Track errors
				errorsCounter.WithLabelValues("user_not_found").Inc()
			}
		case <-m.sendDone:
			logger.Info("LOD send worker stopped")
			return
		}
	}
}

// Stop stops the LOD manager and waits for send worker to finish
func (m *Manager) Stop() {
	logger.Info("Stopping LOD manager")
	close(m.sendDone)
	close(m.lodDone)
	instancesGauge.Dec()
}

// shouldSendUpdate determines if an update should be sent to a receiver
// Logs the reason when returning true
// shouldSendUpdate returns false when update should not be sent
// shouldSendUpdate returns true when update should be sent
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
	// if userEntity doesn't changed after send in previous interval, send update in syncIntervalForFar
	// if userEntity changed after send in previous interval, send update in syncIntervalForNearby
	if distance <= m.maxDistanceForNearby {
		return m.checkForNearby(senderID, senderEntity, receiverID, receiverEntity)
	}

	// between maxDistanceForNearby and maxDistanceForFar -> send based on far rules
	// if userEntity doesn't changed(changedAfterSend == false) after send in previous interval,
	//  send update in syncIntervalForFar
	// if userEntity changed(changedAfterSend == true) after send in previous interval,
	//  send update in dynamic interval
	// dynamic interval is calculated based on distance
	// dynamic interval is between syncIntervalForFar and syncIntervalForNearby
	if distance <= m.maxDistanceForFar {
		return m.checkForBetweenNearbyAndFar(senderID, senderEntity, receiverID, receiverEntity, distance)
	}

	return false
}

func (m *Manager) checkForNearby(
	senderID string,
	senderEntity *UserEntity,
	receiverID string,
	receiverEntity *UserEntity,
) bool {

	// if userEntity changed after send in previous interval, send update in syncIntervalForNearby
	if senderEntity.changedAfterSend {
		if getIntervalFromLastSend(senderEntity, receiverID) > m.syncIntervalForNearby {
			logger.Verbosef("checkForNearby", "senderID", senderID, "receiverID", receiverID, "intervalFromLastSend", getIntervalFromLastSend(senderEntity, receiverID), "syncIntervalForNearby", m.syncIntervalForNearby, "changed", true, "reason", "changed_after_send")
			updatesSentCounter.WithLabelValues("nearby", "changed").Inc()
			return true
		}
		logger.Verbosef("checkForNearby", "senderID", senderID, "receiverID", receiverID, "intervalFromLastSend", getIntervalFromLastSend(senderEntity, receiverID), "syncIntervalForNearby", m.syncIntervalForNearby, "changed", true, "reason", "interval_not_met")
		updatesSkippedCounter.WithLabelValues("interval_not_met").Inc()
		return false
	}

	// if userEntity doesn't changed after send in previous interval, send update in syncIntervalForFar
	if getIntervalFromLastSend(senderEntity, receiverID) > m.syncIntervalForFar {
		logger.Verbosef("checkForNearby", "senderID", senderID, "receiverID", receiverID, "intervalFromLastSend", getIntervalFromLastSend(senderEntity, receiverID), "syncIntervalForFar", m.syncIntervalForFar, "changed", false, "reason", "interval_met")
		updatesSentCounter.WithLabelValues("nearby", "interval").Inc()
		return true
	}
	logger.Verbosef("checkForNearby", "senderID", senderID, "receiverID", receiverID, "intervalFromLastSend", getIntervalFromLastSend(senderEntity, receiverID), "syncIntervalForFar", m.syncIntervalForFar, "changed", false, "reason", "interval_not_met")
	updatesSkippedCounter.WithLabelValues("interval_not_met").Inc()
	return false
}

func (m *Manager) checkForBetweenNearbyAndFar(
	senderID string,
	senderEntity *UserEntity,
	receiverID string,
	receiverEntity *UserEntity,
	distance int32,
) bool {
	interval := time.Duration(m.syncIntervalProportion*int64(distance-m.maxDistanceForNearby)) + m.syncIntervalForNearby
	// syncIntervalForFar and maxDistanceForFar is larger than
	// syncIntervalForNearby and maxDistanceForNearby always
	// cf. func loadLodConfigs()
	if senderEntity.changedAfterSend {
		if getIntervalFromLastSend(senderEntity, receiverID) > interval {
			logger.Verbosef("checkForBetweenNearbyAndFar", "senderID", senderID, "receiverID", receiverID, "intervalFromLastSend", getIntervalFromLastSend(senderEntity, receiverID), "interval", interval, "changed", true, "reason", "changed_after_send")
			updatesSentCounter.WithLabelValues("far", "changed").Inc()
			return true
		}
		logger.Verbosef("checkForBetweenNearbyAndFar", "senderID", senderID, "receiverID", receiverID, "intervalFromLastSend", getIntervalFromLastSend(senderEntity, receiverID), "interval", interval, "changed", true, "reason", "interval_not_met")
		updatesSkippedCounter.WithLabelValues("interval_not_met").Inc()
		return false
	}

	// if userEntity doesn't changed after send in previous interval, send update in syncIntervalForFar
	if getIntervalFromLastSend(senderEntity, receiverID) > m.syncIntervalForFar {
		logger.Verbosef("checkForBetweenNearbyAndFar", "senderID", senderID, "receiverID", receiverID, "intervalFromLastSend", getIntervalFromLastSend(senderEntity, receiverID), "interval", interval, "changed", false, "reason", "interval_met")
		updatesSentCounter.WithLabelValues("far", "interval").Inc()
		return true
	}
	logger.Verbosef("checkForBetweenNearbyAndFar", "senderID", senderID, "receiverID", receiverID, "intervalFromLastSend", getIntervalFromLastSend(senderEntity, receiverID), "interval", interval, "changed", false, "reason", "interval_not_met")
	updatesSkippedCounter.WithLabelValues("interval_not_met").Inc()
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
			select {
			case m.sendBuffer <- sendMessage{
				receiverUserID: receiverUserID,
				ver:            m.ver,
				cmd:            m.cmd,
				payload:        senderUserEntity.payload,
			}:
				senderUserEntity.mu.Lock()
				senderUserEntity.m[receiverUserID] = time.Now()
				senderUserEntity.mu.Unlock()
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
	// Use 1ms ticker for high precision timing
	// This ensures packets are sent at the correct intervals with minimal delay
	tick := time.NewTicker(m.syncIntervalForNearby / 3)
	defer tick.Stop()
	for {
		select {
		case <-m.lodDone:
			logger.Info("LOD loop stopped")
			close(m.sendBuffer)
			return
		case <-tick.C: // if tick is not received, it means that the interval is too short
			m.processAllUsers()
		}
	}

}

// calculateDistance computes Manhattan distance between two entities
func calculateDistance(x1, y1, x2, y2 int32) int32 {
	x := x1 - x2
	y := y1 - y2
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return x + y
}

func getIntervalFromLastSend(senderUserEntity *UserEntity, receiverUserID string) time.Duration {
	senderUserEntity.mu.RLock()
	defer senderUserEntity.mu.RUnlock()
	return time.Since(senderUserEntity.m[receiverUserID])
}
