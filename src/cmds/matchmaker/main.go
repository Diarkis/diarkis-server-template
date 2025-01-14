// © 2019-2024 Diarkis Inc. All rights reserved.

package matchmakercmds

import (
	customcmds "github.com/Diarkis/diarkis-server-template/cmds/custom"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

const sampleTicketType0 uint8 = 0 // max member 2
const sampleTicketType1 uint8 = 1 // max member 4
const sampleTicketType2 uint8 = 2 // max member 8

var logger = log.New("matching")

func Setup() {
	// Set up matching ticket with sampleTicketType
	// This callback is invoked when a new ticket is issued.
	matching.SetOnIssueTicket(sampleTicketType0, func(userData *user.User) *matching.TicketParams {
		// Randomize searchTries and emptySearches to not move to the wait mode at the same time considering all clients issue tickets at the same time.
		searchTries := util.RandomInt(1, 50) // max 50 tries for 5 seconds (searchInterval * searchTries)
		emptySearches := util.RandomInt(1, searchTries)
		return &matching.TicketParams{
			ProfileIDs:     []string{"RankMatch"},
			MaxMembers:     2,
			SearchInterval: 100, // 100ms
			SearchTries:    uint8(searchTries),
			EmptySearches:  uint8(emptySearches),
			TicketDuration: 60, // 1m
			HowMany:        20,
			// Change here as you see fit according to your application needs
			Tags: []string{""},
			// Change here as you see fit according to your application needs
			AddProperties:    map[string]int{"rank": 1},
			SearchProperties: map[string][]int{"rank": {1, 2, 3, 4, 5}},
		}
	})

	// Set up a callback on ticket match with sampleTicketType
	// This callback is invoked when a ticket finds a match.
	// By returning true, the callback allows the found match to be matched.
	// By returning false, the callback rejects the found match and ignores it.
	matching.SetOnTicketMatch(sampleTicketType0,
		func(t *matching.Ticket, matchedUser, ownerUser *user.User, roomID string, memberIDs []string) bool {

			// add custom logic to decide matchmaking completion here
			return false
		})

	// Set up a callback on ticket complete with sampleTicketType
	// This callback is invoked when a ticket completes the matchmaking.
	matching.SetOnTicketComplete(sampleTicketType0, func(ticketProps *matching.TicketProperties, owner *user.User) []byte {
		memberIDs, _ := matching.GetTicketMemberIDs(sampleTicketType0, owner)

		list := make([]string, len(memberIDs)+1)

		// the first element is the owner user ID
		list[0] = owner.ID

		index := 1

		for i := 0; i < len(memberIDs); i++ {
			list[index] = memberIDs[i]
			index++
		}

		return packet.StringListToBytes(list)
	})

	// If OnTicketMemberLeave is not set, notification is NOT sent when matched member leaves.
	matching.SetOnTicketMemberLeaveAnnounce(sampleTicketType0, func(ticket *matching.Ticket, leftUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte) {
		logger.Sys("Matched Member Leave Announce")
		return customcmds.CustomVer, customcmds.MatchedMemberLeaveCmdID, []byte(leftUser.ID)
	})

	matching.SetOnIssueTicket(sampleTicketType1, func(userData *user.User) *matching.TicketParams {
		// Randomize searchTries and emptySearches to not move to the wait mode at the same time considering all clients issue tickets at the same time.
		searchTries := util.RandomInt(1, 50) // max 50 tries for 5 seconds (searchInterval * searchTries)
		emptySearches := util.RandomInt(1, searchTries)
		return &matching.TicketParams{
			ProfileIDs:     []string{"RankMatch20"},
			MaxMembers:     4,
			SearchInterval: 100, // 100ms
			SearchTries:    uint8(searchTries),
			EmptySearches:  uint8(emptySearches),
			TicketDuration: 60, // 1m
			HowMany:        20,
			// Change here as you see fit according to your application needs
			Tags: []string{""},
			// Change here as you see fit according to your application needs
			AddProperties:    map[string]int{"rank": 1},
			SearchProperties: map[string][]int{"rank": {1, 2, 3, 4, 5}},
		}
	})
	matching.SetOnTicketMatch(sampleTicketType1,
		func(t *matching.Ticket, matchedUser, ownerUser *user.User, roomID string, memberIDs []string) bool {
			// add custom logic to decide matchmaking completion here
			return false
		})

	matching.SetOnTicketComplete(sampleTicketType1, func(ticketProps *matching.TicketProperties, owner *user.User) []byte {
		memberIDs, _ := matching.GetTicketMemberIDs(sampleTicketType1, owner)

		list := make([]string, len(memberIDs)+1)

		// the first element is the owner user ID
		list[0] = owner.ID

		index := 1

		for i := 0; i < len(memberIDs); i++ {
			list[index] = memberIDs[i]
			index++
		}

		return packet.StringListToBytes(list)
	})
	// If OnTicketMemberLeave is not set, notification is NOT sent when matched member leaves.
	matching.SetOnTicketMemberLeaveAnnounce(sampleTicketType1, func(ticket *matching.Ticket, leftUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte) {
		logger.Sys("Ticket:1 Matched Member Leave Announce")
		return customcmds.CustomVer, customcmds.MatchedMemberLeaveCmdID, []byte(leftUser.ID)
	})

	matching.SetOnIssueTicket(sampleTicketType2, func(userData *user.User) *matching.TicketParams {
		// Randomize searchTries and emptySearches to not move to the wait mode at the same time considering all clients issue tickets at the same time.
		searchTries := util.RandomInt(1, 50) // max 50 tries for 5 seconds (searchInterval * searchTries)
		emptySearches := util.RandomInt(1, searchTries)
		return &matching.TicketParams{
			ProfileIDs:     []string{"RankMatch20"},
			MaxMembers:     8,
			SearchInterval: 100, // 100ms
			SearchTries:    uint8(searchTries),
			EmptySearches:  uint8(emptySearches),
			TicketDuration: 20, // 20s
			HowMany:        20,
			// Change here as you see fit according to your application needs
			Tags: []string{""},
			// Change here as you see fit according to your application needs
			AddProperties:    map[string]int{"rank": 1},
			SearchProperties: map[string][]int{"rank": {1, 2, 3, 4, 5}},
		}
	})
	matching.SetOnTicketMatch(sampleTicketType2,
		func(t *matching.Ticket, matchedUser, ownerUser *user.User, roomID string, memberIDs []string) bool {
			// add custom logic to decide matchmaking completion here
			return false
		})

	matching.SetOnTicketComplete(sampleTicketType2, func(ticketProps *matching.TicketProperties, owner *user.User) []byte {
		memberIDs, _ := matching.GetTicketMemberIDs(sampleTicketType2, owner)

		list := make([]string, len(memberIDs)+1)

		// the first element is the owner user ID
		list[0] = owner.ID

		index := 1

		for i := 0; i < len(memberIDs); i++ {
			list[index] = memberIDs[i]
			index++
		}

		return packet.StringListToBytes(list)
	})
	// If OnTicketMemberLeave is not set, notification is NOT sent when matched member leaves.
	matching.SetOnTicketMemberLeaveAnnounce(sampleTicketType2, func(ticket *matching.Ticket, leftUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte) {
		logger.Sys("Ticket:2 Matched Member Leave Announce")
		return customcmds.CustomVer, customcmds.MatchedMemberLeaveCmdID, []byte(leftUser.ID)
	})
}
