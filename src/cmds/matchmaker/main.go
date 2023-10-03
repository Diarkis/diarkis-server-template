package matchmakercmds

import (
	"fmt"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

const sampleTicketType uint8 = 0

var logger = log.New("matching")

func Expose(rootpath string) {
	// Set up matching package and load configuration file
	matching.Setup(fmt.Sprintf("%s/configs/shared/matching.json", rootpath))

	// Set up matching ticket
	matching.SetOnIssueTicket(sampleTicketType, func(userData *user.User) *matching.TicketParams {
		return &matching.TicketParams{
			ProfileIDs:     []string{"RankMatch"},
			MaxMembers:     2,
			SearchInterval: 100, // 100ms
			SearchTries: uint8(util.RandomInt(0, 300)),
			EmptySearches: 3,
			TicketDuration: 60,  // 1m
			HowMany:        20,
			// Change here as you see fit according to your application needs
			Tag: "",
			// Change here as you see fit according to your application needs
			AddProperties: map[string]int{"rank": 1},
			SearchProperties: map[string][]int{"rank": []int{1, 2, 3, 4, 5}},
		}
	})
	matching.SetOnTicketMatch(sampleTicketType,
		func(t *matching.Ticket, owner, userData *user.User, roomID string, memberIDs []string) bool {

		// add custom logic to decide matchmaking completion here
		return false
	})
	matching.SetOnTicketComplete(sampleTicketType, func(ticketProps *matching.TicketProperties, owner *user.User) []byte {
		memberIDs, _ := matching.GetTicketMemberIDs(sampleTicketType, owner)

		list := make([]string, len(memberIDs) + 1)

		// the first element is the owner user ID
		list[0] = owner.ID

		index := 1

		for i := 0; i < len(memberIDs); i++ {
			list[index] = memberIDs[i]
			index++
		}

		return packet.StringListToBytes(list)
	})

	// Expose built-in commands to the client
	matching.ExposeCommands()
}
