package rss

import (
	"github.com/mmcdole/gofeed/feed"
)

type Translator interface {
	Translate(rss *Feed) *feed.Feed
}

type DefaultTranslator struct{}

func (t *DefaultTranslator) Translate(rss *Feed) *feed.Feed {
	return nil
}
