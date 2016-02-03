package rss

import (
	"github.com/mmcdole/gofeed/feed"
)

// Translator converts an rss.Feed struct
// into the generic feed.Feed struct
type Translator interface {
	Translate(rss *Feed) *feed.Feed
}

// DefaulTranslator converts an rss.Feed struct
// into the generic feed.Feed struct.
//
// This default implementation defines a set of
// mapping rules between rss.Feed -> feed.Feed
// for each of the fields in feed.Feed.
type DefaultTranslator struct{}

func (t *DefaultTranslator) Translate(rss *Feed) *feed.Feed {
	return nil
}
