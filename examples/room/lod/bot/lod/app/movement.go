// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"math/rand"
)

// moveBot moves the bot randomly within map boundaries
func moveBot(bot *bot) {
	// Random movement in 8 directions (N, NE, E, SE, S, SW, W, NW)
	dx := int32(rand.Intn(3) - 1) // -1, 0, or 1
	dy := int32(rand.Intn(3) - 1) // -1, 0, or 1

	// Apply movement speed
	newX := bot.x + (dx * MovementSpeed)
	newY := bot.y + (dy * MovementSpeed)

	// Clamp to map boundaries
	if newX < MapMinX {
		newX = MapMinX
	} else if newX > MapMaxX {
		newX = MapMaxX
	}

	if newY < MapMinY {
		newY = MapMinY
	} else if newY > MapMaxY {
		newY = MapMaxY
	}

	// Update bot position
	bot.x = newX
	bot.y = newY

	logger.Debug("Bot moved",
		"bot.uid", bot.uid,
		"x", bot.x,
		"y", bot.y)
}
