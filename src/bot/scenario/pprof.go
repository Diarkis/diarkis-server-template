// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	dnet "github.com/Diarkis/diarkis/net"
)

var (
	pprofEnabled bool
	pprofMu      sync.Mutex
	server       *http.Server
)

func togglePprof() {
	pprofMu.Lock()
	defer pprofMu.Unlock()

	if pprofEnabled {
		logger.Info("Stopping pprof server")
		if server != nil {
			_ = server.Close()
		}
		pprofEnabled = false
	} else {
		port := dnet.FindAvailablePort("localhost", "6060")
		addr := "localhost" + ":" + port
		server = &http.Server{Addr: addr, Handler: nil}

		logger.Info("Starting pprof server on %s", addr)
		go func() {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				logger.Error("pprof server error: %v", err)
			}
		}()
		pprofEnabled = true
	}
}

func setupPprof() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)

	go func() {
		for range sigs {
			togglePprof()
		}
	}()

}
