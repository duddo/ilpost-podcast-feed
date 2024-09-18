// cmd/main.go
package main

import (
	"encoding/xml"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/get-ilpost-feed", getIlpostFeedHandler)
	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
