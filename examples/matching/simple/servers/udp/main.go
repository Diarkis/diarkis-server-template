package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/Diarkis/diarkis-server-template/examples/matching/simple/common"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

var logger = log.New("UDP")

const (
	ticketDuration = 15 // 15 seconds
	ticketType     = uint8(1)
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := ""

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json", ExposeCommands: true},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")
	diarkisexec.SetServerCommandHandler(common.AppVersion, common.MatchingStartCmd, handleStartMatching)

	setupMatching()

	diarkisexec.StartDiarkis()

}

func setupMatching() {
	matching.SetOnTicketAllowMatchIf(ticketType, func(ticketProps *matching.TicketProperties, owner, candidate *user.User) bool {
		return true
	})

	matching.SetOnTicketCompleteWithProfileID(ticketType, func(ticketProps *matching.TicketProperties, owner *user.User, tag []string, profileID string) []byte {
		fmt.Printf("**** tag: %q, profileID: %q\n\n", tag, profileID)
		candidates := ticketProps.GetAllCandidates()
		candidateIDs := make([]string, 0, len(candidates))

		for uid := range candidates {
			candidateIDs = append(candidateIDs, uid)
		}

		b, _ := json.Marshal(common.MatchingComplete{
			OwnerID:      owner.ID,
			CandidateIDs: candidateIDs,
			TicketType:   ticketType,
		})

		return b
	})
}

func handleStartMatching(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	var params common.MatchingParams

	err := json.Unmarshal(payload, &params)
	if err != nil {
		err = fmt.Errorf("fail to parse payload. %w", err)
		b, _ := json.Marshal(common.CommandResponse{
			Error: &common.CommandError{
				Message: err.Error(),
			},
		})
		userData.ServerRespond(b, ver, cmd, server.Bad, true)
		next(err)
		return
	}

	// Randomize searchTries and emptySearches to not move to the wait mode at the same time
	// considering all clients issue tickets at the same time.
	searchTries := randomInt(1, 10)
	emptySearches := randomInt(1, max(2, searchTries))
	ticketParams := &matching.TicketParams{
		ProfileIDs:     []string{params.MatchingID},
		MaxMembers:     2,
		SearchInterval: 100, // 100ms
		SearchTries:    uint8(searchTries),
		EmptySearches:  uint8(emptySearches),
		TicketDuration: ticketDuration,
		HowMany:        20,
		// Change here as you see fit according to your application needs.
		Tags: params.Tags,
		// Change here as you see fit according to your application needs.
		AddProperties:    map[string]int{"level": params.Level},
		SearchProperties: map[string][]int{"level": {params.Level}},
	}
	logger.Info("start ticket for user %s. level=%d", userData.ID, params.Level)

	err = matching.StartTicketWithTicketParams(ticketType, userData, ticketParams)
	if err != nil {
		err = fmt.Errorf("fail to start matching ticket. %w", err)
		b, _ := json.Marshal(common.CommandResponse{
			Error: &common.CommandError{
				Message: err.Error(),
			},
		})
		userData.ServerRespond(b, ver, cmd, server.Bad, true)
		next(err)
		return
	}

	userData.ServerRespond(nil, ver, cmd, server.Ok, true)
	next(nil)
}

func randomInt(minValue, maxValue int) int {
	return rand.Intn(maxValue-minValue) + minValue
}
