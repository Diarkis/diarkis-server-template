package common

type MatchingParams struct {
	MatchingID string   `json:"matchingID"`
	Rank       int      `json:"rank"`
	Level      int      `json:"level"`
	Tags       []string `json:"tags"`
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
