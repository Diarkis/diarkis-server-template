// © 2019-2024 Diarkis Inc. All rights reserved.

package scenarios

import (
	"encoding/binary"
	"encoding/json"

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
	// TTL in seconds. This is applied both for session TTL and session scenario.
	SessionTTL uint8 `json:"sessionTTL"`
}

type SessionScenario struct {
	gp             *GlobalParams
	params         *SessionScenarioParams
	metrics        *report.CustomMetrics
	client         *bot_client.UDPClient
	IsSessionOwner bool
	SessionID      *string
	CurrentNum     uint8
	// If TTL is over, the session scenario will be reset.
	TTL int64
}

var _ Scenario = &SessionScenario{}

func NewSessionScenario() Scenario {
	return &SessionScenario{}
}

// // // // // Interface Functions // // // // //
func (s *SessionScenario) GetUserID() string {
	return s.params.UID
}

func (s *SessionScenario) IsInSession() bool {
	return s.SessionID != nil
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

	if s.IsInSession() {
		if s.IsSessionOwner {
			if s.CurrentNum < s.params.SessionMaxMember {
				logger.Infou(s.GetUserID(), "Session has not got maxMember. Invite more members. %d/%d", s.CurrentNum, s.params.SessionMaxMember)
				s.inviteSession()
				return
			}
		}

		// if TTL is over, leave the session
		if s.TTL < util.NowSeconds() {
			logger.Infou(s.GetUserID(), "Session TTL is over. Leaving the session.")
			s.leaveSession()
			return
		}

		// broadcast message if in session and member count is more than 1
		if s.CurrentNum > 1 {
			s.sessionBroadcast()
		}
	} else {
		// restart session scenario
		s.resetSessionScenario()
		s.startSessionScenario()
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
	s.client.RegisterOnResponse(util.CmdBuiltInVer, util.CmdSessionLeave, []uint8{bot_client.ResponseOk}, s.handleOnLeaveSession)
	s.client.RegisterOnResponse(util.CmdBuiltInVer, util.CmdSessionGetSessionInfoBySessionType, []uint8{bot_client.ResponseOk}, s.handleOnGetSessionInfoBySessionType)
}

func (s *SessionScenario) resetSessionScenario() {
	s.IsSessionOwner = false
	s.SessionID = nil
	s.gp.UserState.Set(s.GetUserID(), StateWaitingAsSessionMember, false)
}

func (s *SessionScenario) startSessionScenario() {
	s.TTL = util.NowSeconds() + int64(s.params.SessionTTL)
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
	logger.Debugu(s.GetUserID(), "Someone joined to Session. Session owner?: %v", s.IsSessionOwner)

	s.getSessionInfo()
}

func (s *SessionScenario) handleOnCreateSession(payload []byte) {
	sessionID := string(payload[1:])
	s.SessionID = &sessionID
	s.IsSessionOwner = true
	logger.Debugu(s.GetUserID(), "Created as a Session Owner.")
}

func (s *SessionScenario) handleOnJoinedSession(payload []byte) {
	sessionID := string(payload[1:])
	s.SessionID = &sessionID
	logger.Debugu(s.GetUserID(), "Joined as a Session member.")
}

func (s *SessionScenario) handleOnLeaveSession(_ []byte) {
	s.resetSessionScenario()
	logger.Debugu(s.GetUserID(), "Left the session.")
}

func (s *SessionScenario) handleOnGetSessionInfoBySessionType(payload []byte) {
	sessionType := payload[0]
	sessionIDLen := binary.BigEndian.Uint16(payload[1:3])
	sessionID := string(payload[3 : 3+sessionIDLen])
	currentMembers := binary.BigEndian.Uint16(payload[3+sessionIDLen : 3+sessionIDLen+2])
	maxMembers := binary.BigEndian.Uint16(payload[3+sessionIDLen+2 : 3+sessionIDLen+2+2])
	memberIDsLen := binary.BigEndian.Uint16(payload[3+sessionIDLen+2+2 : 3+sessionIDLen+2+2+2])
	memberIDs := dpacket.BytesToStringList(payload[3+sessionIDLen+2+2+2 : 3+sessionIDLen+2+2+2+memberIDsLen])
	ownerID := string(payload[3+sessionIDLen+2+2+2+memberIDsLen:])

	s.CurrentNum = uint8(currentMembers)

	logger.Debugu(s.GetUserID(), "Get Session Info By Session Type: %d, ID: %s, Members: %d/%d, Members: %v, Owner: %s",
		sessionType, sessionID, currentMembers, maxMembers, memberIDs, ownerID)
}

func (s *SessionScenario) createSession() {
	if s.IsInSession() {
		return
	}
	sessionType := s.params.SessionType
	maxMembers := s.params.SessionMaxMember
	ttl := s.params.SessionTTL
	payload := []byte{sessionType, maxMembers, ttl}
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionCreate, payload)
}

func (s *SessionScenario) joinSession(sessionID string) {
	if s.IsInSession() {
		return
	}
	sessionType := s.params.SessionType
	bytes := []byte{sessionType}
	bytes = append(bytes, []byte(sessionID)...)
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionJoin, bytes)
}

func (s *SessionScenario) leaveSession() {
	if !s.IsInSession() {
		return
	}
	sessionType := s.params.SessionType
	payload := []byte{sessionType}
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionLeave, payload)
}

func (s *SessionScenario) inviteSession() {
	if !s.IsSessionOwner {
		return
	}
	sessionType := s.params.SessionType
	maxMembers := s.params.SessionMaxMember
	// invite members
	memberIDs := []string{}
	userIDs := s.gp.UserState.Search(StateWaitingAsSessionMember, true, int(maxMembers-s.CurrentNum))
	if len(userIDs) == 0 {
		logger.Debugu(s.GetUserID(), "There are no members to invite to the session, therefore leaving the session.")
		s.metrics.Increment("NO_MEMBERS", "")
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
	payload = append(payload, []byte(s.GetUserID()+" Hello")...)
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionBroadcast, payload)
}

func (s *SessionScenario) getSessionInfo() {
	sessionType := s.params.SessionType
	payload := []byte{sessionType}
	s.client.RSend(util.CmdBuiltInVer, util.CmdSessionGetSessionInfoBySessionType, payload)
}
