// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"slices"

	"github.com/Diarkis/diarkis-server-template/examples/matching/custom-criteria/common"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

var logger = log.New("UDP")

// 2 for $teamMaxMembers vs $teamMaxMembers
// 3 for $teamMaxMembers vs $teamMaxMembers vs $teamMaxMembers

const (
	teamVsTeamSize                   = 2
	teamMaxMembers                   = 2
	ticketDuration                   = 15 // 15 seconds
	meshGetTicketMembers      uint16 = 10001
	meshRemoteTicketBroadcast uint16 = 10002
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := ""

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json", ExposeCommands: true},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")

	diarkisexec.SetServerCommandHandler(common.AppVersion, common.MatchingStartCmd, handleStartMatching)
	diarkisexec.SetMeshRPCHandler(meshGetTicketMembers, handleGetTicketMembers)
	diarkisexec.SetMeshRPCHandler(meshRemoteTicketBroadcast, handleRemoteTicketBroadcast)

	setupMatching()

	diarkisexec.StartDiarkis()

}

func setupMatching() {
	// Team matching:
	matching.SetOnTicketAllowMatchIf(common.TeamTicketType, func(ticketProps *matching.TicketProperties, owner, candidate *user.User) bool {
		return true
	})

	matching.SetOnTicketMemberJoined(common.TeamTicketType, func(ticket *matching.Ticket, joinedUser, ownerUser *user.User, memberIDs []string) {
		// Notify ticket members when a member joins the ticket.
		data, _ := json.Marshal(common.MatchingTicketMemberJoined{
			OwnerID:    ownerUser.ID,
			TicketType: common.TeamTicketType,
			MembersIDs: memberIDs,
		})
		matching.TicketBroadcast(common.TeamTicketType, ownerUser, common.AppVersion, common.MatchingTicketMemberJoinedCmd, data)
	})

	// Notify ticket's member when a member leaves the ticket.
	// In this example this is useful to notify a member the owner left the ticket
	// so all the client exists either on success or if the battle ticket time out.
	setOnLeaveTicketRoom := func(ticketType uint8, owner *user.User) {
		matching.SetOnLeaveTicketRoom(ticketType, owner, func(id string, userData *user.User) {
			data, _ := json.Marshal(common.MatchingTicketMemberLeft{
				OwnerID:    owner.ID,
				TicketType: common.TeamTicketType,
				MemberID:   userData.ID,
			})
			matching.TicketBroadcast(common.TeamTicketType, owner, common.AppVersion, common.MatchingTicketMemberLeftCmd, data)
		})
	}

	matching.SetOnTicketCompleteWithProfileID(common.TeamTicketType, func(ticketProps *matching.TicketProperties, owner *user.User, tag []string, profileID string) []byte {
		candidates := ticketProps.GetAllCandidates()
		candidateIDs := make([]string, 0, len(candidates)+1)

		candidateIDs = append(candidateIDs, owner.ID)
		for uid := range candidates {
			candidateIDs = append(candidateIDs, uid)
		}

		slices.Sort(candidateIDs)

		setOnLeaveTicketRoom(common.TeamTicketType, owner)

		// Generate a payload to be sent to all ticket members.
		b, _ := json.Marshal(common.MatchingComplete{
			OwnerID:      owner.ID,
			CandidateIDs: candidateIDs,
			TicketType:   common.TeamTicketType,
		})

		{
			// In case you want to match another team based on some criteria.
			// For example the average level of the team, you could store that information
			// in the owner and use them to adjust the battle criteria.

			// Once the team has been created, the owner starts a new ticket to match
			// against another team.

			searchTries := rand.Intn(10-1) + 1
			emptySearches := rand.Intn(searchTries) + 1
			ticketParams := &matching.TicketParams{
				ProfileIDs:     []string{"battle"},
				MaxMembers:     teamVsTeamSize,
				SearchInterval: 100, // 100ms
				SearchTries:    uint8(searchTries),
				EmptySearches:  uint8(emptySearches),
				TicketDuration: ticketDuration,
				HowMany:        20,
				// Change here as you see fit according to your application needs.
				Tags: nil,
				// The profile we use has no criteria because we focus on implementing
				// a custom check.
				AddProperties:    map[string]int{"level": 1},
				SearchProperties: map[string][]int{"level": {1}},
			}

			if err := matching.StartTicketWithTicketParams(common.BattleTicketType, owner, ticketParams); err != nil {
				logger.Error("failed to start battle ticket. %v", err)
			}
		}

		return b
	})

	// Battle matching:
	matching.SetOnTicketAllowMatchIf(common.BattleTicketType, func(ticketProps *matching.TicketProperties, owner, candidate *user.User) bool {
		return true
	})

	matching.SetOnTicketMemberJoined(common.BattleTicketType, func(ticket *matching.Ticket, joinedUser, ownerUser *user.User, memberIDs []string) {
		// Notify ticket members when a member joins the ticket.
		data, _ := json.Marshal(common.MatchingTicketMemberJoined{
			OwnerID:    ownerUser.ID,
			TicketType: common.BattleTicketType,
			MembersIDs: memberIDs,
		})
		matching.TicketBroadcast(common.BattleTicketType, ownerUser, common.AppVersion, common.MatchingTicketMemberJoinedCmd, data)
	})

	matching.SetOnTicketCompleteWithProfileID(common.BattleTicketType, func(ticketProps *matching.TicketProperties, owner *user.User, tag []string, profileID string) []byte {
		teams := [][]string{}
		// Owner team:
		ownerTeam, _ := matching.GetTicketMemberIDs(common.TeamTicketType, owner)
		slices.Sort(ownerTeam)
		teams = append(teams, ownerTeam)

		// Candidate team:
		for uid, ticketHolder := range ticketProps.GetAllCandidates() {
			meshAddr := ticketHolder.MeshEndPoint
			req := getTicketMembersReq{
				TicketType: common.TeamTicketType,
				OwnerID:    uid,
			}
			data, _ := json.Marshal(req)
			responseData, err := diarkisexec.SendMeshRPC(meshGetTicketMembers, meshAddr, data)
			if err != nil {
				logger.Error("rpc get ticket members failed. %v", err)
			} else {
				var resp getTicketMembersResponse
				_ = json.Unmarshal(responseData, &resp)

				slices.Sort(resp.MemberIDs)
				teams = append(teams, resp.MemberIDs)
			}
		}

		// Generate the payload to be sent to the ticket members.
		b, _ := json.Marshal(common.MatchingComplete{
			OwnerID:    owner.ID,
			TicketType: common.BattleTicketType,
			Teams:      teams,
		})

		// Notify owner's team ticket members a team has been found.
		matching.TicketMultibroadcast([]uint8{common.BattleTicketType, common.TeamTicketType}, owner, common.AppVersion, common.MatchingTicketBattleBroadcastCmd, b)

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
	searchTries := rand.Intn(10-1) + 1
	emptySearches := rand.Intn(searchTries) + 1
	ticketParams := &matching.TicketParams{
		ProfileIDs:     []string{"team"},
		MaxMembers:     teamMaxMembers,
		SearchInterval: 100, // 100ms
		SearchTries:    uint8(searchTries),
		EmptySearches:  uint8(emptySearches),
		TicketDuration: ticketDuration,
		HowMany:        20,
		// Change here as you see fit according to your application needs.
		Tags: nil,
		// The profile we use has no criteria because we focus on implementing
		// a two steps matching.
		AddProperties:    map[string]int{"level": 1},
		SearchProperties: map[string][]int{"level": {1}},
	}

	err = matching.StartTicketWithTicketParams(common.TeamTicketType, userData, ticketParams)
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

// handleGetTicketMembers Returns the members of a matching ticket.
func handleGetTicketMembers(data []byte, senderAddr string) (response []byte, err error) {
	var req getTicketMembersReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	owner := user.GetUserByUID(req.OwnerID)
	if owner == nil {
		return nil, errors.New("owner of ticket not found")
	}

	members, _ := matching.GetTicketMemberIDs(req.TicketType, owner)
	response, err = json.Marshal(getTicketMembersResponse{MemberIDs: members})
	return
}

func handleRemoteTicketBroadcast(data []byte, senderAddr string) (response []byte, err error) {
	var req ticketBroadcastReq
	if err := json.Unmarshal(data, &req); err != nil {
		logger.Error("handleRemoteTicketBroadcast: json.Unmarshal: %v", err)
		return nil, err
	}

	owner := user.GetUserByUID(req.OwnerID)
	if owner == nil {
		logger.Error("handleRemoteTicketBroadcast: owner not found")
		return nil, errors.New("owner of ticket not found")
	}

	err = matching.TicketBroadcast(req.TicketType, owner, req.Version, req.Command, req.Data)

	return
}

type getTicketMembersReq struct {
	TicketType uint8
	OwnerID    string
}

type getTicketMembersResponse struct {
	MemberIDs []string
}

type ticketBroadcastReq struct {
	TicketType uint8
	OwnerID    string
	Data       []byte
	Version    uint8
	Command    uint16
}
