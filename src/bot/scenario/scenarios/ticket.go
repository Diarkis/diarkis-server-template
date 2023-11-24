package scenarios

import (
	"encoding/json"
	"time"

	bot_client "{0}/bot/scenario/lib/client"
	"{0}/bot/scenario/lib/report"
	"{0}/bot/scenario/packets"

	"github.com/Diarkis/diarkis/util"
)

type TicketScenarioParams struct {
	ServerType string `json:"serverType"`
	UID        string `json:"userID"`
	TicketType uint8  `json:"ticketType"`
}

type TicketScenario struct {
	params            *TicketScenarioParams
	client            *bot_client.UDPClient
	metrics           *report.CustomMetrics
	leaveTicketReq    *packets.LeaveMatchingReq
	currentTicketType uint8
}

var _ Scenario = &TicketScenario{}

func NewTicketScenario() Scenario {
	return &TicketScenario{}
}

// // // // // Interface Functions // // // // //
func (s *TicketScenario) ParseParam(index int, params []byte) error {
	var ticketParams *TicketScenarioParams
	err := json.Unmarshal(params, &ticketParams)
	if err != nil {
		logger.Erroru(s.params.UID, "Failed to unmarshal scenario params.", err.Error())
		return err
	}
	s.params = ticketParams
	logger.Debugu(s.params.UID, "Parsed Params. %#v", ticketParams)

	// params for ticket leave
	s.leaveTicketReq = packets.ParseLeaveMatchingReq(params)
	// todo:
	logger.Debugu(s.params.UID, "Parsed Params. %#v", s.leaveTicketReq)

	return nil
}

func (s *TicketScenario) Run(globalParams *GlobalParams) error {
	logger.Noticeu(s.params.UID, "Starting scenario for user %v", s.params.UID)

	_, udpClient, err := bot_client.NewAndConnect(globalParams.Host, s.params.UID, s.params.ServerType, nil, globalParams.ReceiveByteSize, globalParams.UDPSendInterval)
	if err != nil {
		return err
	}

	s.client = udpClient
	s.metrics = report.NewCustomMetrics()
	// ver:1 cmd:218 Issue Ticket - Error case
	udpClient.RegisterOnResponse(util.CmdBuiltInVer, util.CmdMMTicket, []uint8{bot_client.ResponseBad, bot_client.ResponseError}, s.handleTicketIssueError)
	// ver:1 cmd:220 Ticket Complete
	udpClient.RegisterOnPush(util.CmdBuiltInVer, util.CmdMMTicketComplete, s.handleTicketComplete)
	// ver:1 cmd:225 Leave Ticket
	// udpClient.RegisterOnResponse(util.CmdBuiltInVer, util.CmdMMTicketLeave, []uint8{bot_client.ResponseOk}, s.handleTicketLeave)

	// udpClient.OnPush(func(ver uint8, cmd uint16, payload []byte) {
	// 	if ver == 1 && cmd == 220 {
	// 		s.leaveTicket()
	// 	}
	// })
	// udpClient.OnResponse(func(ver uint8, cmd uint16, status uint8, payload []byte) {
	// 	if ver == 1 && cmd == 225 && status == 1 {
	// 		time.Sleep(3 * time.Second)
	// 		s.issueTicket()
	// 	}
	// 	if ver == 1 && cmd == 218 && status != 1 {
	// 		s.issueTicket()
	// 	}
	// })
	s.issueTicket()

	return nil
}

func (s *TicketScenario) OnScenarioEnd() error {
	s.client.Disconnect()
	logger.Infou(s.params.UID, "client disconnected")
	return nil
}

func (s *TicketScenario) WriteReport() (*report.Report, map[string]*report.Report) {
	return nil, nil
}

// // // // //  Scenario Functions // // // // //
func (s *TicketScenario) handleTicketIssueError(payload []byte) {
	s.issueTicket()
}
func (s *TicketScenario) handleTicketComplete(payload []byte) {
	s.leaveTicket()
}
func (s *TicketScenario) handleTicketLeave(payload []byte) {
	time.Sleep(3 * time.Second)
	s.issueTicket()
}

func (s *TicketScenario) issueTicket() {
	// ticketType := uint8(util.RandomInt(0, 1))
	// bytes := packets.CreateLeaveMatchingReq(s.leaveTicketReq)

	// todo
	s.metrics.Increment("ISSUE_TICKET")
	s.client.RSend(1, 218, []byte{s.params.TicketType})
	// s.currentTicketType = ticketType
}

func (s *TicketScenario) leaveTicket() {
	s.client.RSend(1, 225, []byte{s.params.TicketType})
	// s.client.RSend(1, 225, []byte{s.currentTicketType})
}
