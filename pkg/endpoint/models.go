package endpoint

import (
	"encoding/xml"
)

// === RSS Podcast feed structures ===

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Content string   `xml:"xmlns:content,attr"`
	Itunes  string   `xml:"xmlns:itunes,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description CDATA  `xml:"description"`
	Language    string `xml:"language"`
	Generator   string `xml:"generator"`
	Copyright   string `xml:"copyright,omitempty"`

	Subtitle CDATA  `xml:"itunes:subtitle"`
	Summary  CDATA  `xml:"itunes:summary"`
	Author   string `xml:"itunes:author"`
	Block    string `xml:"itunes:block"`
	Explicit string `xml:"itunes:explicit"`
	Image    Image  `xml:"itunes:image"`

	Items []Item `xml:"item"`
}

type Image struct {
	Href string `xml:"href,attr"`
}

type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	PubDate     string    `xml:"pubDate"`
	Guid        string    `xml:"guid"`
	Description string    `xml:"description"`
	Enclosure   Enclosure `xml:"enclosure"`

	ContentEncoded CDATA  `xml:"content:encoded"`
	Duration       string `xml:"itunes:duration"`
	Subtitle       string `xml:"itunes:subtitle"`
	Summary        string `xml:"itunes:summary"`
	Keywords       string `xml:"itunes:keywords"`
	Author         string `xml:"itunes:author"`
	Explicit       string `xml:"itunes:explicit"`
	Block          string `xml:"itunes:block"`
}

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length string `xml:"length,attr,omitempty"`
	Type   string `xml:"type,attr"`
}
