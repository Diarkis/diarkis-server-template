// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
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
		l, err := findAvailablePort("localhost", 6060)
		if err != nil {
			panic(err.Error())
		}
		server = &http.Server{Handler: nil}
		logger.Info("Starting pprof server on %s", l.Addr().String())
		go func() {
			if err := server.Serve(l); err != http.ErrServerClosed {
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

// findAvailablePort Find the first available port starting from port.
func findAvailablePort(host string, port int) (net.Listener, error) {
	var l net.Listener
	var err error

	startingPort := port
	for range 1000 {
		listenAddr := net.JoinHostPort(host, strconv.Itoa(port))
		l, err = net.Listen("tcp", listenAddr)
		if err == nil {
			return l, nil
		}
		port++
	}
	return nil, fmt.Errorf("unable to find available port within %d and %d. %v", startingPort, port, err)
}
