package gofeed

import (
	"errors"
	"io/ioutil"
	"net/http"
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
	if f.AtomTrans != nil {
		return f.AtomTrans
	}
	f.AtomTrans = &DefaultAtomTranslator{}
	return f.AtomTrans
}

func (f *FeedParser) rssTrans() RSSTranslator {
	if f.RSSTrans != nil {
		return f.RSSTrans
	}
	f.RSSTrans = &DefaultRSSTranslator{}
	return f.RSSTrans
}
