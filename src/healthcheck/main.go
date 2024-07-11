// © 2019-2024 Diarkis Inc. All rights reserved.
package main

import (
	"fmt"
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/util"
	"net"
	"os"
	"time"
)

// health check will terminate with error if no response in 5000ms
const timeout = 5000

var logger = log.New("H-CHK")

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("Address and port missing: ./bin/tools/health-check <address:port> <in|out>")
		os.Exit(1)
	}
	rootpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Setup(util.StrConcat(rootpath, "/bin/tools/configs/health-check.json"))
	mesh.Guest(util.StrConcat(rootpath, "/bin/tools/configs/health-check.json"))
	diarkis.OnReady(onReady)
	diarkis.Start()
}

func onReady(next func(error)) {
	go setupTimeout()
	addr := os.Args[1]
	mode := os.Args[2]
	logger.Info("Health (mode:%s) check of %s", mode, addr)
	switch mode {
	case "out":
		resolved, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			logger.Fatal("Error %v", err)
			os.Exit(1)
			return
		}
		conn, err := net.DialUDP("udp", nil, resolved)
		if err != nil {
			logger.Fatal("Error %v", err)
			os.Exit(1)
			return
		}
		ticker := time.NewTicker(3 * time.Second)
		// this is timeout
		go func() {
			select {
			case <-ticker.C:
				conn.Close()
				logger.Fatal("Error timeout")
				os.Exit(1)
				return
			}
		}()
		// we wait for a response packet
		go func() {
			buf := make([]byte, 16)
			size, _, err := conn.ReadFrom(buf)
			conn.Close()
			if err != nil {
				logger.Fatal("Error %v", err)
				os.Exit(1)
				return
			}
			// we received a response
			fmt.Printf("%s", string(buf[0:size]))
			os.Exit(0)
		}()
		_, err = conn.Write([]byte("PING\n"))
		if err != nil {
			logger.Fatal("Error %v", err)
			os.Exit(1)
		}
	case "in":
		data := make(map[string]interface{})
		data["meshOnly"] = true
		mesh.SendRequest(util.MeshHealthCheck, addr, data, onMeshHealthCheck)
	case "mars":
		data := make(map[string]interface{})
		data["meshOnly"] = false
		mesh.SendRequest(util.MeshHealthCheck, addr, data, onMeshHealthCheck)
	}
	next(nil)
}

func onMeshHealthCheck(err error, res map[string]interface{}) {
	if err != nil {
		logger.Fatal("Error %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s", res["message"])
	os.Exit(0)
}

func setupTimeout() {
	time.Sleep(time.Millisecond * time.Duration(timeout))
	logger.Fatal("Error timeout")
	os.Exit(1)
}
