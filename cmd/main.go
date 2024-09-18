// cmd/main.go
package main

import (
	"encoding/xml"
	"log"
	"net/http"

	"github.com/gocolly/colly"
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
	pod := Podcast{
		Url:      "https://www.ilpost.it/podcasts/" + podcastName,
		Episodes: make(map[string]Episode),
	}

	episodesCollector := colly.NewCollector(
	// Visit only domain
	//colly.AllowedDomains("www.ilpost.it"),
	// Parallelism
	//colly.Async(true),
	)

	mediaCollector := colly.NewCollector()

	episodesCollector.OnHTML("h1", func(e *colly.HTMLElement) {
		log.Println("Found podcast: " + e.Text)
		pod.Title = e.Text
	})

	episodesCollector.OnHTML("h3", func(e *colly.HTMLElement) {
		ep := Episode{
			Title: e.Text,
			Url:   e.ChildAttr("a", "href"),
		}
		pod.Episodes[ep.Title] = ep
		log.Println("Found " + ep.Title + " - " + ep.Url)

		err := mediaCollector.Visit(ep.Url)
		if err != nil {
			log.Fatal(err)
		}
	})

	mediaCollector.OnHTML("main.container", func(e *colly.HTMLElement) {
		title := e.ChildText("h1")
		ep := pod.Episodes[title]
		mp3Url := e.ChildText("div._total-duration_dkzl3_145")
		ep.Mp3Url = mp3Url
		log.Println("  Episode " + ep.Title + " url: " + ep.Mp3Url)
	})

	// Handle errors during scraping
	episodesCollector.OnError(func(r *colly.Response, err error) {
		log.Println("Request failed:", r.StatusCode, err)
	})
	mediaCollector.OnError(func(r *colly.Response, err error) {
		log.Println("Request failed:", r.StatusCode, err)
	})

	// Before making a request print "Visiting ..."
	episodesCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})
	mediaCollector.OnRequest(func(r *colly.Request) {
		log.Println("  Visiting", r.URL.String())
	})

	// Visit the website
	err := episodesCollector.Visit(pod.Url)
	if err != nil {
		log.Fatal(err)
	}
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
