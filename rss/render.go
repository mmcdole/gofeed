package rss

import (
	"encoding/xml"
	"io"
)

// rssRoot represents the RSS root element for XML marshaling
type rssRoot struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Feed    `xml:"channel"`
}

// Render writes the RSS feed as RSS 2.0 XML to the provided io.Writer
func (f *Feed) Render(w io.Writer) error {
	if f == nil {
		return nil
	}

	// Create the RSS root structure
	root := &rssRoot{
		Version: "2.0",
		Channel: f,
	}

	// Write XML declaration
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}

	// Create XML encoder and encode the feed
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	return encoder.Encode(root)
}
