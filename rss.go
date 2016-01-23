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
		token, err := d.Token()

		if err != nil {
			return nil, err
		}

		if token == nil {
			return nil, errors.New("No root node found.")
		}

		switch t := token.(type) {
		case xml.StartElement:
			err := parseRoot(t, d, rss)
			if err != nil {
				return nil, err
			}
			return rss, nil
		}
	}
}

func parseRoot(t xml.StartElement, d *xml.Decoder, rss *RSSFeed) error {
	name := strings.ToLower(t.Name.Local)

	// Parse RSS version from root
	switch name {
	case "rss":
		version := attrValue("version", t.Attr)
		if version != "" {
			rss.Version = version
		}
	case "rdf":
		ns := attrValue("xmlns", t.Attr)
		if ns == "http://channel.netscape.com/rdf/simple/0.9/" {
			rss.Version = "0.9"
		} else if ns == "http://purl.org/rss/1.0/" {
			rss.Version = "1.0"
		}
	default:
		return fmt.Errorf("Invalid root element: %s", name)
	}

	//	for {
	//		token, err := d.Token()
	//
	//		if err != nil || token == nil {
	//			fmt.Println(err)
	//			return nil, err
	//		}
	//
	//		switch t := token.(type) {
	//		case xml.StartElement:
	//			err := parseRoot(t, d, rss)
	//			if err != nil {
	//				return nil, error
	//			}
	//		}
	return nil
}

func attrValue(name string, attrs []xml.Attr) string {
	lname := strings.ToLower(name)
	for _, attr := range attrs {
		if strings.ToLower(attr.Name.Local) == lname {
			return attr.Value
		}
	}
	return ""
}

//func extractNamespaces(attrs []xml.Attr) map[xml.Name]string {
//	ns := make(map[xml.Name]string)
//
//	ns[xml.Name{Space: "", Local: "xml"}] = "http://www.w3.org/XML/1998/namespace"
//
//	for i := range attrs {
//		attr := attrs[i].Name
//		val := attrs[i].Value
//
//		if (attr.Local == "xmlns" && attr.Space == "") || attr.Space == "xmlns" {
//			if attr.Local == "xmlns" && attr.Space == "" && val == "" {
//				delete(ns, attr)
//			} else {
//				ns[attr] = val
//			}
//		} else {
//			attrs = append(attrs, &xmlattr.XMLAttr{Attr: ele.Attr[i], Parent: ch})
//		}
//	}
//}
