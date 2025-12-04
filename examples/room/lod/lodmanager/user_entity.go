package lodmanager

// UserEntity represents a user's position and payload data
type UserEntity struct {
	X                int32
	Y                int32
	Payload          []byte
	Remember         *cache
	ChangedAfterSend bool
}

func NewUserEntity(x int32, y int32, payload []byte) *UserEntity {
	return &UserEntity{
		X:                x,
		Y:                y,
		Payload:          payload,
		Remember:         newCache(),
		ChangedAfterSend: true,
	}
}
