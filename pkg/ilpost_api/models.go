package ilpostapi

import "time"

// === list ===
type Podcast struct {
	Author           string `json:"author"`
	Chronological    int    `json:"chronological"`
	Count            int    `json:"count"`
	Cyclicality      string `json:"cyclicality"`
	Description      string `json:"description"`
	Free             int    `json:"free"`
	Gift             int    `json:"gift"`
	GiftAll          int    `json:"gift_all"`
	ID               int    `json:"id"`
	Image            string `json:"image"`
	ImageWeb         string `json:"imageweb"`
	Object           string `json:"object"`
	Order            string `json:"order"`
	PushNotification int    `json:"pushnotification"`
	Robot            string `json:"robot"`
	Slug             string `json:"slug"`
	Sponsored        string `json:"sponsored"`
	Title            string `json:"title"`
	Type             string `json:"type"`
	URL              string `json:"url"`
}

type PodcastListResponse struct {
	Data []Podcast `json:"data"`
}

// === episodes ===
type PodcastEpisodesResponse struct {
	Head Head   `json:"head"`
	Data []Data `json:"data"`
}

type Head struct {
	ExecTime float64  `json:"exec_time"`
	Status   int      `json:"status"`
	Data     MetaData `json:"data"`
}

type MetaData struct {
	Total int `json:"total"`
	Pg    int `json:"pg"`
	Hits  int `json:"hits"`
}

type Data struct {
	ID            int       `json:"id"`
	Author        string    `json:"author"`
	Title         string    `json:"title"`
	Click         string    `json:"_click"`
	Summary       *string   `json:"summary"`
	ContentHTML   string    `json:"content_html"`
	Image         string    `json:"image"`
	ImageWeb      string    `json:"image_web"`
	Object        string    `json:"object"`
	Milliseconds  int       `json:"milliseconds"`
	Minutes       int       `json:"minutes"`
	Special       int       `json:"special"`
	ShareURL      string    `json:"share_url"`
	Slug          string    `json:"slug"`
	FullSlug      string    `json:"full_slug"`
	URL           string    `json:"url"`
	EpisodeRawURL string    `json:"episode_raw_url"`
	Meta          Meta      `json:"meta"`
	AccessLevel   string    `json:"access_level"`
	Timestamp     int64     `json:"timestamp"`
	Date          time.Time `json:"date"`
	DateString    *string   `json:"date_string"`
	Gift          bool      `json:"gift"`
	Parent        Parent    `json:"parent"`
	QueueList     *string   `json:"queue_list"`
}

type Parent struct {
	ID          int    `json:"id"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	ImageWeb    string `json:"image_web"`
	Object      string `json:"object"`
	Slug        string `json:"slug"`
	Meta        Meta   `json:"meta"`
	AccessLevel string `json:"access_level"`
}

type Meta struct {
	Order            int    `json:"order"`
	BackgroundColor  string `json:"background_color"`
	Gift             int    `json:"gift,omitempty"`
	GiftAll          int    `json:"gift_all,omitempty"`
	PushNotification int    `json:"pushnotification,omitempty"`
	Chronological    int    `json:"chronological,omitempty"`
	Robot            string `json:"robot,omitempty"`
	Sponsored        int    `json:"sponsored,omitempty"`
	Cyclicality      string `json:"cyclicality,omitempty"`
	Evidenza         string `json:"evidenza,omitempty"`
	CyclicalityType  string `json:"cyclicalitytype,omitempty"`
}
