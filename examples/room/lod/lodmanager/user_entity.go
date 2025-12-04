package lodmanager

import (
	"fmt"
	"time"
)

// UserEntity represents a user's position and payload data
type UserEntity struct {
	X                int32
	Y                int32
	Payload          []byte
	Remember         map[string]time.Time
	ChangedAfterSend bool
}

func NewUserEntity(x int32, y int32, payload []byte) *UserEntity {
	return &UserEntity{
		X:                x,
		Y:                y,
		Payload:          payload,
		Remember:         make(map[string]time.Time),
		ChangedAfterSend: true,
	}
}

func (ue *UserEntity) String() string {
	return fmt.Sprintf("UserEntity{X: %v, Y: %v, Payload: %v, Remember: %v, ChangedAfterSend: %v}", ue.X, ue.Y, ue.Payload, ue.Remember, ue.ChangedAfterSend)
}
