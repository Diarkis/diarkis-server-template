package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

// Response represents a generic API response
type Response struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// UserRankResponse represents user rank information
type UserRankResponse struct {
	UID  string `json:"uid"`
	Rank int    `json:"rank"`
}

// Main function to run the server if this file is executed directly
func main() {
	StartServer("8080")
}

// StartServer starts the dummy HTTP server
func StartServer(port string) {
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	// This endpoint is used to get authentication information from Diarkis UDP endpoint via API
	mux.HandleFunc("GET /endpoint/type/UDP/user/{uid}", authHandler)

	// This endpoint is used to get user rank
	mux.HandleFunc("GET /users/{uid}", userRankHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Starting API server on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	endpoints := map[string]string{
		"/endpoint/type/UDP/user/{userID}": "Proxy to localhost:7000 UDP user endpoint",
		"/users/{uid}":                     "Get user rank (random 1-200)",
	}

	response := Response{
		Status:    "ok",
		Message:   "Welcome to the dummy API server",
		Data:      endpoints,
		Timestamp: time.Now(),
	}
	sendJSONResponse(w, http.StatusOK, response)
}

// authHandler is used to get authentication information from Diarkis UDP endpoint via API
func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start authentication for user %s", r.PathValue("uid"))
	// Extract user ID from path
	uid := r.PathValue("uid")
	if uid == "" {
		response := Response{
			Status:    "error",
			Message:   "uid is required",
			Timestamp: time.Now(),
		}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// TODO: Pre-processing (authentication, authorization, etc.)

	// Get authentication information from Diarkis UDP endpoint via API
	diarkisAuthURL := fmt.Sprintf("http://localhost:7000/endpoint/type/UDP/user/%s", uid)
	proxyReq, err := http.NewRequest(r.Method, diarkisAuthURL, r.Body)
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		response := Response{
			Status:    "error",
			Message:   "Failed to create proxy request",
			Timestamp: time.Now(),
		}
		sendJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("Error making proxy request: %v", err)
		response := Response{
			Status:    "error",
			Message:   fmt.Sprintf("Failed to connect to target server: %v", err),
			Timestamp: time.Now(),
		}
		sendJSONResponse(w, http.StatusBadGateway, response)
		return
	}
	defer resp.Body.Close()

	// TODO: Post-processing (response creation, etc.)

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
	log.Printf("User %s rank retrieval was successful: %v", uid, resp.StatusCode)
}

// userRankHandler is used to get user rank
func userRankHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	if uid == "" {
		response := Response{
			Status:    "error",
			Message:   "uid is required",
			Timestamp: time.Now(),
		}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Generate random rank between 1-200
	rank := rand.N(200) + 1

	userRank := UserRankResponse{
		UID:  uid,
		Rank: rank,
	}

	response := Response{
		Status:    "ok",
		Message:   "User rank retrieved successfully",
		Data:      userRank,
		Timestamp: time.Now(),
	}

	log.Printf("Generated rank %d for user %s", rank, uid)
	sendJSONResponse(w, http.StatusOK, response)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
