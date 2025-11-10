// © 2019-2025 Diarkis Inc. All rights reserved.

//go:build windows

package main

import (
	"os"
	"sync"
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	// SIGUSR1 USR1 signal constant
	// More invented values for signals
	// See https://cs.opensource.google/go/go/+/refs/tags/go1.23.1:src/syscall/types_windows.go;l=51
	SIGUSR1 = syscall.SIGTERM + 1
	// SIGUSR2 USR2 signal constant
	SIGUSR2 = syscall.SIGTERM + 2
)

type signalHandler struct {
	sync.RWMutex
	m map[syscall.Signal]map[chan<- os.Signal]struct{}
	h map[syscall.Signal]windows.Handle
}

var gSignalHandler = signalHandler{
	m: make(map[syscall.Signal]map[chan<- os.Signal]struct{}),
	h: make(map[syscall.Signal]windows.Handle),
}
