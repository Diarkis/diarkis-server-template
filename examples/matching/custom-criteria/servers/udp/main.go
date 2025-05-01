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

const userDataMatchingGeoCoordKey = "_matching_params_geocoord"
const ticketDuration = 15 // 15 seconds
const maxAllowedDistance = 100

func main() {
	logConfigPath := "configs/shared/log.json"
	meshConfigPath := ""
	log.EnableUnsafeLogging()

	diarkisexec.SetupDiarkis(logConfigPath, meshConfigPath, &diarkisexec.Modules{
		MatchMaker: &diarkisexec.Options{ConfigPath: "configs/shared/matching.json", ExposeCommands: true},
	})
	diarkisexec.SetupDiarkisUDPServer("configs/udp/main.json")

	diarkisexec.SetServerCommandHandler(common.AppVersion, common.MatchingStartCmd, handleStartMatching)

	setupMaching()

	diarkisexec.StartDiarkis()

}

func setupMaching() {
	const ticketType = uint8(1)
	matching.SetOnIssueTicket(ticketType, func(userData *user.User) *matching.TicketParams {
		coords, _ := userData.GetAsFloat64Array(userDataMatchingGeoCoordKey)

		// Pack the user coordinates as an application data.
		// This is a lazy way to pack the user coordinates in a byte array.
		// We recommend to use something else in production.
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, coords)

		// Randomize searchTries and emptySearches to not move to the wait mode at the same time considering all clients issue tickets at the same time.
		searchTries := rand.Intn(10-1) + 1
		emptySearches := rand.Intn(searchTries) + 1
		return &matching.TicketParams{
			ProfileIDs:     []string{"All"},
			MaxMembers:     2,
			SearchInterval: 100, // 100ms
			SearchTries:    uint8(searchTries),
			EmptySearches:  uint8(emptySearches),
			TicketDuration: ticketDuration,
			HowMany:        20,
			// Change here as you see fit according to your application needs
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
	})

	matching.SetOnTicketAllowMatchIf(ticketType, func(ticketProps *matching.TicketProperties, owner, candidate *user.User) bool {
		ownerCoords := []float64{}
		candidateCoords := []float64{}

		candidateProps, _ := ticketProps.GetCandidateByUID(candidate.ID)

		// Retrieve the application data from the ticket props.
		binary.Read(bytes.NewReader(ticketProps.GetOwner().ApplicationData), binary.BigEndian, &ownerCoords)
		binary.Read(bytes.NewReader(candidateProps.ApplicationData), binary.BigEndian, &candidateCoords)

		// Here compute the distance between the owner and candidate.
		distance := computeDistance(ownerCoords, candidateCoords)
		if distance > maxAllowedDistance {
			logger.Debugf("Cannot match candidate", "Distance", distance)
			return false
		}

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

	var latitude, longitude float64
	if false {
		// Here you are supposed to retrieve the client address.
		userAddr := userData.UDPState.GetClientAddr()
		latitude, longitude, _ = geolocalizeAddress(userAddr)
	} else {
		// As an example we use the client data as coordinates.
		latitude = params.Latitude
		longitude = params.Longitude
	}

	// Attach the matching params to the user data
	// so we can retrieve them when creating the ticket.
	userData.SetAsFloat64Array(userDataMatchingGeoCoordKey, []float64{latitude, longitude})
	defer userData.Delete(userDataMatchingGeoCoordKey)

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

func geolocalizeAddress(address string) (latitude float64, longitude float64, err error) {
	// Here you are supposed to lookup the address in a geo DB.
	latitude = 35.65877910898138
	longitude = 139.70128360250285
	return
}

func computeDistance(coords1, coord2 []float64) float64 {
	return maxAllowedDistance
}
