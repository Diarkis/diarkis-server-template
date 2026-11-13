// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/Diarkis/diarkis-server-template/examples/matching/session/protocol"
	"github.com/Diarkis/diarkis/client/go/test/cli"
)

func main() {
	cli.SetupBuiltInSessionCommands()
	// 独自コマンドは標準の "session" とは別の名前で登録する
	// (同じ名前で登録すると "Module [session] is already registered" で panic する)
	cli.RegisterCommands("custom", []cli.Command{
		{CmdName: "show", Desc: "Show all sessions on the server", CmdFunc: showSessions},
	})
	cli.Connect()
	// cli.UDPClient は Connect のあとで使えるようになる
	cli.UDPClient.OnResponse(onResponse)
	cli.Run()
}

func showSessions() {
	cli.UDPClient.RSend(protocol.CustomVer, protocol.ShowSessionsCmd, nil)
}

func onResponse(ver uint8, cmd uint16, status uint8, payload []byte) {
	if ver != protocol.CustomVer || cmd != protocol.ShowSessionsCmd {
		return
	}
	// status 1 が成功(サーバー側の server.Ok)
	if status != 1 {
		fmt.Printf("session show failed: %s\n", string(payload))
		return
	}
	var list []protocol.SessionInfo
	if err := json.Unmarshal(payload, &list); err != nil {
		fmt.Printf("session show: failed to read response: %v\n", err)
		return
	}
	fmt.Printf("Sessions (%d):\n", len(list))
	for _, info := range list {
		fmt.Printf("  type: %d, id: %s, members: %v\n", info.Type, info.ID, info.MemberIDs)
	}
}
