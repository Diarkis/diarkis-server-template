package lodmanager

import (
	"sync"
)

var (
	managerMap      = make(map[string]*Manager)
	managerMapMutex = sync.RWMutex{}
)

// GetRoomManager retrieves the LOD manager for a specific room ID
func GetRoomManager(roomID string) *Manager {
	managerMapMutex.RLock()
	defer managerMapMutex.RUnlock()
	return managerMap[roomID]
}

func SetRoomManager(roomID string, manager *Manager) {
	managerMapMutex.Lock()
	defer managerMapMutex.Unlock()
	managerMap[roomID] = manager
}

func RemoveRoomManager(roomID string) {
	managerMapMutex.Lock()
	defer managerMapMutex.Unlock()
	managerMap[roomID].started.Store(false)
	delete(managerMap, roomID)
}
