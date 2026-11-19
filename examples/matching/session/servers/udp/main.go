// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis-server-template/examples/matching/session/protocol"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mesh"
	psession "github.com/Diarkis/diarkis/proto/go/session"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/session"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

var logger = log.New("SESSION-EXAMPLE")

const (
	builtInVer          uint8  = 1
	createSessionCmd    uint16 = 801
	inviteCmd           uint16 = 802
	kickCmd             uint16 = 806
	rejoinBlockSeconds  int64  = 60
	kickedKeyPrefix            = "kicked:"
	invitedByKeyPrefix         = "invitedBy:"
	meshInviteReplyCmd  uint16 = 10001
	meshListSessionsCmd uint16 = 10002
)

var (
	sessionListMu sync.Mutex
	sessionList   = map[string]uint8{}
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Session: &diarkisexec.Options{ConfigPath: "configs/shared/session.json", ExposeCommands: false},
		Dive:    &diarkisexec.Options{ConfigPath: "configs/shared/dive.json"},
		DM:      &diarkisexec.Options{ConfigPath: "configs/shared/dm.json"},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")
	diarkisexec.SetMeshCommandHandler(meshInviteReplyCmd, handleMeshInviteReply)
	diarkisexec.SetMeshCommandHandler(meshListSessionsCmd, handleMeshListSessions)

	diarkis.OnStart(func() {
		session.ExposePufferCommands()

		if err := server.HandleCommand(builtInVer, createSessionCmd, handleAfterCreateSession); err != nil {
			panic(err)
		}
		if err := server.HandleCommand(builtInVer, inviteCmd, handleAfterInvite); err != nil {
			panic(err)
		}
		if err := server.HandleCommand(builtInVer, kickCmd, handleAfterKick); err != nil {
			panic(err)
		}
		if err := server.HandleCommand(protocol.CustomVer, protocol.ShowSessionsCmd, handleShowSessions); err != nil {
			panic(err)
		}
		if err := server.HandleCommand(protocol.CustomVer, protocol.InviteReplyCmd, handleInviteReply); err != nil {
			panic(err)
		}
	})

	diarkisexec.StartDiarkis()
}

func handleAfterCreateSession(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Infof("handleAfterCreateSession called", "uid", userData.ID)

	sessions := session.GetSessionDataByUser(userData)
	if len(sessions) == 0 {
		logger.Warnf("no session found after create", "uid", userData.ID)
	}
	for _, data := range sessions {
		s, ok := session.GetSessionByID(data.ID)
		if !ok {
			logger.Warnf("session not found", "sessionID", data.ID, "uid", userData.ID)
			continue
		}
		registerCallbacks(s)
	}
	next(nil)
}

func handleAfterKick(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Infof("handleAfterKick called", "uid", userData.ID)

	kick := psession.NewKick()
	if err := kick.Unpack(payload); err != nil {
		logger.Warnf("failed to read kick payload", "uid", userData.ID, "error", err)
		next(nil)
		return
	}
	sessionType := kick.SessionType
	targetUID := kick.TargetUID

	key := kickedKeyPrefix + targetUID
	now := time.Now().Unix()
	session.SetProperty(sessionType, userData, key, now, func(err error, success bool) {
		if err != nil || !success {
			logger.Warnf("failed to save kicked time", "targetUID", targetUID, "error", err)
		} else {
			logger.Infof("kicked time saved", "targetUID", targetUID, "kickedAt", now)
		}
		next(nil)
	})
}

func handleAfterInvite(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	invite := psession.NewInvite()
	if err := invite.Unpack(payload); err != nil {
		logger.Warnf("failed to read invite payload", "uid", userData.ID, "error", err)
		next(nil)
		return
	}
	sessionType := invite.SessionType
	targetUIDs := invite.TargetUIDs

	kvList := make([]*session.PropertyKeyValue, 0, len(targetUIDs))
	for _, targetUID := range targetUIDs {
		kvList = append(kvList, &session.PropertyKeyValue{Key: invitedByKeyPrefix + targetUID, Value: userData.ID})
	}
	session.SetProperties(sessionType, userData, kvList, func(err error) {
		if err != nil {
			logger.Warnf("failed to save inviter", "inviterUID", userData.ID, "error", err)
		} else {
			logger.Infof("inviter saved", "inviterUID", userData.ID, "targetUIDs", targetUIDs)
		}
		next(nil)
	})
}

func isRejoinBlocked(s *session.Session, uid string) bool {
	v, ok := s.GetProperty(kickedKeyPrefix + uid)
	if !ok {
		return false
	}
	kickedAt, ok := util.ToInt64(v)
	if !ok {
		return false
	}
	return time.Now().Unix()-kickedAt < rejoinBlockSeconds
}

func registerCallbacks(s *session.Session) {
	ok := s.SetOnJoin(func(id string, userData *user.User) bool {
		if isRejoinBlocked(s, userData.ID) {
			logger.Infof("OnJoin rejected: kicked recently", "sessionID", id, "uid", userData.ID)
			return false
		}
		logger.Infof("OnJoin", "sessionID", id, "uid", userData.ID)
		return true
	})
	if !ok {
		return
	}

	s.SetOnJoined(func(id string, userData *user.User) {
		logger.Infof("OnJoined", "sessionID", id, "uid", userData.ID)
	})

	s.SetOnLeft(func(id string, userData *user.User) {
		logger.Infof("OnLeft", "sessionID", id, "uid", userData.ID)
	})

	s.SetOnDeleted(func(id string) {
		logger.Infof("OnDeleted", "sessionID", id)
		removeFromSessionList(id)
	})

	addToSessionList(s.GetID(), s.GetType())

	logger.Infof("callbacks registered", "sessionID", s.GetID())
}

func addToSessionList(id string, sessionType uint8) {
	sessionListMu.Lock()
	defer sessionListMu.Unlock()
	sessionList[id] = sessionType
}

func removeFromSessionList(id string) {
	sessionListMu.Lock()
	defer sessionListMu.Unlock()
	delete(sessionList, id)
}

// NOTE: UDPサーバーの台数が増えると、そのぶんmeshの通信も増える。
// 実際のアプリでは、呼ぶ頻度を抑えるか、問い合わせるサーバー数に上限を設けるなどが必要。
func handleShowSessions(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	addrs := mesh.GetNodeAddressesByType(mesh.GetNodeType())

	var mu sync.Mutex
	var wg sync.WaitGroup
	list := []protocol.SessionInfo{}

	for _, addr := range addrs {
		wg.Add(1)
		mesh.SendRequest(meshListSessionsCmd, addr, map[string]any{"uid": userData.ID}, func(err error, res map[string]any) {
			defer wg.Done()
			if err != nil {
				logger.Warnf("failed to get sessions from a server", "addr", addr, "error", err)
				return
			}
			var part []protocol.SessionInfo
			if err := json.Unmarshal(mesh.GetBytes(res, "sessions"), &part); err != nil {
				logger.Warnf("invalid sessions from a server", "addr", addr, "error", err)
				return
			}
			mu.Lock()
			list = append(list, part...)
			mu.Unlock()
		})
	}

	go func() {
		wg.Wait()
		sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })

		res, err := json.Marshal(list)
		if err != nil {
			userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Err, true)
			next(err)
			return
		}
		logger.Infof("session show", "uid", userData.ID, "servers", len(addrs), "count", len(list))
		userData.ServerRespond(res, ver, cmd, server.Ok, true)
		next(nil)
	}()
}

func handleMeshListSessions(req map[string]any) ([]byte, error) {
	sessionListMu.Lock()
	copied := make(map[string]uint8)
	for id, sessionType := range sessionList {
		copied[id] = sessionType
	}
	sessionListMu.Unlock()

	list := make([]protocol.SessionInfo, 0, len(copied))
	for id, sessionType := range copied {
		info := protocol.SessionInfo{ID: id, Type: sessionType, MemberIDs: []string{}}
		if s, ok := session.GetSessionByID(id); ok {
			info.MemberIDs = s.GetMemberIDs()
		}
		list = append(list, info)
	}

	b, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return mesh.CreateReturnBytes(map[string]any{"sessions": b})
}

func handleInviteReply(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	var reply protocol.InviteReply
	if err := json.Unmarshal(payload, &reply); err != nil || reply.SessionID == "" {
		userData.ServerRespond([]byte("invalid invite reply"), ver, cmd, server.Bad, true)
		next(nil)
		return
	}

	if !reply.Accept {
		logger.Infof("invite declined", "sessionID", reply.SessionID, "uid", userData.ID)
		requestInviteReply(reply, userData.ID)
		userData.ServerRespond([]byte("declined"), ver, cmd, server.Ok, true)
		next(nil)
		return
	}

	logger.Infof("invite accepted", "sessionID", reply.SessionID, "uid", userData.ID)
	requestInviteReply(reply, userData.ID)
	session.JoinSessionByID(reply.SessionType, reply.SessionID, userData, func(err error) {
		if err != nil {
			logger.Infof("join by invite failed", "sessionID", reply.SessionID, "uid", userData.ID, "error", err)
			userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Bad, true)
			next(nil)
			return
		}
		userData.ServerRespond([]byte("joined"), ver, cmd, server.Ok, true)
		next(nil)
	})
}

func requestInviteReply(reply protocol.InviteReply, uid string) {
	addr, err := util.GetAddressFromAddressID(reply.SessionID)
	if err != nil {
		logger.Warnf("invalid session ID", "sessionID", reply.SessionID, "error", err)
		return
	}
	req := map[string]any{
		"sessionType": reply.SessionType,
		"sessionID":   reply.SessionID,
		"uid":         uid,
		"accept":      reply.Accept,
	}
	mesh.SendRequest(meshInviteReplyCmd, addr, req, func(err error, res map[string]any) {
		if err != nil {
			logger.Warnf("failed to handle invite reply on the session server", "sessionID", reply.SessionID, "uid", uid, "error", err)
		}
	})
}

func handleMeshInviteReply(req map[string]any) ([]byte, error) {
	sessionType := mesh.GetUint8(req, "sessionType")
	sessionID := mesh.GetString(req, "sessionID")
	uid := mesh.GetString(req, "uid")
	accept := mesh.GetBool(req, "accept")
	s, ok := session.GetSessionByID(sessionID)
	if !ok {
		return nil, errors.New("session not found")
	}

	key := invitedByKeyPrefix + uid
	inviterUID := ""
	if v, ok := s.GetProperty(key); ok {
		inviterUID, _ = util.ToString(v)
	}

	s.SetProperty(key, "")
	logger.Infof("inviter cleared", "sessionID", sessionID, "uid", uid)

	if !accept {
		if inviterUID == "" {
			logger.Warnf("inviter not found", "sessionID", sessionID, "uid", uid)
		} else {
			msg, _ := json.Marshal(protocol.InviteDeclined{SessionID: sessionID, UID: uid})
			if err := session.MessageTo(sessionType, sessionID, uid, []string{inviterUID}, protocol.CustomVer, protocol.InviteDeclinedPushCmd, msg); err != nil {
				logger.Warnf("failed to notify inviter", "inviterUID", inviterUID, "error", err)
			} else {
				logger.Infof("inviter notified", "inviterUID", inviterUID, "uid", uid)
			}
		}
	}

	return mesh.CreateReturnBytes(map[string]any{})
}
