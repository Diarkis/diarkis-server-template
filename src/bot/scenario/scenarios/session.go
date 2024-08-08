// © 2019-2024 Diarkis Inc. All rights reserved.

package scenarios

import (
	"encoding/binary"
	"encoding/json"
	"sync/atomic"

	bot_client "github.com/Diarkis/diarkis-server-template/bot/scenario/lib/client"
	"github.com/Diarkis/diarkis-server-template/bot/scenario/lib/report"
	dpacket "github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/util"
)

const StateWaitingAsSessionMember = "WaitingAsSessionMember"

type SessionScenarioParams struct {
	ServerType       string `json:"serverType"`
	UID              string `json:"userID"`
	SessionType      uint8  `json:"sessionType"`
	SessionMaxMember uint8  `json:"sessionMaxMember"`
	SessionTTL       uint8  `json:"sessionTTL"`
}

type SessionScenario struct {
	gp               *GlobalParams
	params           *SessionScenarioParams
	metrics          *report.CustomMetrics
	client           *bot_client.UDPClient
	IsSessionMember  bool
	IsSessionOwner   bool
	IsInSession      bool
	sessionMemberCnt atomic.Uint32
}

var _ Scenario = &SessionScenario{}

func NewSessionScenario() Scenario {
	return &SessionScenario{}
}

// // // // // Interface Functions // // // // //
func (s *SessionScenario) GetUserID() string {
	return s.params.UID
}

func (s *SessionScenario) ParseParam(index int, params []byte) error {
	var sessionScenarioParams *SessionScenarioParams
	err := json.Unmarshal(params, &sessionScenarioParams)
	if err != nil {
		logger.Error("Failed to unmarshal scenario params.", err.Error())
		return err
	}
	s.params = sessionScenarioParams
	logger.Sys("Parsed Params. %#v", sessionScenarioParams)
	return nil
}

func (s *SessionScenario) Run(gp *GlobalParams) error {
	logger.Info("Starting scenario for user %v", s.GetUserID())

	// store GlobalParams to use it later
	s.gp = gp
	// new report
	s.metrics = report.NewCustomMetrics()

	// connect to UDP server
	_, udpClient, err := bot_client.NewAndConnect(gp.Host, s.GetUserID(), s.params.ServerType, nil, gp.ReceiveByteSize, gp.UDPSendInterval)
	if err != nil {
		return err
	}
	s.client = udpClient

	// set OnPush / OnResponse handlers
	s.RegisterSessionHandlers()

	// start scenario
	s.startSessionScenario()

	return nil
}

func (s *SessionScenario) OnIdle() {
	logger.Infou(s.GetUserID(), "Triggering OnIdle Session...")
	currentMemberCnt := uint8(s.sessionMemberCnt.Load())
	if s.IsSessionOwner && currentMemberCnt < s.params.SessionMaxMember {
		logger.Infou(s.GetUserID(), "Session has not got maxMember. Invite more members. %d/%d", currentMemberCnt, s.params.SessionMaxMember)
		s.inviteSession()
		return
	} else if s.IsSessionOwner && currentMemberCnt == s.params.SessionMaxMember {
		logger.Infou(s.GetUserID(), "Session has got maxMember. Do nothing. %d/%d", currentMemberCnt, s.params.SessionMaxMember)
	}

	// broadcast message if in session and member count is more than 1
	if s.IsInSession && currentMemberCnt > 1 {
		s.sessionBroadcast()
	}
}

func (s *SessionScenario) OnScenarioEnd() error {
	return nil
}

func (s *SessionScenario) RegisterSessionHandlers() {
	s.client.RegisterOnPush(util.CmdBuiltInVer, util.CmdSessionInvite, s.handleSessionInvite)
	s.client.RegisterOnPush(util.CmdBuiltInVer, util.CmdSessionJoin, s.handleSessionJoined)
	s.client.RegisterOnPush(util.CmdBuiltInVer, util.CmdSessionBroadcast, s.handleBroadcast)
	s.client.RegisterOnResponse(util.CmdBuiltInVer, util.CmdSessionCreate, []uint8{bot_client.ResponseOk}, s.handleOnCreateSession)
	s.client.RegisterOnResponse(util.CmdBuiltInVer, util.CmdSessionJoin, []uint8{bot_client.ResponseOk}, s.handleOnJoinedSession)
	s.client.RegisterOnResponse(util.CmdBuiltInVer, util.CmdSessionJoin, []uint8{bot_client.ResponseOk}, s.handleOnLeaveSession)
}

func (s *SessionScenario) startSessionScenario() {
	isOwner := util.RandomInt(1, int(s.params.SessionMaxMember)) == 1
	if isOwner {
		s.createSession()
		return
	}

	// set a waiting flag and wait for party invitation
	s.gp.UserState.Set(s.GetUserID(), StateWaitingAsSessionMember, true)
}

func (s *SessionScenario) handleSessionInvite(payload []byte) {
	logger.Sysu(s.GetUserID(), "Invited to Session")

	sessionIDSize := int(payload[1])
	sessionID := string(payload[2 : 2+sessionIDSize])
	// accept the session invite
	s.joinSession(sessionID)
}

func (s *SessionScenario) handleBroadcast(payload []byte) {
	logger.Verboseu(s.GetUserID(), "Received a broadcast message. %s", string(payload))
}

func (s *SessionScenario) handleSessionJoined(payload []byte) {
	logger.Sysu(s.GetUserID(), "Someone joined to Session")
	if s.IsSessionOwner {
		s.sessionMemberCnt.Add(1)
		if uint8(s.sessionMemberCnt.Load()) == s.params.SessionMaxMember {
			logger.Debugu(s.GetUserID(), "Session has been ready as it's got all members joined.")
		}
	}
}

func (s *SessionScenario) handleOnCreateSession(_ []byte) {
	s.IsInSession = true
	s.IsSessionOwner = true
	s.sessionMemberCnt.Add(1)
	logger.Debugu(s.GetUserID(), "Created as a Session Owner.")
}

func (s *SessionScenario) handleOnJoinedSession(_ []byte) {
	s.IsInSession = true
	s.IsSessionMember = true
	logger.Debugu(s.GetUserID(), "Joined as a Session member.")
}

func (s *SessionScenario) handleOnLeaveSession(_ []byte) {
	s.IsInSession = false
	s.IsSessionOwner = false
	s.IsSessionMember = false
	s.gp.UserState.Set(s.GetUserID(), StateWaitingAsSessionMember, true)
	logger.Debugu(s.GetUserID(), "Left the session.")
}

func (s *SessionScenario) createSession() {
	if s.IsInSession {
		return
	}
	sessionType := s.params.SessionType
	maxMembers := s.params.SessionMaxMember
	ttl := s.params.SessionTTL
	payload := []byte{sessionType, maxMembers, ttl}
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionCreate, payload)
}

func (s *SessionScenario) joinSession(sessionID string) {
	if s.IsInSession {
		return
	}
	sessionType := s.params.SessionType
	bytes := []byte{sessionType}
	bytes = append(bytes, []byte(sessionID)...)
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionJoin, bytes)
}

func (s *SessionScenario) leaveSession() {
	if !s.IsInSession {
		return
	}
	sessionType := s.params.SessionType
	payload := []byte{sessionType}
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionLeave, payload)
}

func (s *SessionScenario) inviteSession() {
	sessionType := s.params.SessionType
	currentMemberCnt := uint8(s.sessionMemberCnt.Load())
	maxMembers := s.params.SessionMaxMember
	// invite members
	memberIDs := []string{}
	userIDs := s.gp.UserState.Search(StateWaitingAsSessionMember, true, int(maxMembers-currentMemberCnt))
	if len(userIDs) == 0 {
		logger.Debugu(s.GetUserID(), "There are no members to invite to the session, change myself as a session member and wait for the invitation.")
		s.leaveSession()
		return
	}

	for _, memberID := range userIDs {
		if s.gp.UserState.Get(memberID, StateWaitingAsSessionMember).(bool) == true {
			s.gp.UserState.Set(memberID, StateWaitingAsSessionMember, false)
			memberIDs = append(memberIDs, memberID)
		}
	}
	message := []byte("invite message")
	payload := []byte{sessionType}
	bMemberIDs := dpacket.StringListToBytes(memberIDs)
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(bMemberIDs)))
	payload = append(payload, size...)
	payload = append(payload, bMemberIDs...)
	payload = append(payload, message...)
	// call session invite
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionInvite, payload)
}

func (s *SessionScenario) sessionBroadcast() {
	sessionType := s.params.SessionType
	payload := []byte{sessionType}
	payload = append(payload, []byte("session broadcast message")...)
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionBroadcast, payload)
}
