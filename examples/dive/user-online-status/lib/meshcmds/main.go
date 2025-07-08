// © 2019-2024 Diarkis Inc. All rights reserved.

package meshCmds

import (
	"errors"

	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
	"github.com/Diarkis/diarkis/uuid/v4"
)

// GetOnlineStatusListCmd is the mesh command ID
const GetOnlineStatusListCmd uint16 = 2100

const createRemoteRoomCmd uint16 = 10001

func Setup() {
	diarkisexec.SetMeshCommandHandler(createRemoteRoomCmd, handleCreateRemoteRoom)
}

func CreateRemoteRoom(serverType string, maxMembers int, ttl, interval int64, cb func(error, string)) {
	data := make(map[string]interface{})
	data["maxMembers"] = maxMembers
	data["ttl"] = ttl
	data["interval"] = interval
	// randomly find an available node to create a room on
	nodes := mesh.GetNodeAddressesByType(serverType)
	seen := make(map[string]bool)
	targetNode := ""
	for {
		rand := util.RandomInt(0, len(nodes)-1)
		node := nodes[rand]
		if !mesh.IsNodeTaken(node) && mesh.IsNodeOnline(node) {
			targetNode = node
			break
		}
		seen[node] = true
		if len(seen) == len(nodes) {
			break
		}
	}
	if targetNode == "" {
		cb(errors.New("No available target node found"), "")
		return
	}
	mesh.SendRequest(createRemoteRoomCmd, targetNode, data, func(err error, res map[string]interface{}) {
		if err != nil {
			cb(err, "")
			return
		}
		roomID := mesh.GetString(res, "roomID")
		cb(nil, roomID)
	})
}

func handleCreateRemoteRoom(req map[string]interface{}) ([]byte, error) {
	maxMembers := mesh.GetInt(req, "maxMembers")
	ttl := mesh.GetInt64(req, "ttl")
	interval := mesh.GetInt64(req, "interval")
	allowEmpty := true
	join := false
	val, err := uuid.New()
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	data["sid"] = val.Bytes
	data["uid"] = val.String
	data["key"] = ""
	data["macKey"] = ""
	_, err = user.New(data, ttl)
	if err != nil {
		return nil, err
	}
	dummy := user.GetUserByID(val.String)
	if dummy == nil {
		return nil, errors.New("Failed to create a dummy user")
	}
	roomID, err := room.NewRoom(dummy, maxMembers, allowEmpty, join, ttl, interval)
	if err != nil {
		return nil, err
	}
	resp := make(map[string]interface{})
	resp["roomID"] = roomID
	ret, err := mesh.CreateReturnBytes(resp)
	return ret, err
}
