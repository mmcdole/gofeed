package feed

import (
	"fmt"
	"time"
)

type RSSFeed struct {
	Title               string
	Link                string
	Description         string
	Language            string
	Copyright           string
	ManagingEditor      string
	WebMaster           string
	PubDate             string
	PubDateParsed       time.Time
	LastBuildDate       string
	LastBuildDateParsed time.Time
	Categories          []RSSCategory
	Generator           string
	Docs                string
	TTL                 string
	Image               RSSImage
	Rating              string
	SkipHours           []string
	SkipDays            []string
	Items               []*RSSItem
	Version             string
	Extensions          map[string]map[string][]Extension
}

func (f *RSSFeed) String() string {
	return fmt.Sprintf("Title: %s\nLink: %s\nDescription: %s\n"+
		"Language: %s\nCopyright: %s\nManagingEditor: %s\n"+
		"WebMaster: %s\nPubDate: %s\nLastBuildDate: %s\n"+
		"Generator: %s\nDocs: %s\nTTL: %s\n"+
		"Rating: %s\nItems: %s\nVersion: %s\n",
		f.Title, f.Link, f.Description,
		f.Language, f.Copyright, f.ManagingEditor,
		f.WebMaster, f.PubDate, f.LastBuildDate,
		f.Generator, f.Docs, f.TTL,
		f.Rating, f.Items, f.Version)
}

type RSSItem struct {
	Title         string
	Link          string
	Description   string
	Author        string
	Categories    []RSSCategory
	Comments      string
	Enclosure     RSSEnclosure
	Guid          RSSGuid
	PubDate       string
	PubDateParsed time.Time
	Source        RSSSource
	Extensions    map[string]map[string][]Extension
}

func (i *RSSItem) String() string {
	return fmt.Sprintf("Title: %s\nLink: %s\nDescription: %s\n"+
		"Author: %s\nComments: %s\nPubDate: %s\n"+
		"Source: %s\n",
		i.Title, i.Link, i.Description,
		i.Author, i.Comments, i.PubDate,
		i.Source)
}

type RSSImage struct {
	URL    string
	Link   string
	Title  string
	Width  string
	Height string
}

type RSSEnclosure struct {
	URL    string
	Length string
	Type   string
}

type RSSGuid struct {
	Value       string
	IsPermalink string
}

type RSSSource struct {
	Title string
	URL   string
}

type RSSCategory struct {
	Domain string
	Value  string
}
