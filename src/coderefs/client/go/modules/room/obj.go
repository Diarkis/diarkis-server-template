package room

const UpdateObjectIncrMode = byte(1)
const UpdateObjectSetMode = byte(2)
const UpdateObjectDelMode = byte(3)

/*
OnObjectUpdateResponse assigns a callback on object update response event
*/
func (room *Room) OnObjectUpdateResponse(callback func(msg []byte)) {
}

/*
OnObjectUpdatePush assigns a callback on object update push event
*/
func (room *Room) OnObjectUpdatePush(callback func(uint8, string, map[string]float64)) {
}

/*
UpdateObject updates a room object by the given properties.
mode: UpdateObjectIncrMode, UpdateObjectSetMode, UpdateObjectDelMode
*/
func (room *Room) UpdateObject(mode byte, name string, obj map[string]float64, reliable bool) {
}
