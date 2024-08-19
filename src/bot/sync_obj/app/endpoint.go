// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type endpointResponse struct {
	ServerHost       string `json:"serverHost"`
	ServerPort       int    `json:"serverPort"`
	Sid              string `json:"sid"`
	EncryptionKey    string `json:"encryptionKey"`
	EncryptionIV     string `json:"encryptionIV"`
	EncryptionMacKey string `json:"encryptionMacKey"`
	ServerType       string `json:"serverType"`
}

func endpoint(host, uid, serverType string) (endpointResponse, error) {
	url := fmt.Sprintf("%s/endpoint/type/%v/user/%v", host, strings.ToUpper(serverType), uid)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return endpointResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("endpoint API returns error", "err", err)
		return endpointResponse{}, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return endpointResponse{}, err
	}
	endpointResp := endpointResponse{}
	err = json.Unmarshal(body, &endpointResp)
	if err != nil {
		logger.Error("endpoint response is invalid")
		return endpointResponse{}, err
	}
	logger.Debug("endpoint response", "url", url, "resp", endpointResp)
	return endpointResp, nil
}
