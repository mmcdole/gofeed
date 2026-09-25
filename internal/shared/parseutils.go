package shared

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	xpp "github.com/mmcdole/goxpp/v2"
)

const CDATA_START = "<![CDATA["
const CDATA_END = "]]>"

// FindRoot iterates through the tokens of an xml document until
// it encounters its first StartTag event.  It returns an error
// if it reaches EndDocument before finding a tag.
func FindRoot(p *xpp.Parser) (event xpp.EventType, err error) {
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
func ParseText(p *xpp.Parser) (string, error) {
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

	return DecodeEntities(result), nil
}

// ParseTextURL is ParseText for an element whose text is a URL: it resolves the
// value against the element's xml:base when one is in scope. The base is
// captured before ParseText runs, because ParseText consumes the element's end
// tag, which pops the base off the stack.
func ParseTextURL(p *xpp.Parser) (string, error) {
	base := p.BaseURL()
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
			dec := DecodeEntities(str[curr:])
			buf.Write([]byte(dec))
			return buf.String()
		}

		// Character data before the CDATA section is still character data:
		// entity-decode it and keep it. Dropping it silently loses any text
		// that precedes or sits between CDATA sections.
		dec := DecodeEntities(str[curr:start])
		buf.Write([]byte(dec))

		end := indexAt(str, CDATA_END, start+len(CDATA_START))

		if end == -1 {
			// Unterminated CDATA (malformed; encoding/xml would have rejected
			// it before this point). Keep the remainder, marker and all, rather
			// than guess where the section was meant to end.
			dec := DecodeEntities(str[start:])
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

// maxEntityLen bounds how far past an '&' the terminating ';' may be. The
// longest named HTML entity is 33 bytes with delimiters
// ("&CounterClockwiseContourIntegral;"); 64 leaves margin without letting a
// distant stray ';' make the scan quadratic on ampersand-heavy text.
const maxEntityLen = 64

// DecodeEntities decodes escaped XML entities
// in a string and returns the unescaped string
func DecodeEntities(str string) string {
	idx := strings.IndexByte(str, '&')
	if idx == -1 {
		return str
	}

	var buf bytes.Buffer
	data := []byte(str)

	for len(data) > 0 {
		idx := bytes.IndexByte(data, '&')
		if idx == -1 {
			buf.Write(data)
			break
		}

		buf.Write(data[:idx])
		data = data[idx:]

		// Find the end of the entity within the bounded window.
		window := data
		if len(window) > maxEntityLen {
			window = window[:maxEntityLen]
		}
		end := bytes.IndexByte(window, ';')

		if end == -1 && len(window) == len(data) {
			// No ';' in the rest of the input; no entity can terminate,
			// so everything left is literal text.
			buf.Write(data)
			break
		}

		// A candidate containing whitespace (or no nearby ';' at all) is not
		// an entity. Emit the '&' literally and keep scanning: entities later
		// in the string must still decode.
		if end == -1 || bytes.ContainsAny(data[1:end], " \t\r\n") {
			buf.WriteByte('&')
			data = data[1:]
			continue
		}

		// Accept the decode only when it consumed the terminating ';'.
		// html.UnescapeString applies HTML legacy rules that decode entity
		// prefixes without a ';' ("&copy=2;" becomes "©=2;"), which corrupts
		// URLs carrying such query parameters. If stripping the ';' from the
		// span decodes to the same text, the ';' was not part of an entity;
		// treat the '&' as literal and keep scanning.
		span := string(data[:end+1])
		dec := html.UnescapeString(span)
		if dec == html.UnescapeString(span[:len(span)-1])+";" {
			buf.WriteByte('&')
			data = data[1:]
			continue
		}

		buf.WriteString(dec)
		data = data[end+1:]
	}

	return buf.String()
}

// ParseNameAddress splits a free-form author string into a name and an email
// address. It recognizes "email (Name)", "Name (email)", "Name <email>",
// "<email>" and a bare email; any other text is returned whole as the name so
// no author is dropped.
func ParseNameAddress(nameAddressText string) (name string, address string) {
	text := strings.TrimSpace(nameAddressText)
	if text == "" {
		return
	}

	switch text[len(text)-1] {
	case '>':
		// "Name <email>" or "<email>"
		if i := strings.LastIndexByte(text, '<'); i >= 0 {
			if addr := strings.TrimSpace(text[i+1 : len(text)-1]); isEmailAddress(addr) {
				return unquoteName(strings.TrimSpace(text[:i])), addr
			}
		}
	case ')':
		// "email (Name)"; the name may itself contain parentheses.
		if i := strings.IndexByte(text, '('); i > 0 {
			if addr := strings.TrimSpace(text[:i]); isEmailAddress(addr) {
				return strings.TrimSpace(text[i+1 : len(text)-1]), addr
			}
		}
		// "Name (email)"
		if i := strings.LastIndexByte(text, '('); i > 0 {
			if addr := strings.TrimSpace(text[i+1 : len(text)-1]); isEmailAddress(addr) {
				return strings.TrimSpace(text[:i]), addr
			}
		}
	}

	if isEmailAddress(text) {
		return "", text
	}
	return text, ""
}

// isEmailAddress reports whether s looks like a single email address: one
// '@' with text on both sides and no whitespace or mailbox punctuation.
func isEmailAddress(s string) bool {
	at := strings.IndexByte(s, '@')
	if at <= 0 || at == len(s)-1 || strings.IndexByte(s[at+1:], '@') >= 0 {
		return false
	}
	return !strings.ContainsAny(s, " \t\r\n<>()\",")
}

// unquoteName strips the double quotes around an RFC 5322 display name.
func unquoteName(name string) string {
	if len(name) >= 2 && name[0] == '"' && name[len(name)-1] == '"' {
		return strings.TrimSpace(name[1 : len(name)-1])
	}
	return name
}

func indexAt(str, substr string, start int) int {
	idx := strings.Index(str[start:], substr)
	if idx > -1 {
		idx += start
	}
	return idx
}
