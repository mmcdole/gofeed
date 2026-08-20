package gofeed

import (
	"bytes"
	"io"
	"strings"

	"github.com/mmcdole/gofeed/internal/shared"
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

// DetectFeedType attempts to determine the type of feed by looking for an XML
// root element or the start of a JSON object. It does not validate the entire
// document; the selected format parser performs full validation. It returns
// FeedTypeUnknown when the reader fails before the type can be determined.
func DetectFeedType(feed io.Reader) FeedType {
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(feed); err != nil {
		return FeedTypeUnknown
	}

	var firstChar byte
loop:
	for {
		ch, err := buffer.ReadByte()
		if err != nil {
			return FeedTypeUnknown
		}
		// ignore leading whitespace & byte order marks
		switch ch {
		case ' ', '\r', '\n', '\t':
		case 0xFE, 0xFF, 0x00, 0xEF, 0xBB, 0xBF: // utf 8-16-32 bom
		default:
			firstChar = ch
			buffer.UnreadByte()
			break loop
		}
	}

	if firstChar == '<' {
		// Check if it's an XML based feed
		p := shared.NewXMLParser(bytes.NewReader(buffer.Bytes()))

		_, err := shared.FindRoot(p)
		if err != nil {
			return FeedTypeUnknown
		}

		name := strings.ToLower(p.Name())
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
	} else if firstChar == '{' {
		return FeedTypeJSON
	}
	return FeedTypeUnknown
}
