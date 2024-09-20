package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	visitPodcast("tienimi-bordone")

	return
	http.HandleFunc("/get-ilpost-feed", getIlpostFeedHandler)
	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Podcast struct {
	Title    string
	Url      string
	Episodes map[string]Episode
}

type Episode struct {
	Title  string
	Url    string
	Mp3Url string
}

func visitPodcast(podcastName string) {
	url := "https://api-prod.ilpost.it/frontend/podcast/list"

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the response body
	fmt.Println(string(body))
}

func getIlpostFeedHandler(w http.ResponseWriter, r *http.Request) {
	podcastName := r.URL.Query().Get("podcast-name")
	if podcastName == "" {
		http.Error(w, "Missing 'podcast-name' parameter", http.StatusBadRequest)
		return
	}

	// Create the XML response
	xmlResponse := RSS{
		Version: "2.0",
	}

	// Set the content type to XML
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)

	// Encode the response as XML
	if err := xml.NewEncoder(w).Encode(xmlResponse); err != nil {
		http.Error(w, "Failed to encode XML", http.StatusInternalServerError)
	}
}
