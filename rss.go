package feed

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/go-xpp"
)

type RSSFeed struct {
	Title               string
	Link                string
	Description         string
	Language            string
	Copyright           string
	ManagingEditor      string
	WebMaster           string
	PubDate             string
	PubDateParsed       time.Time
	LastBuildDate       string
	LastBuildDateParsed time.Time
	Categories          []string
	Generator           string
	Docs                string
	TTL                 string
	Image               RSSImage
	Rating              string
	SkipHours           []string
	SkipDays            []string
	Items               []*RSSItem
	Version             string
	Extensions          map[string]interface{}
}

func (f *RSSFeed) String() string {
	return fmt.Sprintf("Title: %s\nLink: %s\nDescription: %s\n"+
		"Language: %s\nCopyright: %s\nManagingEditor: %s\n"+
		"WebMaster: %s\nPubDate: %s\nLastBuildDate: %s\n"+
		"Generator: %s\nDocs: %s\nTTL: %s\n"+
		"Rating: %s\nItems: %s\nVersion: %s\n",
		f.Title, f.Link, f.Description,
		f.Language, f.Copyright, f.ManagingEditor,
		f.WebMaster, f.PubDate, f.LastBuildDate,
		f.Generator, f.Docs, f.TTL,
		f.Rating, f.Items, f.Version)
}

type RSSItem struct {
	Title         string
	Link          string
	Description   string
	Author        string
	Categories    []string
	Comments      string
	Enclosure     RSSEnclosure
	Guid          RSSGuid
	PubDate       string
	PubDateParsed time.Time
	Source        RSSSource
	Extensions    map[string]interface{}
}

func (i *RSSItem) String() string {
	return fmt.Sprintf("Title: %s\nLink: %s\nDescription: %s\n"+
		"Author: %s\nComments: %s\nPubDate: %s\n"+
		"Source: %s\n",
		i.Title, i.Link, i.Description,
		i.Author, i.Comments, i.PubDate,
		i.Source)
}

type RSSImage struct {
	URL    string
	Link   string
	Width  string
	Height string
}

type RSSEnclosure struct {
	URL    string
	Length string
	Type   string
}

type RSSGuid struct {
	Value       string
	IsPermalink string
}

type RSSSource struct {
	Title string
	URL   string
}

func ParseRSSFeed(feed string) (rss *RSSFeed, err error) {
	p := xpp.NewXMLPullParser(strings.NewReader(feed))

	_, err = p.NextTag()
	if err != nil {
		return
	}

	return parseRoot(p)
}

func parseRoot(p *xpp.XMLPullParser) (rss *RSSFeed, err error) {

	if !p.Matches(xpp.StartTag, "rss") &&
		!p.Matches(xpp.StartTag, "RDF") {
		return nil, errors.New("Unexpected root element")
	}

	items := []*RSSItem{}
	ver := parseVersion(p)

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			if p.Name == "channel" {
				rss, err = parseChannel(p)
				if err != nil {
					return nil, err
				}
			} else if p.Name == "item" {
				// Earlier versions of the RSS spec had "item" elements at the same
				// root level as "channel" elements.  We will merge these items
				// with any channel level items.
				item, err := parseItem(p)
				if err != nil {
					return nil, err
				}
				items = append(items, item)
			} else {
				// Skip any elements that are not "channel" or "item"
				p.Skip()
			}
		}
	}

	if !p.Matches(xpp.EndTag, "rss") &&
		!p.Matches(xpp.EndTag, "RDF") {
		return nil, errors.New("Expected root end tag")
	}

	if rss != nil {
		rss.Items = append(rss.Items, items...)
		rss.Version = ver
		return rss, nil
	} else {
		return nil, errors.New("No channel element found.")
	}
}

func parseChannel(p *xpp.XMLPullParser) (*RSSFeed, error) {
	if !p.Matches(xpp.StartTag, "channel") {
		return nil, errors.New("Expected channel start tag")
	}

	rss := &RSSFeed{}
	rss.Items = []*RSSItem{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			if p.Name == "title" {
				title, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.Title = title
			} else if p.Name == "description" {
				desc, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.Description = desc
			} else if p.Name == "link" {
				link, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.Link = link
			} else if p.Name == "language" {
				lang, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.Language = lang
			} else if p.Name == "copyright" {
				copyright, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.Copyright = copyright
			} else if p.Name == "managingEditor" {
				editor, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.ManagingEditor = editor
			} else if p.Name == "webMaster" {
				web, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.WebMaster = web
			} else if p.Name == "pubDate" {
				pub, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.PubDate = pub
			} else if p.Name == "lastBuildDate" {
				build, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.LastBuildDate = build
			} else if p.Name == "generator" {
				gen, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.Generator = gen
			} else if p.Name == "docs" {
				docs, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.Docs = docs
			} else if p.Name == "ttl" {
				ttl, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.TTL = ttl
			} else if p.Name == "rating" {
				rating, err := p.NextText()
				if err != nil {
					return nil, err
				}
				rss.Rating = rating
			} else if p.Name == "item" {
				item, err := parseItem(p)
				if err != nil {
					return nil, err
				}
				rss.Items = append(rss.Items, item)
			} else {
				// Skip any elements not part of the channel spec
				p.Skip()
			}
		}
	}

	if !p.Matches(xpp.EndTag, "channel") {
		return nil, errors.New("Expected channel end tag")
	}

	return rss, nil
}

func parseItem(p *xpp.XMLPullParser) (*RSSItem, error) {
	if !p.Matches(xpp.StartTag, "item") {
		return nil, errors.New("Expected item start tag")
	}

	item := &RSSItem{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			if p.Name == "title" {
				title, err := p.NextText()
				if err != nil {
					return nil, err
				}
				item.Title = title
			} else if p.Name == "description" {
				desc, err := p.NextText()
				if err != nil {
					return nil, err
				}
				item.Description = desc
			} else if p.Name == "link" {
				link, err := p.NextText()
				if err != nil {
					return nil, err
				}
				item.Link = link
			} else if p.Name == "author" {
				author, err := p.NextText()
				if err != nil {
					return nil, err
				}
				item.Author = author
			} else if p.Name == "comments" {
				comments, err := p.NextText()
				if err != nil {
					return nil, err
				}
				item.Comments = comments
			} else if p.Name == "pubDate" {
				pub, err := p.NextText()
				if err != nil {
					return nil, err
				}
				item.PubDate = pub
			} else if p.Name == "source" {
				source, err := parseSource(p)
				if err != nil {
					return nil, err
				}
				item.Source = source
			} else {
				// Skip any elements not part of the item spec
				p.Skip()
			}
		}
	}

	if !p.Matches(xpp.EndTag, "item") {
		return nil, errors.New("Expected item end tag")
	}

	return item, nil
}

func parseSource(p *xpp.XMLPullParser) (*RSSSource, error) {
	if !p.Matches(xpp.StartTag, "source") {
		return nil, errors.New("Expected source start tag")
	}

	source := &RSSSource{}
	source.URL = p.Attribute("url")

	title, err := p.NextText()
	if err != nil {
		return nil, err
	}

	source.Title = title

	if !p.Matches(xpp.EndTag, "source") {
		return nil, errors.New("Expected source end tag")
	}
	return source, nil
}

func parseVersion(p *xpp.XMLPullParser) (ver string) {
	if p.Name == "rss" {
		ver = p.Attribute("version")
		if ver == "" {
			ver = "2.0"
		}
	} else if p.Name == "RDF" {
		ns := p.Attribute("xmlns")
		if ns == "http://channel.netscape.com/rdf/simple/0.9/" {
			ver = "0.9"
		} else if ns == "http://purl.org/rss/1.0/" {
			ver = "1.0"
		}
	}
	return
}
