package main

import (
	"encoding/json"
	"fmt"

	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/user"

	"github.com/Diarkis/diarkis-server-template/examples/http/json-endpoint/common"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/util"
)

var logger = log.New("UDP")

const userDataMatchingParamKey = "_matching_params"
const ticketDuration = 15 // 15 seconds

func main() {
	logConfigPath := "/configs/shared/log.json"
	meshConfigPath := ""
	log.EnableUnsafeLogging()

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "/configs/shared/matching.json", ExposeCommands: true},
	})
	diarkisexec.SetupDiarkisUDPServer("/configs/udp/main.json")

	diarkisexec.SetServerCommandHandler(common.AppVersion, common.MatchingStartCmd, handleStartMatching)
	// diarkisexec.SetServerCommandHandler(common.AppVersion, common.MatchingTicketBroadcastCmd, handleMatchingTicketBroadcast)

	setupMaching()

	diarkisexec.StartDiarkis()

}

func setupMaching() {
	const ticketType = uint8(1)
	matching.SetOnIssueTicket(ticketType, func(userData *user.User) *matching.TicketParams {
		b, _ := userData.GetAsBytes(userDataMatchingParamKey)
		var params common.MatchingParams
		err := json.Unmarshal(b, &params)
		if err != nil {
			logger.Error("fail to parse matching params. %v", err)
			return nil
		}

		// Randomize searchTries and emptySearches to not move to the wait mode at the same time considering all clients issue tickets at the same time.
		searchTries := util.RandomInt(1, 10)
		emptySearches := util.RandomInt(1, searchTries)
		return &matching.TicketParams{
			ProfileIDs:     []string{params.MatchingID},
			MaxMembers:     2,
			SearchInterval: 100, // 100ms
			SearchTries:    uint8(searchTries),
			EmptySearches:  uint8(emptySearches),
			TicketDuration: ticketDuration,
			HowMany:        20,
			// Change here as you see fit according to your application needs
			Tags: params.Tags,
			// Change here as you see fit according to your application needs
			AddProperties:    map[string]int{"level": params.Level},
			SearchProperties: map[string][]int{"level": {params.Level}},
		}
	})

	matching.SetOnTicketAllowMatchIf(ticketType, func(ticketProps *matching.TicketProperties, owner, candidate *user.User) bool {
		return true
	})

	/*
		matching.SetOnTicketComplete(ticketType, func(ticketProps *matching.TicketProperties, owner *user.User) []byte {
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
	*/
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

	// Attach the matching params to the user data
	// so we can retrieve them when creating the ticket.
	userData.Set(userDataMatchingParamKey, payload)
	defer userData.Delete(userDataMatchingParamKey)

	err = matching.StartTicket(1, userData)
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
