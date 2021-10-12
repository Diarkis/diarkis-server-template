package customcmds

import (
	"errors"
	dpayload "{0}/lib/payload"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

const p2pAddrList = "p2pAddrList"

func reportP2PAddr(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	addr := string(payload)
	if addr == "" {
		userData.ServerRespond([]byte("Invalid payload"), ver, cmd, server.Bad, true)
		next(errors.New("Invalid payload"))
		return
	}
	// we assume the user reporting is in a room
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		userData.ServerRespond([]byte("User not in room"), ver, cmd, server.Bad, true)
		next(errors.New("User not in room"))
		return
	}
	updated := true
	room.UpdateProperties(roomID, func (props map[string]interface{}) bool {
		if _, ok := props[p2pAddrList]; !ok {
			props[p2pAddrList] = make([]string, 0)
		}
		if _, ok := props[p2pAddrList].([]string); !ok {
			updated = false
			return false
		}
		props[p2pAddrList] = append(props[p2pAddrList].([]string), addr)
		return true
	})
	if !updated {
		userData.ServerRespond([]byte("Invalid room property"), ver, cmd, server.Bad, true)
		next(errors.New("Invalid room property"))
		return
	}
	userData.ServerRespond([]byte("OK"), ver, cmd, server.Ok, true)
	next(nil)
}

func initP2P(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	// we assume the P2P starting user is in a room
	roomID := room.GetRoomID(userData)
	if roomID == "" {
		userData.ServerRespond([]byte("User not in room"), ver, cmd, server.Bad, true)
		next(errors.New("User not in room"))
		return
	}
	addrList := room.GetProperty(roomID, p2pAddrList)
	if addrList == nil {
		userData.ServerRespond([]byte("No address list"), ver, cmd, server.Bad, true)
		next(errors.New("No address list"))
		return
	}
	bytes := dpayload.PackP2PInit(addrList.([]string))
	room.Broadcast(roomID, userData, util.CmdBuiltInVer, util.CmdBroadcastRoom, bytes, true)
	userData.ServerRespond([]byte("OK"), ver, cmd, server.Ok, true)
	next(nil)
}
