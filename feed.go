package gofeed

import (
	"time"

	"github.com/mmcdole/gofeed/extensions"
)

type Feed struct {
	Title           string            `json:"title,omitempty"`
	Description     string            `json:"description,omitempty"`
	Link            string            `json:"link,omitempty"`
	FeedLink        string            `json:"feedLink,omitempty"`
	Updated         string            `json:"updated,omitempty"`
	UpdatedParsed   *time.Time        `json:"updatedParsed,omitempty"`
	Published       string            `json:"published,omitempty"`
	PublishedParsed *time.Time        `json:"publishedParsed,omitempty"`
	Author          *Person           `json:"author,omitempty"`
	Language        string            `json:"language,omitempty"`
	Image           *Image            `json:"image,omitempty"`
	Copyright       string            `json:"copyright,omitempty"`
	Generator       string            `json:"generator,omitempty"`
	Categories      []string          `json:"categories,omitempty"`
	Extensions      ext.Extensions    `json:"extensions,omitempty"`
	Custom          map[string]string `json:"custom,omitempty"`
	Items           []*Item           `json:"items"`
	FeedType        string            `json:"feedType"`
	FeedVersion     string            `json:"feedVersion"`
}

type Item struct {
	Title           string            `json:"title,omitempty"`
	Description     string            `json:"description,omitempty"`
	Content         string            `json:"content,omitempty"`
	Link            string            `json:"link,omitempty"`
	Updated         string            `json:"updated,omitempty"`
	UpdatedParsed   *time.Time        `json:"updatedParsed,omitempty"`
	Published       string            `json:"published,omitempty"`
	PublishedParsed *time.Time        `json:"publishedParsed,omitempty"`
	Author          *Person           `json:"author,omitempty"`
	Guid            string            `json:"guid,omitempty"`
	Image           *Image            `json:"image,omitempty"`
	Categories      []string          `json:"categories,omitempty"`
	Enclosures      []*Enclosure      `json:"enclosures,omitempty"`
	Extensions      ext.Extensions    `json:"extensions,omitempty"`
	Custom          map[string]string `json:"custom,omitempty"`
}

type Person struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type Image struct {
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}

type Enclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}
