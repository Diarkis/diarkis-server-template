// © 2019-2024 Diarkis Inc. All rights reserved.

package httpcmds

// rootpath is defined in cmds/main.go

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Diarkis/diarkis/matching"
	"github.com/Diarkis/diarkis/server/http"
)

func exposeMatchMaker() {
	http.Post("/mm/add/:mmID/:uniqueID/:ttl", addToMatchMaker)
	http.Delete("/mm/rm/:mmID", removeFromMatchMaker)
	http.Post("/mm/search/:mmIDs/:limit", searchMatchMaker)
	defineMatchMakerRules()
}

func defineMatchMakerRules() {
	// rank will match rank range of 10. i.e. 0 to 10, 11 to 20, 21 to 30...
	rankRule := make(map[string]int)
	rankRule["rank"] = 10
	matching.Define("rank", rankRule)
	// score will match score range of 1000. i.e. 0 to 1000, 1001 to 2000...
	scoreRule := make(map[string]int)
	scoreRule["score"] = 1000
	matching.Define("score", scoreRule)
	// sample for UDP and TCP server
	rateAndPlayRule := make(map[string]int)
	rateAndPlayRule["rate"] = 1
	rateAndPlayRule["play"] = 1
	matching.Define("RateAndPlay", rateAndPlayRule)
	rankProps := make(map[string]int)
	rankProps["rank"] = 10
	matching.Define("RankMatch", rankProps)
	rank20Props := make(map[string]int)
	rank20Props["rank"] = 20
	matching.Define("RankMatch20", rank20Props)
	rank50Props := make(map[string]int)
	rank50Props["rank"] = 50
	matching.Define("RankMatch50", rank50Props)
}

func addToMatchMaker(res *http.Response, req *http.Request, params *http.Params, next func(error)) {
	mmID, err := params.GetAsString("mmID")
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	uniqueID, err := params.GetAsString("uniqueID")
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	ttl, err := params.GetAsInt64("ttl")
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	dataJSON := req.Req.FormValue("props")
	if dataJSON == "" {
		err := errors.New("MatchMaker data is missing")
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	// MatchMaker data is expected to be JSON
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(dataJSON), &m)
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	data := make(map[string]int)
	for k, v := range m {
		if _, ok := v.(float64); !ok {
			// invalid data
			err := errors.New("Invalid MatchMaker data. The all values must be numbers")
			res.Respond(err.Error(), http.Bad)
			next(err)
			return
		}
		data[k] = int(v.(float64))
	}
	// meta data to be stored is expected to be JSON
	metadataJSON := req.Req.FormValue("metadata")
	if metadataJSON == "" {
		err := errors.New("Missing metadata")
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	metadata := make(map[string]interface{})
	err = json.Unmarshal([]byte(metadataJSON), &metadata)
	matching.Add(mmID, uniqueID, "", data, metadata, ttl, 2)
	res.Respond("OK", http.Ok)
	next(nil)
}

func removeFromMatchMaker(res *http.Response, req *http.Request, params *http.Params, next func(error)) {
	mmID, err := params.GetAsString("mmID")
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	// uniqueIDList is expected to be a JSON array
	uniqueIDList, err := params.GetAsString("uniqueIDList")
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	idlist := make([]string, 0)
	err = json.Unmarshal([]byte(uniqueIDList), &idlist)
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	matching.Remove(mmID, idlist, 2)
	res.Respond("OK", http.Ok)
	next(nil)
}

func searchMatchMaker(res *http.Response, req *http.Request, params *http.Params, next func(error)) {
	// mmIDs is comma separated list of matching IDs
	mmIDs, err := params.GetAsString("mmIDs")
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	mmIDList := strings.Split(mmIDs, ",")
	// search properties are expected to be JSON
	propsJSON := req.Req.FormValue("props")
	if propsJSON == "" {
		err := errors.New("MatchMaker data is missing")
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	// MatchMaker data is expected to be JSON
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(propsJSON), &m)
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	props := make(map[string]int)
	for k, v := range m {
		if _, ok := v.(float64); !ok {
			// invalid data
			err := errors.New("Invalid MatchMaker data. The all values must be numbers")
			res.Respond(err.Error(), http.Bad)
			next(err)
			return
		}
		props[k] = int(v.(float64))
	}
	// search result limit
	limit, err := params.GetAsInt("limit")
	if err != nil {
		res.Respond(err.Error(), http.Bad)
		next(err)
		return
	}
	matching.Search(mmIDList, "", props, limit, func(err error, results []interface{}) {
		if err != nil {
			res.Respond(err.Error(), http.Bad)
			next(err)
			return
		}
		response := make(map[string]interface{})
		response["props"] = props
		response["results"] = results
		encodedResponse, err := json.Marshal(response)
		if err != nil {
			res.Respond(err.Error(), http.Bad)
			next(err)
			return
		}
		res.Respond(string(encodedResponse), http.Ok)
		next(nil)
	})
}
