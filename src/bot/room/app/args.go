// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"log/slog"
	"os"
	"strconv"
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
	hostKey                          = "HOST"
	defaultHost                      = "localhost:7000"
	botsKey                          = "BOTS"
	defaultBots                      = 10
	authIntervalKey                  = "AUTH_INTERVAL"
	defaultAuthInterval              = 100
	roomSizeKey                      = "ROOM_SIZE"
	defaultRoomSize                  = 10
	packetIntervalMillisecondKey     = "PACKET_INTERVAL"
	defaultPacketIntervalMilliSecond = 100
	packetSizeKey                    = "PACKET_SIZE"
	defaultPacketSize                = 100
	logLevelKey                      = "LOG_LEVEL"
	defaultLogLevel                  = "info"
	protocolKey                      = "PROTOCOL"
	defaultProtocol                  = "udp"
	udpClientSendIntervalKey         = "UDP_CLIENT_SEND_INTERVAL"
	defaultUdpClientSendInterval     = 5
)

func parseArgs() {
	programLevel := new(slog.LevelVar)
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
	logger.Info("bot args",
		"host", hostKey,
		"bots", botsKey,
		"authInterval", authIntervalKey,
		"roomSize", roomSizeKey,
		"packetInterval", packetIntervalMillisecondKey,
		"packetSize", packetSizeKey,
		"logLevel", logLevelKey,
		"protocol", protocolKey,
	)

	var ok bool
	host, ok = getParams(hostKey)
	if !ok {
		logger.Info("host is set to default value", "host", defaultHost)
		host = defaultHost
	}

	bots, ok = getParamsInt(botsKey)
	if !ok {
		logger.Info("bots is set to default value", "bots", defaultBots)
		bots = defaultBots
	}

	authInterval, ok = getParamsInt(authIntervalKey)
	if !ok {
		logger.Info("endpoint interval is set to default value", "authInterval", defaultAuthInterval)
		authInterval = defaultAuthInterval
	}

	roomSize, ok = getParamsInt(roomSizeKey)
	if !ok {
		logger.Info("room size is set to default value", "roomSize", defaultRoomSize)
		roomSize = defaultRoomSize
	}

	packetInterval, ok = getParamsInt(packetIntervalMillisecondKey)
	if !ok {
		logger.Info("packet interval is set to default value", "packetInternal", defaultPacketIntervalMilliSecond)
		packetInterval = defaultPacketIntervalMilliSecond
	}

	packetSize, ok = getParamsInt(packetSizeKey)
	if !ok {
		logger.Info("packet size is set to default value", "packetSize", defaultPacketSize)
		packetSize = defaultPacketSize
	}

	logLevel, ok = getParams(logLevelKey)
	if !ok {
		logger.Info("log level is set to default value", "logLevel", defaultLogLevel)
		logLevel = defaultLogLevel
	}

	protocol, ok = getParams(protocolKey)
	if !ok {
		logger.Info("protocol is set to default value", "protocol", defaultProtocol)
		protocol = defaultProtocol
	}

	udpClientSendInterval, ok = getParamsInt(udpClientSendIntervalKey)
	if !ok {
		logger.Info("udp client send interval is set to default value: %v", "udpClientSendInterval", defaultUdpClientSendInterval)
		udpClientSendInterval = defaultUdpClientSendInterval
	}

	switch logLevel {
	case "debug":
		programLevel.Set(slog.LevelDebug)
	case "info":
		programLevel.Set(slog.LevelInfo)
	case "warn":
		programLevel.Set(slog.LevelWarn)
	case "error":
		programLevel.Set(slog.LevelError)
	default:
		programLevel.Set(slog.LevelInfo)
	}
}

func getParams(key string) (string, bool) {
	v, ok := os.LookupEnv(key)
	if ok {
		return v, true
	}
	return "", false
}

func getParamsInt(key string) (int, bool) {
	v, ok := os.LookupEnv(key)
	if ok {
		retV, err := strconv.Atoi(v)
		if err != nil {
			slog.Error("Error converting env var: %v", key)
			os.Exit(1)
			return 0, false
		}
		return retV, true
	}

	return 0, false
}
