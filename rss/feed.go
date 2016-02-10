package rss

import (
	"time"

	"github.com/mmcdole/gofeed/feed"
)

type Feed struct {
	Title               string              `json:"title,omitempty"`
	Link                string              `json:"link,omitempty"`
	Description         string              `json:"description,omitempty"`
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
	Cloud               *Cloud              `json:"cloud,omitempty"`
	TextInput           *TextInput          `json:"textInput,omitempty"`
	Items               []*Item             `json:"items"`
	Extensions          feed.FeedExtensions `json:"extensions"`
	Version             string              `json:"version"`
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
	Extensions    feed.FeedExtensions `json:"extensions"`
}

type Image struct {
	URL         string `json:"url,omitempty"`
	Link        string `json:"link,omitempty"`
	Title       string `json:"title,omitempty"`
	Width       string `json:"width,omitempty"`
	Height      string `json:"height,omitempty"`
	Description string `json:"description,omitempty"`
}

type Enclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}

type Guid struct {
	Value       string `json:"value,omitempty"`
	IsPermalink string `json:"isPermalink,omitempty"`
}

type Source struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type Category struct {
	Domain string `json:"domain,omitempty"`
	Value  string `json:"value,omitempty"`
}

type TextInput struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Link        string `json:"link,omitempty"`
}

type Cloud struct {
	Domain            string `json:"domain,omitempty"`
	Port              string `json:"port,omitempty"`
	Path              string `json:"path,omitempty"`
	RegisterProcedure string `json:"registerProcedure,omitempty"`
	Protocol          string `json:"protocol,omitempty"`
}
