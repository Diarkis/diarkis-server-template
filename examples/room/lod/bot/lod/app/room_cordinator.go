package app

import "sync"

// RoomCoordinator manages room distribution among bots
type RoomCoordinator struct {
	CreatedRoomMap sync.Map
}

func NewRoomCoordinator() *RoomCoordinator {
	return &RoomCoordinator{}
}

func (rc *RoomCoordinator) AddRoom(index int, roomID string) {
	rc.CreatedRoomMap.Store(index, roomID)
}

func (rc *RoomCoordinator) GetRoom(index int) (string, bool) {
	val, ok := rc.CreatedRoomMap.Load(index)
	if !ok {
		return "", false
	}
	return val.(string), true
}
