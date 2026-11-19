// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Diarkis/diarkis-server-template/examples/matching/session/protocol"
	"github.com/Diarkis/diarkis/client/go/modules/session"
	"github.com/Diarkis/diarkis/client/go/test/cli"
)

var (
	lastInviteMu sync.Mutex
	lastInvite   *protocol.InviteReply
)

func main() {
	cli.SetupBuiltInSessionCommands()
	cli.RegisterCommands("custom", []cli.Command{
		{CmdName: "show", Desc: "Show all sessions on the server", CmdFunc: showSessions},
	})
	cli.RegisterCommands("invite", []cli.Command{
		{CmdName: "yes", Desc: "Accept the last invitation and join the Session", CmdFunc: acceptInvite},
		{CmdName: "no", Desc: "Decline the last invitation and notify the user who sent it", CmdFunc: declineInvite},
	})
	cli.Connect()
	cli.UDPClient.OnResponse(onResponse)
	cli.UDPClient.OnPush(onPush)
	session.NewSessionAsUDP(cli.UDPClient).OnSessionInvite(onSessionInvite)
	cli.Run()
}

func showSessions() {
	cli.UDPClient.RSend(protocol.CustomVer, protocol.ShowSessionsCmd, nil)
}

func acceptInvite() {
	sendInviteReply(true)
}

func declineInvite() {
	sendInviteReply(false)
}

func onSessionInvite(sessionType uint8, sessionID string, message []byte) {
	lastInviteMu.Lock()
	defer lastInviteMu.Unlock()
	lastInvite = &protocol.InviteReply{SessionType: sessionType, SessionID: sessionID}
	fmt.Println("Reply with `invite yes` or `invite no`")
}

func sendInviteReply(accept bool) {
	lastInviteMu.Lock()
	invite := lastInvite
	lastInvite = nil
	lastInviteMu.Unlock()

	if invite == nil {
		fmt.Println("No invitation to reply to")
		return
	}

	payload, err := json.Marshal(protocol.InviteReply{
		SessionType: invite.SessionType,
		SessionID:   invite.SessionID,
		Accept:      accept,
	})
	if err != nil {
		fmt.Printf("Failed to create invite reply: %v\n", err)
		return
	}
	cli.UDPClient.RSend(protocol.CustomVer, protocol.InviteReplyCmd, payload)
}

func onResponse(ver uint8, cmd uint16, status uint8, payload []byte) {
	if ver != protocol.CustomVer {
		return
	}
	switch cmd {
	case protocol.ShowSessionsCmd:
		if status != 1 {
			fmt.Printf("custom show failed: %s\n", string(payload))
			return
		}
		var list []protocol.SessionInfo
		if err := json.Unmarshal(payload, &list); err != nil {
			fmt.Printf("custom show: failed to read response: %v\n", err)
			return
		}
		fmt.Printf("Sessions (%d):\n", len(list))
		for _, info := range list {
			fmt.Printf("  type: %d, id: %s, members: %v\n", info.Type, info.ID, info.MemberIDs)
		}

	case protocol.InviteReplyCmd:
		if status != 1 {
			fmt.Printf("Invite reply failed: %s\n", string(payload))
			return
		}
		fmt.Printf("Invite reply sent: %s\n", string(payload))
	}
}

func onPush(ver uint8, cmd uint16, payload []byte) {
	if ver != protocol.CustomVer || cmd != protocol.InviteDeclinedPushCmd {
		return
	}
	var declined protocol.InviteDeclined
	if err := json.Unmarshal(payload, &declined); err != nil {
		fmt.Printf("Failed to read invite declined notification: %v\n", err)
		return
	}
	fmt.Printf("%s declined the invitation. SessionID: %s\n", declined.UID, declined.SessionID)
}
