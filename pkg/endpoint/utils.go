package endpoint

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"

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

func testHandler(w http.ResponseWriter, r *http.Request) *appError {
	ilpostData, err := unmarshalFromFile("./pkg/endpoint/bordone.json")
	if err != nil {
		return &appError{err, "Can't unmarshal test json", http.StatusInternalServerError}
	}

	episodes := Convert_ilpost_to_RSS(*ilpostData)

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
