package gofeed

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mmcdole/go-xpp"
	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/rss"
)

type FeedType int

const (
	FeedTypeUnknown FeedType = iota
	FeedTypeAtom
	FeedTypeRSS
)

type FeedParser struct {
	ATrans AtomTranslator
	RTrans RSSTranslator
	rp     *RSSParser
	ap     *AtomParser
}

func NewFeedParser() *FeedParser {
	fp := FeedParser{
		rp: &RSSParser{},
		ap: &AtomParser{},
	}
	return &fp
}

func (f *FeedParser) ParseFeedURL(feedURL string) (*Feed, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return f.ParseFeed(string(body))
}

func (f *FeedParser) ParseFeed(feed string) (*Feed, error) {
	ft := f.DetectFeedType(feed)
	switch ft {
	case FeedTypeAtom:
		return f.parseFeedFromAtom(feed)
	case FeedTypeRSS:
		return f.parseFeedFromRSS(feed)
	}
	return nil, errors.New("Failed to detect feed type")
}

func (f *FeedParser) DetectFeedType(feed string) FeedType {
	p := xpp.NewXMLPullParser(strings.NewReader(feed))

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

func (f *FeedParser) parseFeedFromAtom(feed string) (*Feed, error) {
	af, err := f.ap.ParseFeed(feed)
	if err != nil {
		return nil, err
	}
	result := f.atomTrans().Translate(af)
	return result, nil
}

func (f *FeedParser) parseFeedFromRSS(feed string) (*Feed, error) {
	rf, err := f.rp.ParseFeed(feed)
	if err != nil {
		return nil, err
	}

	result := f.rssTrans().Translate(rf)
	return result, nil
}

func (f *FeedParser) atomTrans() AtomTranslator {
	if f.ATrans != nil {
		return f.ATrans
	}
	f.ATrans = &DefaultAtomTranslator
}

func (f *FeedParser) rssTrans() RSSTranslator {
	if f.RTrans != nil {
		return f.RTrans
	}
	f.RTrans = &DefaultRSSTranslator
}
