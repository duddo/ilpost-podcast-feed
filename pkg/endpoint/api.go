package endpoint

import (
	"encoding/xml"
	"fmt"
	ilpostapi "ilpost-podcast-feed/pkg/ilpost_api"
	rssfeed "ilpost-podcast-feed/pkg/rss_feed"
	"log"
	"net/http"
)

func PodcastListHandler(w http.ResponseWriter, r *http.Request) {
	response, err := ilpostapi.FetchPodcastList(nil)
	if err != nil {
		log.Println("Error fetching podcasts:", err)
		return
	}

	_ = response

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	listItems := ""
	for i, podcast := range response.Data {
		if podcast.Title == "" {
			continue
		}

		listItems += fmt.Sprintf(`<li>%d - %s <a href="%s">Il Post</a></li> <a href="%s">Feed</a></li>`, i, podcast.Title, podcast.URL, "/feed?podcast-name="+podcast.Slug)
	}

	// HTML response with the list
	html := fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Hello World</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f0f0f0;
                    display: flex;
                    justify-content: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    padding: 20px;
                    background-color: #fff;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                }
                h1 {
                    color: #333;
                    font-size: 2.5rem;
                }
                ul {
                    list-style-type: none;
                    padding: 0;
                }
                li {
                    margin: 10px 0;
                }
                a {
                    color: #1e90ff;
                    text-decoration: none;
                }
                a:hover {
                    text-decoration: underline;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>Hello, World!</h1>
                <ul>
                    %s
                </ul>
            </div>
        </body>
        </html>
    `, listItems)

	w.Write([]byte(html))
}

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	podcastName := r.URL.Query().Get("podcast-name")
	if podcastName == "" {
		http.Error(w, "Missing 'podcast-name' parameter", http.StatusBadRequest)
		return
	}

	// Create the XML response
	xmlResponse := rssfeed.RSS{
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
