package main

import "C"

import (
	"fmt"
	"github.com/Diarkis/diarkis/server/connector"
	"os"
	"strconv"
)

func main() {

}

const host = "127.0.0.1"
const anchorPort = "9000"

var peerAddr = ""

// Initialize must be called when the client server process starts.
//
//export Initialize
func Initialize() C.int {
	diarkisEP := connector.Initialize()
	if diarkisEP == "" {
		return C.int(-1)
	}
	peerAddr = diarkisEP
	return C.int(1)
}

// GetOpenPort returns the port for the client server to bind for client communication.
//
//export GetOpenPort
func GetOpenPort() C.int {
	res, value := connector.SendToDiarkis(peerAddr, "port", -1)
	if value == nil {
		return C.int(res)
	}
	portString := string(value)
	port, err := strconv.Atoi(portString)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v\n", err))
		return -1
	}
	return C.int(port)
}

// GetAddress returns public endpoint address that the client uses to connect to client server.
// Allocated address string memory must be released manually.
//
//export GetAddress
func GetAddress() *C.char {
	res, value := connector.SendToDiarkis(peerAddr, "addr", -1)
	if res == -1 {
		return C.CString("")
	}
	if value == nil {
		return C.CString("")
	}
	return C.CString(string(value))
}

// Ready must be called when client server is ready to receive clients.
//
//export Ready
func Ready() C.int {
	res, _ := connector.SendToDiarkis(peerAddr, "ready", -1)
	return C.int(res)
}

// Health must be called every 1 second to report CCU and "healthiness".
//
//export Health
func Health(ccu C.int) C.int {
	res, _ := connector.SendToDiarkis(peerAddr, "health", int(ccu))
	return C.int(res)
}

// Allocate should be called when client server needs to stop revcieving new clients.
// Ready must be called again to start receive new clients.
//
//export Allocate
func Allocate() C.int {
	res, _ := connector.SendToDiarkis(peerAddr, "allocate", -1)
	return C.int(res)
}

// Shutdown should be called when the client server is shutting down.
//
//export Shutdown
func Shutdown() C.int {
	res, _ := connector.SendToDiarkis(peerAddr, "shutdown", -1)
	return C.int(res)
}
