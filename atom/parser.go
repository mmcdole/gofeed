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
	p := xpp.NewXMLPullParser(strings.NewReader(feed))

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
				result, err := ap.ParseText(p)
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
			} else if p.Name == "updated" {
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
			} else if p.Name == "subtitle" {
				result, err := ap.ParseText(p)
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
			} else if p.Name == "rights" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				atom.Rights = result
			} else if p.Name == "contributor" {
				result, err := ap.ParseText(p)
				if err != nil {
					return nil, err
				}
				atom.Rights = result
			} else {
				p.Skip()
			}
		}
	}

	if err := p.Expect(xpp.EndTag, "feed"); err != nil {
		return nil, err
	}

	return atom, nil
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

	if _, err := p.NextTag(); err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, "link"); err != nil {
		return nil, err
	}
	return l, nil
}

func (ap *Parser) parseGenerator(p *xpp.XMLPullParser) (*Generator, error) {
	if err := p.Expect(xpp.StartTag, "generator"); err != nil {
		return nil, err
	}
	g := &Generator{}
	g.Term = p.Attribute("term")
	g.Scheme = p.Attribute("scheme")
	g.Label = p.Attribute("label")

	if _, err := p.NextTag(); err != nil {
		return nil, err
	}

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
