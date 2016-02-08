package atom

import (
	"strings"

	"github.com/mmcdole/gofeed/feed"
	"github.com/mmcdole/gofeed/shared"
	"github.com/mmcdole/goxpp"
)

type Parser struct {
	shared.BaseParser
}

func (ap *Parser) ParseFeed(feed string) (*Feed, error) {
	p := xpp.NewXMLPullParser(strings.NewReader(feed), false)

	_, err := p.NextTag()
	if err != nil {
		return nil, err
	}

	return ap.parseRoot(p)
}

func (ap *Parser) parseRoot(p *xpp.XMLPullParser) (*Feed, error) {
	if err := p.Expect(xpp.StartTag, "feed"); err != nil {
		return nil, err
	}

	atom := &Feed{}
	atom.Entries = []*Entry{}
	atom.Extensions = feed.FeedExtensions{}
	atom.Version = ap.parseVersion(p)

	contributors := []*Person{}
	authors := []*Person{}
	categories := []*Category{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {

			if ap.IsExtension(p) {
				ext, err := ap.ParseExtension(atom.Extensions, p)
				if err != nil {
					return nil, err
				}
				atom.Extensions = ext
			} else if p.Name == "title" {
				result, err := ap.parseAtomText(p)
				if err != nil {
					return nil, err
				}
				atom.Title = result
			} else if p.Name == "id" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				atom.ID = result
			} else if p.Name == "updated" ||
				p.Name == "modified" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				atom.Updated = result
				date, err := ap.ParseDate(result)
				if err == nil {
					utcDate := date.UTC()
					atom.UpdatedParsed = &utcDate
				}
			} else if p.Name == "subtitle" ||
				p.Name == "tagline" {
				result, err := ap.parseAtomText(p)
				if err != nil {
					return nil, err
				}
				atom.Subtitle = result
			} else if p.Name == "link" {
				result, err := ap.parseLink(p)
				if err != nil {
					return nil, err
				}
				atom.Link = result
			} else if p.Name == "generator" {
				result, err := ap.parseGenerator(p)
				if err != nil {
					return nil, err
				}
				atom.Generator = result
			} else if p.Name == "icon" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				atom.Icon = result
			} else if p.Name == "logo" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				atom.Logo = result
			} else if p.Name == "rights" ||
				p.Name == "copyright" {
				result, err := ap.parseAtomText(p)
				if err != nil {
					return nil, err
				}
				atom.Rights = result
			} else if p.Name == "contributor" {
				result, err := ap.parsePerson("contributor", p)
				if err != nil {
					return nil, err
				}
				contributors = append(contributors, result)
			} else if p.Name == "author" {
				result, err := ap.parsePerson("author", p)
				if err != nil {
					return nil, err
				}
				authors = append(authors, result)
			} else if p.Name == "category" {
				result, err := ap.parseCategory(p)
				if err != nil {
					return nil, err
				}
				categories = append(categories, result)
			} else if p.Name == "source" {
				result, err := ap.parseSource(p)
				if err != nil {
					return nil, err
				}
				atom.Source = result
			} else if p.Name == "entry" {
				result, err := ap.parseEntry(p)
				if err != nil {
					return nil, err
				}
				atom.Entries = append(atom.Entries, result)
			} else {
				p.Skip()
			}
		}
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

	if err := p.Expect(xpp.EndTag, "feed"); err != nil {
		return nil, err
	}

	return atom, nil
}

func (ap *Parser) parseEntry(p *xpp.XMLPullParser) (*Entry, error) {
	if err := p.Expect(xpp.StartTag, "entry"); err != nil {
		return nil, err
	}
	entry := &Entry{}
	entry.Extensions = feed.FeedExtensions{}

	contributors := []*Person{}
	authors := []*Person{}
	categories := []*Category{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {

			if ap.IsExtension(p) {
				ext, err := ap.ParseExtension(entry.Extensions, p)
				if err != nil {
					return nil, err
				}
				entry.Extensions = ext
			} else if p.Name == "title" {
				result, err := ap.parseAtomText(p)
				if err != nil {
					return nil, err
				}
				entry.Title = result
			} else if p.Name == "id" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				entry.ID = result
			} else if p.Name == "updated" ||
				p.Name == "modified" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				entry.Updated = result
				date, err := ap.ParseDate(result)
				if err == nil {
					utcDate := date.UTC()
					entry.UpdatedParsed = &utcDate
				}
			} else if p.Name == "contributor" {
				result, err := ap.parsePerson("contributor", p)
				if err != nil {
					return nil, err
				}
				contributors = append(contributors, result)
			} else if p.Name == "author" {
				result, err := ap.parsePerson("author", p)
				if err != nil {
					return nil, err
				}
				authors = append(authors, result)
			} else if p.Name == "category" {
				result, err := ap.parseCategory(p)
				if err != nil {
					return nil, err
				}
				categories = append(categories, result)
			} else if p.Name == "link" {
				result, err := ap.parseLink(p)
				if err != nil {
					return nil, err
				}
				entry.Link = result
			} else if p.Name == "published" ||
				p.Name == "issued" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				entry.Published = result
				date, err := ap.ParseDate(result)
				if err == nil {
					utcDate := date.UTC()
					entry.PublishedParsed = &utcDate
				}
			} else if p.Name == "content" {
				result, err := ap.parseContent(p)
				if err != nil {
					return nil, err
				}
				entry.Content = result
			} else {
				p.Skip()
			}
		}
	}

	if len(categories) > 0 {
		entry.Categories = categories
	}

	if len(authors) > 0 {
		entry.Authors = authors
	}

	if len(contributors) > 0 {
		entry.Contributors = contributors
	}

	if err := p.Expect(xpp.EndTag, "entry"); err != nil {
		return nil, err
	}

	return entry, nil
}

func (ap *Parser) parseSource(p *xpp.XMLPullParser) (*Source, error) {

	if err := p.Expect(xpp.StartTag, "source"); err != nil {
		return nil, err
	}

	source := &Source{}
	source.Extensions = feed.FeedExtensions{}

	contributors := []*Person{}
	authors := []*Person{}
	categories := []*Category{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {

			if ap.IsExtension(p) {
				ext, err := ap.ParseExtension(source.Extensions, p)
				if err != nil {
					return nil, err
				}
				source.Extensions = ext
			} else if p.Name == "title" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				source.Title = result
			} else if p.Name == "id" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				source.ID = result
			} else if p.Name == "updated" ||
				p.Name == "modified" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				source.Updated = result
				date, err := ap.ParseDate(result)
				if err == nil {
					utcDate := date.UTC()
					source.UpdatedParsed = &utcDate
				}
			} else if p.Name == "subtitle" ||
				p.Name == "tagline" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				source.Subtitle = result
			} else if p.Name == "link" {
				result, err := ap.parseLink(p)
				if err != nil {
					return nil, err
				}
				source.Link = result
			} else if p.Name == "generator" {
				result, err := ap.parseGenerator(p)
				if err != nil {
					return nil, err
				}
				source.Generator = result
			} else if p.Name == "icon" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				source.Icon = result
			} else if p.Name == "logo" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				source.Logo = result
			} else if p.Name == "rights" ||
				p.Name == "copyright" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				source.Rights = result
			} else if p.Name == "contributor" {
				result, err := ap.parsePerson("contributor", p)
				if err != nil {
					return nil, err
				}
				contributors = append(contributors, result)
			} else if p.Name == "author" {
				result, err := ap.parsePerson("author", p)
				if err != nil {
					return nil, err
				}
				authors = append(authors, result)
			} else if p.Name == "category" {
				result, err := ap.parseCategory(p)
				if err != nil {
					return nil, err
				}
				categories = append(categories, result)
			} else {
				p.Skip()
			}
		}
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

	if err := p.Expect(xpp.EndTag, "source"); err != nil {
		return nil, err
	}

	return source, nil
}

func (ap *Parser) parseAtomText(p *xpp.XMLPullParser) (string, error) {

	var text struct {
		Type     string `xml:"type,attr"`
		Body     string `xml:",chardata"`
		InnerXML string `xml:",innerxml"`
	}

	err := p.DecodeElement(&text)
	if err != nil {
		return "", err
	}

	// TODO: unwrap XHTML surrounding div
	// and handle other rules based on type
	result := ""
	if len(text.InnerXML) > 0 {
		result = text.InnerXML
	} else if len(text.Body) > 0 {
		result = text.Body
	}

	return result, nil
}

func (ap *Parser) parseContent(p *xpp.XMLPullParser) (*Content, error) {

	var content struct {
		Src      string `xml:"src,attr"`
		Type     string `xml:"type,attr"`
		Body     string `xml:",chardata"`
		InnerXML string `xml:",innerxml"`
	}

	err := p.DecodeElement(&content)
	if err != nil {
		return nil, err
	}

	c := &Content{}
	c.Type = content.Type
	c.Src = content.Src

	// TODO: base64 decode?
	if len(content.InnerXML) > 0 {
		c.Value = strings.TrimSpace(content.InnerXML)
	} else if len(content.Body) > 0 {
		c.Value = strings.TrimSpace(content.Body)
	}

	return c, nil
}

func (ap *Parser) parsePerson(name string, p *xpp.XMLPullParser) (*Person, error) {

	if err := p.Expect(xpp.StartTag, name); err != nil {
		return nil, err
	}

	person := &Person{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			if p.Name == "name" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				person.Name = result
			} else if name == "email" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				person.Email = result
			} else if name == "uri" ||
				name == "url" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				person.URI = result
			} else {
				p.Skip()
			}
		}
	}

	if err := p.Expect(xpp.EndTag, name); err != nil {
		return nil, err
	}

	return person, nil
}

func (ap *Parser) parseLink(p *xpp.XMLPullParser) (*Link, error) {
	if err := p.Expect(xpp.StartTag, "link"); err != nil {
		return nil, err
	}

	l := &Link{}
	l.Href = p.Attribute("href")
	l.Rel = p.Attribute("rel")
	l.Hreflang = p.Attribute("hreflang")
	l.Type = p.Attribute("type")
	l.Length = p.Attribute("length")

	if err := p.Skip(); err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, "link"); err != nil {
		return nil, err
	}
	return l, nil
}

func (ap *Parser) parseCategory(p *xpp.XMLPullParser) (*Category, error) {
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

func (ap *Parser) parseGenerator(p *xpp.XMLPullParser) (*Generator, error) {

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

	result, err := ap.ParseText(p)
	if err != nil {
		return nil, err
	}

	g.Value = result

	if err := p.Expect(xpp.EndTag, "generator"); err != nil {
		return nil, err
	}

	return g, nil
}

func (ap *Parser) parseVersion(p *xpp.XMLPullParser) string {
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
