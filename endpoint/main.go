package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	baseURL = "https://api.github.com"
)

type Repository struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Define the routes
	router.HandleFunc("/repos/{owner}/{repo}", getRepo).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getRepo(w http.ResponseWriter, r *http.Request) {
	// Get the owner and repo parameters from the request URL
	params := mux.Vars(r)
	owner := params["owner"]
	repo := params["repo"]

	// Build the request URL
	url := fmt.Sprintf("%s/repos/%s/%s", baseURL, owner, repo)

	// Create a new HTTP client with a timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Create a new request with a custom user agent header
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("User-Agent", "my-github-api-client")

	// Send the request and measure the elapsed time
	startTime := time.Now()
	resp, err := client.Do(req)
	elapsedTime := time.Since(startTime)

	if err != nil {
		log.Printf("Error sending request: %v", err)
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response status code: %d %s", resp.StatusCode, resp.Status)
		http.Error(w, "Error response status code", resp.StatusCode)
		return
	}

	// Parse the response body into a Repository object
	var repoObj Repository
	err = json.NewDecoder(resp.Body).Decode(&repoObj)
	if err != nil {
		log.Printf("Error decoding response body: %v", err)
		http.Error(w, "Error decoding response body", http.StatusInternalServerError)
		return
	}

	// Log the request and response metadata
	logData := map[string]interface{}{
		"method":      req.Method,
		"url":         req.URL.String(),
		"headers":     req.Header,
		"elapsedTime": elapsedTime.Seconds(),
		"status":      resp.Status,
		"headersOut":  resp.Header,
		"body":        repoObj,
	}
	logJSON, err := json.MarshalIndent(logData, "", "    ")
	if err != nil {
		log.Printf("Error encoding log data: %v", err)
		http.Error(w, "Error encoding log data", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(logJSON))

	// Write the Repository object as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(repoObj)
	if err != nil {
		log.Printf("Error encoding response body: %v", err)
		http.Error(w, "Error encoding response body", http.StatusInternalServerError)
		return
	}
}

