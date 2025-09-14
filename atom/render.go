package atom

import (
	"encoding/xml"
	"io"
)

// Render writes the Atom feed as Atom 1.0 XML to the provided io.Writer
func (f *Feed) Render(w io.Writer) error {
	if f == nil {
		return nil
	}

	// Write XML declaration
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}

	// Create XML encoder and encode the feed
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	return encoder.Encode(f)
}
