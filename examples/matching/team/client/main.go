package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Diarkis/diarkis-server-template/examples/matching/custom-criteria/common"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/util"
)

var logger = log.New("CLI")

const rcvByteSize = 8000
const udpSendInterval int64 = 200

func main() {
	os.Exit(run())
}

func run() int {
	var (
		flagHost      = flag.String("host", "127.0.0.1:7000", "the address of the HTTP server")
		flagUID       = flag.String("uid", "", "the unique identifier of the client like user ID")
		flagClientKey = flag.String("clientKey", "", "the client key to authenticate with the server")
	)
	context, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	flag.Parse()

	if *flagUID == "" {
		flag.Usage()
		return 1
	}

	cli := new(client)
	cli.uid = *flagUID

	done := make(chan struct{})
	go func() {
		defer close(done)
		cli.connect(*flagHost, *flagClientKey)
	}()

	select {
	case <-context.Done():
		cli.uc.Disconnect()
		time.Sleep(time.Second)
	case <-done:
		// client disconnected
	}
	return 0
}

type client struct {
	uc  *udp.Client
	uid string
}

func (c *client) connect(host string, clientKey string) {
	udpURL := fmt.Sprintf("http://%s/endpoint/type/UDP/user/%s", host, c.uid)
	fmt.Printf("Connecting to HTTP server first: %s - clientKey = %v\n", udpURL, clientKey)
	udpDecoded := getAuthInfo(udpURL, clientKey)

	c.startConnection(udpDecoded, clientKey)
}

func validateAuthResponse(d map[string]any) error {
	var errList []error
	for _, k := range []string{"serverHost", "serverPort", "sid", "encryptionKey", "encryptionIV", "encryptionMacKey"} {
		_, ok := d[k]
		if !ok {
			errList = append(errList, fmt.Errorf("missing property %q", k))
		}
	}

	for _, k := range []string{"serverHost", "sid", "encryptionKey", "encryptionIV", "encryptionMacKey"} {
		_, ok := d[k].(string)
		if !ok {
			errList = append(errList, fmt.Errorf("property %q is not a string", k))
		}
	}

	{
		_, ok := d["serverPort"].(float64)
		if !ok {
			errList = append(errList, fmt.Errorf("property %q is not a number", "serverPort"))
		}
	}

	return errors.Join(errList...)
}

func (c *client) startConnection(udpData map[string]interface{}, clientKey string) error {
	if err := validateAuthResponse(udpData); err != nil {
		return err
	}
	port := udpData["serverPort"].(float64)
	udpAddr := fmt.Sprintf("%s:%d", udpData["serverHost"].(string), int(port))

	fmt.Printf("UDP address = %s\n", udpAddr)

	udpSid, err := hex.DecodeString(udpData["sid"].(string))
	if err != nil {
		panic("Failed to decode hex encoded string")
	}
	udpEncKey, err := hex.DecodeString(udpData["encryptionKey"].(string))
	if err != nil {
		panic("Failed to decode hex encoded string")
	}
	udpEncIV, err := hex.DecodeString(udpData["encryptionIV"].(string))
	if err != nil {
		panic("Failed to decode hex encoded string")
	}
	udpEncMacKey, err := hex.DecodeString(udpData["encryptionMacKey"].(string))
	if err != nil {
		panic("Failed to decode hex encoded string")
	}
	fmt.Printf("UDP sid         = %s\n", udpData["sid"].(string))
	fmt.Printf("UDP key         = %s\n", udpData["encryptionKey"].(string))
	fmt.Printf("UDP iv          = %s\n", udpData["encryptionIV"].(string))
	fmt.Printf("UDP mac         = %s\n", udpData["encryptionMacKey"].(string))
	c.connectUDP(udpAddr, udpSid, udpEncKey, udpEncIV, udpEncMacKey, clientKey)

	return nil
}

func (c *client) connectUDP(addr string, sid []byte, key []byte, iv []byte, mackey []byte, clientKey string) {
	UDPClient := udp.New(rcvByteSize, udpSendInterval)
	UDPClient.SetID(c.uid)

	c.uc = UDPClient

	c.uc.SetClientKey(clientKey)
	c.uc.SetEncryptionKeys(sid, key, iv, mackey)
	c.uc.OnResponse(c.handleOnResponseUDP)
	c.uc.OnPush(c.handleOnPushUDP)
	c.uc.OnConnect(c.handleOnConnectUDP)

	done := make(chan struct{})
	c.uc.OnDisconnect(func() {
		close(done)
	})

	c.uc.Connect(addr)
	<-done
}

func (c *client) handleOnConnectUDP() {
	logger.Info("Connected UDP")
	c.startMatching()
}

func (c *client) handleOnResponseUDP(ver uint8, cmd uint16, status uint8, payload []byte) {
	logger.Debugf("UDP onResponse", "ver", ver, "cmd", cmd, "status", status, "payload", string(payload))
	c.handleResponse(ver, cmd, status, payload)
}

func (c *client) handleOnPushUDP(ver uint8, cmd uint16, payload []byte) {
	logger.Debugf("UDP onPush", "ver", ver, "cmd", cmd, "payload", string(payload))
	if ver == util.CmdBuiltInVer {
		switch cmd {
		case util.CmdMMTicketComplete:
			logger.Infof("UDP onPush", "ver", ver, "cmd", cmd, "payload", string(payload))
			c.handleTicketComplete(payload)
		case util.CmdMMTicketErr:
			logger.Infof("UDP onPush", "ver", ver, "cmd", cmd, "payload", hex.EncodeToString(payload))
			c.handleTicketError(payload)
		case util.CmdMMTicketBroadcast:
			logger.Infof("received ticket broadcast", "ver", ver, "cmd", cmd, "payload", string(payload))
			c.handleTicketBroadcast(payload)
		}
	} else if ver == common.AppVersion {
		switch cmd {
		case common.MatchingTicketBattleBroadcastCmd:
			logger.Info("Team matched: %s", payload)
			logger.Info("test is finished, disconnect")
			c.uc.Disconnect()
		case common.MatchingTicketMemberLeftCmd:
			c.handleCustomMatchingTicketLeft(payload)
		}
	}
}

func (c *client) handleResponse(ver uint8, cmd uint16, status uint8, payload []byte) {
	if ver != common.AppVersion {
		return
	}
	if cmd == common.MatchingStartCmd {
		// We received the matching start response.
		if status == server.Ok {
			logger.Info("Matching successfully started")
		} else {
			logger.Warnf("Matching start error", "status", status, "payload", payload)
			// We cannot call Disconnect from handleResponse.
			go func() {
				time.Sleep(time.Second)
				logger.Info("test is finished, disconnect")
				c.uc.Disconnect()
			}()
			return
		}
	}
}

func (c *client) handleTicketComplete(payload []byte) {
	var e common.MatchingComplete
	err := json.Unmarshal(payload, &e)
	if err != nil {
		logger.Error("fail to parse matching complete payload. %v (%#v)", err, payload)
		return
	}

	logger.Info("matching complete: %+v", e)
	switch e.TicketType {
	case common.TeamTicketType:
		logger.Info("Team created. Members %q", e.CandidateIDs)
	}
}

func (c *client) handleTicketError(payload []byte) {
	code := binary.BigEndian.Uint16(payload[:2])
	message := string(payload[2:])
	logger.Info("matching ticket error code %d message %s", code, message)
	if code == 4009 {
		// timeout
		logger.Info("received matching timeout")
		logger.Info("test is finished, disconnect")
		c.uc.Disconnect()
		return
	}
}

func (c *client) handleTicketBroadcast(payload []byte) {
	_ = payload
	logger.Info("received matching ticket broadcast %s", payload)
	logger.Info("test is finished, disconnect")
	c.uc.Disconnect()
}

func (c *client) handleCustomMatchingTicketLeft(payload []byte) {
	var evt common.MatchingTicketMemberLeft
	json.Unmarshal(payload, &evt)
	if evt.MemberID == evt.OwnerID {
		logger.Info("The owner of the ticket left, finish the run.")
		logger.Info("test is finished, disconnect")
		c.uc.Disconnect()
	}
}

func (c *client) startMatching() error {
	var params common.MatchingParams

	logger.Info("start matching")

	b, _ := json.Marshal(params)
	c.uc.RSend(common.AppVersion, common.MatchingStartCmd, b)
	return nil
}

func getAuthInfo(url string, clientKey string) map[string]interface{} {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Sprintf("Error %v", err))
	}
	req.Header.Add("ClientKey", clientKey)
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Error %v", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read the HTTP response: %v", err))
	}
	decoded := make(map[string]interface{})
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode JSON response: %v - %v", err, string(body)))
	}
	return decoded
}
