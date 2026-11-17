// © 2019-2025 Diarkis Inc. All rights reserved.

package protocol

const CustomVer uint8 = 2

const ShowSessionsCmd uint16 = 100

type SessionInfo struct {
	ID        string   `json:"id"`
	Type      uint8    `json:"type"`
	MemberIDs []string `json:"memberIDs"`
}

const InviteReplyCmd uint16 = 101

const InviteDeclinedPushCmd uint16 = 102

type InviteReply struct {
	SessionType uint8  `json:"sessionType"`
	SessionID   string `json:"sessionID"`
	Accept      bool   `json:"accept"`
}

type InviteDeclined struct {
	SessionID string `json:"sessionID"`
	UID       string `json:"uid"`
}
