// © 2019-2024 Diarkis Inc. All rights reserved.

package app

import (
	"flag"
	"log/slog"
	"os"
)

var (
	host           = flag.String("host", "localhost:7000", "host")
	bots           = flag.Int("bots", 200, "bots")
	authInterval   = flag.Int("authInterval", 1000, "authInterval milliseconds")
	roomSize       = flag.Int("roomSize", 200, "roomSize")
	packetInterval = flag.Int("packetInterval", 100, "packetInterval milliseconds")
	logLevel       = flag.String("logLevel", "info", "logLevel")
	mode           = flag.String("mode", "incr", "object update mode[update, incr]")
	changePercent  = flag.Int("changePercent", 2, "obj value changePercent")
)

func parseArgs() {
	flag.Parse()

	programLevel := new(slog.LevelVar)
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
	logger.Info("bot args",
		"host", *host,
		"bots", *bots,
		"authInterval", *authInterval,
		"roomSize", *roomSize,
		"packetInterval", *packetInterval,
		"logLevel", *logLevel,
		"mode", *mode,
	)

	switch *logLevel {
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
