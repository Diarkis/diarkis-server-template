package customcmds

import (
	"encoding/binary"
	"errors"
	dpayload "{0}/lib/payload"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

// add to matching and create a room
func addToMatchMaker(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	mmAdd := dpayload.UnpackMMAdd(payload)
	if mmAdd == nil {
		userData.ServerRespond([]byte("Invalid payload"), ver, cmd, server.Bad, true)
		next(errors.New("Invalid payload"))
		return
	}
	// raise Room.OnCreate event
	maxMembers := 10
	allowEmpty := false
	join := true
	ttl := int64(120)
	interval := int64(0)
	roomID, err := room.NewRoom(userData, maxMembers, allowEmpty, join, ttl, interval)
	if err != nil {
		userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Bad, true)
		next(err)
		return
	}
	created := room.GetCreatedTime(roomID)
	timeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(timeBytes, uint32(created))
	roomIDBytes := []byte(roomID)
	bytes := append(timeBytes, roomIDBytes...)
	userData.ServerRespond(bytes, util.CmdBuiltInVer, util.CmdCreateRoom, server.Ok, true)
	metadata := make(map[string]interface{})
	metadata["serialized"] = mmAdd.Metadata
	metadata["uniqueID"] = mmAdd.UID
	metadata["roomID"] = roomID
	metadata["maxMembers"] = maxMembers
	matching.Add(mmAdd.ID, mmAdd.UID, mmAdd.Props, metadata, mmAdd.TTL, 2)
	// update MatchMaker add every 40 seconds
	timeCnt := util.NowSeconds()
	room.SetOnTick(roomID, func(roomID_ string) {
		if util.NowSeconds() - timeCnt < mmAddInterval {
			return
		}
		if len(room.GetMemberIDs(roomID_)) == maxMembers {
			return
		}
		timeCnt = util.NowSeconds()
		metadata["serialized"] = mmAdd.Metadata
		metadata["uniqueID"] = mmAdd.UID
		metadata["roomID"] = roomID
		metadata["maxMembers"] = maxMembers
		matching.Add(mmAdd.ID, mmAdd.UID, mmAdd.Props, metadata, mmAdd.TTL, 2)
	})
	// response for matchmaking
	userData.ServerRespond([]byte("OK"), ver, cmd, server.Ok, true)
	next(nil)
}

func searchMatchMaker(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	mmSearch := dpayload.UnpackMMSearch(payload)
	if mmSearch == nil {
		userData.ServerRespond([]byte("Invalid payload"), ver, cmd, server.Bad, true)
		next(errors.New("Invalid payload"))
		return
	}
	howmany := 10
	matching.Search(mmSearch.IDs, mmSearch.Props, howmany, func(err error, results []interface{}) {
		if err != nil {
			userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Bad, true)
			next(err)
			return
		}
		logger.Sys("MatchMaker search results %v", results)
		if len(results) == 0 {
			userData.ServerRespond([]byte("Matching not found"), ver, cmd, server.Bad, true)
			next(errors.New("Matching not found"))
			return
		}
		list := make([]map[string]interface{}, len(results))
		operations := make([]func(func(error)), len(results))
		index := 0
		joinedRoomID := ""
		ct := int64(0)
		isRoomFull := false
		done := func(err error) {
			if err != nil {
				userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Bad, true)
			}
			if joinedRoomID == "" {
				userData.ServerRespond([]byte("No room found to join"), ver, cmd, server.Bad, true)
				next(errors.New("No room found to join"))
				return
			}
			timeBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(timeBytes, uint32(ct))
			roomIDBytes := []byte(joinedRoomID)
			bytes := append(timeBytes, roomIDBytes...)
			userData.ServerRespond(bytes, util.CmdBuiltInVer, util.CmdJoinRoom, server.Ok, true)
			userData.ServerRespond([]byte(joinedRoomID), ver, cmd, server.Ok, true)
			if isRoomFull {
				userData.ServerPush(ver, matchmakerComplete, []byte(joinedRoomID), true)
			}
			next(err)
		}
		join := func(moveon func(error)) {
			item := list[index]
			uniqueID := item["uniqueID"].(string)
			roomID := item["roomID"].(string)
			logger.Sys("Try to join room %v", roomID)
			room.Join(roomID, userData, util.CmdBuiltInVer, util.CmdJoinRoom, []byte(userData.ID),
			func(err error, memberIDs []string, ownerID string, createdTime int64, props map[string]interface{}) {
				if err != nil {
					// discard stale room ID
					removeFromMM(mmSearch.IDs, uniqueID)
					// try the next room
					index++
					moveon(nil)
					return
				}
				logger.Sys("Successfully joined a from by MatchMaker %v", roomID)
				joinedRoomID = roomID
				ct = createdTime
				isRoomFull = len(memberIDs) == int(item["maxMembers"].(float64))
				// joined successfully
				done(nil)
			})
		}
		for i, v := range results {
			list[i] = v.(map[string]interface{})
			operations[i] = join
		}
		util.Waterfall(operations, done)
	})
}

func removeFromMM(mmIDs []string, uniqueID string) {
	uids := make([]string, 1)
	uids[0] = uniqueID
	for _, id := range mmIDs {
		matching.Remove(id, uids, 2)
	}
}
