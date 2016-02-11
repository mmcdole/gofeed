package gofeed

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/feed"
	"github.com/mmcdole/gofeed/rss"
	"github.com/mmcdole/goxpp"
)

type FeedType int

const (
	FeedTypeUnknown FeedType = iota
	FeedTypeAtom
	FeedTypeRSS
)

type FeedParser struct {
	AtomTrans atom.Translator
	RSSTrans  rss.Translator
	rp        *rss.Parser
	ap        *atom.Parser
}

func NewFeedParser() *FeedParser {
	fp := FeedParser{
		rp: &rss.Parser{},
		ap: &atom.Parser{},
	}
	return &fp
}

func (f *FeedParser) ParseFeedURL(feedURL string) (*feed.Feed, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return f.ParseFeed(string(body))
}

func (f *FeedParser) ParseFeed(feed string) (*feed.Feed, error) {
	fmt.Println(feed)
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
	p := xpp.NewXMLPullParser(strings.NewReader(feed), false)

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

func (f *FeedParser) parseFeedFromAtom(feed string) (*feed.Feed, error) {
	af, err := f.ap.ParseFeed(feed)
	if err != nil {
		return nil, err
	}
	result := f.atomTrans().Translate(af)
	return result, nil
}

func (f *FeedParser) parseFeedFromRSS(feed string) (*feed.Feed, error) {
	rf, err := f.rp.ParseFeed(feed)
	if err != nil {
		return nil, err
	}

	result := f.rssTrans().Translate(rf)
	return result, nil
}

func (f *FeedParser) atomTrans() atom.Translator {
	if f.AtomTrans != nil {
		return f.AtomTrans
	}
	f.AtomTrans = &atom.DefaultTranslator{}
	return f.AtomTrans
}

func (f *FeedParser) rssTrans() rss.Translator {
	if f.RSSTrans != nil {
		return f.RSSTrans
	}
	f.RSSTrans = &rss.DefaultTranslator{}
	return f.RSSTrans
}
