package gofeed

import (
	"fmt"
	"io"
)

func (f *Feed) RenderJSON(w io.Writer, converter JSONConverter) error {
	if converter == nil {
		converter = &DefaultJSONConverter{}
	}

	jsonFeed, err := converter.Convert(f)
	if err != nil {
		return fmt.Errorf("failed converting to JSON Feed: %v", err)
	}
	return jsonFeed.Render(w)
}

func (f *Feed) RenderAtom(w io.Writer, converter AtomConverter) error {
	if converter == nil {
		converter = &DefaultAtomConverter{}
	}

	atomFeed, err := converter.Convert(f)
	if err != nil {
		return fmt.Errorf("failed converting to Atom Feed: %v", err)
	}
	return atomFeed.Render(w)
}

func (f *Feed) RenderRSS(w io.Writer, converter RSSConverter) error {
	if converter == nil {
		converter = &DefaultRSSConverter{}
	}

	rssFeed, err := converter.Convert(f)
	if err != nil {
		return fmt.Errorf("failed converting to RSS Feed: %v", err)
	}
	return rssFeed.Render(w)
}
