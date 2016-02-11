package rss

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed/feed"
	"github.com/mmcdole/gofeed/shared"
	"github.com/mmcdole/goxpp"
)

type Parser struct {
	shared.BaseParser
}

func (rp *Parser) ParseFeed(feed string) (*Feed, error) {
	p := xpp.NewXMLPullParser(strings.NewReader(feed), false)

	_, err := p.NextTag()
	if err != nil {
		return nil, err
	}

	return rp.parseRoot(p)
}

func (rp *Parser) parseRoot(p *xpp.XMLPullParser) (*Feed, error) {
	rssErr := rp.Expect(p, xpp.StartTag, "rss")
	rdfErr := rp.Expect(p, xpp.StartTag, "rdf")
	if rssErr != nil && rdfErr != nil {
		return nil, fmt.Errorf("%s or %s", rssErr.Error(), rdfErr.Error())
	}

	// Items found in feed root
	var channel *Feed
	var textinput *TextInput
	var image *Image
	items := []*Item{}

	ver := rp.parseVersion(p)

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {

			// Skip any extensions found in the feed root.
			if rp.IsExtension(p) {
				p.Skip()
				continue
			}

			name := strings.ToLower(p.Name)

			if name == "channel" {
				channel, err = rp.parseChannel(p)
				if err != nil {
					fmt.Printf("error in channel: %s\n", err)
					return nil, err
				}
			} else if name == "item" {
				item, err := rp.parseItem(p)
				if err != nil {
					return nil, err
				}
				items = append(items, item)
			} else if name == "textinput" {
				textinput, err = rp.parseTextInput(p)
				if err != nil {
					return nil, err
				}
			} else if name == "image" {
				image, err = rp.parseImage(p)
				if err != nil {
					return nil, err
				}
			} else {
				p.Skip()
			}
		}
	}

	rssErr = rp.Expect(p, xpp.EndTag, "rss")
	rdfErr = rp.Expect(p, xpp.EndTag, "rdf")
	if rssErr != nil && rdfErr != nil {
		return nil, fmt.Errorf("%s or %s", rssErr.Error(), rdfErr.Error())
	}

	if channel == nil {
		channel = &Feed{}
		channel.Items = []*Item{}
		channel.Extensions = feed.FeedExtensions{}
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

func (rp *Parser) parseChannel(p *xpp.XMLPullParser) (rss *Feed, err error) {

	if err = rp.Expect(p, xpp.StartTag, "channel"); err != nil {
		return nil, err
	}

	rss = &Feed{}
	rss.Items = []*Item{}
	rss.Extensions = feed.FeedExtensions{}
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

			name := strings.ToLower(p.Name)

			if rp.IsExtension(p) {
				ext, err := rp.ParseExtension(rss.Extensions, p)
				if err != nil {
					return nil, err
				}
				rss.Extensions = ext
			} else if name == "title" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.Title = result
			} else if name == "description" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.Description = result
			} else if name == "link" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.Link = result
			} else if name == "language" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.Language = result
			} else if name == "copyright" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.Copyright = result
			} else if name == "managingeditor" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.ManagingEditor = result
			} else if name == "webmaster" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.WebMaster = result
			} else if name == "pubdate" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.PubDate = result
				date, err := rp.ParseDate(result)
				if err == nil {
					utcDate := date.UTC()
					rss.PubDateParsed = &utcDate
				}
			} else if name == "lastbuilddate" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.LastBuildDate = result
				date, err := rp.ParseDate(result)
				if err == nil {
					utcDate := date.UTC()
					rss.LastBuildDateParsed = &utcDate
				}
			} else if name == "generator" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.Generator = result
			} else if name == "docs" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.Docs = result
			} else if name == "ttl" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.TTL = result
			} else if name == "rating" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				rss.Rating = result
			} else if name == "skiphours" {
				result, err := rp.parseSkipHours(p)
				if err != nil {
					return nil, err
				}
				rss.SkipHours = result
			} else if name == "skipdays" {
				result, err := rp.parseSkipDays(p)
				if err != nil {
					return nil, err
				}
				rss.SkipDays = result
			} else if name == "item" {
				result, err := rp.parseItem(p)
				if err != nil {
					return nil, err
				}
				rss.Items = append(rss.Items, result)
			} else if name == "cloud" {
				result, err := rp.parseCloud(p)
				if err != nil {
					return nil, err
				}
				rss.Cloud = result
			} else if name == "category" {
				result, err := rp.parseCategory(p)
				if err != nil {
					return nil, err
				}
				categories = append(categories, result)
			} else if name == "image" {
				result, err := rp.parseImage(p)
				if err != nil {
					return nil, err
				}
				rss.Image = result
			} else if name == "textinput" {
				result, err := rp.parseTextInput(p)
				if err != nil {
					return nil, err
				}
				rss.TextInput = result
			} else {
				// Skip element as it isn't an extension and not
				// part of the spec
				p.Skip()
			}
		}
	}

	if err = rp.Expect(p, xpp.EndTag, "channel"); err != nil {
		return nil, err
	}

	if len(categories) > 0 {
		rss.Categories = categories
	}

	return rss, nil
}

func (rp *Parser) parseItem(p *xpp.XMLPullParser) (item *Item, err error) {

	if err = rp.Expect(p, xpp.StartTag, "item"); err != nil {
		return nil, err
	}

	item = &Item{}
	item.Extensions = feed.FeedExtensions{}

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

			name := strings.ToLower(p.Name)

			if rp.IsExtension(p) {
				ext, err := rp.ParseExtension(item.Extensions, p)
				if err != nil {
					return nil, err
				}
				item.Extensions = ext
			} else if name == "title" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				item.Title = result
			} else if name == "description" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				item.Description = result
			} else if name == "link" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				item.Link = result
			} else if name == "author" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				item.Author = result
			} else if name == "comments" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				item.Comments = result
			} else if name == "pubdate" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				item.PubDate = result
				date, err := rp.ParseDate(result)
				if err == nil {
					utcDate := date.UTC()
					item.PubDateParsed = &utcDate
				}
			} else if name == "source" {
				result, err := rp.parseSource(p)
				if err != nil {
					return nil, err
				}
				item.Source = result
			} else if name == "enclosure" {
				result, err := rp.parseEnclosure(p)
				if err != nil {
					return nil, err
				}
				item.Enclosure = result
			} else if name == "guid" {
				result, err := rp.parseGuid(p)
				if err != nil {
					return nil, err
				}
				item.Guid = result
			} else if name == "category" {
				result, err := rp.parseCategory(p)
				if err != nil {
					return nil, err
				}
				categories = append(categories, result)
			} else {
				// Skip any elements not part of the item spec
				p.Skip()
			}
		}
	}

	if len(categories) > 0 {
		item.Categories = categories
	}

	if err = rp.Expect(p, xpp.EndTag, "item"); err != nil {
		return nil, err
	}

	return item, nil
}

func (rp *Parser) parseSource(p *xpp.XMLPullParser) (source *Source, err error) {
	if err = rp.Expect(p, xpp.StartTag, "source"); err != nil {
		return nil, err
	}

	source = &Source{}
	source.URL = p.Attribute("url")

	result, err := rp.ParseText(p)
	if err != nil {
		return source, err
	}
	source.Title = result

	if err = rp.Expect(p, xpp.EndTag, "source"); err != nil {
		return nil, err
	}
	return source, nil
}

func (rp *Parser) parseEnclosure(p *xpp.XMLPullParser) (enclosure *Enclosure, err error) {
	if err = rp.Expect(p, xpp.StartTag, "enclosure"); err != nil {
		return nil, err
	}

	enclosure = &Enclosure{}
	enclosure.URL = p.Attribute("url")
	enclosure.Length = p.Attribute("length")
	enclosure.Type = p.Attribute("type")

	p.NextTag()

	if err = rp.Expect(p, xpp.EndTag, "enclosure"); err != nil {
		return nil, err
	}

	return enclosure, nil
}

func (rp *Parser) parseImage(p *xpp.XMLPullParser) (image *Image, err error) {
	if err = rp.Expect(p, xpp.StartTag, "image"); err != nil {
		return nil, err
	}

	image = &Image{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return image, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			name := strings.ToLower(p.Name)

			if name == "url" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				image.URL = result
			} else if name == "title" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				image.Title = result
			} else if name == "link" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				image.Link = result
			} else if name == "width" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				image.Width = result
			} else if name == "height" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				image.Height = result
			} else if name == "description" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				image.Description = result
			} else {
				p.Skip()
			}
		}
	}

	if err = rp.Expect(p, xpp.EndTag, "image"); err != nil {
		return nil, err
	}

	return image, nil
}

func (rp *Parser) parseGuid(p *xpp.XMLPullParser) (guid *Guid, err error) {
	if err = rp.Expect(p, xpp.StartTag, "guid"); err != nil {
		return nil, err
	}

	guid = &Guid{}
	guid.IsPermalink = p.Attribute("isPermalink")

	result, err := rp.ParseText(p)
	if err != nil {
		return
	}
	guid.Value = result

	if err = rp.Expect(p, xpp.EndTag, "guid"); err != nil {
		return nil, err
	}

	return guid, nil
}

func (rp *Parser) parseCategory(p *xpp.XMLPullParser) (cat *Category, err error) {

	if err = rp.Expect(p, xpp.StartTag, "category"); err != nil {
		return nil, err
	}

	cat = &Category{}
	cat.Domain = p.Attribute("domain")

	result, err := rp.ParseText(p)
	if err != nil {
		return nil, err
	}

	cat.Value = result

	if err = rp.Expect(p, xpp.EndTag, "category"); err != nil {
		return nil, err
	}
	return cat, nil
}

func (rp *Parser) parseTextInput(p *xpp.XMLPullParser) (*TextInput, error) {
	if err := rp.Expect(p, xpp.StartTag, "textinput"); err != nil {
		return nil, err
	}

	ti := &TextInput{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			name := strings.ToLower(p.Name)

			if name == "title" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				ti.Title = result
			} else if name == "description" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				ti.Description = result
			} else if name == "name" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				ti.Name = result
			} else if name == "link" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				ti.Link = result
			} else {
				p.Skip()
			}
		}
	}

	if err := rp.Expect(p, xpp.EndTag, "textinput"); err != nil {
		return nil, err
	}

	return ti, nil
}

func (rp *Parser) parseSkipHours(p *xpp.XMLPullParser) ([]string, error) {
	if err := rp.Expect(p, xpp.StartTag, "skiphours"); err != nil {
		return nil, err
	}

	hours := []string{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			name := strings.ToLower(p.Name)
			if name == "hour" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				hours = append(hours, result)
			} else {
				p.Skip()
			}
		}
	}

	if err := rp.Expect(p, xpp.EndTag, "skiphours"); err != nil {
		return nil, err
	}

	return hours, nil
}

func (rp *Parser) parseSkipDays(p *xpp.XMLPullParser) ([]string, error) {
	if err := rp.Expect(p, xpp.StartTag, "skipdays"); err != nil {
		return nil, err
	}

	days := []string{}

	for {
		tok, err := p.NextTag()
		if err != nil {
			return nil, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			name := strings.ToLower(p.Name)
			if name == "day" {
				result, err := rp.ParseText(p)
				if err != nil {
					return nil, err
				}
				days = append(days, result)
			} else {
				p.Skip()
			}
		}
	}

	if err := rp.Expect(p, xpp.EndTag, "skipdays"); err != nil {
		return nil, err
	}

	return days, nil
}

func (rp *Parser) parseCloud(p *xpp.XMLPullParser) (*Cloud, error) {
	if err := rp.Expect(p, xpp.StartTag, "cloud"); err != nil {
		return nil, err
	}

	cloud := &Cloud{}
	cloud.Domain = p.Attribute("domain")
	cloud.Port = p.Attribute("port")
	cloud.Path = p.Attribute("path")
	cloud.RegisterProcedure = p.Attribute("registerProcedure")
	cloud.Protocol = p.Attribute("protocol")

	p.NextTag()

	if err := rp.Expect(p, xpp.EndTag, "cloud"); err != nil {
		return nil, err
	}

	return cloud, nil
}

func (rp *Parser) parseVersion(p *xpp.XMLPullParser) (ver string) {
	name := strings.ToLower(p.Name)
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
