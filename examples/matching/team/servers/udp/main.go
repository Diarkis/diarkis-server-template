// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"slices"

	"github.com/Diarkis/diarkis-server-template/examples/matching/custom-criteria/common"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/session"
	"github.com/Diarkis/diarkis/user"
)

var logger = log.New("UDP")

// 2 for $teamMaxMembers vs $teamMaxMembers
// 3 for $teamMaxMembers vs $teamMaxMembers vs $teamMaxMembers

const (
	teamVsTeamSize = 2
	teamMaxMembers = 3
	ticketDuration = 15 // 15 seconds
	meshAddrKey = "meshAddr"
	teamMembersKey = "teamMembers"
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json", ExposeCommands: true},
		Session:    &diarkisexec.Options{ConfigPath: "configs/shared/session.json", ExposeCommands: true},
		Dive:       &diarkisexec.Options{ConfigPath: "configs/shared/dive.json"},
		DM:         &diarkisexec.Options{ConfigPath: "configs/shared/dm.json"},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")

	diarkisexec.SetServerCommandHandler(common.AppVersion, common.MatchingStartCmd, handleStartMatching)

	user.OnNewConnection(func(userData *user.User) {
		userData.SetAsString(meshAddrKey, mesh.GetMyEndPoint())
	})

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
		clients, _ := matching.GetTicketMemberClients(common.TeamTicketType, owner)
		memberIDs := make([]string, 0, len(clients))
		otherIDs := []string{}
		otherAddrs := []string{}
		for _, c := range clients {
			memberIDs = append(memberIDs, c.ID)
			if c.ID == owner.ID {
				continue
			}
			addr, _ := c.UserData[meshAddrKey].(string)
			otherIDs = append(otherIDs, c.ID)
			otherAddrs = append(otherAddrs, addr)
		}
		slices.Sort(memberIDs)
		logger.Infof("team ticket complete", "owner", owner.ID, "members", memberIDs)
		setTeamMembers(owner, memberIDs)

		setOnLeaveTicketRoom(common.TeamTicketType, owner)

		// Generate a payload to be sent to all ticket members.
		b, _ := json.Marshal(common.MatchingComplete{
			OwnerID:      owner.ID,
			CandidateIDs: memberIDs,
			TicketType:   common.TeamTicketType,
		})

		// Once the team has been created, the owner starts a new ticket to match
		// against another team. The other team members join the battle ticket with the owner.
		if err := startBattleTicket(owner, otherIDs, otherAddrs); err != nil {
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

	if err := json.Unmarshal(payload, &params); err != nil {
		respondError(userData, ver, cmd, fmt.Errorf("fail to parse payload. %w", err), next)
		return
	}

	if session.GetSessionID(common.PartySessionType, userData) == "" {
		startMatching(userData, nil, nil, ver, cmd, next)
		return
	}

	session.GetAllSessionMemberData(common.PartySessionType, userData, func(err error, members []*session.UserData) {
		if err != nil {
			respondError(userData, ver, cmd, fmt.Errorf("fail to get party members. %w", err), next)
			return
		}
		memberIDs := []string{}
		meshAddrs := []string{}
		for _, m := range members {
			if m.UserID == userData.ID {
				continue
			}
			memberIDs = append(memberIDs, m.UserID)
			meshAddrs = append(meshAddrs, m.MeshAddr)
		}
		startMatching(userData, memberIDs, meshAddrs, ver, cmd, next)
	})
}

func startMatching(owner *user.User, memberIDs []string, meshAddrs []string, ver uint8, cmd uint16, next func(error)) {
	var err error
	if 1+len(memberIDs) >= teamMaxMembers {
		logger.Infof("party is full, skip team ticket", "owner", owner.ID, "members", memberIDs)
		setTeamMembers(owner, append([]string{owner.ID}, memberIDs...))
		err = startBattleTicket(owner, memberIDs, meshAddrs)
	} else {
		err = startTeamTicket(owner, memberIDs, meshAddrs)
	}
	if err != nil {
		respondError(owner, ver, cmd, fmt.Errorf("fail to start matching ticket. %w", err), next)
		return
	}
	owner.ServerRespond(nil, ver, cmd, server.Ok, true)
	next(nil)
}

func respondError(userData *user.User, ver uint8, cmd uint16, err error, next func(error)) {
	b, _ := json.Marshal(common.CommandResponse{
		Error: &common.CommandError{
			Message: err.Error(),
		},
	})
	userData.ServerRespond(b, ver, cmd, server.Bad, true)
	next(err)
}

func setTeamMembers(owner *user.User, memberIDs []string) {
	b, _ := json.Marshal(memberIDs)
	owner.SetAsString(teamMembersKey, string(b))
}

func getTeamMembers(owner *user.User) []string {
	var memberIDs []string
	if v, ok := owner.GetAsString(teamMembersKey); ok {
		_ = json.Unmarshal([]byte(v), &memberIDs)
	}
	return memberIDs
}

func newTicketMembers(memberIDs []string, meshAddrs []string) *matching.TicketMembers {
	if len(memberIDs) == 0 {
		return nil
	}
	return &matching.TicketMembers{IDs: memberIDs, MeshAddrs: meshAddrs}
}

func startTeamTicket(owner *user.User, memberIDs []string, meshAddrs []string) error {
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
		TicketMembers: newTicketMembers(memberIDs, meshAddrs),
	}
	logger.Infof("start team ticket", "owner", owner.ID, "members", memberIDs, "meshAddrs", meshAddrs)
	return matching.StartTicketWithOptions(common.TeamTicketType, owner, options)
}

func startBattleTicket(owner *user.User, memberIDs []string, meshAddrs []string) error {
	searchTries := rand.Intn(10-1) + 1
	emptySearches := rand.Intn(searchTries) + 1
	options := &matching.TicketOptions{
		TicketParams: &matching.TicketParams{
			ProfileIDs:       []string{"battle"},
			MaxMembers:       teamMaxMembers * teamVsTeamSize,
			SearchInterval:   100, // 100ms
			SearchTries:      uint8(searchTries),
			EmptySearches:    uint8(emptySearches),
			TicketDuration:   ticketDuration,
			HowMany:          20,
			AddProperties:    map[string]int{"level": 1},
			SearchProperties: map[string][]int{"level": {1}},
		},
		TicketMembers: newTicketMembers(memberIDs, meshAddrs),
	}
	logger.Infof("start battle ticket", "owner", owner.ID, "members", memberIDs, "meshAddrs", meshAddrs)
	return matching.StartTicketWithOptions(common.BattleTicketType, owner, options)
}
