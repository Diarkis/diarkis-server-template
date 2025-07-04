// © 2019-2025 Diarkis Inc. All rights reserved.

package common

type MatchingParams struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CommandResponse struct {
	Error *CommandError `json:"error"`
	Data  []byte        `json:"data"`
}

type CommandError struct {
	Message string `json:"msg"`
	Code    int    `json:"code"`
}

type MatchingComplete struct {
	OwnerID      string   `json:"ownerID"`
	CandidateIDs []string `json:"candidateIDs"`
	TicketType   uint8    `json:"ticketType"`
}
