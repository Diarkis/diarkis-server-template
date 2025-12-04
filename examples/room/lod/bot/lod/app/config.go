// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"github.com/Diarkis/diarkis/config"
)

// Bot config values
var (
	BotStartX        int32 // Starting X position for bots
	BotStartY        int32 // Starting Y position for bots
	MapMinX          int32 // Minimum X boundary
	MapMaxX          int32 // Maximum X boundary
	MapMinY          int32 // Minimum Y boundary
	MapMaxY          int32 // Maximum Y boundary
	MovementInterval int32 // Movement update interval in milliseconds
	MovementSpeed    int32 // Movement distance per update
)

const (
	configPath              = "bot/lod/bot.json"
	botStartXKey            = "BotStartX"
	botStartYKey            = "BotStartY"
	mapMinXKey              = "MapMinX"
	mapMaxXKey              = "MapMaxX"
	mapMinYKey              = "MapMinY"
	mapMaxYKey              = "MapMaxY"
	movementIntervalKey     = "MovementInterval"
	movementSpeedKey        = "MovementSpeed"
	defaultBotStartX        = 0
	defaultBotStartY        = 0
	defaultMapMinX          = -25000
	defaultMapMaxX          = 25000
	defaultMapMinY          = -25000
	defaultMapMaxY          = 25000
	defaultMovementInterval = 100
	defaultMovementSpeed    = 50
)

// loadBotLodConfig loads bot configuration from file
func loadBotLodConfig() {
	config.Load("BotConfig", configPath)

	BotStartX = config.GetAsInt32("BotConfig", botStartXKey, defaultBotStartX)
	BotStartY = config.GetAsInt32("BotConfig", botStartYKey, defaultBotStartY)
	MapMinX = config.GetAsInt32("BotConfig", mapMinXKey, defaultMapMinX)
	MapMaxX = config.GetAsInt32("BotConfig", mapMaxXKey, defaultMapMaxX)
	MapMinY = config.GetAsInt32("BotConfig", mapMinYKey, defaultMapMinY)
	MapMaxY = config.GetAsInt32("BotConfig", mapMaxYKey, defaultMapMaxY)
	MovementInterval = config.GetAsInt32("BotConfig", movementIntervalKey, defaultMovementInterval)
	MovementSpeed = config.GetAsInt32("BotConfig", movementSpeedKey, defaultMovementSpeed)

	logger.Info("Bot config loaded",
		"BotStartX", BotStartX,
		"BotStartY", BotStartY,
		"MapMinX", MapMinX,
		"MapMaxX", MapMaxX,
		"MapMinY", MapMinY,
		"MapMaxY", MapMaxY,
		"MovementInterval", MovementInterval,
		"MovementSpeed", MovementSpeed)
}
