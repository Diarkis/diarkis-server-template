// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"time"

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
	})

	logger.Infof("callbacks registered", "sessionID", s.GetID())
}
