package feed

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
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
	Items               []RSSItem
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

func ParseRSSFeed(feed string) (*RSSFeed, error) {
	rss := &RSSFeed{}

	d := xml.NewDecoder(strings.NewReader(feed))
	d.Strict = false

	for {
		if tok, err = d.Token(); err != nil {
			if err == io.EOF {
				return errors.New("No root node found.")
			}
			return err
		}

		switch tt := tok.(type) {
		case xml.SyntaxError:
			return nil, errors.New(tt.Error())
		case xml.StartElement:
			name := strings.ToLower(tok.Name.Local)
			if name != "rdf" && name != "rss" {
				return nil, fmt.Errorf("Invalid root node: %s", name)
			} else {
				err := parseRoot(t, d, rss)
				if err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			return rss, nil
		}
	}
}

func parseRoot(t xml.StartElement, d *xml.Decoder, rss *RSSFeed) error {
	rss.Version = parseVersion(t)

	for {
		if tok, err = d.Token(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch tt := token.(type) {
		case xml.StartElement:
			name := strings.ToLower(tt.Name.Local)
			if name == "channel" {
				parseChannel(tt, d, rss)
			} else if name == "item" {
				parseItem(tt, d, rss)
			} else {
				d.Skip()
			}
		}
	}
}

func parseChannel(t xml.StartElement, d *xml.Decoder, rss *RSSFeed) error {
	for {
		if tok, err = d.Token(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch tt := token.(type) {
		case xml.StartElement:

		case xml.CharData:

		case xml.EndElement:

		}
	}
}

func parseItem(t xml.StartElement, d *xml.Decoder, rss *RSSFeed) error {
	for {
		token, err := d.Token()

		if err != nil {
			return err
		}

		if token == nil {
			return nil
		}

		switch t := token.(type) {
		case xml.StartElement:
			name := strings.ToLower(t.Name.Local)
			switch name {
			case "channel":
				parseChannel(t, d, rss)
			case "item":
				parseItem(t, d, rss)
			default:
				d.Skip()
			}
		}
	}
	return nil
}

func parseVersion(root xml.StartElement) string {
	var result string
	name := strings.ToLower(root.Name.Local)
	if name == "rss" {
		version := attrValue("version", root.Attr)
		if version != "" {
			result = version
		}
	} else if name == "rdf" {
		ns := attrValue("xmlns", root.Attr)
		if ns == "http://channel.netscape.com/rdf/simple/0.9/" {
			result = "0.9"
		} else if ns == "http://purl.org/rss/1.0/" {
			result = "1.0"
		}
	}
	return result
}

func attrValue(name string, attrs []xml.Attr) string {
	n := strings.ToLower(name)
	for _, attr := range attrs {
		if strings.ToLower(attr.Name.Local) == n {
			return attr.Value
		}
	}
	return ""
}
