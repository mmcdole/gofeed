package gofeed

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed/atom"
	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/mmcdole/gofeed/internal/shared"
	"github.com/mmcdole/gofeed/json"
	"github.com/mmcdole/gofeed/rss"
	"golang.org/x/net/html"
)

// Translator converts a particular feed (atom.Feed or rss.Feed of json.Feed)
// into the generic Feed struct
type Translator interface {
	Translate(feed interface{}) (*Feed, error)
}

// DefaultRSSTranslator converts an rss.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between rss.Feed -> Feed
// for each of the fields in Feed.
type DefaultRSSTranslator struct{}

// Translate converts an RSS feed into the universal
// feed type.
func (t *DefaultRSSTranslator) Translate(feed interface{}) (*Feed, error) {
	rss, found := feed.(*rss.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *rss.Feed")
	}

	result := &Feed{}
	result.Title = t.translateFeedTitle(rss)
	result.Description = t.translateFeedDescription(rss)
	result.Link = t.translateFeedLink(rss)
	result.Links = t.translateFeedLinks(rss)
	result.FeedLink = t.translateFeedFeedLink(rss)
	result.Updated = t.translateFeedUpdated(rss)
	result.UpdatedParsed = t.translateFeedUpdatedParsed(rss)
	result.Published = t.translateFeedPublished(rss)
	result.PublishedParsed = t.translateFeedPublishedParsed(rss)
	result.Author = t.translateFeedAuthor(rss)
	result.Authors = t.translateFeedAuthors(rss)
	result.Language = t.translateFeedLanguage(rss)
	result.Image = t.translateFeedImage(rss)
	result.Copyright = t.translateFeedCopyright(rss)
	result.Generator = t.translateFeedGenerator(rss)
	result.Categories = t.translateFeedCategories(rss)
	result.Items = t.translateFeedItems(rss)
	result.ITunesExt = rss.ITunesExt
	result.DublinCoreExt = rss.DublinCoreExt
	result.Extensions = rss.Extensions
	result.FeedVersion = rss.Version
	result.FeedType = "rss"
	return result, nil
}

func (t *DefaultRSSTranslator) translateFeedItem(rssItem *rss.Item) (item *Item) {
	item = &Item{}
	item.Title = t.translateItemTitle(rssItem)
	item.Description = t.translateItemDescription(rssItem)
	item.Content = t.translateItemContent(rssItem)
	item.Link = t.translateItemLink(rssItem)
	item.Links = t.translateItemLinks(rssItem)
	item.Published = t.translateItemPublished(rssItem)
	item.PublishedParsed = t.translateItemPublishedParsed(rssItem)
	item.Author = t.translateItemAuthor(rssItem)
	item.Authors = t.translateItemAuthors(rssItem)
	item.GUID = t.translateItemGUID(rssItem)
	item.Image = t.translateItemImage(rssItem)
	item.Categories = t.translateItemCategories(rssItem)
	item.Enclosures = t.translateItemEnclosures(rssItem)
	item.DublinCoreExt = rssItem.DublinCoreExt
	item.ITunesExt = rssItem.ITunesExt
	item.Extensions = rssItem.Extensions
	return
}

func (t *DefaultRSSTranslator) translateFeedTitle(rss *rss.Feed) (title string) {
	if rss.Title != "" {
		title = rss.Title
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Title != nil {
		title = t.firstEntry(rss.DublinCoreExt.Title)
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedDescription(rss *rss.Feed) (desc string) {
	if rss.Description != "" {
		desc = rss.Description
	} else if rss.ITunesExt != nil && rss.ITunesExt.Summary != "" {
		desc = rss.ITunesExt.Summary
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedLink(rss *rss.Feed) (link string) {
	if rss.Link != "" {
		link = rss.Link
	} else if rss.ITunesExt != nil && rss.ITunesExt.Subtitle != "" {
		link = rss.ITunesExt.Subtitle
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedFeedLink(rss *rss.Feed) (link string) {
	atomExtensions := t.extensionsForKeys([]string{"atom", "atom10", "atom03"}, rss.Extensions)
	for _, ex := range atomExtensions {
		if links, ok := ex["link"]; ok {
			for _, l := range links {
				if l.Attrs["rel"] == "self" {
					link = l.Attrs["href"]
				}
			}
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedLinks(rss *rss.Feed) (links []string) {
	if len(rss.Links) > 0 {
		links = append(links, rss.Links...)
	}
	atomExtensions := t.extensionsForKeys([]string{"atom", "atom10", "atom03"}, rss.Extensions)
	for _, ex := range atomExtensions {
		if lks, ok := ex["link"]; ok {
			for _, l := range lks {
				if l.Attrs["rel"] == "" || l.Attrs["rel"] == "alternate" || l.Attrs["rel"] == "self" {
					links = append(links, l.Attrs["href"])
				}
			}
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedUpdated(rss *rss.Feed) (updated string) {
	if rss.LastBuildDate != "" {
		updated = rss.LastBuildDate
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Date != nil {
		updated = t.firstEntry(rss.DublinCoreExt.Date)
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedUpdatedParsed(rss *rss.Feed) (updated *time.Time) {
	if rss.LastBuildDateParsed != nil {
		updated = rss.LastBuildDateParsed
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Date != nil {
		dateText := t.firstEntry(rss.DublinCoreExt.Date)
		date, err := shared.ParseDate(dateText)
		if err == nil {
			updated = &date
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedPublished(rss *rss.Feed) (published string) {
	return rss.PubDate
}

func (t *DefaultRSSTranslator) translateFeedPublishedParsed(rss *rss.Feed) (published *time.Time) {
	return rss.PubDateParsed
}

func (t *DefaultRSSTranslator) translateFeedAuthor(rss *rss.Feed) (author *Person) {
	if rss.ManagingEditor != "" {
		name, address := shared.ParseNameAddress(rss.ManagingEditor)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rss.WebMaster != "" {
		name, address := shared.ParseNameAddress(rss.WebMaster)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Author != nil {
		dcAuthor := t.firstEntry(rss.DublinCoreExt.Author)
		name, address := shared.ParseNameAddress(dcAuthor)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Creator != nil {
		dcCreator := t.firstEntry(rss.DublinCoreExt.Creator)
		name, address := shared.ParseNameAddress(dcCreator)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rss.ITunesExt != nil && rss.ITunesExt.Author != "" {
		name, address := shared.ParseNameAddress(rss.ITunesExt.Author)
		author = &Person{}
		author.Name = name
		author.Email = address
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedAuthors(rss *rss.Feed) (authors []*Person) {
	if author := t.translateFeedAuthor(rss); author != nil {
		authors = []*Person{author}
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedLanguage(rss *rss.Feed) (language string) {
	if rss.Language != "" {
		language = rss.Language
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Language != nil {
		language = t.firstEntry(rss.DublinCoreExt.Language)
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedImage(rss *rss.Feed) *Image {
	if rss.Image != nil {
		return &Image{
			Title: rss.Image.Title,
			URL:   rss.Image.URL,
		}
	}
	if rss.ITunesExt != nil && rss.ITunesExt.Image != "" {
		return &Image{URL: rss.ITunesExt.Image}
	}
	if media, ok := rss.Extensions["media"]; ok {
		if content, ok := media["content"]; ok {
			for _, c := range content {
				if strings.HasPrefix(c.Attrs["type"], "image/") || c.Attrs["medium"] == "image" {
					return &Image{URL: c.Attrs["url"]}
				}
			}
		}
	}
	return firstImageFromHtmlDocument(rss.Description)
}

func (t *DefaultRSSTranslator) translateFeedCopyright(rss *rss.Feed) (rights string) {
	if rss.Copyright != "" {
		rights = rss.Copyright
	} else if rss.DublinCoreExt != nil && rss.DublinCoreExt.Rights != nil {
		rights = t.firstEntry(rss.DublinCoreExt.Rights)
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedGenerator(rss *rss.Feed) (generator string) {
	return rss.Generator
}

func (t *DefaultRSSTranslator) translateFeedCategories(rss *rss.Feed) (categories []string) {
	var cats []string
	if rss.Categories != nil {
		cats = make([]string, 0, len(rss.Categories))
		for _, c := range rss.Categories {
			cats = append(cats, c.Value)
		}
	}

	if rss.ITunesExt != nil && rss.ITunesExt.Keywords != "" {
		keywords := strings.Split(rss.ITunesExt.Keywords, ",")
		cats = append(cats, keywords...)
	}

	if rss.ITunesExt != nil && rss.ITunesExt.Categories != nil {
		for _, c := range rss.ITunesExt.Categories {
			cats = append(cats, c.Text)
			if c.Subcategory != nil {
				cats = append(cats, c.Subcategory.Text)
			}
		}
	}

	if rss.DublinCoreExt != nil && rss.DublinCoreExt.Subject != nil {
		cats = append(cats, rss.DublinCoreExt.Subject...)
	}

	if len(cats) > 0 {
		categories = cats
	}

	return
}

func (t *DefaultRSSTranslator) translateFeedItems(rss *rss.Feed) (items []*Item) {
	items = make([]*Item, 0, len(rss.Items))
	for _, i := range rss.Items {
		items = append(items, t.translateFeedItem(i))
	}
	return
}

func (t *DefaultRSSTranslator) translateItemTitle(rssItem *rss.Item) (title string) {
	if rssItem.Title != "" {
		title = rssItem.Title
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Title != nil {
		title = t.firstEntry(rssItem.DublinCoreExt.Title)
	}
	return
}

func (t *DefaultRSSTranslator) translateItemDescription(rssItem *rss.Item) (desc string) {
	if rssItem.Description != "" {
		desc = rssItem.Description
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Description != nil {
		desc = t.firstEntry(rssItem.DublinCoreExt.Description)
	} else if rssItem.ITunesExt != nil && rssItem.ITunesExt.Summary != "" {
		desc = rssItem.ITunesExt.Summary
	}
	return
}

func (t *DefaultRSSTranslator) translateItemContent(rssItem *rss.Item) (content string) {
	return rssItem.Content
}

func (t *DefaultRSSTranslator) translateItemLink(rssItem *rss.Item) (link string) {
	return rssItem.Link
}

func (t *DefaultRSSTranslator) translateItemLinks(rssItem *rss.Item) (links []string) {
	if len(rssItem.Links) > 0 {
		links = append(links, rssItem.Links...)
	}
	return links
}

func (t *DefaultRSSTranslator) translateItemUpdated(rssItem *rss.Item) (updated string) {
	if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		updated = t.firstEntry(rssItem.DublinCoreExt.Date)
	}
	return updated
}

func (t *DefaultRSSTranslator) translateItemUpdatedParsed(rssItem *rss.Item) (updated *time.Time) {
	if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		updatedText := t.firstEntry(rssItem.DublinCoreExt.Date)
		updatedDate, err := shared.ParseDate(updatedText)
		if err == nil {
			updated = &updatedDate
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateItemPublished(rssItem *rss.Item) (pubDate string) {
	if rssItem.PubDate != "" {
		return rssItem.PubDate
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		return t.firstEntry(rssItem.DublinCoreExt.Date)
	}
	return
}

func (t *DefaultRSSTranslator) translateItemPublishedParsed(rssItem *rss.Item) (pubDate *time.Time) {
	if rssItem.PubDateParsed != nil {
		return rssItem.PubDateParsed
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		pubDateText := t.firstEntry(rssItem.DublinCoreExt.Date)
		pubDateParsed, err := shared.ParseDate(pubDateText)
		if err == nil {
			pubDate = &pubDateParsed
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateItemAuthor(rssItem *rss.Item) (author *Person) {
	if rssItem.Author != "" {
		name, address := shared.ParseNameAddress(rssItem.Author)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Author != nil {
		dcAuthor := t.firstEntry(rssItem.DublinCoreExt.Author)
		name, address := shared.ParseNameAddress(dcAuthor)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Creator != nil {
		dcCreator := t.firstEntry(rssItem.DublinCoreExt.Creator)
		name, address := shared.ParseNameAddress(dcCreator)
		author = &Person{}
		author.Name = name
		author.Email = address
	} else if rssItem.ITunesExt != nil && rssItem.ITunesExt.Author != "" {
		name, address := shared.ParseNameAddress(rssItem.ITunesExt.Author)
		author = &Person{}
		author.Name = name
		author.Email = address
	}
	return
}

func (t *DefaultRSSTranslator) translateItemAuthors(rssItem *rss.Item) (authors []*Person) {
	if author := t.translateItemAuthor(rssItem); author != nil {
		authors = []*Person{author}
	}
	return
}

func (t *DefaultRSSTranslator) translateItemGUID(rssItem *rss.Item) (guid string) {
	if rssItem.GUID != nil {
		guid = rssItem.GUID.Value
	}
	return
}

func (t *DefaultRSSTranslator) translateItemImage(rssItem *rss.Item) *Image {
	if rssItem.ITunesExt != nil && rssItem.ITunesExt.Image != "" {
		return &Image{URL: rssItem.ITunesExt.Image}
	}
	if media, ok := rssItem.Extensions["media"]; ok {
		if content, ok := media["content"]; ok {
			for _, c := range content {
				if strings.Contains(c.Attrs["type"], "image") || strings.Contains(c.Attrs["medium"], "image") {
					return &Image{URL: c.Attrs["url"]}
				}
			}
		}
	}
	for _, enc := range rssItem.Enclosures {
		if strings.HasPrefix(enc.Type, "image/") {
			return &Image{URL: enc.URL}
		}
	}
	if img := firstImageFromHtmlDocument(rssItem.Content); img != nil {
		return img
	}
	if img := firstImageFromHtmlDocument(rssItem.Description); img != nil {
		return img
	}
	return nil
}

func firstImageFromHtmlDocument(document string) *Image {
	if doc, err := html.Parse(bytes.NewBufferString(document)); err == nil {
		doc := goquery.NewDocumentFromNode(doc)
		for _, node := range doc.FindMatcher(goquery.Single("img[src]")).Nodes {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					return &Image{
						URL: attr.Val,
					}
				}
			}
		}
	}
	return nil
}

func (t *DefaultRSSTranslator) translateItemCategories(rssItem *rss.Item) (categories []string) {
	var cats []string
	if rssItem.Categories != nil {
		cats = make([]string, 0, len(rssItem.Categories))
		for _, c := range rssItem.Categories {
			cats = append(cats, c.Value)
		}
	}

	if rssItem.ITunesExt != nil && rssItem.ITunesExt.Keywords != "" {
		keywords := strings.Split(rssItem.ITunesExt.Keywords, ",")
		cats = append(cats, keywords...)
	}

	if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Subject != nil {
		cats = append(cats, rssItem.DublinCoreExt.Subject...)
	}

	if len(cats) > 0 {
		categories = cats
	}

	return
}

func (t *DefaultRSSTranslator) translateItemEnclosures(rssItem *rss.Item) (enclosures []*Enclosure) {
	if rssItem.Enclosures != nil && len(rssItem.Enclosures) > 0 {
		// Accumulate the enclosures
		for _, enc := range rssItem.Enclosures {
			e := &Enclosure{}
			e.URL = enc.URL
			e.Type = enc.Type
			e.Length = enc.Length
			enclosures = append(enclosures, e)
		}
	}

	if len(enclosures) == 0 {
		enclosures = nil
	}

	return
}

func (t *DefaultRSSTranslator) extensionsForKeys(keys []string, extensions ext.Extensions) (matches []map[string][]ext.Extension) {
	matches = []map[string][]ext.Extension{}

	if extensions == nil {
		return
	}

	for _, key := range keys {
		if match, ok := extensions[key]; ok {
			matches = append(matches, match)
		}
	}
	return
}

func (t *DefaultRSSTranslator) firstEntry(entries []string) (value string) {
	if entries == nil {
		return
	}

	if len(entries) == 0 {
		return
	}

	return entries[0]
}

// DefaultAtomTranslator converts an atom.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between atom.Feed -> Feed
// for each of the fields in Feed.
type DefaultAtomTranslator struct{}

// Translate converts an Atom feed into the universal
// feed type.
func (t *DefaultAtomTranslator) Translate(feed interface{}) (*Feed, error) {
	atom, found := feed.(*atom.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *atom.Feed")
	}

	result := &Feed{}
	result.Title = t.translateFeedTitle(atom)
	result.Description = t.translateFeedDescription(atom)
	result.Link = t.translateFeedLink(atom)
	result.FeedLink = t.translateFeedFeedLink(atom)
	result.Links = t.translateFeedLinks(atom)
	result.Updated = t.translateFeedUpdated(atom)
	result.UpdatedParsed = t.translateFeedUpdatedParsed(atom)
	result.Author = t.translateFeedAuthor(atom)
	result.Authors = t.translateFeedAuthors(atom)
	result.Language = t.translateFeedLanguage(atom)
	result.Image = t.translateFeedImage(atom)
	result.Copyright = t.translateFeedCopyright(atom)
	result.Categories = t.translateFeedCategories(atom)
	result.Generator = t.translateFeedGenerator(atom)
	result.Items = t.translateFeedItems(atom)
	result.Extensions = atom.Extensions
	result.FeedVersion = atom.Version
	result.FeedType = "atom"
	return result, nil
}

func (t *DefaultAtomTranslator) translateFeedItem(entry *atom.Entry) (item *Item) {
	item = &Item{}
	item.Title = t.translateItemTitle(entry)
	item.Description = t.translateItemDescription(entry)
	item.Content = t.translateItemContent(entry)
	item.Link = t.translateItemLink(entry)
	item.Links = t.translateItemLinks(entry)
	item.Updated = t.translateItemUpdated(entry)
	item.UpdatedParsed = t.translateItemUpdatedParsed(entry)
	item.Published = t.translateItemPublished(entry)
	item.PublishedParsed = t.translateItemPublishedParsed(entry)
	item.Author = t.translateItemAuthor(entry)
	item.Authors = t.translateItemAuthors(entry)
	item.GUID = t.translateItemGUID(entry)
	item.Image = t.translateItemImage(entry)
	item.Categories = t.translateItemCategories(entry)
	item.Enclosures = t.translateItemEnclosures(entry)
	item.Extensions = entry.Extensions
	return
}

func (t *DefaultAtomTranslator) translateFeedTitle(atom *atom.Feed) (title string) {
	return atom.Title
}

func (t *DefaultAtomTranslator) translateFeedDescription(atom *atom.Feed) (desc string) {
	return atom.Subtitle
}

func (t *DefaultAtomTranslator) translateFeedLink(atom *atom.Feed) (link string) {
	l := t.firstLinkWithType("alternate", atom.Links)
	if l != nil {
		link = l.Href
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedFeedLink(atom *atom.Feed) (link string) {
	feedLink := t.firstLinkWithType("self", atom.Links)
	if feedLink != nil {
		link = feedLink.Href
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedLinks(atom *atom.Feed) (links []string) {
	for _, l := range atom.Links {
		if l.Rel == "" || l.Rel == "alternate" || l.Rel == "self" {
			links = append(links, l.Href)
		}
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedUpdated(atom *atom.Feed) (updated string) {
	return atom.Updated
}

func (t *DefaultAtomTranslator) translateFeedUpdatedParsed(atom *atom.Feed) (updated *time.Time) {
	return atom.UpdatedParsed
}

func (t *DefaultAtomTranslator) translateFeedAuthor(atom *atom.Feed) (author *Person) {
	a := t.firstPerson(atom.Authors)
	if a != nil {
		feedAuthor := Person{}
		feedAuthor.Name = a.Name
		feedAuthor.Email = a.Email
		author = &feedAuthor
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedAuthors(atom *atom.Feed) (authors []*Person) {
	if atom.Authors != nil {
		authors = make([]*Person, 0, len(atom.Authors))

		for _, a := range atom.Authors {
			authors = append(authors, &Person{
				Name:  a.Name,
				Email: a.Email,
			})
		}
	}

	return
}

func (t *DefaultAtomTranslator) translateFeedLanguage(atom *atom.Feed) (language string) {
	return atom.Language
}

func (t *DefaultAtomTranslator) translateFeedImage(atom *atom.Feed) (image *Image) {
	if atom.Logo != "" {
		feedImage := Image{}
		feedImage.URL = atom.Logo
		image = &feedImage
	} else if atom.Icon != "" {
		feedImage := Image{}
		feedImage.URL = atom.Icon
		image = &feedImage
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedCopyright(atom *atom.Feed) (rights string) {
	return atom.Rights
}

func (t *DefaultAtomTranslator) translateFeedGenerator(atom *atom.Feed) (generator string) {
	if atom.Generator != nil {
		if atom.Generator.Value != "" {
			generator += atom.Generator.Value
		}
		if atom.Generator.Version != "" {
			generator += " v" + atom.Generator.Version
		}
		if atom.Generator.URI != "" {
			generator += " " + atom.Generator.URI
		}
		generator = strings.TrimSpace(generator)
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedCategories(atom *atom.Feed) (categories []string) {
	if atom.Categories != nil {
		categories = make([]string, 0, len(atom.Categories))
		for _, c := range atom.Categories {
			if c.Label != "" {
				categories = append(categories, c.Label)
			} else {
				categories = append(categories, c.Term)
			}
		}
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedItems(atom *atom.Feed) (items []*Item) {
	items = make([]*Item, 0, len(atom.Entries))
	for _, entry := range atom.Entries {
		items = append(items, t.translateFeedItem(entry))
	}
	return
}

func (t *DefaultAtomTranslator) translateItemTitle(entry *atom.Entry) (title string) {
	return entry.Title
}

func (t *DefaultAtomTranslator) translateItemDescription(entry *atom.Entry) (desc string) {
	return entry.Summary
}

func (t *DefaultAtomTranslator) translateItemContent(entry *atom.Entry) (content string) {
	if entry.Content != nil {
		content = entry.Content.Value
	}
	return
}

func (t *DefaultAtomTranslator) translateItemLink(entry *atom.Entry) (link string) {
	l := t.firstLinkWithType("alternate", entry.Links)
	if l != nil {
		link = l.Href
	}
	return
}

func (t *DefaultAtomTranslator) translateItemLinks(entry *atom.Entry) (links []string) {
	for _, l := range entry.Links {
		if l.Rel == "" || l.Rel == "alternate" || l.Rel == "self" {
			links = append(links, l.Href)
		}
	}
	return
}

func (t *DefaultAtomTranslator) translateItemUpdated(entry *atom.Entry) (updated string) {
	return entry.Updated
}

func (t *DefaultAtomTranslator) translateItemUpdatedParsed(entry *atom.Entry) (updated *time.Time) {
	return entry.UpdatedParsed
}

func (t *DefaultAtomTranslator) translateItemPublished(entry *atom.Entry) (published string) {
	published = entry.Published
	if published == "" {
		published = entry.Updated
	}
	return
}

func (t *DefaultAtomTranslator) translateItemPublishedParsed(entry *atom.Entry) (published *time.Time) {
	published = entry.PublishedParsed
	if published == nil {
		published = entry.UpdatedParsed
	}
	return
}

func (t *DefaultAtomTranslator) translateItemAuthor(entry *atom.Entry) (author *Person) {
	a := t.firstPerson(entry.Authors)
	if a != nil {
		author = &Person{}
		author.Name = a.Name
		author.Email = a.Email
	}
	return
}

func (t *DefaultAtomTranslator) translateItemAuthors(entry *atom.Entry) (authors []*Person) {
	if entry.Authors != nil {
		authors = make([]*Person, 0, len(entry.Authors))
		for _, a := range entry.Authors {
			authors = append(authors, &Person{
				Name:  a.Name,
				Email: a.Email,
			})
		}
	}
	return
}

func (t *DefaultAtomTranslator) translateItemGUID(entry *atom.Entry) (guid string) {
	return entry.ID
}

func (t *DefaultAtomTranslator) translateItemImage(entry *atom.Entry) (image *Image) {
	return nil
}

func (t *DefaultAtomTranslator) translateItemCategories(entry *atom.Entry) (categories []string) {
	if entry.Categories != nil {
		categories = make([]string, 0, len(entry.Categories))
		for _, c := range entry.Categories {
			if c.Label != "" {
				categories = append(categories, c.Label)
			} else {
				categories = append(categories, c.Term)
			}
		}
	}
	return
}

func (t *DefaultAtomTranslator) translateItemEnclosures(entry *atom.Entry) (enclosures []*Enclosure) {
	if entry.Links != nil {
		enclosures = make([]*Enclosure, 0, len(entry.Links))
		for _, e := range entry.Links {
			if e.Rel == "enclosure" {
				enclosure := &Enclosure{}
				enclosure.URL = e.Href
				enclosure.Length = e.Length
				enclosure.Type = e.Type
				enclosures = append(enclosures, enclosure)
			}
		}

		if len(enclosures) == 0 {
			enclosures = nil
		}
	}
	return
}

func (t *DefaultAtomTranslator) firstLinkWithType(linkType string, links []*atom.Link) *atom.Link {
	if links == nil {
		return nil
	}

	for _, link := range links {
		if link.Rel == linkType {
			return link
		}
	}
	return nil
}

func (t *DefaultAtomTranslator) firstPerson(persons []*atom.Person) (person *atom.Person) {
	if persons == nil || len(persons) == 0 {
		return
	}

	person = persons[0]
	return
}

// DefaultJSONTranslator converts an json.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between json.Feed -> Feed
// for each of the fields in Feed.
type DefaultJSONTranslator struct{}

// Translate converts an JSON feed into the universal
// feed type.
func (t *DefaultJSONTranslator) Translate(feed interface{}) (*Feed, error) {
	json, found := feed.(*json.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *json.Feed")
	}

	result := &Feed{}
	result.FeedVersion = json.Version
	result.Title = t.translateFeedTitle(json)
	result.Link = t.translateFeedLink(json)
	result.FeedLink = t.translateFeedFeedLink(json)
	result.Links = t.translateFeedLinks(json)
	result.Description = t.translateFeedDescription(json)
	result.Image = t.translateFeedImage(json)
	result.Author = t.translateFeedAuthor(json)
	result.Authors = t.translateFeedAuthors(json)
	result.Language = t.translateFeedLanguage(json)
	result.Items = t.translateFeedItems(json)
	result.Updated = t.translateFeedUpdated(json)
	result.UpdatedParsed = t.translateFeedUpdatedParsed(json)
	result.Published = t.translateFeedPublished(json)
	result.PublishedParsed = t.translateFeedPublishedParsed(json)
	result.FeedType = "json"
	// TODO UserComment is missing in global Feed
	// TODO NextURL is missing in global Feed
	// TODO Favicon is missing in global Feed
	// TODO Exipred is missing in global Feed
	// TODO Hubs is not supported in json.Feed
	// TODO Extensions is not supported in json.Feed
	return result, nil
}

func (t *DefaultJSONTranslator) translateFeedItem(jsonItem *json.Item) (item *Item) {
	item = &Item{}
	item.GUID = t.translateItemGUID(jsonItem)
	item.Link = t.translateItemLink(jsonItem)
	item.Links = t.translateItemLinks(jsonItem)
	item.Title = t.translateItemTitle(jsonItem)
	item.Content = t.translateItemContent(jsonItem)
	item.Description = t.translateItemDescription(jsonItem)
	item.Image = t.translateItemImage(jsonItem)
	item.Published = t.translateItemPublished(jsonItem)
	item.PublishedParsed = t.translateItemPublishedParsed(jsonItem)
	item.Updated = t.translateItemUpdated(jsonItem)
	item.UpdatedParsed = t.translateItemUpdatedParsed(jsonItem)
	item.Author = t.translateItemAuthor(jsonItem)
	item.Authors = t.translateItemAuthors(jsonItem)
	item.Categories = t.translateItemCategories(jsonItem)
	item.Enclosures = t.translateItemEnclosures(jsonItem)
	// TODO ExternalURL is missing in global Feed
	// TODO BannerImage is missing in global Feed
	return
}

func (t *DefaultJSONTranslator) translateFeedTitle(json *json.Feed) (title string) {
	if json.Title != "" {
		title = json.Title
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedDescription(json *json.Feed) (desc string) {
	return json.Description
}

func (t *DefaultJSONTranslator) translateFeedLink(json *json.Feed) (link string) {
	if json.HomePageURL != "" {
		link = json.HomePageURL
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedFeedLink(json *json.Feed) (link string) {
	if json.FeedURL != "" {
		link = json.FeedURL
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedLinks(json *json.Feed) (links []string) {
	if json.HomePageURL != "" {
		links = append(links, json.HomePageURL)
	}
	if json.FeedURL != "" {
		links = append(links, json.FeedURL)
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedUpdated(json *json.Feed) (updated string) {
	if len(json.Items) > 0 {
		updated = json.Items[0].DateModified
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedUpdatedParsed(json *json.Feed) (updated *time.Time) {
	if len(json.Items) > 0 {
		updateTime, err := shared.ParseDate(json.Items[0].DateModified)
		if err == nil {
			updated = &updateTime
		}
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedPublished(json *json.Feed) (published string) {
	if len(json.Items) > 0 {
		published = json.Items[0].DatePublished
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedPublishedParsed(json *json.Feed) (published *time.Time) {
	if len(json.Items) > 0 {
		publishTime, err := shared.ParseDate(json.Items[0].DatePublished)
		if err == nil {
			published = &publishTime
		}
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedAuthor(json *json.Feed) (author *Person) {
	if json.Author != nil {
		name, address := shared.ParseNameAddress(json.Author.Name)
		author = &Person{}
		author.Name = name
		author.Email = address
	}
	// Author.URL is missing in global feed
	// Author.Avatar is missing in global feed
	return
}

func (t *DefaultJSONTranslator) translateFeedAuthors(json *json.Feed) (authors []*Person) {
	if json.Authors != nil {
		authors = make([]*Person, 0, len(json.Authors))
		for _, a := range json.Authors {
			name, address := shared.ParseNameAddress(a.Name)
			author := &Person{}
			author.Name = name
			author.Email = address

			authors = append(authors, author)
		}
	} else if author := t.translateFeedAuthor(json); author != nil {
		authors = []*Person{author}
	}
	// Author.URL is missing in global feed
	// Author.Avatar is missing in global feed
	return
}

func (t *DefaultJSONTranslator) translateFeedLanguage(json *json.Feed) (language string) {
	language = json.Language
	return
}

func (t *DefaultJSONTranslator) translateFeedImage(json *json.Feed) (image *Image) {
	// Using the Icon rather than the image
	// icon (optional, string) is the URL of an image for the feed suitable to be used in a timeline. It should be square and relatively large — such as 512 x 512
	if json.Icon != "" {
		image = &Image{}
		image.URL = json.Icon
	}
	return
}

func (t *DefaultJSONTranslator) translateFeedItems(json *json.Feed) (items []*Item) {
	items = make([]*Item, 0, len(json.Items))
	for _, i := range json.Items {
		items = append(items, t.translateFeedItem(i))
	}
	return
}

func (t *DefaultJSONTranslator) translateItemTitle(jsonItem *json.Item) (title string) {
	if jsonItem.Title != "" {
		title = jsonItem.Title
	}
	return
}

func (t *DefaultJSONTranslator) translateItemDescription(jsonItem *json.Item) (desc string) {
	if jsonItem.Summary != "" {
		desc = jsonItem.Summary
	}
	return
}

func (t *DefaultJSONTranslator) translateItemContent(jsonItem *json.Item) (content string) {
	if jsonItem.ContentHTML != "" {
		content = jsonItem.ContentHTML
	} else if jsonItem.ContentText != "" {
		content = jsonItem.ContentText
	}
	return
}

func (t *DefaultJSONTranslator) translateItemLink(jsonItem *json.Item) (link string) {
	return jsonItem.URL
}

func (t *DefaultJSONTranslator) translateItemLinks(jsonItem *json.Item) (links []string) {
	if jsonItem.URL != "" {
		links = append(links, jsonItem.URL)
	}
	if jsonItem.ExternalURL != "" {
		links = append(links, jsonItem.ExternalURL)
	}
	return
}

func (t *DefaultJSONTranslator) translateItemUpdated(jsonItem *json.Item) (updated string) {
	if jsonItem.DateModified != "" {
		updated = jsonItem.DateModified
	}
	return updated
}

func (t *DefaultJSONTranslator) translateItemUpdatedParsed(jsonItem *json.Item) (updated *time.Time) {
	if jsonItem.DateModified != "" {
		updatedTime, err := shared.ParseDate(jsonItem.DateModified)
		if err == nil {
			updated = &updatedTime
		}
	}
	return
}

func (t *DefaultJSONTranslator) translateItemPublished(jsonItem *json.Item) (pubDate string) {
	if jsonItem.DatePublished != "" {
		pubDate = jsonItem.DatePublished
	}
	return
}

func (t *DefaultJSONTranslator) translateItemPublishedParsed(jsonItem *json.Item) (pubDate *time.Time) {
	if jsonItem.DatePublished != "" {
		publishTime, err := shared.ParseDate(jsonItem.DatePublished)
		if err == nil {
			pubDate = &publishTime
		}
	}
	return
}

func (t *DefaultJSONTranslator) translateItemAuthor(jsonItem *json.Item) (author *Person) {
	if jsonItem.Author != nil {
		name, address := shared.ParseNameAddress(jsonItem.Author.Name)
		author = &Person{}
		author.Name = name
		author.Email = address
	}
	// Author.URL is missing in global feed
	// Author.Avatar is missing in global feed
	return
}

func (t *DefaultJSONTranslator) translateItemAuthors(jsonItem *json.Item) (authors []*Person) {
	if jsonItem.Authors != nil {
		authors = make([]*Person, 0, len(jsonItem.Authors))
		for _, a := range jsonItem.Authors {
			name, address := shared.ParseNameAddress(a.Name)
			author := &Person{}
			author.Name = name
			author.Email = address

			authors = append(authors, author)
		}
	} else if author := t.translateItemAuthor(jsonItem); author != nil {
		authors = []*Person{author}
	}
	// Author.URL is missing in global feed
	// Author.Avatar is missing in global feed
	return
}

func (t *DefaultJSONTranslator) translateItemGUID(jsonItem *json.Item) (guid string) {
	if jsonItem.ID != "" {
		guid = jsonItem.ID
	}
	return
}

func (t *DefaultJSONTranslator) translateItemImage(jsonItem *json.Item) (image *Image) {
	if jsonItem.Image != "" {
		image = &Image{}
		image.URL = jsonItem.Image
	} else if jsonItem.BannerImage != "" {
		image = &Image{}
		image.URL = jsonItem.BannerImage
	}
	return
}

func (t *DefaultJSONTranslator) translateItemCategories(jsonItem *json.Item) (categories []string) {
	if len(jsonItem.Tags) > 0 {
		categories = jsonItem.Tags
	}
	return
}

func (t *DefaultJSONTranslator) translateItemEnclosures(jsonItem *json.Item) (enclosures []*Enclosure) {
	if jsonItem.Attachments != nil {
		for _, attachment := range *jsonItem.Attachments {
			e := &Enclosure{}
			e.URL = attachment.URL
			e.Type = attachment.MimeType
			e.Length = fmt.Sprintf("%d", attachment.DurationInSeconds)
			// Title is not defined in global enclosure
			// SizeInBytes is not defined in global enclosure
			enclosures = append(enclosures, e)
		}
	}
	return
}
