package shared

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"

	xpp "github.com/mmcdole/goxpp"
)

var (
	emailNameRgx = regexp.MustCompile(`^([^@]+@[^\s]+)\s+\(([^@]+)\)$`)
	nameEmailRgx = regexp.MustCompile(`^([^@]+)\s+\(([^@]+@[^)]+)\)$`)
	nameOnlyRgx  = regexp.MustCompile(`^([^@()]+)$`)
	emailOnlyRgx = regexp.MustCompile(`^([^@()]+@[^@()]+)$`)

	ErrTruncatedEntity         = errors.New("truncated entity")
	ErrInvalidNumericReference = errors.New("invalid numeric reference")
)

const CDATA_START = "<![CDATA["
const CDATA_END = "]]>"

// FindRoot iterates through the tokens of an xml document until
// it encounters its first StartTag event.  It returns an error
// if it reaches EndDocument before finding a tag.
func FindRoot(p *xpp.XMLPullParser) (event xpp.XMLEventType, err error) {
	for {
		event, err = p.Next()
		if err != nil {
			return event, err
		}
		if event == xpp.StartTag {
			break
		}

		if event == xpp.EndDocument {
			return event, fmt.Errorf("failed to find root node before document end")
		}
	}
	return
}

// ParseText is a helper function for parsing the text
// from the current element of the XMLPullParser.
// This function can handle parsing naked XML text from
// an element.
func ParseText(p *xpp.XMLPullParser) (string, error) {
	var text struct {
		Type     string `xml:"type,attr"`
		InnerXML string `xml:",innerxml"`
	}

	err := p.DecodeElement(&text)
	if err != nil {
		return "", err
	}

	result := text.InnerXML
	result = strings.TrimSpace(result)

	if strings.Contains(result, CDATA_START) {
		return StripCDATA(result), nil
	}

	return DecodeEntities(result)
}

// ParseTextURL is ParseText for an element whose text is a URL: it resolves the
// value against the element's xml:base when one is in scope. The base is
// captured before ParseText runs, because ParseText consumes the element's end
// tag, which pops the base off the stack.
func ParseTextURL(p *xpp.XMLPullParser) (string, error) {
	base := p.BaseStack.Top()
	s, err := ParseText(p)
	if err != nil {
		return "", err
	}
	return ResolveURLIfBase(base, s), nil
}

// StripCDATA removes CDATA tags from the string
// content outside of CDATA tags is passed via DecodeEntities
func StripCDATA(str string) string {
	buf := bytes.NewBuffer([]byte{})

	curr := 0

	for curr < len(str) {

		start := indexAt(str, CDATA_START, curr)

		if start == -1 {
			dec, _ := DecodeEntities(str[curr:])
			buf.Write([]byte(dec))
			return buf.String()
		}

		// Character data before the CDATA section is still character data:
		// entity-decode it and keep it. Dropping it silently loses any text
		// that precedes or sits between CDATA sections.
		dec, _ := DecodeEntities(str[curr:start])
		buf.Write([]byte(dec))

		end := indexAt(str, CDATA_END, start+len(CDATA_START))

		if end == -1 {
			// Unterminated CDATA (malformed; encoding/xml would have rejected
			// it before this point). Keep the remainder, marker and all, rather
			// than guess where the section was meant to end.
			dec, _ := DecodeEntities(str[start:])
			buf.Write([]byte(dec))
			return buf.String()
		}

		// CDATA content is taken verbatim (no entity decoding).
		buf.Write([]byte(str[start+len(CDATA_START) : end]))

		// end is an absolute index; advance past the closing ]]> without
		// re-adding curr (doing so overshoots and drops trailing content).
		curr = end + len(CDATA_END)
	}

	return buf.String()
}

// DecodeEntities decodes escaped XML entities
// in a string and returns the unescaped string
func DecodeEntities(str string) (string, error) {
	data := []byte(str)
	buf := bytes.NewBuffer([]byte{})

	for len(data) > 0 {
		// Find the next entity
		idx := bytes.IndexByte(data, '&')
		if idx == -1 {
			buf.Write(data)
			break
		}

		buf.Write(data[:idx])
		data = data[idx:]

		// If there is only the '&' left here
		if len(data) == 1 {
			buf.Write(data)
			return buf.String(), nil
		}

		// Find the end of the entity
		end := bytes.IndexByte(data, ';')
		if end == -1 {
			// it's not an entitiy. just a plain old '&' possibly with extra bytes
			buf.Write(data)
			return buf.String(), nil
		}

		// Check if there is a space somewhere within the 'entitiy'.
		// If there is then skip the whole thing since it's not a real entity.
		if strings.Contains(string(data[1:end]), " ") {
			buf.Write(data)
			return buf.String(), nil
		} else {
			buf.WriteString(html.UnescapeString(string(data[0 : end+1])))
		}

		// Skip the entity
		data = data[end+1:]
	}

	return buf.String(), nil
}

// ParseNameAddress parses name/email strings commonly
// found in RSS feeds of the format "Example Name (example@site.com)"
// and other variations of this format.
func ParseNameAddress(nameAddressText string) (name string, address string) {
	if nameAddressText == "" {
		return
	}

	if emailNameRgx.MatchString(nameAddressText) {
		result := emailNameRgx.FindStringSubmatch(nameAddressText)
		address = result[1]
		name = result[2]
	} else if nameEmailRgx.MatchString(nameAddressText) {
		result := nameEmailRgx.FindStringSubmatch(nameAddressText)
		name = result[1]
		address = result[2]
	} else if nameOnlyRgx.MatchString(nameAddressText) {
		result := nameOnlyRgx.FindStringSubmatch(nameAddressText)
		name = result[1]
	} else if emailOnlyRgx.MatchString(nameAddressText) {
		result := emailOnlyRgx.FindStringSubmatch(nameAddressText)
		address = result[1]
	}
	return
}

func indexAt(str, substr string, start int) int {
	idx := strings.Index(str[start:], substr)
	if idx > -1 {
		idx += start
	}
	return idx
}
