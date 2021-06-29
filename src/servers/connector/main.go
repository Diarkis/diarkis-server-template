package main

import (
	"fmt"
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/server"
	"os"
)

/**
* Use DIARKIS_NODE_TYPE env to set node type
*/

var logger = log.New("CONN")

func main() {
	rootpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Setup(fmt.Sprintf("%s/configs/shared/log.json", rootpath))
	mesh.Setup(fmt.Sprintf("%s/configs/shared/mesh.json", rootpath))
	// change node type as you see fit
	mesh.SetNodeType("CONNECTOR")
	server.SetupAsConnector(fmt.Sprintf("%s/configs/connector/main.json", rootpath))
	diarkis.Start()
}
