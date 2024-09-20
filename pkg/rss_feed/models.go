package rssfeed

type RSS struct {
	//XMLName xml.Name `xml:"rss"`
	Version string `xml:"version,attr"`
	//Channel Channel  `xml:"channel"`
}

// Channel represents the <channel> element in the RSS feed.
type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	PubDate       string `xml:"pubDate"`
	LastBuildDate string `xml:"lastBuildDate"`
	TTL           string `xml:"ttl"`
	//Items       []Item  `xml:"item"`
}
