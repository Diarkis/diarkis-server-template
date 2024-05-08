package room

/*
OnPropertyUpdateSync assigns a callback on property sync event
*/
func (room *Room) OnPropertyUpdateSync(callback func(rops map[string][]byte)) {
}

/*
OnIncrProperty assigns a callback on property incr event
*/
func (room *Room) OnIncrProperty(callback func(success bool, numProps int64)) {
}

/*
OnIncrPropertySync assigns a callback on property incr sync event
*/
func (room *Room) OnIncrPropertySync(callback func(success bool, key string, numProps uint64)) {
}

/*
OnGetProperty assigns a callback on get property event
*/
func (room *Room) OnGetProperty(callback func(success bool, props [][]byte)) {
}

/*
OnSyncProperties assigns a callback on sync properties event
*/
func (room *Room) OnSyncProperties(callback func(success bool, props map[string][]byte)) {
}

/*
SyncProperties etrieve properties of the room.

	Raises OnPropertySync event.

Only properties with byte array value will be synchronized.
*/
func (room *Room) SyncProperties() {
}

/*
GetProperty retrieves properties of the room.
*/
func (room *Room) GetProperty(keys []string) {
}

/*
UpdateProperty updates properties of the room.
*/
func (room *Room) UpdateProperty(props map[string][]byte, sync bool) {
}

/*
IncrementProperty increments a property of the room.
*/
func (room *Room) IncrementProperty(key string, delta int64, sync bool) {
}
