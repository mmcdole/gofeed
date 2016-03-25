package gofeed

import (
	"time"

	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/internal/shared"
	"github.com/mmcdole/gofeed/rss"
)

// Translator converts an atom.Feed struct
// into the generic Feed struct
type AtomTranslator interface {
	Translate(atom *atom.Feed) *Feed
}

// Translator converts an rss.Feed struct
// into the generic Feed struct
type RSSTranslator interface {
	Translate(rss *rss.Feed) *Feed
}

// DefaulRSSTranslator converts an rss.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between rss.Feed -> Feed
// for each of the fields in Feed.
type DefaultRSSTranslator struct{}

func (t *DefaultRSSTranslator) Translate(rss *rss.Feed) *Feed {
	feed := &Feed{}
	feed.Title = t.translateFeedTitle(rss)
	feed.Description = t.translateFeedDescription(rss)
	feed.Link = t.translateFeedLink(rss)
	feed.FeedLink = t.translateFeedFeedLink(rss)
	feed.Published = t.translateFeedPublished(rss)
	feed.PublishedParsed = t.translateFeedPublishedParsed(rss)
	feed.Author = t.translateFeedAuthor(rss)
	feed.Language = t.translateFeedLanguage(rss)
	feed.Image = t.translateFeedImage(rss)
	feed.Copyright = t.translateFeedCopyright(rss)
	feed.Generator = t.translateFeedGenerator(rss)
	feed.Categories = t.translateFeedCategories(rss)
	feed.Items = t.translateFeedItems(rss)
	feed.Extensions = rss.Extensions
	feed.FeedVersion = rss.Version
	feed.FeedType = "rss"
	return feed
}

func (t *DefaultRSSTranslator) translateFeedItem(rssItem *rss.Item) (item *Item) {
	item = &Item{}
	item.Title = t.translateItemTitle(rssItem)
	item.Description = t.translateItemDescription(rssItem)
	item.Link = t.translateItemLink(rssItem)
	item.Published = t.translateItemPublished(rssItem)
	item.PublishedParsed = t.translateItemPublishedParsed(rssItem)
	item.Author = t.translateItemAuthor(rssItem)
	item.Guid = t.translateItemGuid(rssItem)
	item.Image = t.translateItemImage(rssItem)
	item.Categories = t.translateItemCategories(rssItem)
	item.Enclosures = t.translateItemEnclosures(rssItem)
	item.Extensions = rssItem.Extensions
	return item
}

func (t *DefaultRSSTranslator) translateFeedTitle(rss *rss.Feed) (title string) {
	return rss.Title
}

func (t *DefaultRSSTranslator) translateFeedDescription(rss *rss.Feed) (desc string) {
	return rss.Description
}

func (t *DefaultRSSTranslator) translateFeedLink(rss *rss.Feed) (link string) {
	return rss.Link
}

func (t *DefaultRSSTranslator) translateFeedFeedLink(rss *rss.Feed) (link string) {
	if rss.Extensions == nil {
		return
	}
	if atom, ok := rss.Extensions["atom"]; ok {
		if links, ok := atom["link"]; ok {
			for _, l := range links {
				if l.Attrs["Rel"] == "self" {
					return l.Value
				}
			}
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedPublished(rss *rss.Feed) (updated string) {
	return rss.PubDate
}

func (t *DefaultRSSTranslator) translateFeedPublishedParsed(rss *rss.Feed) (updated *time.Time) {
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
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedLanguage(rss *rss.Feed) (language string) {
	return rss.Language
}

func (t *DefaultRSSTranslator) translateFeedImage(rss *rss.Feed) (image *Image) {
	if rss.Image != nil {
		image = &Image{}
		image.Title = rss.Image.Title
		image.URL = rss.Image.URL
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedCopyright(rss *rss.Feed) (rights string) {
	return rss.Copyright
}

func (t *DefaultRSSTranslator) translateFeedGenerator(rss *rss.Feed) (generator string) {
	return rss.Generator
}

func (t *DefaultRSSTranslator) translateFeedCategories(rss *rss.Feed) (categories []string) {
	if rss.Categories != nil {
		categories = []string{}
		for _, c := range rss.Categories {
			categories = append(categories, c.Value)
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateFeedItems(rss *rss.Feed) (items []*Item) {
	items = []*Item{}
	for _, i := range rss.Items {
		items = append(items, t.translateFeedItem(i))
	}
	return
}

func (t *DefaultRSSTranslator) translateItemTitle(rssItem *rss.Item) (title string) {
	return rssItem.Title
}

func (t *DefaultRSSTranslator) translateItemDescription(rssItem *rss.Item) (desc string) {
	return rssItem.Description
}

func (t *DefaultRSSTranslator) translateItemLink(rssItem *rss.Item) (link string) {
	return rssItem.Link
}

func (t *DefaultRSSTranslator) translateItemPublished(rssItem *rss.Item) (updated string) {
	return rssItem.PubDate
}

func (t *DefaultRSSTranslator) translateItemPublishedParsed(rssItem *rss.Item) (updated *time.Time) {
	return rssItem.PubDateParsed
}

func (t *DefaultRSSTranslator) translateItemAuthor(rssItem *rss.Item) (author *Person) {
	return // TODO: dc and itunes
}

func (t *DefaultRSSTranslator) translateItemGuid(rssItem *rss.Item) (guid string) {
	if rssItem.Guid != nil {
		guid = rssItem.Guid.Value
	}
	return
}

func (t *DefaultRSSTranslator) translateItemImage(rssItem *rss.Item) (image *Image) {
	return // TODO: itunes
}

func (t *DefaultRSSTranslator) translateItemCategories(rssItem *rss.Item) (categories []string) {
	if rssItem.Categories != nil {
		categories = []string{}
		for _, c := range rssItem.Categories {
			categories = append(categories, c.Value)
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateItemEnclosures(rssItem *rss.Item) (enclosures []*Enclosure) {
	if rssItem.Enclosure != nil {
		e := &Enclosure{}
		e.URL = rssItem.Enclosure.URL
		e.Type = rssItem.Enclosure.Type
		e.Length = rssItem.Enclosure.Length
		enclosures = []*Enclosure{e}
	}
	return
}

// DefaultAtomTranslator converts an atom.Feed struct
// into the generic Feed struct.
//
// This default implementation defines a set of
// mapping rules between atom.Feed -> Feed
// for each of the fields in Feed.
type DefaultAtomTranslator struct{}

func (t *DefaultAtomTranslator) Translate(atom *atom.Feed) *Feed {
	feed := &Feed{}
	feed.Title = t.translateFeedTitle(atom)
	feed.Description = t.translateFeedDescription(atom)
	feed.Link = t.translateFeedLink(atom)
	feed.FeedLink = t.translateFeedFeedLink(atom)
	feed.Updated = t.translateFeedUpdated(atom)
	feed.UpdatedParsed = t.translateFeedUpdatedParsed(atom)
	feed.Author = t.translateFeedAuthor(atom)
	feed.Language = t.translateFeedLanguage(atom)
	feed.Image = t.translateFeedImage(atom)
	feed.Copyright = t.translateFeedCopyright(atom)
	feed.Categories = t.translateFeedCategories(atom)
	feed.Items = t.translateFeedItems(atom)
	feed.Extensions = atom.Extensions
	feed.FeedVersion = atom.Version
	feed.FeedType = "atom"
	return feed
}

func (t *DefaultAtomTranslator) translateFeedItem(entry *atom.Entry) (item *Item) {
	item = &Item{}
	item.Title = t.translateItemTitle(entry)
	item.Description = t.translateItemDescription(entry)
	item.Content = t.translateItemContent(entry)
	item.Link = t.translateItemLink(entry)
	item.Updated = t.translateItemUpdated(entry)
	item.UpdatedParsed = t.translateItemUpdatedParsed(entry)
	item.Published = t.translateItemPublished(entry)
	item.PublishedParsed = t.translateItemPublishedParsed(entry)
	item.Author = t.translateItemAuthor(entry)
	item.Guid = t.translateItemGuid(entry)
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

func (t *DefaultAtomTranslator) translateFeedUpdated(atom *atom.Feed) (updated string) {
	return atom.Updated
}

func (t *DefaultAtomTranslator) translateFeedUpdatedParsed(atom *atom.Feed) (updated *time.Time) {
	return atom.UpdatedParsed
}

func (t *DefaultAtomTranslator) translateFeedAuthor(atom *atom.Feed) (author *Person) {
	a := t.firstPerson(atom.Authors)
	if author != nil {
		feedAuthor := Person{}
		feedAuthor.Name = a.Name
		feedAuthor.Email = a.Email
		author = &feedAuthor
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
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedCopyright(atom *atom.Feed) (rights string) {
	return atom.Rights
}

func (t *DefaultAtomTranslator) translateFeedGenerator(atom *atom.Feed) (generator string) {
	if atom.Generator != nil {
		generator = atom.Generator.Value
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedCategories(atom *atom.Feed) (categories []string) {
	if atom.Categories != nil {
		categories := []string{}
		for _, c := range atom.Categories {
			categories = append(categories, c.Term)
		}
	}
	return
}

func (t *DefaultAtomTranslator) translateFeedItems(atom *atom.Feed) (items []*Item) {
	items = []*Item{}
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

func (t *DefaultAtomTranslator) translateItemUpdated(entry *atom.Entry) (updated string) {
	return entry.Updated
}

func (t *DefaultAtomTranslator) translateItemUpdatedParsed(entry *atom.Entry) (updated *time.Time) {
	return entry.UpdatedParsed
}

func (t *DefaultAtomTranslator) translateItemPublished(entry *atom.Entry) (updated string) {
	return entry.Published
}

func (t *DefaultAtomTranslator) translateItemPublishedParsed(entry *atom.Entry) (updated *time.Time) {
	return entry.PublishedParsed
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

func (t *DefaultAtomTranslator) translateItemGuid(entry *atom.Entry) (guid string) {
	return entry.ID
}

func (t *DefaultAtomTranslator) translateItemImage(entry *atom.Entry) (image *Image) {
	return nil
}

func (t *DefaultAtomTranslator) translateItemCategories(entry *atom.Entry) (categories []string) {
	if entry.Categories != nil {
		categories := []string{}
		for _, c := range entry.Categories {
			categories = append(categories, c.Term)
		}
	}
	return
}

func (t *DefaultAtomTranslator) translateItemEnclosures(entry *atom.Entry) (enclosures []*Enclosure) {
	if entry.Links != nil {
		enclosures := []*Enclosure{}
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
		if link.Rel == "alternate" {
			return link
		}
	}
	return nil
}

func (t *DefaultAtomTranslator) firstPerson(persons []*atom.Person) *atom.Person {
	if persons == nil {
		return nil
	}

	if len(persons) == 0 {
		return nil
	}

	return persons[0]
}
