package endpoint

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"

	ilpostapi "ilpost-podcast-feed/pkg/ilpost_api"
)

type key int

const userKey key = 0

func basicAuth(cookieCache *CookieCache, next appHandler) appHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)

			return &appError{nil, "unauthorized, missing Authorization header", http.StatusUnauthorized}
		}

		// Check if the header contains "Basic" prefix
		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return &appError{nil, "missing Basic prefix in header", http.StatusBadRequest}
		}

		// Decode the base64-encoded username:password
		payload, err := base64.StdEncoding.DecodeString(authHeaderParts[1])
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return &appError{err, "can't decode username:password", http.StatusBadRequest}
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Enter username and password"`)
			return &appError{nil, "missing username or password, passed: " + string(payload), http.StatusBadRequest}
		}

		username := pair[0]
		password := pair[1]
		fromChache, cookies := cookieCache.TryGetValidCookie(username)
		if !fromChache {
			cks, err := ilpostapi.Login(username, password)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+err.Error()+`"`)
				return &appError{err, "can't decode username:password", http.StatusBadRequest}
			}
			cookies = cks
			cookieCache.Add(username, cookies)
		}

		// Store the username in the context to pass it to the next handler
		ctx := context.WithValue(r.Context(), userKey, cookies)

		return next(w, r.WithContext(ctx))
	}
}

func podcastListHandler(w http.ResponseWriter, _ *http.Request) *appError {
	ilpostResponse, err := ilpostapi.FetchPodcastList(nil)
	if err != nil {
		return &appError{err, "Can't fetch podcasts", http.StatusInternalServerError}
	}

	response := Podcasts{
		Items: []Podcast{},
	}

	for i, podcast := range ilpostResponse.Data {
		if podcast.Title == "" {
			continue
		}

		p := Podcast{
			ID:    i,
			Title: podcast.Title,
			URL:   podcast.URL,
			Feed:  "/feed?podcast-name=" + podcast.Slug,
		}

		response.Items = append(response.Items, p)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		return &appError{err, "Can't write response", http.StatusInternalServerError}
	}

	return nil
}

func feedHandler(w http.ResponseWriter, r *http.Request) *appError {
	podcastName := r.URL.Query().Get("podcast-name")
	if podcastName == "" {
		return &appError{nil, "missing 'podcast-name' parameter", http.StatusBadRequest}
	}

	// retrieve cookies from context
	cookies, ok := r.Context().Value(userKey).([]*http.Cookie)
	if !ok {
		return &appError{nil, "user not authenticated", http.StatusForbidden}
	}

	episodes, err := ilpostapi.FetchPodcastEpisodes(cookies, podcastName, 1, 20)
	if err != nil {
		return &appError{err, "failed to retrieve episodes", http.StatusBadGateway}
	}

	// create XML response
	xmlResponse := Convert_ilpost_to_RSS(episodes)

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return &appError{err, "failed to write XML header", http.StatusInternalServerError}
	}

	// Encode the response as XML
	if err := xml.NewEncoder(w).Encode(xmlResponse); err != nil {
		return &appError{err, "failed to encode XML data", http.StatusInternalServerError}
	}

	return nil
}
