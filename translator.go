package gofeed

type AtomTranslator interface {
	Translate(feed *AtomFeed) *Feed
}

type RSSTranslator interface {
	Translate(feed *RSSFeed) *Feed
}

type DefaultAtomTranslator struct{}

func (t *DefaultAtomTranslator) Translate(feed *AtomFeed) *Feed {
	return nil
}

type DefaultRSSTranslator struct{}

func (t *DefaultRSSTranslator) Translate(feed *RSSFeed) *Feed {
	return nil
}
