// © 2019-2026 Diarkis Inc. All rights reserved.

package scenarios

import (
	"encoding/binary"
	"encoding/json"
	"math/rand/v2"
	"strconv"
	"sync/atomic"
	"time"

	bot_client "github.com/Diarkis/diarkis-server-template/bot/scenario/lib/client"
	"github.com/Diarkis/diarkis-server-template/bot/scenario/lib/report"
	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/util"
)

const StateDMReady = "DMReady"

// DM_LATENCY の単位は秒。
var dmLatencyBuckets = []float64{0.001, 0.002, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10}

func init() {
	report.SetCustomMetricsAsHistogram("DM_LATENCY", dmLatencyBuckets)
}

type DMScenarioParams struct {
	ServerType string `json:"serverType"`
	// userID は連番の数値にする。index が 2n と 2n+1 のユーザーがペアになる。
	UID string `json:"userID"`
	// メッセージのバイト数。先頭 8 バイトに送信時刻を入れる。
	MessageSize int `json:"messageSize"`
	// 1 ユーザーの送信間隔 (ミリ秒)。
	SendInterval int `json:"sendInterval"`
}

type DMScenario struct {
	gp      *GlobalParams
	params  *DMScenarioParams
	peerUID string
	client  *bot_client.UDPClient
	stopped atomic.Bool
}

var _ Scenario = &DMScenario{}

func NewDMScenario() Scenario {
	return &DMScenario{}
}

// // // // // Interface Functions // // // // //
func (s *DMScenario) GetUserID() string {
	return s.params.UID
}

func (s *DMScenario) ParseParam(index int, params []byte) error {
	var dmParams *DMScenarioParams
	err := json.Unmarshal(params, &dmParams)
	if err != nil {
		logger.Error("Failed to unmarshal scenario params.", err.Error())
		return err
	}
	if dmParams.MessageSize < 8 {
		dmParams.MessageSize = 8
	}
	if dmParams.SendInterval <= 0 {
		dmParams.SendInterval = 1000
	}
	s.params = dmParams

	uid, err := strconv.Atoi(dmParams.UID)
	if err != nil {
		logger.Error("DM scenario requires a sequential numeric userID. UID:%s", dmParams.UID)
		return err
	}
	s.peerUID = strconv.Itoa(uid - index + (index ^ 1))
	logger.Sys("Parsed Params. %#v PeerUID:%s", dmParams, s.peerUID)
	return nil
}

func (s *DMScenario) Run(gp *GlobalParams) error {
	s.gp = gp
	if s.peerUID == "" {
		return util.NewError("DM scenario requires a sequential numeric userID. UID:%s", s.GetUserID())
	}

	_, udpClient, err := bot_client.New(gp.Host, s.GetUserID(), s.params.ServerType, nil, gp.ReceiveByteSize, gp.UDPSendInterval)
	if err != nil {
		report.IncrementCustomMetrics("DM_CONNECT_ERROR", "")
		return err
	}
	s.client = udpClient
	s.registerHandlers()
	s.client.Connect()

	go s.sendLoop()
	return nil
}

func (s *DMScenario) OnIdle() {
	// do nothing
}

func (s *DMScenario) OnScenarioEnd() error {
	s.stopped.Store(true)
	// Disconnect は送信ループの終了を待つ。全ユーザーを同時に止めるため、切断は並行で行う。
	go s.client.Disconnect()
	return nil
}

func (s *DMScenario) registerHandlers() {
	uid := s.GetUserID()
	s.client.OnConnect(func() {
		report.IncrementCustomMetrics("DM_CONNECTED", "")
		s.gp.UserState.Set(uid, StateDMReady, true)
	})
	s.client.OnDisconnect(func() {
		report.IncrementCustomMetrics("DM_DISCONNECTED", "")
		s.gp.UserState.Set(uid, StateDMReady, false)
	})
	s.client.RegisterOnResponse(util.CmdBuiltInVer, util.CmdDMSend, []uint8{bot_client.ResponseOk}, func(_ []byte) {
		report.IncrementCustomMetrics("DM_SEND_OK", "")
	})
	s.client.RegisterOnResponse(util.CmdBuiltInVer, util.CmdDMSend, []uint8{bot_client.ResponseBad, bot_client.ResponseError}, func(payload []byte) {
		report.IncrementCustomMetrics("DM_SEND_ERROR", "")
		logger.Warnu(uid, "DM send failed. %s", string(payload))
	})
	s.client.RegisterOnPush(util.CmdBuiltInVer, util.CmdDMSend, s.handleMessage)
}

// handleMessage は受信した DM の先頭 8 バイトから遅延を計算する。送信者は同じ bot プロセスにいる。
func (s *DMScenario) handleMessage(payload []byte) {
	list := packet.BytesToBytesList(payload)
	if len(list) != 2 || len(list[1]) < 8 {
		report.IncrementCustomMetrics("DM_RECEIVED_INVALID", "")
		return
	}
	sentAt := int64(binary.BigEndian.Uint64(list[1][:8]))
	report.IncrementCustomMetrics("DM_RECEIVED", "")
	report.AddCustomMetrics("DM_LATENCY", "", time.Since(time.Unix(0, sentAt)).Seconds())
}

func (s *DMScenario) isPeerReady() bool {
	ready, _ := s.gp.UserState.Get(s.peerUID, StateDMReady).(bool)
	return ready
}

func (s *DMScenario) sendLoop() {
	// ペアの相手が接続するまで送らない。
	for !s.isPeerReady() {
		if s.stopped.Load() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	interval := time.Duration(s.params.SendInterval) * time.Millisecond
	// ユーザーごとに送信時刻をずらす。
	time.Sleep(rand.N(interval))

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	message := make([]byte, s.params.MessageSize)
	for !s.stopped.Load() {
		binary.BigEndian.PutUint64(message, uint64(time.Now().UnixNano()))
		s.send(message)
		<-ticker.C
	}
}

// send はクライアントの dm モジュールと同じ形式で送る。先頭 1 バイトは互換用のフラグ。
func (s *DMScenario) send(message []byte) {
	payload := append([]byte{0}, packet.BytesListToBytes([][]byte{[]byte(s.peerUID), message})...)
	s.client.RSend(util.CmdBuiltInVer, util.CmdDMSend, payload)
	report.IncrementCustomMetrics("DM_SENT", "")
}
