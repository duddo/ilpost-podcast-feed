package main

import (
	"log"
	"net/http"

	"ilpost-podcast-feed/pkg/endpoint"
)

func main() {
	http.HandleFunc("/podcast-list", endpoint.PodcastListHandler)
	http.HandleFunc("/feed", endpoint.BasicAuth(endpoint.FeedHandler))

	log.Printf("IlPost Podcast Feed %s", Version)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
