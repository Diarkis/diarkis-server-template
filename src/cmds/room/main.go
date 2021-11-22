package roomcmds

import (
	"encoding/binary"
	"{0}/lib/meshCmds"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/roomSupport"
	"strconv"
)

var logger = log.New("room")

func Expose() {
	room.ExposeCommands()
	roomSupport.ExposeCommands()
	// we broadcast user ID to room members when the client leaves unexpectedly
	// The message will be propagated by OnMemberLeave even on the client
	room.SetOnDiscardCustomMessage(onDiscardCustomMessage)
	// set up custom mesh commands
	meshCmds.Setup()
}

func onDiscardCustomMessage(roomID string, userID string, sid string) []byte {
	logger.Debug("OnDiscardCustomMessage roomID:%v userID:%v sid:%v", roomID, userID, sid)
	// UE4 sample client uses uin64 as the data type for userID
	conv, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return []byte(userID)
	}
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(conv))
	logger.Debug("OnDiscardCustomMessage message user ID %v bytes %v", conv, bytes)
	return bytes
}
