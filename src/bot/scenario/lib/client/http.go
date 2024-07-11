// © 2019-2024 Diarkis Inc. All rights reserved.

package bot_client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Diarkis/diarkis/util"
)

type AuthRequest interface {
	// Rank          string `json:"Rank"`
	// Latency       string `json:"Latency"`
	// Region        string `json:"Region"`
	// QuestProgress string `json:"QuestProgress"`
}
type AuthResponse struct {
	SID        string `json:"sid"`
	Key        string `json:"encryptionKey"`
	Iv         string `json:"encryptionIV"`
	MacKey     string `json:"encryptionMacKey"`
	Port       int    `json:"serverPort"`
	ServerType string `json:"serverType"`
	Host       string `json:"serverHost"`
	Custom     interface{}
	// Token   string `json:"token"`
	// Address string `json:"UDP"`
}

type Credentials struct {
	SID    []byte
	Key    []byte
	Iv     []byte
	MacKey []byte
}

type HTTPClient struct {
	uid     string
	host    string
	servers map[string]struct {
		endpoint    *AuthResponse
		credentials *Credentials
	}
}

func NewHTTPClient(host string, uid string) *HTTPClient {
	httpClient := HTTPClient{
		host: host,
		uid:  uid,
		servers: map[string]struct {
			endpoint    *AuthResponse
			credentials *Credentials
		}{},
	}
	return &httpClient
}

func (hc *HTTPClient) Connect(serverType string, requestBody any) (int, error) {

	var req *http.Request
	var err error

	uri := fmt.Sprintf("http://%s/endpoint/type/%s/user/%s", hc.host, serverType, hc.uid)
	if requestBody == nil {
		req, err = http.NewRequest("GET", uri, nil)
	} else {
		requestBodyJSON, err := json.Marshal(requestBody)
		if err != nil {
			return 0, err
		}

		req, err = http.NewRequest("POST", uri, bytes.NewBuffer(requestBodyJSON))
		if err != nil {
			return 0, err
		}
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		if resp == nil {
			return 0, err
		}
		return resp.StatusCode, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return resp.StatusCode, err
	}
	if resp.StatusCode != 200 {
		return resp.StatusCode, util.NewError("Failed to get endpoint. %v", string(body))
	}
	logger.Debug("Got response from HTTP server. %v", string(body))

	var ar AuthResponse
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return resp.StatusCode, err
	}

	sid, _ := hex.DecodeString(ar.SID)
	key, _ := hex.DecodeString(ar.Key)
	iv, _ := hex.DecodeString(ar.Iv)
	macKey, _ := hex.DecodeString(ar.MacKey)

	server := hc.servers[serverType]
	credentials := &Credentials{
		SID:    sid,
		Key:    key,
		Iv:     iv,
		MacKey: macKey,
	}
	server.credentials = credentials
	server.endpoint = &ar
	hc.servers[serverType] = server
	return resp.StatusCode, nil
}

func (hc *HTTPClient) GetCredentials(serverType string) *Credentials {
	server, ok := hc.servers[serverType]
	if !ok {
		return nil
	}
	return server.credentials
}

func (hc *HTTPClient) GetEndpoint(serverType string) *AuthResponse {
	server, ok := hc.servers[serverType]
	if !ok {
		return nil
	}
	return server.endpoint
}

func (hc *HTTPClient) GetUserID() string {
	return hc.uid
}

func (hc *HTTPClient) SetCustomAuthResponse(cb func()) bool {

	return false
}
