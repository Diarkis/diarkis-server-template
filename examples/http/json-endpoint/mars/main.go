// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mars"
)

var logger = log.New("MARS")

func main() {
	mars.Setup()
	diarkis.OnReady(onReady)
	diarkis.Start()
}

func onReady(next func(error)) {
	logger.Info("Mars server is now ready")
	next(nil)
}
