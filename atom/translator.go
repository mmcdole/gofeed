package atom

import (
	"github.com/mmcdole/gofeed/feed"
)

// Translator converts an atom.Feed struct
// into the generic feed.Feed struct
type Translator interface {
	Translate(atom *Feed) *feed.Feed
}

// DefaultTranslator converts an atom.Feed struct
// into the generic feed.Feed struct.
//
// This default implementation defines a set of
// mapping rules between atom.Feed -> feed.Feed
// for each of the fields in feed.Feed.
type DefaultTranslator struct{}

func (t *DefaultTranslator) Translate(atom *Feed) *feed.Feed {
	feed := &feed.Feed{}
	feed.Title = atom.Title
	feed.Subtitle = atom.Subtitle
	feed.Rights = atom.Rights
	//feed.Links =
	if atom.Generator != nil {
		feed.Generator = atom.Generator.Value
	}

	return feed
}
