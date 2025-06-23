package main

import (
	"context"
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

	"github.com/Diarkis/diarkis/client/go/modules/room"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/log"
)

var logger = log.New("CLI")

const rcvByteSize = 8000
const udpSendInterval int64 = 200

func main() {
	os.Exit(run())
}

func run() int {
	var (
		flagHost       = flag.String("host", "127.0.0.1:7000", "The address of the HTTP server.")
		flagUID        = flag.String("uid", "", "The unique identifier of the client like user ID.")
		flagCreateRoom = flag.Bool("create-room", false, "Create a room upon connection.")
	)

	context, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	flag.Parse()

	cli := new(client)
	cli.uid = *flagUID
	cli.createRoom = *flagCreateRoom

	done := make(chan struct{})
	go func() {
		defer close(done)
		cli.connect(*flagHost, "")
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
	uc         *udp.Client
	uid        string
	createRoom bool
	room       *room.Room
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
	c.room = &room.Room{}
	c.room.SetupAsUDP(c.uc)

	c.uc.SetClientKey(clientKey)
	c.uc.SetEncryptionKeys(sid, key, iv, mackey)
	c.uc.OnConnect(func() {
		if c.createRoom {
			c.room.Create(2, false, true, 10, 100)
		}
	})

	done := make(chan struct{})
	c.uc.OnDisconnect(func() {
		close(done)
	})

	c.uc.Connect(addr)
	<-done
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
