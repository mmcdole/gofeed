package rss

import (
	"fmt"
	"io"
	"strings"
	"time"

	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/mmcdole/gofeed/internal/shared"
	xpp "github.com/mmcdole/goxpp/v2"
)

// Parser is a RSS Parser
type Parser struct{}

// Parse parses an xml feed into an rss.Feed
func (rp *Parser) Parse(feed io.Reader) (*Feed, error) {
	feed = shared.NewControlCharFilterReader(feed)
	p := shared.NewXMLParser(feed)

	_, err := shared.FindRoot(p)
	if err != nil {
		return nil, err
	}

	return rp.parseRoot(p)
}

func (rp *Parser) parseRoot(p *xpp.Parser) (*Feed, error) {
	rssErr := p.Expect(xpp.StartTag, "rss")
	rdfErr := p.Expect(xpp.StartTag, "rdf")
	if rssErr != nil && rdfErr != nil {
		return nil, fmt.Errorf("%s or %s", rssErr.Error(), rdfErr.Error())
	}

	// Items found in feed root
	var channel *Feed
	var textinput *TextInput
	var image *Image
	items := []*Item{}

	ver := rp.parseVersion(p)

	err := shared.ForEachChild(p, func(name string) error {
		// Skip any extensions found in the feed root.
		if shared.IsExtension(p) {
			return p.Skip()
		}
		var err error
		switch name {
		case "channel":
			channel, err = rp.parseChannel(p)
		case "item":
			var item *Item
			if item, err = rp.parseItem(p); err == nil {
				items = append(items, item)
			}
		case "textinput":
			textinput, err = rp.parseTextInput(p)
		case "image":
			image, err = rp.parseImage(p)
		default:
			err = p.Skip()
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	rssErr = p.Expect(xpp.EndTag, "rss")
	rdfErr = p.Expect(xpp.EndTag, "rdf")
	if rssErr != nil && rdfErr != nil {
		return nil, fmt.Errorf("%s or %s", rssErr.Error(), rdfErr.Error())
	}

	if channel == nil {
		channel = &Feed{}
		channel.Items = []*Item{}
	}

	if len(items) > 0 {
		channel.Items = append(channel.Items, items...)
	}

	if textinput != nil {
		channel.TextInput = textinput
	}

	if image != nil {
		channel.Image = image
	}

	channel.Version = ver
	return channel, nil
}

func (rp *Parser) parseChannel(p *xpp.Parser) (rss *Feed, err error) {
	if err = p.Expect(xpp.StartTag, "channel"); err != nil {
		return nil, err
	}

	rss = &Feed{}
	rss.Items = []*Item{}

	extensions := ext.Extensions{}
	categories := []*Category{}
	links := []string{}

	// parseDate mirrors the historical behavior for pubDate and
	// lastBuildDate: the raw text is kept even when unparseable, and the
	// parsed form is normalized to UTC.
	parseDate := func(text string) *time.Time {
		if date, err := shared.ParseDate(text); err == nil {
			utc := date.UTC()
			return &utc
		}
		return nil
	}

	err = shared.ForEachChild(p, func(name string) error {
		if shared.IsExtension(p) {
			extensions, err = shared.ParseExtension(extensions, p)
			return err
		}
		var err error
		switch name {
		case "title":
			rss.Title, err = shared.ParseText(p)
		case "description":
			rss.Description, err = shared.ParseText(p)
		case "link":
			if rss.Link, err = rp.parseLink(p); err == nil {
				links = append(links, rss.Link)
			}
		case "language":
			rss.Language, err = shared.ParseText(p)
		case "copyright":
			rss.Copyright, err = shared.ParseText(p)
		case "managingeditor":
			rss.ManagingEditor, err = shared.ParseText(p)
		case "webmaster":
			rss.WebMaster, err = shared.ParseText(p)
		case "pubdate":
			if rss.PubDate, err = shared.ParseText(p); err == nil {
				rss.PubDateParsed = parseDate(rss.PubDate)
			}
		case "lastbuilddate":
			if rss.LastBuildDate, err = shared.ParseText(p); err == nil {
				rss.LastBuildDateParsed = parseDate(rss.LastBuildDate)
			}
		case "generator":
			rss.Generator, err = shared.ParseText(p)
		case "docs":
			rss.Docs, err = shared.ParseTextURL(p)
		case "ttl":
			rss.TTL, err = shared.ParseText(p)
		case "rating":
			rss.Rating, err = shared.ParseText(p)
		case "skiphours":
			rss.SkipHours, err = rp.parseSkipHours(p)
		case "skipdays":
			rss.SkipDays, err = rp.parseSkipDays(p)
		case "item":
			var item *Item
			if item, err = rp.parseItem(p); err == nil {
				rss.Items = append(rss.Items, item)
			}
		case "cloud":
			rss.Cloud, err = rp.parseCloud(p)
		case "category":
			var cat *Category
			if cat, err = rp.parseCategory(p); err == nil {
				categories = append(categories, cat)
			}
		case "image":
			rss.Image, err = rp.parseImage(p)
		case "textinput":
			rss.TextInput, err = rp.parseTextInput(p)
		case "items":
			// The RSS 1.0 <items> rdf:Seq is a structural list of item
			// references, not content; skip it rather than capture it.
			err = p.Skip()
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

	if err = p.Expect(xpp.EndTag, "channel"); err != nil {
		return nil, err
	}

	if len(categories) > 0 {
		rss.Categories = categories
	}

	if len(links) > 0 {
		rss.Links = links
	}

	if len(extensions) > 0 {
		rss.Extensions = extensions

		if itunes, ok := rss.Extensions["itunes"]; ok {
			rss.ITunesExt = ext.NewITunesFeedExtension(itunes)
		}

		if dc, ok := rss.Extensions["dc"]; ok {
			rss.DublinCoreExt = ext.NewDublinCoreExtension(dc)
		}
	}

	return rss, nil
}

func (rp *Parser) parseItem(p *xpp.Parser) (item *Item, err error) {
	if err = p.Expect(xpp.StartTag, "item"); err != nil {
		return nil, err
	}

	item = &Item{}
	extensions := ext.Extensions{}
	categories := []*Category{}
	enclosures := []*Enclosure{}
	links := []string{}

	// captureCustom files an unrecognized child in the extension map under
	// the _custom pseudo namespace, preserving nesting and attributes, and
	// keeps the flat Custom map populated for elements without children so
	// existing Custom reads are unchanged.
	captureCustom := func(p *xpp.Parser) error {
		var e ext.Extension
		var err error
		extensions, e, err = shared.ParseCustom(extensions, p)
		if err != nil {
			return err
		}
		if len(e.Children) == 0 {
			if item.Custom == nil {
				item.Custom = make(map[string]string)
			}
			item.Custom[e.Name] = shared.DecodeEntities(e.Value)
		}
		return nil
	}

	err = shared.ForEachChild(p, func(name string) error {
		if shared.IsExtension(p) {
			extensions, err = shared.ParseExtension(extensions, p)
			return err
		}
		var err error
		switch name {
		case "title":
			item.Title, err = shared.ParseText(p)
		case "description":
			item.Description, err = shared.ParseText(p)
		case "encoded":
			if shared.PrefixForNamespace(p.Space(), p) == "content" {
				item.Content, err = shared.ParseText(p)
			} else {
				err = captureCustom(p)
			}
		case "link":
			if item.Link, err = rp.parseLink(p); err == nil {
				links = append(links, item.Link)
			}
		case "author":
			item.Author, err = shared.ParseText(p)
		case "comments":
			item.Comments, err = shared.ParseTextURL(p)
		case "pubdate":
			if item.PubDate, err = shared.ParseText(p); err == nil {
				if date, derr := shared.ParseDate(item.PubDate); derr == nil {
					utc := date.UTC()
					item.PubDateParsed = &utc
				}
			}
		case "source":
			item.Source, err = rp.parseSource(p)
		case "enclosure":
			var enc *Enclosure
			if enc, err = rp.parseEnclosure(p); err == nil {
				item.Enclosure = enc
				enclosures = append(enclosures, enc)
			}
		case "guid":
			item.GUID, err = rp.parseGUID(p)
		case "category":
			var cat *Category
			if cat, err = rp.parseCategory(p); err == nil {
				categories = append(categories, cat)
			}
		default:
			err = captureCustom(p)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if len(enclosures) > 0 {
		item.Enclosures = enclosures
	}

	if len(categories) > 0 {
		item.Categories = categories
	}

	if len(links) > 0 {
		item.Links = links
	}

	if len(extensions) > 0 {
		item.Extensions = extensions

		if itunes, ok := item.Extensions["itunes"]; ok {
			item.ITunesExt = ext.NewITunesItemExtension(itunes)
		}

		if dc, ok := item.Extensions["dc"]; ok {
			item.DublinCoreExt = ext.NewDublinCoreExtension(dc)
		}
	}

	if err = p.Expect(xpp.EndTag, "item"); err != nil {
		return nil, err
	}

	return item, nil
}

func (rp *Parser) parseLink(p *xpp.Parser) (url string, err error) {
	base := p.BaseURL()
	href := p.Attribute("href")
	url, err = shared.ParseText(p)
	if err != nil {
		return
	}
	if url == "" && href != "" {
		url = href
	}
	url = shared.ResolveURLIfBase(base, url)
	return url, err
}

func (rp *Parser) parseSource(p *xpp.Parser) (source *Source, err error) {
	if err = p.Expect(xpp.StartTag, "source"); err != nil {
		return nil, err
	}

	source = &Source{}
	source.URL = p.Attribute("url")

	result, err := shared.ParseText(p)
	if err != nil {
		return source, err
	}
	source.Title = result

	if err = p.Expect(xpp.EndTag, "source"); err != nil {
		return nil, err
	}
	return source, nil
}

func (rp *Parser) parseEnclosure(p *xpp.Parser) (enclosure *Enclosure, err error) {
	if err = p.Expect(xpp.StartTag, "enclosure"); err != nil {
		return nil, err
	}

	enclosure = &Enclosure{}
	enclosure.URL = p.Attribute("url")
	enclosure.Length = p.Attribute("length")
	enclosure.Type = p.Attribute("type")

	// The spec defines <enclosure> as an empty element; skip to the
	// matching end tag so stray children can't desync the parse.
	if err = p.Skip(); err != nil {
		return nil, err
	}

	if err = p.Expect(xpp.EndTag, "enclosure"); err != nil {
		return nil, err
	}

	return enclosure, nil
}

func (rp *Parser) parseImage(p *xpp.Parser) (image *Image, err error) {
	if err = p.Expect(xpp.StartTag, "image"); err != nil {
		return nil, err
	}

	image = &Image{}

	err = shared.ForEachChild(p, func(name string) error {
		var err error
		switch name {
		case "url":
			image.URL, err = shared.ParseTextURL(p)
		case "title":
			image.Title, err = shared.ParseText(p)
		case "link":
			image.Link, err = shared.ParseTextURL(p)
		case "width":
			image.Width, err = shared.ParseText(p)
		case "height":
			image.Height, err = shared.ParseText(p)
		case "description":
			image.Description, err = shared.ParseText(p)
		default:
			err = p.Skip()
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if err = p.Expect(xpp.EndTag, "image"); err != nil {
		return nil, err
	}

	return image, nil
}

func (rp *Parser) parseGUID(p *xpp.Parser) (guid *GUID, err error) {
	if err = p.Expect(xpp.StartTag, "guid"); err != nil {
		return nil, err
	}

	guid = &GUID{}
	// The RSS 2.0 attribute is "isPermaLink" (note the capital L). XML
	// attribute names are case sensitive, so read that spelling first and fall
	// back to the lowercase form some feeds use. When absent it defaults to
	// true per the spec.
	isPermalink := p.Attribute("isPermaLink")
	if isPermalink == "" {
		isPermalink = p.Attribute("isPermalink")
	}
	if isPermalink == "" {
		isPermalink = "true"
	}
	guid.IsPermalink = isPermalink

	result, err := shared.ParseText(p)
	if err != nil {
		return
	}
	guid.Value = result

	if err = p.Expect(xpp.EndTag, "guid"); err != nil {
		return nil, err
	}

	return guid, nil
}

func (rp *Parser) parseCategory(p *xpp.Parser) (cat *Category, err error) {
	if err = p.Expect(xpp.StartTag, "category"); err != nil {
		return nil, err
	}

	cat = &Category{}
	cat.Domain = p.Attribute("domain")

	result, err := shared.ParseText(p)
	if err != nil {
		return nil, err
	}

	cat.Value = result

	if err = p.Expect(xpp.EndTag, "category"); err != nil {
		return nil, err
	}
	return cat, nil
}

func (rp *Parser) parseTextInput(p *xpp.Parser) (*TextInput, error) {
	if err := p.Expect(xpp.StartTag, "textinput"); err != nil {
		return nil, err
	}

	ti := &TextInput{}

	err := shared.ForEachChild(p, func(name string) error {
		var err error
		switch name {
		case "title":
			ti.Title, err = shared.ParseText(p)
		case "description":
			ti.Description, err = shared.ParseText(p)
		case "name":
			ti.Name, err = shared.ParseText(p)
		case "link":
			ti.Link, err = shared.ParseText(p)
		default:
			err = p.Skip()
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, "textinput"); err != nil {
		return nil, err
	}

	return ti, nil
}

func (rp *Parser) parseSkipHours(p *xpp.Parser) ([]string, error) {
	if err := p.Expect(xpp.StartTag, "skiphours"); err != nil {
		return nil, err
	}

	hours := []string{}

	err := shared.ForEachChild(p, func(name string) error {
		if name != "hour" {
			return p.Skip()
		}
		hour, err := shared.ParseText(p)
		if err == nil {
			hours = append(hours, hour)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, "skiphours"); err != nil {
		return nil, err
	}

	return hours, nil
}

func (rp *Parser) parseSkipDays(p *xpp.Parser) ([]string, error) {
	if err := p.Expect(xpp.StartTag, "skipdays"); err != nil {
		return nil, err
	}

	days := []string{}

	err := shared.ForEachChild(p, func(name string) error {
		if name != "day" {
			return p.Skip()
		}
		day, err := shared.ParseText(p)
		if err == nil {
			days = append(days, day)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, "skipdays"); err != nil {
		return nil, err
	}

	return days, nil
}

func (rp *Parser) parseCloud(p *xpp.Parser) (*Cloud, error) {
	if err := p.Expect(xpp.StartTag, "cloud"); err != nil {
		return nil, err
	}

	cloud := &Cloud{}
	cloud.Domain = p.Attribute("domain")
	cloud.Port = p.Attribute("port")
	cloud.Path = p.Attribute("path")
	cloud.RegisterProcedure = p.Attribute("registerProcedure")
	cloud.Protocol = p.Attribute("protocol")

	// The spec defines <cloud> as an empty element, but feeds exist with text
	// or child elements inside it. Skip to the matching end tag rather than
	// assuming the next tag is it.
	if err := p.Skip(); err != nil {
		return nil, err
	}

	if err := p.Expect(xpp.EndTag, "cloud"); err != nil {
		return nil, err
	}

	return cloud, nil
}

func (rp *Parser) parseVersion(p *xpp.Parser) (ver string) {
	name := strings.ToLower(p.Name())
	if name == "rss" {
		ver = p.Attribute("version")
	} else if name == "rdf" {
		ns := p.Attribute("xmlns")
		if ns == "http://channel.netscape.com/rdf/simple/0.9/" ||
			ns == "http://my.netscape.com/rdf/simple/0.9/" {
			ver = "0.9"
		} else if ns == "http://purl.org/rss/1.0/" {
			ver = "1.0"
		}
	}
	return
}
