package main

import (
	"ilpost-podcast-feed/pkg/endpoint"
	"log"
)

func main() {
	log.Printf("IlPost Podcast Feed %s", Version)

	endpoint.StartServer()
}
