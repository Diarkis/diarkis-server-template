// © 2019-2024 Diarkis Inc. All rights reserved.

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type EndpointResponse struct {
	ServerHost       string `json:"serverHost"`
	ServerPort       int    `json:"serverPort"`
	Sid              string `json:"sid"`
	EncryptionKey    string `json:"encryptionKey"`
	EncryptionIV     string `json:"encryptionIV"`
	EncryptionMacKey string `json:"encryptionMacKey"`
	ServerType       string `json:"serverType"`
}

func Endpoint(host, uid, serverType string) (EndpointResponse, error) {
	url := fmt.Sprintf("http://%s/endpoint/type/%v/user/%v", host, strings.ToUpper(serverType), uid)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return EndpointResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return EndpointResponse{}, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return EndpointResponse{}, err
	}
	endpointResp := EndpointResponse{}
	err = json.Unmarshal(body, &endpointResp)
	if err != nil {

		return EndpointResponse{}, err
	}

	return endpointResp, nil
}
