package main

import (
	"fmt"
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis-server-template/cmds"
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
	server.SetupAsTCPServer(fmt.Sprintf("%s/configs/tcp/main.json", rootpath))
	cmds.Setup(rootpath)
	cmds.ExposeServer()
	diarkis.Start()
}
