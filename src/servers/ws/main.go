package main

import (
	"fmt"
	"github.com/Diarkis/diarkis"
	cmds "{0}/ws_cmds"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/server"
	"os"
)

func main() {
	rootpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Setup(fmt.Sprintf("%s/configs/shared/log.json", rootpath))
	mesh.Setup(fmt.Sprintf("%s/configs/shared/mesh.json", rootpath))
	server.SetupAsWebSocketServer(fmt.Sprintf("%s/configs/ws/main.json", rootpath))
	cmds.Setup(rootpath)
	cmds.ExposeServer()
	diarkis.Start()
}
