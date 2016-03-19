package feed

import (
	"time"

	"github.com/mmcdole/gofeed/extensions"
)

type Feed struct {
	Title           string     `json:"title"`
	Subtitle        string     `json:"subtitle"`
	Published       string     `json:"published"`
	PublishedParsed *time.Time `json:"publishedParsed"`
	Rights          string     `json:"rights"`
	//	Image               *Image                   `json:"image"`
	Items               []Item                   `json:"items"`
	Generator           string                   `json:"generator"`
	ITunesExtension     *ext.ITunesFeedExtension `json:"itunesExt,omitempty"`
	DublinCoreExtension *ext.DublinCoreExtension `json:"dcExt,omitempty"`
	Extensions          ext.Extensions           `json:"extensions"`
	Custom              map[string]string        `json:"custom"`
	FeedType            string                   `json:"feedType"`
	FeedVersion         string                   `json:"feedVersion"`
}

type Item struct {
	Title               string                    `json:"title"`
	Description         string                    `json:"description"`
	ITunesExtension     *ext.ITunesEntryExtension `json:"itunesExt,omitempty"`
	DublinCoreExtension *ext.DublinCoreExtension  `json:"dcExt,omitempty"`
	Extensions          ext.Extensions            `json:"extensions"`
	Custom              map[string]string         `json:"custom"`
}
