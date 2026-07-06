package atom

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"
	"time"

	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/mmcdole/gofeed/internal/shared"
	xpp "github.com/mmcdole/goxpp/v2"
	"golang.org/x/net/html"
)

var (
	// Atom elements which contain URIs
	// https://tools.ietf.org/html/rfc4287
	atomUriElements = map[string]bool{
		"icon": true,
		"id":   true,
		"logo": true,
		"uri":  true,
		"url":  true, // atom 0.3
	}
)

// Parser is an Atom Parser
type Parser struct{}

// Parse parses an xml feed into an atom.Feed
func (ap *Parser) Parse(feed io.Reader) (*Feed, error) {
	feed = shared.NewControlCharFilterReader(feed)
	p := shared.NewXMLParser(feed)

	_, err := shared.FindRoot(p)
	if err != nil {
		return nil, err
	}

	return ap.parseRoot(p)
}

func (ap *Parser) parseRoot(p *xpp.Parser) (*Feed, error) {
	if err := p.Expect(xpp.StartTag, "feed"); err != nil {
		return nil, err
	}

	atom := &Feed{}
	atom.Entries = []*Entry{}
	atom.Version = ap.parseVersion(p)
	atom.Language = ap.parseLanguage(p)

	contributors := []*Person{}
	authors := []*Person{}
	categories := []*Category{}
	links := []*Link{}
	extensions := ext.Extensions{}

	err := shared.ForEachChild(p, func(name string) error {
		if shared.IsExtension(p) {
			var err error
			extensions, err = shared.ParseExtension(extensions, p)
			return err
		}
		var err error
		switch name {
		case "title":
			atom.Title, err = ap.parseAtomText(p)
		case "id":
			atom.ID, err = ap.parseAtomText(p)
		case "updated", "modified":
			if atom.Updated, err = ap.parseAtomText(p); err == nil {
				atom.UpdatedParsed = parseDateUTC(atom.Updated)
			}
		case "subtitle", "tagline":
			atom.Subtitle, err = ap.parseAtomText(p)
		case "link":
			var link *Link
			if link, err = ap.parseLink(p); err == nil {
				links = append(links, link)
			}
		case "generator":
			atom.Generator, err = ap.parseGenerator(p)
		case "icon":
			atom.Icon, err = ap.parseAtomText(p)
		case "logo":
			atom.Logo, err = ap.parseAtomText(p)
		case "rights", "copyright":
			atom.Rights, err = ap.parseAtomText(p)
		case "contributor":
			var person *Person
			if person, err = ap.parsePerson("contributor", p); err == nil {
				contributors = append(contributors, person)
			}
		case "author":
			var person *Person
			if person, err = ap.parsePerson("author", p); err == nil {
				authors = append(authors, person)
			}
		case "category":
			var cat *Category
			if cat, err = ap.parseCategory(p); err == nil {
				categories = append(categories, cat)
			}
		case "entry":
			var entry *Entry
			if entry, err = ap.parseEntry(p); err == nil {
				atom.Entries = append(atom.Entries, entry)
			}
		default:
			// Not part of the spec: capture it into the extension map
			// under the _custom pseudo namespace instead of dropping it.
			extensions, _, err = shared.ParseCustom(extensions, p)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if len(categories) > 0 {
		atom.Categories = categories
	}

	if len(authors) > 0 {
		atom.Authors = authors
	}

	if len(contributors) > 0 {
		atom.Contributors = contributors
	}

	if len(links) > 0 {
		atom.Links = links
	}

	if len(extensions) > 0 {
		atom.Extensions = extensions
	}

	if err := p.Expect(xpp.EndTag, "feed"); err != nil {
		return nil, err
	}

	return atom, nil
}

// parseDateUTC parses a date the historical way: the raw text is kept by the
// caller even when unparseable, and the parsed form is normalized to UTC.
func parseDateUTC(text string) *time.Time {
	if date, err := shared.ParseDate(text); err == nil {
		utc := date.UTC()
		return &utc
	}
	return nil
}

func (ap *Parser) parseEntry(p *xpp.Parser) (*Entry, error) {
	if err := p.Expect(xpp.StartTag, "entry"); err != nil {
		return nil, err
	}
	entry := &Entry{}

	contributors := []*Person{}
	authors := []*Person{}
	categories := []*Category{}
	links := []*Link{}
	extensions := ext.Extensions{}

	err := shared.ForEachChild(p, func(name string) error {
		if shared.IsExtension(p) {
			var err error
			extensions, err = shared.ParseExtension(extensions, p)
			return err
		}
		var err error
		switch name {
		case "title":
			entry.Title, err = ap.parseAtomText(p)
		case "id":
			entry.ID, err = ap.parseAtomText(p)
		case "rights", "copyright":
			entry.Rights, err = ap.parseAtomText(p)
		case "summary":
			entry.Summary, err = ap.parseAtomText(p)
		case "source":
			entry.Source, err = ap.parseSource(p)
		case "updated", "modified":
			if entry.Updated, err = ap.parseAtomText(p); err == nil {
				entry.UpdatedParsed = parseDateUTC(entry.Updated)
			}
		case "contributor":
			var person *Person
			if person, err = ap.parsePerson("contributor", p); err == nil {
				contributors = append(contributors, person)
			}
		case "author":
			var person *Person
			if person, err = ap.parsePerson("author", p); err == nil {
				authors = append(authors, person)
			}
		case "category":
			var cat *Category
			if cat, err = ap.parseCategory(p); err == nil {
				categories = append(categories, cat)
			}
		case "link":
			var link *Link
			if link, err = ap.parseLink(p); err == nil {
				links = append(links, link)
			}
		case "published", "issued":
			if entry.Published, err = ap.parseAtomText(p); err == nil {
				entry.PublishedParsed = parseDateUTC(entry.Published)
			}
		case "content":
			entry.Content, err = ap.parseContent(p)
		default:
			// Not part of the spec: capture it into the extension map
			// under the _custom pseudo namespace instead of dropping it.
			extensions, _, err = shared.ParseCustom(extensions, p)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if len(categories) > 0 {
		entry.Categories = categories
	}

	if len(authors) > 0 {
		entry.Authors = authors
	}

	if len(links) > 0 {
		entry.Links = links
	}

	if len(contributors) > 0 {
		entry.Contributors = contributors
	}

	if len(extensions) > 0 {
		entry.Extensions = extensions
	}

	if err := p.Expect(xpp.EndTag, "entry"); err != nil {
		return nil, err
	}

	return entry, nil
}

func (ap *Parser) parseSource(p *xpp.Parser) (*Source, error) {
	if err := p.Expect(xpp.StartTag, "source"); err != nil {
		return nil, err
	}

	source := &Source{}

	contributors := []*Person{}
	authors := []*Person{}
	categories := []*Category{}
	links := []*Link{}
	extensions := ext.Extensions{}

	err := shared.ForEachChild(p, func(name string) error {
		if shared.IsExtension(p) {
			var err error
			extensions, err = shared.ParseExtension(extensions, p)
			return err
		}
		var err error
		switch name {
		case "title":
			source.Title, err = ap.parseAtomText(p)
		case "id":
			source.ID, err = ap.parseAtomText(p)
		case "updated", "modified":
			if source.Updated, err = ap.parseAtomText(p); err == nil {
				source.UpdatedParsed = parseDateUTC(source.Updated)
			}
		case "subtitle", "tagline":
			source.Subtitle, err = ap.parseAtomText(p)
		case "link":
			var link *Link
			if link, err = ap.parseLink(p); err == nil {
				links = append(links, link)
			}
		case "generator":
			source.Generator, err = ap.parseGenerator(p)
		case "icon":
			source.Icon, err = ap.parseAtomText(p)
		case "logo":
			source.Logo, err = ap.parseAtomText(p)
		case "rights", "copyright":
			source.Rights, err = ap.parseAtomText(p)
		case "contributor":
			var person *Person
			if person, err = ap.parsePerson("contributor", p); err == nil {
				contributors = append(contributors, person)
			}
		case "author":
			var person *Person
			if person, err = ap.parsePerson("author", p); err == nil {
				authors = append(authors, person)
			}
		case "category":
			var cat *Category
			if cat, err = ap.parseCategory(p); err == nil {
				categories = append(categories, cat)
			}
		default:
			// Not part of the spec: capture it into the extension map
			// under the _custom pseudo namespace instead of dropping it.
			extensions, _, err = shared.ParseCustom(extensions, p)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if len(categories) > 0 {
		source.Categories = categories
	}

	if len(authors) > 0 {
		source.Authors = authors
	}

	if len(contributors) > 0 {
		source.Contributors = contributors
	}

	if len(links) > 0 {
		source.Links = links
	}

	if len(extensions) > 0 {
		source.Extensions = extensions
	}

	if err := p.Expect(xpp.EndTag, "source"); err != nil {
		return nil, err
	}

	return source, nil
}

func (ap *Parser) parseContent(p *xpp.Parser) (*Content, error) {
	c := &Content{}
	c.Type = p.Attribute("type")
	c.Src = p.Attribute("src")

	text, err := ap.parseAtomText(p)
	if err != nil {
		return nil, err
	}
	c.Value = text

	return c, nil
}

func (ap *Parser) parsePerson(name string, p *xpp.Parser) (*Person, error) {
	if err := p.Expect(xpp.StartTag, name); err != nil {
		return nil, err
	}

	person := &Person{}

	err := shared.ForEachChild(p, func(child string) error {
		var err error
		switch child {
		case "name":
			person.Name, err = ap.parseAtomText(p)
		case "email":
			person.Email, err = ap.parseAtomText(p)
		case "uri", "url", "homepage":
			person.URI, err = ap.parseAtomText(p)
		default:
			err = p.Skip()
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, name); err != nil {
		return nil, err
	}

	return person, nil
}

func (ap *Parser) parseLink(p *xpp.Parser) (*Link, error) {
	if err := p.Expect(xpp.StartTag, "link"); err != nil {
		return nil, err
	}

	l := &Link{}
	l.Href = p.Attribute("href")
	l.Hreflang = p.Attribute("hreflang")
	l.Type = p.Attribute("type")
	l.Length = p.Attribute("length")
	l.Title = p.Attribute("title")
	l.Rel = p.Attribute("rel")
	if l.Rel == "" {
		l.Rel = "alternate"
	}

	if err := p.Skip(); err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, "link"); err != nil {
		return nil, err
	}
	return l, nil
}

func (ap *Parser) parseCategory(p *xpp.Parser) (*Category, error) {
	if err := p.Expect(xpp.StartTag, "category"); err != nil {
		return nil, err
	}

	c := &Category{}
	c.Term = p.Attribute("term")
	c.Scheme = p.Attribute("scheme")
	c.Label = p.Attribute("label")

	if err := p.Skip(); err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, "category"); err != nil {
		return nil, err
	}
	return c, nil
}

func (ap *Parser) parseGenerator(p *xpp.Parser) (*Generator, error) {

	if err := p.Expect(xpp.StartTag, "generator"); err != nil {
		return nil, err
	}

	g := &Generator{}

	uri := p.Attribute("uri") // Atom 1.0
	url := p.Attribute("url") // Atom 0.3

	if uri != "" {
		g.URI = uri
	} else if url != "" {
		g.URI = url
	}

	g.Version = p.Attribute("version")

	result, err := ap.parseAtomText(p)
	if err != nil {
		return nil, err
	}

	g.Value = result

	if err := p.Expect(xpp.EndTag, "generator"); err != nil {
		return nil, err
	}

	return g, nil
}

func (ap *Parser) parseAtomText(p *xpp.Parser) (string, error) {

	var text struct {
		Type     string `xml:"type,attr"`
		Mode     string `xml:"mode,attr"`
		InnerXML string `xml:",innerxml"`
	}

	// get current base URL before it is clobbered by DecodeElement
	base := p.BaseURL()
	err := p.DecodeElement(&text)
	if err != nil {
		return "", err
	}

	result := text.InnerXML
	result = strings.TrimSpace(result)

	lowerType := strings.ToLower(text.Type)
	lowerMode := strings.ToLower(text.Mode)

	if strings.Contains(result, "<![CDATA[") {
		result = shared.StripCDATA(result)
		if lowerType == "html" || strings.Contains(lowerType, "xhtml") {
			result, _ = shared.ResolveHTML(base, result)
		}
	} else {
		// decode non-CDATA contents depending on type

		if lowerType == "text" ||
			strings.HasPrefix(lowerType, "text/") ||
			(lowerType == "" && lowerMode == "") {
			result = shared.DecodeEntities(result)
		} else if strings.Contains(lowerType, "xhtml") {
			result = ap.stripWrappingDiv(result)
			result, _ = shared.ResolveHTML(base, result)
		} else if lowerType == "html" {
			result = ap.stripWrappingDiv(result)
			result = shared.DecodeEntities(result)
			result, _ = shared.ResolveHTML(base, result)
		} else if lowerMode == "base64" || isBinaryMediaType(lowerType) {
			// Decode base64 only when the content says so: an explicit Atom 0.3
			// mode="base64", or a binary media type. Decoding by default
			// corrupts ordinary text whose type happens to be valid base64
			// (e.g. "test").
			if decoded, derr := base64.StdEncoding.DecodeString(result); derr == nil {
				result = string(decoded)
			}
		}
		// else: text with an unknown/non-binary type, leave it as parsed.
	}

	// resolve relative URIs in URI-containing elements according to xml:base
	name := strings.ToLower(p.Name())
	if atomUriElements[name] {
		resolved, err := shared.XmlBaseResolveUrl(base, result)
		if resolved != nil && err == nil {
			result = resolved.String()
		}
	}

	return result, err
}

// isBinaryMediaType reports whether an Atom content type should be treated as
// base64-encoded binary. Text and XML types never are.
func isBinaryMediaType(t string) bool {
	if t == "" || strings.HasPrefix(t, "text/") || strings.Contains(t, "xml") {
		return false
	}
	if strings.HasPrefix(t, "image/") ||
		strings.HasPrefix(t, "audio/") ||
		strings.HasPrefix(t, "video/") {
		return true
	}
	switch t {
	case "application/octet-stream", "application/pdf", "application/zip",
		"application/gzip", "application/x-gzip", "application/ogg":
		return true
	}
	return false
}

func (ap *Parser) parseLanguage(p *xpp.Parser) string {
	return p.Attribute("lang")
}

func (ap *Parser) parseVersion(p *xpp.Parser) string {
	ver := p.Attribute("version")
	if ver != "" {
		return ver
	}

	ns := p.Attribute("xmlns")
	if ns == "http://purl.org/atom/ns#" {
		return "0.3"
	}

	if ns == "http://www.w3.org/2005/Atom" {
		return "1.0"
	}

	return ""
}

// stripWrappingDiv removes the wrapping <div> an xhtml text construct
// carries per RFC 4287 section 3.1.1.3: when the parsed body holds exactly
// one element child and it is a div, the div's inner HTML is returned.
// Anything else comes back unchanged.
func (ap *Parser) stripWrappingDiv(content string) string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return content
	}
	body := findElement(doc, "body")
	if body == nil {
		return content
	}

	var div *html.Node
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if div != nil || c.Data != "div" {
			return content
		}
		div = c
	}
	if div == nil {
		return content
	}

	var buf bytes.Buffer
	for c := div.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(&buf, c); err != nil {
			return content
		}
	}
	return buf.String()
}

// findElement returns the first element with the given name in a depth-first
// walk of the parsed document.
func findElement(n *html.Node, name string) *html.Node {
	if n.Type == html.ElementNode && n.Data == name {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findElement(c, name); found != nil {
			return found
		}
	}
	return nil
}
