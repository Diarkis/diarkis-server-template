// © 2019-2024 Diarkis Inc. All rights reserved.
package bot_client

import (
	"github.com/Diarkis/diarkis-server-template/bot/scenario/lib/log"
	"strconv"
	"strings"

	"github.com/Diarkis/diarkis/util"
)

var logger = log.New("BOT/CLI")

// NewAndConnect gets an endpoint for the transport and connect it.
func NewAndConnect(host string, userID string, serverType string, httpRequestBody any, rcvMaxSize int, interval int64) (*HTTPClient, *UDPClient, error) {

	httpClient := NewHTTPClient(host, userID)
	logger.Infou(userID, "HTTP client created %#v", httpClient)

	_, err := httpClient.Connect(serverType, httpRequestBody)
	if err != nil {
		return httpClient, nil, util.StackError(util.NewError("Failed to connect to HTTP."), err)
	}

	credentials := httpClient.GetCredentials(serverType)
	if credentials == nil {
		return httpClient, nil, util.NewError("Cannot get credentials from HTTP response")
	}
	ep := httpClient.GetEndpoint(serverType)
	if ep.Host == "" || ep.Port == 0 {
		return httpClient, nil, util.NewError("Cannot get host and port from HTTP response. Response could be v1...")
	}
	endpoint := strings.Join([]string{ep.Host, strconv.Itoa(ep.Port)}, ":")
	udpClient, err := NewUDPClient(userID, endpoint, credentials, rcvMaxSize, interval)
	if err != nil {
		return httpClient, udpClient, util.StackError(util.NewError("Failed to connect to UDP."), err)
	}
	logger.Infou(userID, "UDP client created %#v", udpClient.client)

	udpClient.Connect()
	return httpClient, udpClient, nil

}
