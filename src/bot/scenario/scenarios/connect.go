// © 2019-2024 Diarkis Inc. All rights reserved.

package scenarios

import (
	"encoding/json"

	bot_client "github.com/Diarkis/diarkis-server-template/bot/scenario/lib/client"
)

type ConnectParams struct {
	ServerType string `json:"serverType"`
	UID        string `json:"userID"`
}

type ConnectScenario struct {
	params *ConnectParams
}

var _ Scenario = &ConnectScenario{}

func NewConnectScenario() Scenario {
	return &ConnectScenario{}
}

// // // // // Interface Functions // // // // //
func (s *ConnectScenario) GetUserID() string {
	return s.params.UID
}

func (s *ConnectScenario) ParseParam(index int, params []byte) error {
	var connectParams *ConnectParams
	err := json.Unmarshal(params, &connectParams)
	if err != nil {
		logger.Error("Failed to unmarshal scenario params.", err.Error())
		return err
	}
	s.params = connectParams
	logger.Sys("Parsed Params. %#v", connectParams)
	return nil
}

func (s *ConnectScenario) Run(globalParams *GlobalParams) error {
	logger.Notice("Starting scenario for user %v", s.params.UID)

	_, udpClient, err := bot_client.NewAndConnect(globalParams.Host, s.params.UID, s.params.ServerType, nil, globalParams.ReceiveByteSize, globalParams.UDPSendInterval)
	if err != nil {
		return err
	}

	udpClient.Connect()
	udpClient.Disconnect()
	return nil
}

func (s *ConnectScenario) OnIdle() {
	// do nothing
	return
}

func (s *ConnectScenario) OnScenarioEnd() error {
	return nil
}
