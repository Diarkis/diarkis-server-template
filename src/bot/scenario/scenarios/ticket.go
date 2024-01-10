package scenarios

import (
	"encoding/json"
	"fmt"
	"time"

	bot_client "{0}/bot/scenario/lib/client"
	"{0}/bot/scenario/lib/report"
	"{0}/bot/scenario/packets"

	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/util"
)

type TicketScenarioParams struct {
	ServerTypeMM   string `json:"serverTypeMM"`
	ServerTypeTurn string `json:"serverTypeTurn"`
	UID            string `json:"userID"`
	TicketType     uint8  `json:"ticketType"`
	BattleDuration int    `json:"battleDuration"`
}

type TicketScenario struct {
	gp                *GlobalParams
	params            *TicketScenarioParams
	client            *bot_client.UDPClient
	trnClient         *bot_client.UDPClient
	metrics           *report.CustomMetrics
	createRoomReq     *packets.CreateRoomReq
	matchingStartedAt time.Time
	isOwner           bool
	roomID            string
}

var _ Scenario = &TicketScenario{}

func NewTicketScenario() Scenario {
	return &TicketScenario{}
}

// // // // // Interface Functions // // // // //
func (s *TicketScenario) ParseParam(index int, params []byte) error {
	// parse scenario params
	var ticketParams *TicketScenarioParams
	err := json.Unmarshal(params, &ticketParams)
	if err != nil {
		logger.Erroru(s.params.UID, "Failed to unmarshal scenario params.", err.Error())
		return err
	}
	s.params = ticketParams
	logger.Debugu(s.params.UID, "Scenario Params. %#v", ticketParams)

	// params for create room
	s.createRoomReq = &packets.CreateRoomReq{}
	json.Unmarshal(params, s.createRoomReq)
	logger.Debugu(s.params.UID, "Params for create Room. %#v", s.createRoomReq)

	return nil
}

func (s *TicketScenario) Run(gp *GlobalParams) error {
	logger.Infou(s.params.UID, "Starting scenario for user %v", s.params.UID)

	// store GlobalParams to use it later
	s.gp = gp
	// new report
	s.metrics = report.NewCustomMetrics()

	// connect to udp server
	_, udpClient, err := bot_client.NewAndConnect(gp.Host, s.params.UID, s.params.ServerTypeMM, nil, gp.ReceiveByteSize, gp.UDPSendInterval)
	if err != nil {
		return err
	}
	s.client = udpClient

	// ver:1 cmd:218 Issue Ticket - Error case
	udpClient.RegisterOnResponse(util.CmdBuiltInVer, util.CmdMMTicket, []uint8{bot_client.ResponseBad, bot_client.ResponseError}, s.handleTicketIssueError)
	// ver:1 cmd:220 Ticket Complete
	udpClient.RegisterOnPush(util.CmdBuiltInVer, util.CmdMMTicketComplete, s.handleTicketComplete)
	// ver:1 cmd:225 Leave Ticket
	udpClient.RegisterOnResponse(util.CmdBuiltInVer, util.CmdMMTicketLeave, []uint8{bot_client.ResponseOk}, s.handleTicketLeave)
	// ver:1 cmd:224 Ticket Broadcast
	udpClient.RegisterOnPush(util.CmdBuiltInVer, util.CmdMMTicketBroadcast, s.handleTicketBroadcast)

	// start issuing ticket
	s.issueTicket()

	return nil
}

func (s *TicketScenario) OnScenarioEnd() error {
	// stop update metrics loop
	s.metrics.Stop()

	// print if the client is active and print the last command if not
	isActive := report.IsActive(s.params.UID)
	if !isActive {
		kind, ver, cmd := s.client.GetLastActivity()
		logger.Warnu(s.params.UID, "I did not have any activities more than %d seconds, last command was ver: %d cmd: %d type: %s", report.Interval, ver, cmd, kind)
	}

	// disconnect from all servers
	logger.Infou(s.params.UID, "disconnecting client...")
	s.client.Disconnect()
	if s.params.ServerTypeMM != s.params.ServerTypeTurn && s.trnClient != nil {
		s.trnClient.Disconnect()
	}

	// print report
	logger.Noticeu(s.params.UID, "result per client   === \\")
	s.metrics.Print()
	logger.Noticeu(s.params.UID, "result per client   === /")
	return nil
}

// // // // //  Scenario Functions // // // // //
func (s *TicketScenario) handleTicketIssueError(payload []byte) {
	time.Sleep(5 * time.Second)
	s.issueTicket()
}
func (s *TicketScenario) handleTicketComplete(payload []byte) {
	// add matching duration metrics
	duration := time.Since(s.matchingStartedAt)
	s.metrics.Add("MATCHING_DURATION", "", duration.Seconds())

	// parse response
	res := packet.BytesToBytesList(payload)

	// connect to turn server
	s.connectTurnServer()
	// create Room if I am the ticket owner
	s.isOwner = string(res[0]) == s.params.UID
	if s.isOwner {
		logger.Sysu(s.params.UID, "Room Owner")
		s.createRoom()
	}
}
func (s *TicketScenario) handleTicketLeave(payload []byte) {
	// get new randomised ticketType
	s.regenerateParams()

	// issue ticket again to loop the scenario
	s.issueTicket()
	logger.Infou(s.params.UID, "Restarting ticket...")
}

func (s *TicketScenario) handleTicketBroadcast(payload []byte) {
	if !s.isOwner {
		// store roomID
		s.roomID = string(payload)
		// create payload for join room
		req := payload
		req = append(req, []byte("hello")...)
		s.trnClient.RSend(util.CmdBuiltInVer, util.CmdJoinRoom, req)
		logger.Sysu(s.params.UID, "Joining room... roomID: %s", s.roomID)
	}
}

func (s *TicketScenario) issueTicket() {
	// set timestamp for matching duration metrics
	s.matchingStartedAt = time.Now()
	// count issued ticket by type
	s.metrics.Increment("ISSUE_TICKET", fmt.Sprintf("TYPE%d", s.params.TicketType))
	// ver:1 cmd:218 issue ticket
	s.client.RSend(util.CmdBuiltInVer, util.CmdMMTicket, []byte{s.params.TicketType})
}

func (s *TicketScenario) leaveTicket() {
	// ver:1 cmd:225 leave ticket
	s.client.RSend(util.CmdBuiltInVer, util.CmdMMTicketLeave, []byte{s.params.TicketType})
}

func (s *TicketScenario) connectTurnServer() {
	// connect to turn server
	if s.params.ServerTypeMM == s.params.ServerTypeTurn {
		s.trnClient = s.client
	} else {
		_, trnClient, err := bot_client.NewAndConnect(s.gp.Host, s.params.UID, s.params.ServerTypeTurn, nil, s.gp.ReceiveByteSize, s.gp.UDPSendInterval)
		if err != nil {
			logger.Erroru(s.params.UID, "Failed to get Turn server")
			return
		}
		s.trnClient = trnClient
	}
	// ver:1 cmd: 100 create room
	s.trnClient.RegisterOnResponse(util.CmdBuiltInVer, util.CmdCreateRoom, []uint8{bot_client.ResponseOk}, s.onCreateRoom)
	// ver:1 cmd: 101 join room
	s.trnClient.RegisterOnPush(util.CmdBuiltInVer, util.CmdJoinRoom, s.battle)

	// connect to turn server
	if s.params.ServerTypeMM != s.params.ServerTypeTurn {
		s.trnClient.Connect()
	}
}

func (s *TicketScenario) createRoom() {
	// ver:1 cmd: 100 create room
	bytes := packets.CreateCreateRoomReq(s.createRoomReq)
	s.trnClient.RSend(util.CmdBuiltInVer, util.CmdCreateRoom, bytes)
}

func (s *TicketScenario) leaveRoom() {
	// ver:1 cmd: 102 create room
	s.trnClient.RSend(util.CmdBuiltInVer, util.CmdLeaveRoom, []byte(s.roomID))
}

func (s *TicketScenario) onCreateRoom(payload []byte) {
	s.roomID = string(payload[4:])
	// ver:1 cmd: 224 ticket broadcast
	req := []byte{s.params.TicketType}
	req = append(req, []byte(s.roomID)...)
	s.client.RSend(util.CmdBuiltInVer, util.CmdMMTicketBroadcast, req)
}

func (s *TicketScenario) battle(payload []byte) {
	// simulate battle...
	for i := 0; i < s.params.BattleDuration; i++ {
		// use Send as battle command is normally unreliable
		s.trnClient.Send(util.CmdBuiltInVer, util.CmdBroadcastRoom, []byte("some battle command"))
		time.Sleep(time.Second)
	}

	// disconnect from turn server
	if s.params.ServerTypeMM != s.params.ServerTypeTurn {
		s.trnClient.Disconnect()
		s.trnClient = nil
	}

	// finish battle
	s.leaveRoom()
	s.leaveTicket()

}

func (s *TicketScenario) regenerateParams() {
	// generate params again to get a random value for ticketType
	bytes, _ := GenerateParams(0, s.gp.Raw.ScenarioParams)
	var ticketParams *TicketScenarioParams
	json.Unmarshal(bytes, &ticketParams)
	// update with new ticket type
	s.params.TicketType = ticketParams.TicketType
}
