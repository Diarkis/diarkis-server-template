// © 2019-2024 Diarkis Inc. All rights reserved.

package scenarios

import (
	"encoding/json"

	bot_client "github.com/Diarkis/diarkis-server-template/bot/scenario/lib/client"
)

type SessionScenarioParams struct {
	ServerType string `json:"serverType"`
	UID        string `json:"userID"`
}

type SessionScenario struct {
	params *SessionScenarioParams
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

func (s *SessionScenario) Run(globalParams *GlobalParams) error {
	logger.Notice("Starting scenario for user %v", s.params.UID)

	_, udpClient, err := bot_client.NewAndConnect(globalParams.Host, s.params.UID, s.params.ServerType, nil, globalParams.ReceiveByteSize, globalParams.UDPSendInterval)
	if err != nil {
		return err
	}

	udpClient.Connect()
	udpClient.Disconnect()
	return nil
}

func (s *SessionScenario) OnIdle() {
	// do nothing
	return
}

func (s *SessionScenario) OnScenarioEnd() error {
	return nil
}
