package atom

import (
	"github.com/mmcdole/gofeed/feed"
)

type Translator interface {
	Translate(atom *Feed) *feed.Feed
}

type DefaultTranslator struct{}

func (t *DefaultTranslator) Translate(atom *Feed) *feed.Feed {
	return nil
}
