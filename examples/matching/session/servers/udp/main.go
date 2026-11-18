// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis-server-template/examples/matching/session/protocol"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	psession "github.com/Diarkis/diarkis/proto/go/session"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/session"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

var logger = log.New("SESSION-EXAMPLE")

const (
	builtInVer         uint8  = 1
	createSessionCmd   uint16 = 801
	inviteCmd          uint16 = 802
	kickCmd            uint16 = 806
	rejoinBlockSeconds int64  = 60
	kickedKeyPrefix           = "kicked:"
	invitedByKeyPrefix        = "invitedBy:"
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

func handleShowSessions(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
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
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })

	res, err := json.Marshal(list)
	if err != nil {
		userData.ServerRespond([]byte(err.Error()), ver, cmd, server.Err, true)
		next(err)
		return
	}
	logger.Infof("session show", "uid", userData.ID, "count", len(list))
	userData.ServerRespond(res, ver, cmd, server.Ok, true)
	next(nil)
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
		notifyInviter(reply, userData)
		clearInviter(reply.SessionID, userData.ID)
		userData.ServerRespond([]byte("declined"), ver, cmd, server.Ok, true)
		next(nil)
		return
	}

	logger.Infof("invite accepted", "sessionID", reply.SessionID, "uid", userData.ID)
	clearInviter(reply.SessionID, userData.ID)
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

// clearInviter 招待への返事 (yes / no) が来たら、保存しておいた招待者 (invitedBy:<UID>) を消す。
// Session のプロパティを削除する API はないので、空文字で上書きする (notifyInviter は空文字を「招待者なし」として扱う)。
// 招待を断ったユーザーは Session のメンバーではないため、メンバー用の session.SetProperty は使えない。
// そのため Session があるサーバーで *Session を直接書き換える (UDP サーバー1台が前提)。
func clearInviter(sessionID string, uid string) {
	s, ok := session.GetSessionByID(sessionID)
	if !ok {
		return
	}
	s.SetProperty(invitedByKeyPrefix+uid, "")
	logger.Infof("inviter cleared", "sessionID", sessionID, "uid", uid)
}

// GetSessionByID は Session があるサーバーでしか使えないので、このサンプル(UDP サーバー1台)が前提。
func notifyInviter(reply protocol.InviteReply, userData *user.User) {
	s, ok := session.GetSessionByID(reply.SessionID)
	if !ok {
		logger.Warnf("session not found to notify inviter", "sessionID", reply.SessionID)
		return
	}
	v, ok := s.GetProperty(invitedByKeyPrefix + userData.ID)
	if !ok {
		logger.Warnf("inviter not found", "sessionID", reply.SessionID, "uid", userData.ID)
		return
	}
	inviterUID, ok := util.ToString(v)
	if !ok || inviterUID == "" {
		logger.Warnf("invalid inviter", "sessionID", reply.SessionID, "uid", userData.ID)
		return
	}

	msg, _ := json.Marshal(protocol.InviteDeclined{SessionID: reply.SessionID, UID: userData.ID})
	err := session.MessageTo(reply.SessionType, reply.SessionID, userData.ID, []string{inviterUID}, protocol.CustomVer, protocol.InviteDeclinedPushCmd, msg)
	if err != nil {
		logger.Warnf("failed to notify inviter", "inviterUID", inviterUID, "error", err)
		return
	}
	logger.Infof("inviter notified", "inviterUID", inviterUID, "uid", userData.ID)
}
