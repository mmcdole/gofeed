package feed

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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

func ParseRSSFeed(feed string) (*RSSFeed, error) {
	fmt.Println("Parsing feed...")
	var rss *RSSFeed
	d := xml.NewDecoder(strings.NewReader(feed))
	d.Strict = false

TokenLoop:
	for {
		tok, err := getToken(d)
		if err != nil {
			return nil, err
		}

		switch tt := tok.(type) {
		case xml.StartElement:
			name := strings.ToLower(tt.Name.Local)
			if name == "rdf" || name == "rss" {
				rss = parseRoot(tt, d)
				break TokenLoop
			} else {
				// Shouldn't be necessary since an XML document
				// can only have a single root element.
				d.Skip()
			}
		}
	}
}

func parseRoot(root xml.StartElement, d *xml.Decoder) (*RSSFeed, error) {
	fmt.Println("Parsing root...")
	var rss *RSSFeed
	rootItems := []*RSSItem{}

TokenLoop:
	for {
		tok, err := getToken(d)
		if err != nil {
			return nil, err
		}

		switch tt := tok.(type) {
		case xml.StartElement:
			name := strings.ToLower(tt.Name.Local)
			if name == "channel" {
				if rss != nil {
					// Skip any subsequent "channel" elements after we have already
					// parsed one.
					d.Skip()
				} else {
					if rss, err = parseChannel(tt, d); err != nil {
						return nil, err
					}
				}
			} else if name == "item" {
				// Earlier versions of the RSS spec had "item" elements at the same
				// root level as "channel" elements.
				if item, err := parseItem(tt, d); err != nil {
					return nil, err
				} else {
					rootItems = append(rootItems, item)
				}
			} else {
				// Skip any elements that are not "channel" or "item"
				d.Skip()
			}
		case xml.EndElement:
			// End the root element
			break TokenLoop
		}
	}

	if rss == nil {
		return nil, errors.New("No channel element found.")
	}

	rss.Items = append(rss.Items, rootItems...)
	rss.Version = parseVersion(root)
	return rss, nil
}

func parseChannel(t xml.StartElement, d *xml.Decoder) (*RSSFeed, error) {
	fmt.Println("Parsing channel...")
	rss := &RSSFeed{}
	rss.Items = []*RSSItem{}

TokenLoop:
	for {
		tok, err := getToken(d)
		if err != nil {
			return nil, err
		}

		switch tt := tok.(type) {
		case xml.StartElement:
			name := strings.ToLower(tt.Name.Local)

			if name == "title" {
				if title, err := getElementText(tt, d); err != nil {
					return nil, err
				} else {
					rss.Title = title
				}
				d.Skip()
			} else {
				d.Skip()
			}
		case xml.EndElement:
			break TokenLoop
		}
	}

	return rss, nil
}

func parseItem(t xml.StartElement, d *xml.Decoder) (*RSSItem, error) {
	return &RSSItem{}, nil
}

func parseVersion(root xml.StartElement) string {
	var result string
	name := strings.ToLower(root.Name.Local)
	if name == "rss" {
		version := getAttrValue("version", root.Attr)
		if version != "" {
			result = version
		}
	} else if name == "rdf" {
		ns := getAttrValue("xmlns", root.Attr)
		if ns == "http://channel.netscape.com/rdf/simple/0.9/" {
			result = "0.9"
		} else if ns == "http://purl.org/rss/1.0/" {
			result = "1.0"
		}
	}
	return result
}

func getElementText(e xml.StartElement, d *xml.Decoder) (string, error) {
	var result string

TokenLoop:
	for {
		tok, err := getToken(d)
		if err != nil {
			return "", err
		}

		switch tt := tok.(type) {
		case xml.StartElement:
			d.Skip()
		case xml.EndElement:
			break TokenLoop
		case xml.CharData:
			result = string([]byte(tt))
		}
	}
	return result, nil
}

func getToken(d *xml.Decoder) (xml.Token, error) {
	tok, err := d.Token()
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("Unexpected end of feed")
		}
		return nil, err
	}
	return tok, nil
}

func getAttrValue(name string, attrs []xml.Attr) string {
	n := strings.ToLower(name)
	for _, attr := range attrs {
		if strings.ToLower(attr.Name.Local) == n {
			return attr.Value
		}
	}
	return ""
}
