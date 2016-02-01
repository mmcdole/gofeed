package rss

import (
	"time"
)

type RSSFeed struct {
	Title               string         `json:"title"`
	Link                string         `json:"link"`
	Description         string         `json:"description"`
	Language            string         `json:"langauge,omitempty"`
	Copyright           string         `json:"copyright,omitempty"`
	ManagingEditor      string         `json:"managingEditor,omitempty"`
	WebMaster           string         `json:"webMaster,omitempty"`
	PubDate             string         `json:"pubDate,omitempty"`
	PubDateParsed       *time.Time     `json:"pubDateParsed,omitempty"`
	LastBuildDate       string         `json:"lastBuildDate,omitempty"`
	LastBuildDateParsed *time.Time     `json:"lastBuildDateParsed,omitempty"`
	Categories          []*RSSCategory `json:"categories,omitempty"`
	Generator           string         `json:"generator,omitempty"`
	Docs                string         `json:"docs,omitempty"`
	TTL                 string         `json:"ttl,omitempty"`
	Image               *RSSImage      `json:"image,omitempty"`
	Rating              string         `json:"rating,omitempty"`
	SkipHours           []string       `json:"skipHours,omitempty"`
	SkipDays            []string       `json:"skipDays,omitempty"`
	TextInput           *RSSTextInput  `json:"textInput,omitempty"`
	Items               []*RSSItem     `json:"items,omitempty"`
	Version             string         `json:"version,omitempty"`
	Extensions          FeedExtensions `json:"extensions,omitempty"`
}

type RSSItem struct {
	Title         string         `json:"title,omitempty"`
	Link          string         `json:"link,omitempty"`
	Description   string         `json:"description,omitempty"`
	Author        string         `json:"author,omitempty"`
	Categories    []*RSSCategory `json:"categories,omitempty"`
	Comments      string         `json:"comments,omitempty"`
	Enclosure     *RSSEnclosure  `json:"enclosure,omitempty"`
	Guid          *RSSGuid       `json:"guid,omitempty"`
	PubDate       string         `json:"pubDate,omitempty"`
	PubDateParsed *time.Time     `json:"pubDateParsed,omitempty"`
	Source        *RSSSource     `json:"source,omitempty"`
	Extensions    FeedExtensions `json:"extensions,omitempty"`
}

type RSSImage struct {
	URL    string `json:"url"`
	Link   string `json:"link"`
	Title  string `json:"title"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

type RSSEnclosure struct {
	URL    string `json:"url"`
	Length string `json:"length"`
	Type   string `json:"type"`
}

type RSSGuid struct {
	Value       string `json:"value"`
	IsPermalink string `json:"isPermalink"`
}

type RSSSource struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type RSSCategory struct {
	Domain string `json:"domain"`
	Value  string `json:"value"`
}

type RSSTextInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Link        string `json:"link"`
}
