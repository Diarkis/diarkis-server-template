package main

import (
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/config"
	ddebug "github.com/Diarkis/diarkis/debug"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mars"
)

var logger = log.New("MARS")

func main() {
	mars.Setup()
	if config.GetAsBool("Mars", "debug", false) {
		ddebug.Enable()
		ddebug.EnableMutexLogging()
	}
	diarkis.OnReady(onReady)
	diarkis.Start()
}

func onReady(next func(error)) {
	logger.Info("Mars server is now ready")
	next(nil)
}
