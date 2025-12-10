// © 2019-2025 Diarkis Inc. All rights reserved.

package lodmanager

import (
	"fmt"
	"sync"
	"time"
)

// UserEntity represents a user's position and payload data
type UserEntity struct {
	x                int32
	y                int32
	payload          []byte
	m                map[string]time.Time
	mu               sync.RWMutex
	changedAfterSend bool
}

func NewUserEntity(x int32, y int32, payload []byte) *UserEntity {
	return &UserEntity{
		x:                x,
		y:                y,
		payload:          payload,
		m:                make(map[string]time.Time),
		changedAfterSend: true,
	}
}

func (ue *UserEntity) String() string {
	return fmt.Sprintf("UserEntity{X: %v, Y: %v, Payload: %v, Remember: %v, ChangedAfterSend: %v}", ue.x, ue.y, ue.payload, ue.m, ue.changedAfterSend)
}
