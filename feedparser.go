package gofeed

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mmcdole/go-xpp"
	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/rss"
	"github.com/mmcdole/gofeed/util"
)

type FeedParser struct {
	AtomTrans AtomTranslator
	RSSTrans  RSSTranslator
	rp        *RSSParser
	ap        *AtomParser
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
	ft := DetectFeedType(feed)
	switch ft {
	case FeedTypeAtom:
		return f.parseFeedFromAtom(feed)
	case FeedTypeRSS:
		return f.parseFeedFromRSS(feed)
	}
	return nil, errors.New("Failed to detect feed type")
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
	f.AtomTrans = &DefaultAtomTranslator
}

func (f *FeedParser) rssTrans() RSSTranslator {
	if f.RTrans != nil {
		return f.RTrans
	}
	f.RSSTrans = &DefaultRSSTranslator
}
