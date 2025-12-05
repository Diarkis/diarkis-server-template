// © 2019-2025 Diarkis Inc. All rights reserved.

package lodmanager

import (
	"sync"

	"github.com/Diarkis/diarkis/log"
)

var (
	managerMap      = make(map[string]*Manager)
	managerMapMutex = sync.RWMutex{}
)

var logger = log.New("LOD_MANAGER")

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
	if _, ok := managerMap[roomID]; !ok {
		return
	}
	managerMap[roomID].Stop()
	delete(managerMap, roomID)
}
