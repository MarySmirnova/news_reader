package rss

import "time"

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Generator string `xml:"generator"`
	Link      string `xml:"link"`
	Items     []Item `xml:"item"`
}

type Item struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	PubDate     time.Time `xml:"pubDate"`
}
