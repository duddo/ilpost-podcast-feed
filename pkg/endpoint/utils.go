package endpoint

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	ilpostapi "ilpost-podcast-feed/pkg/ilpost_api"
)

func StartServer() {
	var cookieCache = make(CookieCache)

	http.Handle("/podcast-list", appHandler(podcastListHandler))
	http.Handle("/feed", basicAuth(&cookieCache, feedHandler))
	http.Handle("/test", appHandler(testHandler))

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func BuildFeed(episodes ilpostapi.PodcastEpisodesResponse) RSS {
	var items []Item

	for _, episode := range episodes.Data {
		items = append(items, Item{
			Title:       episode.Title,
			Link:        episode.URL,
			PubDate:     episode.Date.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
			Guid:        strconv.Itoa(episode.ID),
			Description: episode.ShareURL,
			Enclosure: Enclosure{
				Url:    episode.EpisodeRawURL,
				Length: "", // TODO provide a length without dowloading the mp3?
				Type:   "audio/mpeg",
			},
			ContentEncoded: CDATA(episode.ContentHTML),
			Duration:       strconv.Itoa(episode.Milliseconds * 1000),
			Subtitle:       "",
			Summary:        "",
			Keywords:       "",
			Author:         episode.Author,
			Explicit:       "no",
			Block:          "no",
		})
	}

	// TODO do something if no episode returned
	channelDescr := episodes.Data[0].Parent

	return RSS{
		Version: "2.0",
		Content: "http://purl.org/rss/1.0/modules/content/",
		Itunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
		Channel: Channel{
			Title:       channelDescr.Title,
			Description: CDATA(channelDescr.Description),
			Link:        "https://www.ilpost.it/podcasts/" + channelDescr.Slug,
			Language:    "it",
			Generator:   "https://github.com/duddo/ilpost-podcast-feed",
			Subtitle:    CDATA(channelDescr.Description),
			Summary:     CDATA(channelDescr.Description),
			Author:      channelDescr.Author,
			Block:       "no",
			Explicit:    "no",
			Image:       Image{Href: channelDescr.Image},

			Items: items,
		},
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) *appError {
	ilpostData, err := unmarshalFromFile("./pkg/endpoint/bordone.json")
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

func unmarshalFromFile(filename string) (*ilpostapi.PodcastEpisodesResponse, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// return nil, fmt.Errorf("this is very bad")

	filedata, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data ilpostapi.PodcastEpisodesResponse
	err = json.Unmarshal(filedata, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
