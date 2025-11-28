// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/diarkisexec"

	proom "github.com/Diarkis/diarkis-server-template/examples/room/lod/puffer/go/room"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := ""

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		Room: &diarkisexec.Options{ConfigPath: "configs/shared/room.json", ExposeCommands: true},
		Dive: &diarkisexec.Options{ConfigPath: "configs/shared/dive.json"},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")
	diarkisexec.SetServerCommandHandler(proom.BroadcastLoDVer, proom.BroadcastLoDCmd, handleRoomBroadcastLoD)
	diarkisexec.SetServerCommandHandler(proom.GetLoDInfoVer, proom.GetLoDInfoCmd, handleRoomGetLoDInfo)

	diarkisexec.StartDiarkis()
}
