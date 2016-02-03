package rss

import (
	"time"

	"github.com/mmcdole/gofeed/feed"
)

type Feed struct {
	Title               string              `json:"title"`
	Link                string              `json:"link"`
	Description         string              `json:"description"`
	Language            string              `json:"langauge,omitempty"`
	Copyright           string              `json:"copyright,omitempty"`
	ManagingEditor      string              `json:"managingEditor,omitempty"`
	WebMaster           string              `json:"webMaster,omitempty"`
	PubDate             string              `json:"pubDate,omitempty"`
	PubDateParsed       *time.Time          `json:"pubDateParsed,omitempty"`
	LastBuildDate       string              `json:"lastBuildDate,omitempty"`
	LastBuildDateParsed *time.Time          `json:"lastBuildDateParsed,omitempty"`
	Categories          []*Category         `json:"categories,omitempty"`
	Generator           string              `json:"generator,omitempty"`
	Docs                string              `json:"docs,omitempty"`
	TTL                 string              `json:"ttl,omitempty"`
	Image               *Image              `json:"image,omitempty"`
	Rating              string              `json:"rating,omitempty"`
	SkipHours           []string            `json:"skipHours,omitempty"`
	SkipDays            []string            `json:"skipDays,omitempty"`
	TextInput           *TextInput          `json:"textInput,omitempty"`
	Items               []*Item             `json:"items,omitempty"`
	Version             string              `json:"version,omitempty"`
	Extensions          feed.FeedExtensions `json:"extensions,omitempty"`
}

type Item struct {
	Title         string              `json:"title,omitempty"`
	Link          string              `json:"link,omitempty"`
	Description   string              `json:"description,omitempty"`
	Author        string              `json:"author,omitempty"`
	Categories    []*Category         `json:"categories,omitempty"`
	Comments      string              `json:"comments,omitempty"`
	Enclosure     *Enclosure          `json:"enclosure,omitempty"`
	Guid          *Guid               `json:"guid,omitempty"`
	PubDate       string              `json:"pubDate,omitempty"`
	PubDateParsed *time.Time          `json:"pubDateParsed,omitempty"`
	Source        *Source             `json:"source,omitempty"`
	Extensions    feed.FeedExtensions `json:"extensions,omitempty"`
}

type Image struct {
	URL    string `json:"url"`
	Link   string `json:"link"`
	Title  string `json:"title"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

type Enclosure struct {
	URL    string `json:"url"`
	Length string `json:"length"`
	Type   string `json:"type"`
}

type Guid struct {
	Value       string `json:"value"`
	IsPermalink string `json:"isPermalink"`
}

type Source struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Category struct {
	Domain string `json:"domain"`
	Value  string `json:"value"`
}

type TextInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Link        string `json:"link"`
}
