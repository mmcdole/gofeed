package shared

import (
	"strings"

	"github.com/mmcdole/goxpp"
)

// Parse the next text value from the document,
// trim any whitespace and return the value.
func ParseTrimText(p *xpp.XMLPullParser) (text string, err error) {
	text, err = p.NextText()
	if err != nil {
		return text, err
	}

	text = strings.TrimSpace(text)
	return text, nil
}

// IsExtension returns wether or not the current
// XML element is an extension element (if it has an
// empty prefix)
func IsExtension(p *xpp.XMLPullParser) bool {
	space := strings.TrimSpace(p.Space)
	if prefix, ok := p.Spaces[space]; ok {
		return prefix != ""
	}

	return p.Space != ""
}
