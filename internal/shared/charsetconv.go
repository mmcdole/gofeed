package shared

import (
	"encoding/xml"
	"io"

	xpp "github.com/mmcdole/goxpp/v2"
	"golang.org/x/net/html/charset"
)

func NewReaderLabel(label string, input io.Reader) (io.Reader, error) {
	conv, err := charset.NewReaderLabel(label, input)

	if err != nil {
		return nil, err
	}

	return conv, nil
}

// NewXMLParser returns a pull parser configured the way every gofeed parser
// needs it: non-strict, so real-world feeds with unescaped entities and other
// common mistakes still tokenize, and with charset conversion for feeds that
// declare a non-UTF-8 encoding.
func NewXMLParser(r io.Reader) *xpp.Parser {
	d := xml.NewDecoder(r)
	d.Strict = false
	d.CharsetReader = NewReaderLabel
	return xpp.New(d)
}
