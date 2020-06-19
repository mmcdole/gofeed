package gofeed

import (
	"bytes"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/mmcdole/gofeed/internal/shared"
	xpp "github.com/mmcdole/goxpp"
)

// FeedType represents one of the possible feed
// types that we can detect.
type FeedType int

const (
	// FeedTypeUnknown represents a feed that could not have its
	// type determiend.
	FeedTypeUnknown FeedType = iota
	// FeedTypeAtom repesents an Atom feed
	FeedTypeAtom
	// FeedTypeRSS represents an RSS feed
	FeedTypeRSS
	// FeedTypeJSON represents a JSON feed
	FeedTypeJSON
)

// DetectFeedType attempts to determine the type of feed
// by looking for specific xml elements unique to the
// various feed types.
func DetectFeedType(feed io.Reader) FeedType {

	// Check if document is valid JSON
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(feed)

	if jsoniter.Valid(buffer.Bytes()) {
		return FeedTypeJSON
	}

	// If not, check if it's an XML based feed
	p := xpp.NewXMLPullParser(bytes.NewReader(buffer.Bytes()), false, shared.NewReaderLabel)

	xmlBase := shared.XMLBase{}
	_, err := xmlBase.FindRoot(p)
	if err != nil {
		return FeedTypeUnknown
	}

	name := strings.ToLower(p.Name)
	switch name {
	case "rdf":
		return FeedTypeRSS
	case "rss":
		return FeedTypeRSS
	case "feed":
		return FeedTypeAtom
	default:
		return FeedTypeUnknown
	}
}
