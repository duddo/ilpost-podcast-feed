package ilpostapi

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
