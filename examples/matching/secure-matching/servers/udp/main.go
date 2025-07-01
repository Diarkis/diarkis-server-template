// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/Diarkis/diarkis-server-template/examples/matching/secure-matching/common"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/user"
)

var logger = log.New("UDP")

const ticketType = uint8(1)

// UserResponse represents the response from the user API
type UserResponse struct {
	Data struct {
		Rank int `json:"rank"`
	} `json:"data"`
}

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := ""

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json", ExposeCommands: true},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")

	setupMatching()

	diarkisexec.StartDiarkis()
}

func setupMatching() {
	matching.SetOnIssueTicket(ticketType, func(userData *user.User) *matching.TicketParams {

		// Get rank from API
		rank, err := getRankFromAPI(userData.ID)
		if err != nil {
			logger.Errorf("Failed to get rank from API for user %s: %v, using fallback random rank", userData.ID, err)
			rank = randomInt(1, 200) // fallback to random if API fails
		}
		logger.Sysf("onIssueTicket: You got rank via API:", "rank", rank)

		searchRankProperties := generateSearchProperties(rank)
		logger.Sysf("onIssueTicket: searchRankProperties:", "searchRankProperties", searchRankProperties)

		// Randomize searchTries and emptySearches to not move to the wait mode at the same time considering all clients issue tickets at the same time.
		maxSearchTries := 50
		searchTries := randomInt(1, maxSearchTries) // max 50 tries for 5 seconds (searchInterval * searchTries)
		emptySearches := randomInt(1, max(2, searchTries))

		return &matching.TicketParams{
			ProfileIDs:       []string{"RankMatch"},
			MaxMembers:       2,
			SearchInterval:   100, // 100ms
			SearchTries:      uint8(searchTries),
			EmptySearches:    uint8(emptySearches),
			TicketDuration:   60,
			HowMany:          20,
			AddProperties:    map[string]int{"rank": rank},
			SearchProperties: map[string][]int{"rank": searchRankProperties},
		}
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

// getRankFromAPI fetches user rank from localhost:8080/users/{uid}
// This is intended to simulate an API server provided by the service side.
// In a real-world application, this would be replaced with an actual API server.
func getRankFromAPI(uid string) (int, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:8080/users/%s", uid)
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return userResp.Data.Rank, nil
}

// randomInt generates a random integer between minValue and maxValue
func randomInt(minValue, maxValue int) int {
	return rand.N(maxValue-minValue) + minValue
}

// generateSearchProperties generates search range properties.
// The result is as follows:
//
//	rank 0: [0, 10, 20]: It means that the user can match with users in rank from 0 to 20.
//	rank 1-10: [0, 10, 20, 30]: It means that the user can match with users in rank from 0 to 30.
//	rank 11-20: [0, 10, 20, 30, 40]: It means that the user can match with users in rank from 0 to 40.
//	rank 21-30: [10, 20, 30, 40, 50]: It means that the user can match with users in rank from 1 to 50.
//	rank 31-40: [20, 30, 40, 50, 60]: It means that the user can match with users in rank from 11 to 60.
func generateSearchProperties(rank int) []int {
	searchStep := 10 // This should be the same as 'rank' in the matching profile
	currentBucket := int(math.Ceil(float64(rank) / float64(searchStep)))
	searchBucketRange := 2 // Search ±2 buckets from the current bucket
	var searchRankProperties []int
	minBucket := currentBucket - searchBucketRange
	if minBucket < 0 {
		minBucket = 0
	}
	maxBucket := currentBucket + searchBucketRange
	for bucket := minBucket; bucket <= maxBucket; bucket++ {
		searchRankProperties = append(searchRankProperties, bucket*searchStep)
	}
	return searchRankProperties
}
