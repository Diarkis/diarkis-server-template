// © 2019-2024 Diarkis Inc. All rights reserved.

package matchmakercmds

import (
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

const sampleTicketType0 uint8 = 0 // max member 2
const sampleTicketType1 uint8 = 1 // max member 4

var logger = log.New("matching")

func Setup() {
	// Set up matching ticket with sampleTicketType
	// This callback is invoked when a new ticket is issued.
	matching.SetOnIssueTicket(sampleTicketType0, func(userData *user.User) *matching.TicketParams {
		return &matching.TicketParams{
			ProfileIDs:     []string{"RankMatch"},
			MaxMembers:     2,
			SearchInterval: 100, // 100ms
			SearchTries:    uint8(util.RandomInt(0, 300)),
			EmptySearches:  3,
			TicketDuration: 60, // 1m
			HowMany:        20,
			// Change here as you see fit according to your application needs
			Tag: "",
			// Change here as you see fit according to your application needs
			AddProperties:       map[string]int{"rank": 1},
			SearchProperties:    map[string][]int{"rank": {1, 2, 3, 4, 5}},
			IsLeaveNotification: true,
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

	matching.SetOnIssueTicket(sampleTicketType1, func(userData *user.User) *matching.TicketParams {
		return &matching.TicketParams{
			ProfileIDs:     []string{"RankMatch"},
			MaxMembers:     4,
			SearchInterval: 100, // 100ms
			SearchTries:    uint8(util.RandomInt(0, 300)),
			EmptySearches:  3,
			TicketDuration: 180, // 1m
			HowMany:        20,
			// Change here as you see fit according to your application needs
			Tag: "",
			// Change here as you see fit according to your application needs
			AddProperties:       map[string]int{"rank": 1},
			SearchProperties:    map[string][]int{"rank": {1, 2, 3, 4, 5}},
			IsLeaveNotification: true,
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
}
