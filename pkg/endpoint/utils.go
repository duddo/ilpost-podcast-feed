package endpoint

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	ilpostapi "ilpost-podcast-feed/pkg/ilpost_api"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func StartServer() {
	var cookieCache = make(CookieCache)

	http.HandleFunc("/podcast-list", podcastListHandler)
	http.HandleFunc("/feed", basicAuth(&cookieCache, feedHandler))
	http.Handle("/test", appHandler(testHandler))

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func areCookieValid(cookies []*http.Cookie) bool {
	if cookies == nil {
		return false
	}

	for _, cookie := range cookies {
		if !cookie.Expires.IsZero() && cookie.Expires.Before(time.Now()) {
			return false
		}
	}

	return true
}

func BuildFeed(episodes ilpostapi.PodcastEpisodesResponse) RSS {
	var items []Item

	for _, episode := range episodes.Data {
		items = append(items, Item{
			Title:       episode.Title,
			Link:        episode.URL,
			PubDate:     episode.Date.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
			Guid:        fmt.Sprintf("%d", episode.ID),
			Description: episode.ShareURL,
			Enclosure: Enclosure{
				Url:    episode.EpisodeRawURL,
				Length: "200",
				Type:   "audio/mpeg",
			},
		})
	}

	return RSS{
		Version: "2.0",
		Content: "http://purl.org/rss/1.0/modules/content/",
		Channel: Channel{
			Title:       "BOH MAH",
			Link:        "http...",
			Description: "descrizione...",
			Language:    "it",
			Generator:   "https://github.com/duddo/ilpost-podcast-feed",
			Items:       items,
		},
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) *appError {
	ilpostData, err := unmarshal("./pkg/endpoint/bordone.json")
	if err != nil {
		return &appError{err, "Can't unmarshal test json", http.StatusInternalServerError}
	}

	episodes := BuildFeed(*ilpostData)

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return &appError{err, "Can't write XML header", http.StatusInternalServerError}
	}

	err = xml.NewEncoder(w).Encode(episodes)
	if err != nil {
		return &appError{err, "Can't encode data to XML", http.StatusInternalServerError}
	}

	return nil
}

func unmarshal(filename string) (*ilpostapi.PodcastEpisodesResponse, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return nil, fmt.Errorf("brutta storia")

	filedata, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data ilpostapi.PodcastEpisodesResponse
	json.Unmarshal(filedata, &data)

	return &data, nil
}
