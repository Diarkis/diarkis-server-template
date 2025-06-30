package common

type MatchingComplete struct {
	OwnerID      string   `json:"ownerID"`
	CandidateIDs []string `json:"candidateIDs"`
	TicketType   uint8    `json:"ticketType"`
}
