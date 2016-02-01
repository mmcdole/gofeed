package feed

// FeedNormalizer is responsible for taking either
// an RSS or Atom feed and outputting a normalized
// Feed.
//
// If multiple fields from Atom/RSS or their extensions
// map to a single field in the `Feed` object, then the
// normalizer decides the precedence.
type FeedNormalizer interface {

	// NormalizeAtom maps an AtomFeed to a Feed type.
	NormalizeAtom(feed *AtomFeed) *Feed

	// NormalizeRSS maps an RSSFeed to a Feed type.
	NormalizeRSS(feed *RSSFeed) *Feed
}

// DefaultNormalizer is the default normalizer
// that maps RSS/Atom feeds to the `Feed` object.
type DefaultNormalizer struct{}

func (n *DefaultNormalizer) NormalizeRSS(feed *RSSFeed) *Feed {
	return &Feed{}
}

func (n *DefaultNormalizer) NormalizeAtom(feed *AtomFeed) *Feed {
	return &Feed{}
}
