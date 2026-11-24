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
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

var logger = log.New("UDP")

// 2 for $teamMaxMembers vs $teamMaxMembers
// 3 for $teamMaxMembers vs $teamMaxMembers vs $teamMaxMembers

const (
	teamVsTeamSize                   = 2
	teamMaxMembers                   = 3
	ticketDuration                   = 15 // 15 seconds
	meshGetTicketMembers      uint16 = 10001
	meshRemoteTicketBroadcast uint16 = 10002
)

var hardcodedParties = map[string][]string{
	"user-1": {"user-2", "user-3"},
	"user-4": {"user-5"},
}

func isPartyMember(uid string) bool {
	for _, members := range hardcodedParties {
		if slices.Contains(members, uid) {
			return true
		}
	}
	return false
}

func getTeamMembers(owner *user.User) []string {
	if ids, _ := matching.GetTicketMemberIDs(common.TeamTicketType, owner); len(ids) > 0 {
		return ids
	}
	return append([]string{owner.ID}, hardcodedParties[owner.ID]...) // TEMP: use Session members later
}

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
		memberIDs, _ := matching.GetTicketMemberIDs(common.TeamTicketType, owner)
		slices.Sort(memberIDs)
		logger.Infof("team ticket complete", "owner", owner.ID, "members", memberIDs)

		setOnLeaveTicketRoom(common.TeamTicketType, owner)

		// Generate a payload to be sent to all ticket members.
		b, _ := json.Marshal(common.MatchingComplete{
			OwnerID:      owner.ID,
			CandidateIDs: memberIDs,
			TicketType:   common.TeamTicketType,
		})

		// Once the team has been created, the owner starts a new ticket to match
		// against another team. The other team members join the battle ticket with the owner.
		others := slices.DeleteFunc(slices.Clone(memberIDs), func(id string) bool { return id == owner.ID })
		if err := startBattleTicket(owner, others); err != nil {
			logger.Error("failed to start battle ticket. %v", err)
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
		// All members of both teams are in the battle ticket.
		battleMembers, _ := matching.GetTicketMemberIDs(common.BattleTicketType, owner)

		// Owner team, and the other team = everyone else in the battle ticket.
		ownerTeam := getTeamMembers(owner)
		otherTeam := slices.DeleteFunc(slices.Clone(battleMembers), func(id string) bool { return slices.Contains(ownerTeam, id) })
		slices.Sort(ownerTeam)
		slices.Sort(otherTeam)
		teams := [][]string{ownerTeam, otherTeam}
		logger.Infof("battle ticket complete", "owner", owner.ID, "battleMembers", battleMembers, "teams", teams)

		// Generate the payload to be sent to the ticket members.
		b, _ := json.Marshal(common.MatchingComplete{
			OwnerID:    owner.ID,
			TicketType: common.BattleTicketType,
			Teams:      teams,
		})

		// Notify all battle members (both teams) the battle has been found.
		if err := matching.TicketBroadcast(common.BattleTicketType, owner, common.AppVersion, common.MatchingTicketBattleBroadcastCmd, b); err != nil {
			logger.Error("failed to broadcast battle result. %v", err)
		}

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

	// Party members do not start a ticket. The party owner adds them to its ticket.
	if isPartyMember(userData.ID) {
		logger.Infof("waiting for the party owner to start a ticket", "uid", userData.ID)
		userData.ServerRespond(nil, ver, cmd, server.Ok, true)
		next(nil)
		return
	}
	partyMembers := hardcodedParties[userData.ID] // nil for solo users

	if 1+len(partyMembers) >= teamMaxMembers {
		// The party already fills a team: skip the team ticket and go straight to battle.
		logger.Infof("party is full, skip team ticket", "owner", userData.ID, "members", partyMembers)
		err = startBattleTicket(userData, partyMembers)
	} else {
		err = startTeamTicket(userData, partyMembers)
	}
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

func newTicketMembers(memberIDs []string) *matching.TicketMembers {
	if len(memberIDs) == 0 {
		return nil
	}
	meshAddrs := make([]string, len(memberIDs))
	for i := range memberIDs {
		meshAddrs[i] = mesh.GetMyEndPoint()
	}
	return &matching.TicketMembers{IDs: memberIDs, MeshAddrs: meshAddrs}
}

func startTeamTicket(owner *user.User, memberIDs []string) error {
	searchTries := rand.Intn(10-1) + 1
	emptySearches := rand.Intn(searchTries) + 1
	options := &matching.TicketOptions{
		TicketParams: &matching.TicketParams{
			ProfileIDs:       []string{"team"},
			MaxMembers:       teamMaxMembers,
			SearchInterval:   100, // 100ms
			SearchTries:      uint8(searchTries),
			EmptySearches:    uint8(emptySearches),
			TicketDuration:   ticketDuration,
			HowMany:          20,
			AddProperties:    map[string]int{"level": 1},
			SearchProperties: map[string][]int{"level": {1}},
		},
		TicketMembers: newTicketMembers(memberIDs),
	}
	logger.Infof("start team ticket", "owner", owner.ID, "members", memberIDs)
	return matching.StartTicketWithOptions(common.TeamTicketType, owner, options)
}

func startBattleTicket(owner *user.User, memberIDs []string) error {
	searchTries := rand.Intn(10-1) + 1
	emptySearches := rand.Intn(searchTries) + 1
	options := &matching.TicketOptions{
		TicketParams: &matching.TicketParams{
			ProfileIDs:       []string{"battle"},
			MaxMembers:       teamMaxMembers * teamVsTeamSize, // all members of all teams
			SearchInterval:   100,                             // 100ms
			SearchTries:      uint8(searchTries),
			EmptySearches:    uint8(emptySearches),
			TicketDuration:   ticketDuration,
			HowMany:          20,
			AddProperties:    map[string]int{"level": 1},
			SearchProperties: map[string][]int{"level": {1}},
		},
		TicketMembers: newTicketMembers(memberIDs),
	}
	logger.Infof("start battle ticket", "owner", owner.ID, "members", memberIDs)
	return matching.StartTicketWithOptions(common.BattleTicketType, owner, options)
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
