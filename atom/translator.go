package atom

import (
	"github.com/mmcdole/gofeed/feed"
)

// Translator converts an atom.Feed struct
// into the generic feed.Feed struct
type Translator interface {
	Translate(atom *Feed) *feed.Feed
}

// DefaulTranslator converts an atom.Feed struct
// into the generic feed.Feed struct.
//
// This default implementation defines a set of
// mapping rules between atom.Feed -> feed.Feed
// for each of the fields in feed.Feed.
type DefaultTranslator struct{}

func (t *DefaultTranslator) Translate(atom *Feed) *feed.Feed {
	return nil
}
