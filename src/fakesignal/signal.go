//go:build windows

package main

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

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

// FormatFakeSignalEventName generate a unique string associated to the
// signal to be used to create win32 Event.
func formatFakeSignalEventName(sig syscall.Signal) string {
	return formatFakeSignalEventNameWithPID(sig, os.Getpid())
}

// FormatFakeSignalEventNameWithPID generate a unique string associated to the
// signal and pid to be used to create win32 Event.
func formatFakeSignalEventNameWithPID(sig syscall.Signal, pid int) string {
	// name the event with something that cannot be guessed or be already used
	// by another process.
	const format = "DIARKIS#SIGNAL_%s#PID_%d"
	var sigStr string
	switch sig {
	case SIGUSR1:
		sigStr = "SIGUSR1"
	case SIGUSR2:
		sigStr = "SIGUSR2"
	case syscall.SIGHUP:
		sigStr = "SIGHUP"
	default:
		fmt.Printf("FormatFakeSignalEventNameWithPID: unsupported signal %s\n", sig)
		return ""
	}

	return fmt.Sprintf(format, sigStr, pid)
}

func (s *signalHandler) add(c chan os.Signal, sig syscall.Signal, name string) {
	s.Lock()
	defer s.Unlock()

	// Check if the win32 Handle already exists
	if _, ok := s.h[sig]; !ok {
		namep, err := windows.UTF16PtrFromString(name)
		if err != nil {
			fmt.Printf("failed to create utf16 string. %v\n", err)
			return
		}

		handle, err := windows.CreateEvent(nil, 0, 0, namep)
		if err != nil {
			fmt.Printf("failed to create handle with name %q. %v\n", name, err)
			return
		}

		s.h[sig] = handle
		go s.waitForSignal(sig)
	}

	// append the signal
	// use a map to make multiple call idempotent
	if s.m[sig] == nil {
		s.m[sig] = make(map[chan<- os.Signal]struct{})
	}
	s.m[sig][c] = struct{}{}
}

func (s *signalHandler) waitForSignal(sig syscall.Signal) {
	var ok bool
	var handle windows.Handle

	s.Lock()
	handle, ok = s.h[sig]
	s.Unlock()

	if !ok {
		// Should never happen
		return
	}

	for {
		evt, err := windows.WaitForSingleObject(handle, 30_000)
		if err != nil {
			fmt.Printf("waitForSignal: %s: windows.WaitForSingleObject error: %v", sig, err)
			time.Sleep(time.Second * 30)
			continue
		}

		if evt == windows.WAIT_OBJECT_0 {
			// signaled, notify the channel
			s.RLock()
			for c := range s.m[sig] {
				// send but do not block for it
				// Same as signal.Notify.
				// See https://github.com/golang/go/blob/e9a500f47dadcd73c970649a1072d28997617610/src/os/signal/signal.go#L243
				select {
				case c <- sig:
				default:
				}
			}
			s.RUnlock()
		}
	}
}
