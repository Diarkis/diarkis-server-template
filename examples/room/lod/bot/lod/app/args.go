// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"github.com/Diarkis/diarkis/config"
)

// args
var (
	host                  string
	bots                  int
	authInterval          int
	roomSize              int
	packetInterval        int
	packetSize            int
	logLevel              string
	protocol              string
	udpClientSendInterval int
)

const (
	hostKey                          = "Host"
	defaultHost                      = "localhost:7000"
	botsKey                          = "Bots"
	defaultBots                      = 10
	authIntervalKey                  = "AuthInterval"
	defaultAuthInterval              = 100
	roomSizeKey                      = "RoomSize"
	defaultRoomSize                  = 10
	packetIntervalMillisecondKey     = "PacketIntervalMillisecond"
	defaultPacketIntervalMilliSecond = 100
	packetSizeKey                    = "PacketSize"
	defaultPacketSize                = 100
	logLevelKey                      = "LogLevel"
	defaultLogLevel                  = "info"
	protocolKey                      = "Protocol"
	defaultProtocol                  = "udp"
	udpClientSendIntervalKey         = "UDPClientSendInterval"
	defaultUdpClientSendInterval     = 5
)

func loadBotConfig() {
	const configPath = "bot/lod/bot.json"
	const configName = "BotConfig"

	config.Load(configName, configPath)

	host = config.GetAsString(configName, hostKey, defaultHost)
	bots = int(config.GetAsInt32(configName, botsKey, int32(defaultBots)))
	authInterval = int(config.GetAsInt32(configName, authIntervalKey, int32(defaultAuthInterval)))
	roomSize = int(config.GetAsInt32(configName, roomSizeKey, int32(defaultRoomSize)))
	packetInterval = int(config.GetAsInt32(configName, packetIntervalMillisecondKey, int32(defaultPacketIntervalMilliSecond)))
	packetSize = int(config.GetAsInt32(configName, packetSizeKey, int32(defaultPacketSize)))
	logLevel = config.GetAsString(configName, logLevelKey, defaultLogLevel)
	protocol = config.GetAsString(configName, protocolKey, defaultProtocol)
	udpClientSendInterval = int(config.GetAsInt32(configName, udpClientSendIntervalKey, int32(defaultUdpClientSendInterval)))

	logger.Info("Bot configuration loaded",
		"configPath", configPath,
		"host", host,
		"bots", bots,
		"authInterval", authInterval,
		"roomSize", roomSize,
		"packetInterval", packetInterval,
		"packetSize", packetSize,
		"logLevel", logLevel,
		"protocol", protocol,
		"udpClientSendInterval", udpClientSendInterval)
}
