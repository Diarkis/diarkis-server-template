// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/session"
	"github.com/Diarkis/diarkis/user"
)

var logger = log.New("SESSION-EXAMPLE")

const (
	builtInVer       uint8  = 1
	createSessionCmd uint16 = 801
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

func registerCallbacks(s *session.Session) {
	ok := s.SetOnJoin(func(id string, userData *user.User) bool {
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
