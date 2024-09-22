package endpoint

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	ilpostapi "ilpost-podcast-feed/pkg/ilpost_api"
)

type key int

const userKey key = 0

func basicAuth(cookieCache *CookieCache, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the value of the "Authorization" header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the header contains "Basic" prefix
		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Decode the base64-encoded username:password
		payload, _ := base64.StdEncoding.DecodeString(authHeaderParts[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Try to retrieve valid cookies from cache
		var username = pair[0]
		var password = pair[1]
		var cookies = (*cookieCache)[username]
		if !areCookieValid(cookies) {
			// Ask for actual authorization to ilpost endpoint, otherwise
			cks, err := ilpostapi.Login(username, password)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+err.Error()+`"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			} else {
				cookies = cks
				// Store the cookies in cache
				(*cookieCache)[username] = cookies
			}
		}

		// Store the username in the context to pass it to the next handler
		ctx := context.WithValue(r.Context(), userKey, cookies)

		// Pass the request with the context to the next handler
		next(w, r.WithContext(ctx))
	}
}

func podcastListHandler(w http.ResponseWriter, _ *http.Request) *appError {
	response, err := ilpostapi.FetchPodcastList(nil)
	if err != nil {
		return &appError{err, "Can't fetch podcasts", http.StatusInternalServerError}
	}

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

	_, err = w.Write([]byte(html))
	if err != nil {
		return &appError{err, "Can't write response", http.StatusInternalServerError}
	}

	return nil
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	podcastName := r.URL.Query().Get("podcast-name")
	if podcastName == "" {
		http.Error(w, fmt.Sprintf("ERROR %d - Missing 'podcast-name' parameter", http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Retrieve the cookies from the context
	cookies, ok := r.Context().Value(userKey).([]*http.Cookie)
	if !ok {
		http.Error(w, fmt.Sprintf("ERROR %d - User not authenticated", http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	episodes, err := ilpostapi.FetchPodcastEpisodes(cookies, podcastName, 1, 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR %d - Failed to retrieve episodes: %s", http.StatusInternalServerError, err.Error()), http.StatusInternalServerError)
		return
	}

	// Create the XML response
	xmlResponse := BuildFeed(episodes)

	// Set the content type to XML
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)

	// Encode the response as XML
	if err := xml.NewEncoder(w).Encode(xmlResponse); err != nil {
		http.Error(w, fmt.Sprintf("ERROR %d - Failed to encode XML", http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
