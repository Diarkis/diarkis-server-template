// © 2019-2025 Diarkis Inc. All rights reserved.

//go:build windows

// This is a utility to send a signal (SIGUSR1, SIGUSR2 or SIGHUP)
// to diarkis process.
package main

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/sys/windows"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <PID> <SIGNAL>\n", os.Args[0])
		os.Exit(1)
	}

	pidStr := os.Args[1]
	signalStr := os.Args[2]

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Printf("invalid pid %q\n", pidStr)
		os.Exit(1)
	}
	switch signalStr {
	case "SIGUSR1":
	case "SIGUSR2":
	case "SIGHUP":
	default:
		fmt.Printf("invalid signal name %q\n", signalStr)
		os.Exit(1)
	}

	name := fmt.Sprintf("DIARKIS#SIGNAL_%s#PID_%d", signalStr, pid)

	namep, err := windows.UTF16PtrFromString(name)
	if err != nil {
		fmt.Printf("failed to create utf16 string. %v\n", err)
		os.Exit(1)
	}

	handle, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, namep)
	if err != nil {
		fmt.Printf("failed to create handle. %v\n", err)
		fmt.Printf("the process with pid %d might not exist anymore\n", pid)
		os.Exit(1)
	}

	err = windows.SetEvent(handle)
	if err != nil {
		fmt.Printf("failed to set event. %v", err)
		os.Exit(1)
	}
}
