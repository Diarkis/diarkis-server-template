package mesh

import (
	"sync"
)

const MaxSize = 1300
const MaxMessageSize = 10000000
const MaxPacketSize = 1400
const Shrink = 10

/*
Msg Data structure of Msg
*/
type Msg struct {
	UUID      []byte
	Packet    []byte
	AddedSize int
	TTL       int64
	sync.RWMutex
}
