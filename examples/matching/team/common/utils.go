// © 2019-2025 Diarkis Inc. All rights reserved.

package common

// MatchingParams matching parameters the client
// sends to the server.
// For this example we do not need to configure the criteria.
type MatchingParams struct {
}

// CommandResponse ...
type CommandResponse struct {
	Error *CommandError `json:"error"`
	Data  []byte        `json:"data"`
}

// CommandError ...
type CommandError struct {
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

// MatchingComplete matching complete custom event.
type MatchingComplete struct {
	OwnerID      string     `json:"ownerID"`
	CandidateIDs []string   `json:"candidateIDs"`
	TicketType   uint8      `json:"ticketType"`
	Teams        [][]string `json:"teams"`
}

// MatchingTicketMemberJoined event to notify a member joined
// a matching ticket.
type MatchingTicketMemberJoined struct {
	OwnerID    string   `json:"ownerID"`
	MembersIDs []string `json:"membersIDs"`
	TicketType uint8    `json:"ticketType"`
}

// MatchingTicketMemberJoined event to notify a member joined
// a matching ticket.
type MatchingTicketMemberLeft struct {
	OwnerID    string `json:"ownerID"`
	MemberID   string `json:"memberID"`
	TicketType uint8  `json:"ticketType"`
}
