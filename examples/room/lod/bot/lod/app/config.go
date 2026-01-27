// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Diarkis/diarkis/config"
)

// Settings holds all bot configuration
type Settings struct {
	// Connection
	Host                  string
	Protocol              string
	AuthInterval          int
	UDPClientSendInterval int

	// Bot logic
	BotsCount         int
	RoomSize          int
	AverageRoomMember int
	PacketInterval    int
	PacketSize        int
	LogLevel          string

	// Movement & Map
	BotStartX        int32
	BotStartY        int32
	MapMinX          int32
	MapMaxX          int32
	MapMinY          int32
	MapMaxY          int32
	MovementInterval int32
	MovementSpeed    int32
}

func (s *Settings) String() string {
	return fmt.Sprintf("Settings{Host: %s, Protocol: %s, AuthInterval: %d, UDPClientSendInterval: %d, BotsCount: %d, RoomSize: %d, AverageRoomMember: %d, PacketInterval: %d, PacketSize: %d, LogLevel: %s, BotStartX: %d, BotStartY: %d, MapMinX: %d, MapMaxX: %d, MapMinY: %d, MapMaxY: %d, MovementInterval: %d, MovementSpeed: %d}",
		s.Host,
		s.Protocol,
		s.AuthInterval,
		s.UDPClientSendInterval,
		s.BotsCount,
		s.RoomSize,
		s.AverageRoomMember,
		s.PacketInterval,
		s.PacketSize,
		s.LogLevel,
		s.BotStartX,
		s.BotStartY,
		s.MapMinX,
		s.MapMaxX,
		s.MapMinY,
		s.MapMaxY,
		s.MovementInterval,
		s.MovementSpeed)
}

func LoadSettings() *Settings {
	const configPath = "bot/lod/bot.json"
	const configName = "BotConfig"
	config.Load(configName, configPath)

	s := &Settings{}

	s.Host = config.GetAsString(configName, "Host", "localhost:7000")
	s.BotsCount = int(config.GetAsInt32(configName, "Bots", 10))
	s.AuthInterval = int(config.GetAsInt32(configName, "AuthInterval", 100))
	s.RoomSize = int(config.GetAsInt32(configName, "RoomSize", 10))
	s.AverageRoomMember = int(config.GetAsInt32(configName, "AverageRoomMember", 0))
	s.PacketInterval = int(config.GetAsInt32(configName, "PacketIntervalMillisecond", 100))
	s.PacketSize = int(config.GetAsInt32(configName, "PacketSize", 100))
	s.LogLevel = config.GetAsString(configName, "LogLevel", "info")
	s.Protocol = config.GetAsString(configName, "Protocol", "udp")
	s.UDPClientSendInterval = int(config.GetAsInt32(configName, "UDPClientSendInterval", 5))
	s.BotStartX = config.GetAsInt32(configName, "BotStartX", 0)
	s.BotStartY = config.GetAsInt32(configName, "BotStartY", 0)
	s.MapMinX = config.GetAsInt32(configName, "MapMinX", -25000)
	s.MapMaxX = config.GetAsInt32(configName, "MapMaxX", 25000)
	s.MapMinY = config.GetAsInt32(configName, "MapMinY", -25000)
	s.MapMaxY = config.GetAsInt32(configName, "MapMaxY", 25000)
	s.MovementInterval = config.GetAsInt32(configName, "MovementInterval", 100)
	s.MovementSpeed = config.GetAsInt32(configName, "MovementSpeed", 50)
	return s
}

func setupLogger(settings *Settings) *slog.Logger {
	programLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))

	switch settings.LogLevel {
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
	return logger
}
