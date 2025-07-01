// © 2019-2025 Diarkis Inc. All rights reserved.

package common

type MatchingComplete struct {
	OwnerID      string   `json:"ownerID"`
	CandidateIDs []string `json:"candidateIDs"`
	TicketType   uint8    `json:"ticketType"`
}
