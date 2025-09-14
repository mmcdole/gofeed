package rss

import (
	"encoding/json"
	"time"

	ext "github.com/mmcdole/gofeed/extensions"
)

// Feed is an RSS Feed
type Feed struct {
	Title               string                   `json:"title,omitempty" xml:"title"`
	Link                string                   `json:"link,omitempty" xml:"link"`
	Links               []string                 `json:"links,omitempty" xml:"-"`
	Description         string                   `json:"description,omitempty" xml:"description"`
	Language            string                   `json:"language,omitempty" xml:"language,omitempty"`
	Copyright           string                   `json:"copyright,omitempty" xml:"copyright,omitempty"`
	ManagingEditor      string                   `json:"managingEditor,omitempty" xml:"managingEditor,omitempty"`
	WebMaster           string                   `json:"webMaster,omitempty" xml:"webMaster,omitempty"`
	PubDate             string                   `json:"pubDate,omitempty" xml:"pubDate,omitempty"`
	PubDateParsed       *time.Time               `json:"pubDateParsed,omitempty" xml:"-"`
	LastBuildDate       string                   `json:"lastBuildDate,omitempty" xml:"lastBuildDate,omitempty"`
	LastBuildDateParsed *time.Time               `json:"lastBuildDateParsed,omitempty" xml:"-"`
	Categories          []*Category              `json:"categories,omitempty" xml:"category,omitempty"`
	Generator           string                   `json:"generator,omitempty" xml:"generator,omitempty"`
	Docs                string                   `json:"docs,omitempty" xml:"docs,omitempty"`
	TTL                 string                   `json:"ttl,omitempty" xml:"ttl,omitempty"`
	Image               *Image                   `json:"image,omitempty" xml:"image"`
	Rating              string                   `json:"rating,omitempty" xml:"rating,omitempty"`
	SkipHours           []string                 `json:"skipHours,omitempty" xml:"skipHours>hour"`
	SkipDays            []string                 `json:"skipDays,omitempty" xml:"skipDays>day"`
	Cloud               *Cloud                   `json:"cloud,omitempty" xml:"cloud"`
	TextInput           *TextInput               `json:"textInput,omitempty" xml:"textInput"`
	DublinCoreExt       *ext.DublinCoreExtension `json:"dcExt,omitempty" xml:"-"`
	ITunesExt           *ext.ITunesFeedExtension `json:"itunesExt,omitempty" xml:"-"`
	Extensions          ext.Extensions           `json:"extensions,omitempty" xml:"-"`
	Items               []*Item                  `json:"items" xml:"item"`
	Version             string                   `json:"version" xml:"-"`
}

func (f Feed) String() string {
	json, _ := json.MarshalIndent(f, "", "    ")
	return string(json)
}

// Item is an RSS Item
type Item struct {
	Title         string                   `json:"title,omitempty" xml:"title"`
	Link          string                   `json:"link,omitempty" xml:"link"`
	Links         []string                 `json:"links,omitempty" xml:"-"`
	Description   string                   `json:"description,omitempty" xml:"description"`
	Content       string                   `json:"content,omitempty" xml:"-"`
	Author        string                   `json:"author,omitempty" xml:"author"`
	Categories    []*Category              `json:"categories,omitempty" xml:"category"`
	Comments      string                   `json:"comments,omitempty" xml:"comments"`
	Enclosure     *Enclosure               `json:"enclosure,omitempty" xml:"enclosure"`
	Enclosures    []*Enclosure             `json:"enclosures,omitempty" xml:"-"`
	GUID          *GUID                    `json:"guid,omitempty" xml:"guid"`
	PubDate       string                   `json:"pubDate,omitempty" xml:"pubDate"`
	PubDateParsed *time.Time               `json:"pubDateParsed,omitempty" xml:"-"`
	Source        *Source                  `json:"source,omitempty" xml:"source"`
	DublinCoreExt *ext.DublinCoreExtension `json:"dcExt,omitempty" xml:"-"`
	ITunesExt     *ext.ITunesItemExtension `json:"itunesExt,omitempty" xml:"-"`
	Extensions    ext.Extensions           `json:"extensions,omitempty" xml:"-"`
	Custom        map[string]string        `json:"custom,omitempty" xml:"-"`
}

// Image is an image that represents the feed
type Image struct {
	URL         string `json:"url,omitempty" xml:"url"`
	Link        string `json:"link,omitempty" xml:"link"`
	Title       string `json:"title,omitempty" xml:"title"`
	Width       string `json:"width,omitempty" xml:"width"`
	Height      string `json:"height,omitempty" xml:"height"`
	Description string `json:"description,omitempty" xml:"description"`
}

// Enclosure is a media object that is attached to
// the item
type Enclosure struct {
	URL    string `json:"url,omitempty" xml:"url,attr"`
	Length string `json:"length,omitempty" xml:"length,attr"`
	Type   string `json:"type,omitempty" xml:"type,attr"`
}

// GUID is a unique identifier for an item
type GUID struct {
	Value       string `json:"value,omitempty" xml:",chardata"`
	IsPermalink string `json:"isPermalink,omitempty" xml:"isPermaLink,attr"`
}

// Source contains feed information for another
// feed if a given item came from that feed
type Source struct {
	Title string `json:"title,omitempty" xml:",chardata"`
	URL   string `json:"url,omitempty" xml:"url,attr"`
}

// Category is category metadata for Feeds and Entries
type Category struct {
	Domain string `json:"domain,omitempty" xml:"domain,attr,omitempty"`
	Value  string `json:"value,omitempty" xml:",chardata"`
}

// TextInput specifies a text input box that
// can be displayed with the channel
type TextInput struct {
	Title       string `json:"title,omitempty" xml:"title"`
	Description string `json:"description,omitempty" xml:"description"`
	Name        string `json:"name,omitempty" xml:"name"`
	Link        string `json:"link,omitempty" xml:"link"`
}

// Cloud allows processes to register with a
// cloud to be notified of updates to the channel,
// implementing a lightweight publish-subscribe protocol
// for RSS feeds
type Cloud struct {
	Domain            string `json:"domain,omitempty" xml:"domain,attr"`
	Port              string `json:"port,omitempty" xml:"port,attr"`
	Path              string `json:"path,omitempty" xml:"path,attr"`
	RegisterProcedure string `json:"registerProcedure,omitempty" xml:"registerProcedure,attr"`
	Protocol          string `json:"protocol,omitempty" xml:"protocol,attr"`
}
