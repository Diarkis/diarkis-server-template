// © 2019-2025 Diarkis Inc. All rights reserved.

package common

import (
	"os"
	"path/filepath"
)

type MatchingComplete struct {
	OwnerID      string   `json:"ownerID"`
	CandidateIDs []string `json:"candidateIDs"`
	TicketType   uint8    `json:"ticketType"`
}

func GetConfigPath(path string) string {
	if len(path) == 0 {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	rootpath, _ := os.Getwd()

	return filepath.Join(rootpath, path)
}
