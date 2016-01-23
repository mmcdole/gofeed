package feed

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type FeedParser struct {
	// Normalizer is the object that is responsible
	// for transforming RSSFeed or AtomFeed into the
	// common Feed object.
	//
	// Defaults to `DefaultNormalizer` which tries to make
	// a best effort at mapping between the two formats.
	// Replace this FeedNormalizer if you need to define custom
	// mappings.
	Normalizer FeedNormalizer
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
	af, err := ParseAtomFeed(feed)
	if err != nil {
		return nil, err
	}
	nf := f.normalizer().NormalizeAtom(af)
	return nf, nil
}

func (f *FeedParser) parseFeedFromRSS(feed string) (*Feed, error) {
	rf, err := ParseRSSFeed(feed)
	if err != nil {
		return nil, err
	}

	nf := f.Normalizer.NormalizeRSS(rf)
	return nf, nil
}

func (f *FeedParser) normalizer() FeedNormalizer {
	if f.Normalizer == nil {
		f.Normalizer = &DefaultNormalizer{}
	}
	return f.Normalizer
}
