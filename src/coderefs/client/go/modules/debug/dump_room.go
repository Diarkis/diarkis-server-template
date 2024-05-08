package debug

import (
	"github.com/Diarkis/diarkis/room"
)

/*
DumpRoom dumps the content of the room
*/
func (d *Debug) DumpRoom() {
}

/*
OnDumpRoomResponse assigns a callback on DumpRoom
*/
func (d *Debug) OnDumpRoomResponse(callback func(success bool, roomData *room.DataDump)) {
}
