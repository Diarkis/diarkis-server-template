package mars

import (
	"sync"
)

/*
SharedData is shared and updated by any server node in Diarkis cluster.
The number of shared data entry is limited and dictated by util.SharedDataLimitLength
*/
type SharedData struct{ sync.RWMutex }

/*
NewSharedData creates a new SharedData instance.
*/
func NewSharedData() *SharedData {
	return nil
}

/*
Update changes SharedData value and propagate the change to all pods in the cluster.
Update does not change the value atomically therefore it may suffer from race condition.
*/
func (sd *SharedData) Update(key string, value interface{}, set bool) bool {
	return false
}

/*
FlushRemovedKeys flushes out removed keys from the memory.
*/
func (sd *SharedData) FlushRemovedKeys() map[string]int16 {
	return nil
}

/*
GetSharedData returns the value of SharedData.
*/
func (sd *SharedData) GetSharedData() map[string]interface{} {
	return nil
}

/*
EncodeSharedData encodes the SharedData value into a byte array.
*/
func (sd *SharedData) EncodeSharedData() []byte {
	return nil
}
