package gofeed

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/rss"
	"github.com/mmcdole/goxpp"
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
)

// DetectFeedType takes a feed reader and attempts
// to detect its feed type.
func DetectFeedType(feed io.Reader) FeedType {
	p := xpp.NewXMLPullParser(feed, false)

	_, err := p.NextTag()
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

// FeedParser is a universal feed parser that detects
// a given feed type, parsers it, and translates it
// to the universal feed type.
type FeedParser struct {
	AtomTranslator Translator
	RSSTranslator  Translator
	Client         *http.Client
	rp             *rss.Parser
	ap             *atom.Parser
}

// NewFeedParser creates a FeedParser.
func NewFeedParser() *FeedParser {
	fp := FeedParser{
		rp: &rss.Parser{},
		ap: &atom.Parser{},
	}
	return &fp
}

func (f *FeedParser) ParseFeed(feed io.Reader) (*Feed, error) {
	// Wrap the feed io.Reader in a io.TeeReader
	// so we can capture all the bytes read by the
	// DetectFeedType function and construct a new
	// reader with those bytes intact for when we
	// attempt to parse the feeds.
	var buf bytes.Buffer
	tee := io.TeeReader(feed, &buf)
	feedType := DetectFeedType(tee)

	// Glue the read bytes from the detect function
	// back into a new reader
	r := io.MultiReader(&buf, feed)

	switch feedType {
	case FeedTypeAtom:
		return f.parseAtomFeed(r)
	case FeedTypeRSS:
		return f.parseRSSFeed(r)
	}
	return nil, errors.New("Failed to detect feed type")
}

// ParseFeedURL fetches the contents of a given feed url and
// parses the feed into the universal feed type.
func (f *FeedParser) ParseFeedURL(feedURL string) (*Feed, error) {
	client := f.httpClient()
	resp, err := client.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return f.ParseFeed(resp.Body)
}

// ParseFeed takes a feed XML string and parses it into the
// universal feed type.
func (f *FeedParser) ParseFeedString(feed string) (*Feed, error) {
	return f.ParseFeed(strings.NewReader(feed))
}

func (f *FeedParser) parseAtomFeed(feed io.Reader) (*Feed, error) {
	af, err := f.ap.ParseFeed(feed)
	if err != nil {
		return nil, err
	}
	return f.atomTrans().Translate(af)
}

func (f *FeedParser) parseRSSFeed(feed io.Reader) (*Feed, error) {
	rf, err := f.rp.ParseFeed(feed)
	if err != nil {
		return nil, err
	}

	return f.rssTrans().Translate(rf)
}

func (f *FeedParser) atomTrans() Translator {
	if f.AtomTranslator != nil {
		return f.AtomTranslator
	}
	f.AtomTranslator = &DefaultAtomTranslator{}
	return f.AtomTranslator
}

func (f *FeedParser) rssTrans() Translator {
	if f.RSSTranslator != nil {
		return f.RSSTranslator
	}
	f.RSSTranslator = &DefaultRSSTranslator{}
	return f.RSSTranslator
}

func (f *FeedParser) httpClient() *http.Client {
	if f.Client != nil {
		return f.Client
	}
	f.Client = &http.Client{}
	return f.Client
}
