package feed

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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
	Source        string
	Extensions    map[string]interface{}
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

func ParseRSSFeed(feed string) (rss *RSSFeed, err error) {
	fmt.Println("Parsing feed...")

	p := xpp.NewXMLPullParser(strings.NewReader(feed))

	tok, err := p.NextTag()
	if err != nil {
		return
	}

	return parseRoot(p)
}

func parseRoot(p *xpp.XMLPullParser) (rss *RSSFeed, err error) {

	if !p.Matches(xpp.StartTag, nil, "rss") &&
		!p.Matches(xpp.StartTag, nil, "RDF") {
		return nil, errors.New("Unexpected root element")
	}

	items := []*RSSItem{}
	ver = parseVersion(p)

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
				item, err = parseItem(i)
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

	if !p.Matches(xpp.EndTag, nil, "rss") &&
		!p.Matches(xpp.EndTag, nil, "RDF") {
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

func parseChannel(t xml.StartElement, d *xml.Decoder) (*RSSFeed, error) {
	if !p.Matches(xpp.StartTag, nil, "channel") {
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
			} else {
				// Skip any elements not part of the channel spec
				p.Skip()
			}
		}
	}

	if !p.Matches(xpp.EndTag, nil, "channel") {
		return nil, errors.New("Expected channel end tag")
	}

	return rss, nil
}

func parseItem(t xml.StartElement, d *xml.Decoder) (*RSSItem, error) {
	return &RSSItem{}, nil
}

func parseVersion(p *xpp.XMLPullParser) (ver string) {
	if p.Name == "rss" {
		ver = p.Attribute("version")
	} else if p.Name == "RDF" {
		if p.Space == "http://channel.netscape.com/rdf/simple/0.9/" {
			ver = "0.9"
		} else if p.Space == "http://purl.org/rss/1.0/" {
			ver = "1.0"
		}
	}
	return
}
