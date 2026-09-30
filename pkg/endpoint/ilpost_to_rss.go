package endpoint

import (
	ilpostapi "ilpost-podcast-feed/pkg/ilpost_api"
	"strconv"
)

func Convert_ilpost_to_RSS(episodes ilpostapi.PodcastEpisodesResponse) RSS {
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
