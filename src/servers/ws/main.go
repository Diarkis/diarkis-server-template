package main

import (
	"fmt"
	"github.com/Diarkis/diarkis"
	cmds "{0}/ws_cmds"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/dairkis/ws"
	"io/ioutil"
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
	// this assumes there is nginx running along side and bin/tools/findport outputs /tmp/NGINX_PORT
	diarkis.OnReady(func(next func(error)) {
		bytes, err := ioutil.ReadFile("/tmp/NGINX_PORT")
		if err != nil {
			// we failed to read nginx port, we ignore it
			next(nil)
			return
		}
		nginxport := string(bytes)
		switch os.Getenv("DIARKIS_CLOUD_ENV") {
			case "AZURE_WS":
				pubEP, _ := util.GetPublicEndPointMSLB()
				ws.SetPublicEndPointWithPort(pubEP, nginxport)
			// TODO: add GCP and AWS version here too
		}
		next(nil)
	})
	diarkis.Start()
}
