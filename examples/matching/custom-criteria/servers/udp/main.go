package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/Diarkis/diarkis-server-template/examples/matching/custom-criteria/common"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

var logger = log.New("UDP")

const (
	// maxAllowedDistance is the maximum distance allowed for two users to match together.
	maxAllowedDistance = 5000
	ticketType         = uint8(1)
	ticketDuration     = 15 // 15 seconds
)

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := ""

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json", ExposeCommands: true},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")

	diarkisexec.SetServerCommandHandler(common.AppVersion, common.MatchingStartCmd, handleStartMatching)

	setupMaching()

	diarkisexec.StartDiarkis()

}

func setupMaching() {

	matching.SetOnTicketAllowMatchIf(ticketType, func(ticketProps *matching.TicketProperties, owner, candidate *user.User) bool {
		var ownerCoords, candidateCoords common.Coordinates

		candidateProps, _ := ticketProps.GetCandidateByUID(candidate.ID)

		// Retrieve the application data from the ticket properties.
		binary.Read(bytes.NewReader(ticketProps.GetOwner().ApplicationData), binary.BigEndian, &ownerCoords)
		binary.Read(bytes.NewReader(candidateProps.ApplicationData), binary.BigEndian, &candidateCoords)

		// Here we compute the distance between the owner and candidate.
		distance := common.ComputeDistanceHubeny(ownerCoords, candidateCoords)
		if distance > maxAllowedDistance {
			logger.Debugf("Cannot match candidates", "Owner", owner.ID, "Candidate", candidate.ID,
				"Distance", distance, "MaxAllowedDistance", maxAllowedDistance)
			return false
		}
		logger.Debugf("Candidates matched successfully", "Owner", owner.ID, "Candidate", candidate.ID,
			"Distance", distance, "MaxAllowedDistance", maxAllowedDistance)
		return true
	})

	matching.SetOnTicketCompleteWithProfileID(ticketType, func(ticketProps *matching.TicketProperties, owner *user.User, tag []string, profileID string) []byte {
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

	var coords *common.Coordinates
	if false {
		// As stated by the note in the README.md, if you were to implement a geolocation database,
		// you could attempt to retrieve the coordinates of the client by their IP address.
		//
		// We do not provide this, but we opt to leave it as a functional stub to encourage creative
		// use of Diarkis Matchmaker.
		userAddr := userData.UDPState.GetClientAddr()
		latitude, longitude, _ := geolocalizeAddress(userAddr)
		coords = &common.Coordinates{Latitude: latitude, Longitude: longitude}
	} else {
		// As an example we use the client data as coordinates:
		coords = &common.Coordinates{Latitude: params.Latitude, Longitude: params.Longitude}
	}

	// Pack the user coordinates as an application data.
	// This is a lazy way to pack the user coordinates in a byte array.
	// We recommend to use something else in production.
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, coords)

	// Randomize searchTries and emptySearches to not move to the wait mode at the same time considering all clients issue tickets at the same time.
	searchTries := randomInt(1, 10)
	emptySearches := randomInt(1, max(2, searchTries))
	ticketParams := &matching.TicketParams{
		ProfileIDs:     []string{"All"},
		MaxMembers:     2,
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
		// Associated some metadata related to the client to the ticket.
		// We will retrieve them in `matching.SetOnTicketAllowMatchIf` to check
		// if two candidates can match together.
		ApplicationData: buf.Bytes(),
	}

	err = matching.StartTicketWithTicketParams(1, userData, ticketParams)
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

func geolocalizeAddress(address string) (latitude float64, longitude float64, err error) {
	// Here you may look up the address in a geolocation database.
	_ = address
	latitude = 35.65877910898138
	longitude = 139.70128360250285
	return
}

func randomInt(minValue, maxValue int) int {
	return rand.Intn(maxValue-minValue) + minValue
}
