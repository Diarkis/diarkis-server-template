// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/Diarkis/diarkis-server-template/examples/matching/session/protocol"

	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/session"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

var logger = log.New("SESSION-EXAMPLE")

const (
	builtInVer       uint8  = 1
	createSessionCmd uint16 = 801
	kickCmd          uint16 = 806

	// kick された人が再参加できるようになるまでの秒数
	rejoinBlockSeconds int64 = 60
	// kick した時刻を Session のプロパティに保存するときのキーの頭につける文字
	kickedKeyPrefix = "kicked:"
)

// sessionList このサーバーで作られた Session の一覧(Session ID → Session type)。
// コマンドの処理と OnDeleted(別 goroutine)から同時に触られるので、Mutex で守る。
// このサンプルは UDP サーバーが1台なので、このサーバーの一覧だけを返す。
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
		session.ExposeCommands()

		if err := server.HandleCommand(builtInVer, createSessionCmd, handleAfterCreateSession); err != nil {
			panic(err)
		}
		if err := server.HandleCommand(builtInVer, kickCmd, handleAfterKick); err != nil {
			panic(err)
		}
		if err := server.HandleCommand(protocol.CustomVer, protocol.ShowSessionsCmd, handleShowSessions); err != nil {
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

// handleAfterKick 標準の session kick が成功したあとに呼ばれる。
// kick された人と時刻を Session のプロパティに保存する。
func handleAfterKick(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Infof("handleAfterKick called", "uid", userData.ID)

	// payload は「1バイト目が Session type、残りが kick された人の UID」
	// (例: user2 を kick すると 00 75 73 65 72 32 になることをログで確認済み)
	if len(payload) < 2 {
		logger.Warnf("kick payload is too short", "uid", userData.ID)
		next(nil)
		return
	}
	sessionType := payload[0]
	targetUID := string(payload[1:])

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

// isRejoinBlocked kick されてから rejoinBlockSeconds 秒たっていなければ true を返す。
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

// handleShowSessions 独自コマンド session show。Session の一覧を JSON で返す。
func handleShowSessions(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	// Mutex を持っている時間を短くするため、まず一覧をコピーする
	sessionListMu.Lock()
	copied := make(map[string]uint8, len(sessionList))
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
	// 表示の順番が毎回変わらないように並べる
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
