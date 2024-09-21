package endpoint_test

import (
	"encoding/json"
	"ilpost-podcast-feed/pkg/endpoint"
	ilpostapi "ilpost-podcast-feed/pkg/ilpost_api"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestBuildFeed(t *testing.T) {

	tests := []struct {
		name     string
		episodes ilpostapi.PodcastEpisodesResponse
		want     endpoint.RSS
	}{
		{
			name:     "Test with valid podcast name",
			episodes: Unmarshal("bordone.json"),
			want: endpoint.RSS{
				Version: "1.0",
			},
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := endpoint.BuildFeed(tt.episodes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildFeed(...) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Unmarshal(filename string) ilpostapi.PodcastEpisodesResponse {
	file, _ := os.Open(filename)
	defer file.Close()

	filedata, _ := io.ReadAll(file)

	var data ilpostapi.PodcastEpisodesResponse
	json.Unmarshal(filedata, &data)

	return data
}
