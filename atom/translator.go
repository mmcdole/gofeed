package atom

import (
	"time"

	"github.com/mmcdole/gofeed/feed"
)

// Translator converts an atom.Feed struct
// into the generic feed.Feed struct
type Translator interface {
	Translate(atom *Feed) *feed.Feed
}

// DefaultTranslator converts an atom.Feed struct
// into the generic feed.Feed struct.
//
// This default implementation defines a set of
// mapping rules between atom.Feed -> feed.Feed
// for each of the fields in feed.Feed.
type DefaultTranslator struct{}

func (t *DefaultTranslator) Translate(atom *Feed) *feed.Feed {
	feed := &feed.Feed{}
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

func (t *DefaultTranslator) translateFeedItem(entry *Entry) (item *feed.Item) {
	item = &feed.Item{}
	item.Title = t.translateItemTitle(entry)
	item.Description = t.translateItemDescription(entry)
	item.Summary = t.translateItemSummary(entry)
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

func (t *DefaultTranslator) translateFeedTitle(atom *Feed) (title string) {
	return atom.Title
}

func (t *DefaultTranslator) translateFeedDescription(atom *Feed) (desc string) {
	return atom.Subtitle
}

func (t *DefaultTranslator) translateFeedLink(atom *Feed) (link string) {
	l := t.firstLinkWithType("alternate", atom.Links)
	if l != nil {
		link = l.Href
	}
	return
}

func (t *DefaultTranslator) translateFeedFeedLink(atom *Feed) (link string) {
	feedLink := t.firstLinkWithType("self", atom.Links)
	if feedLink != nil {
		link = feedLink.Href
	}
	return
}

func (t *DefaultTranslator) translateFeedUpdated(atom *Feed) (updated string) {
	return atom.Updated
}

func (t *DefaultTranslator) translateFeedUpdatedParsed(atom *Feed) (updated *time.Time) {
	return atom.UpdatedParsed
}

func (t *DefaultTranslator) translateFeedAuthor(atom *Feed) (author *feed.Person) {
	a := t.firstPerson(atom.Authors)
	if author != nil {
		feedAuthor := feed.Person{}
		feedAuthor.Name = a.Name
		feedAuthor.Email = a.Email
		author = &feedAuthor
	}
	return
}

func (t *DefaultTranslator) translateFeedLanguage(atom *Feed) (language string) {
	return atom.Language
}

func (t *DefaultTranslator) translateFeedImage(atom *Feed) (image *feed.Image) {
	if atom.Logo != "" {
		feedImage := feed.Image{}
		feedImage.URL = atom.Logo
		image = &feedImage
	}
	return
}

func (t *DefaultTranslator) translateFeedCopyright(atom *Feed) (rights string) {
	return atom.Rights
}

func (t *DefaultTranslator) translateFeedGenerator(atom *Feed) (generator string) {
	if atom.Generator != nil {
		generator = atom.Generator.Value
	}
	return
}

func (t *DefaultTranslator) translateFeedCategories(atom *Feed) (categories []string) {
	if atom.Categories != nil {
		categories := []string{}
		for _, c := range atom.Categories {
			categories = append(categories, c.Term)
		}
	}
	return
}

func (t *DefaultTranslator) translateFeedItems(atom *Feed) (items []*feed.Item) {
	items = []*feed.Item{}
	for _, entry := range atom.Entries {
		items = append(items, t.translateFeedItem(entry))
	}
	return
}

func (t *DefaultTranslator) translateItemTitle(entry *Entry) (title string) {
	return entry.Title
}

func (t *DefaultTranslator) translateItemDescription(entry *Entry) (desc string) {
	if entry.Content != nil {
		desc = entry.Content.Value
	}
	return
}

func (t *DefaultTranslator) translateItemSummary(entry *Entry) (summary string) {
	return entry.Summary
}

func (t *DefaultTranslator) translateItemLink(entry *Entry) (link string) {
	l := t.firstLinkWithType("alternate", entry.Links)
	if l != nil {
		link = l.Href
	}
	return
}

func (t *DefaultTranslator) translateItemUpdated(entry *Entry) (updated string) {
	return entry.Updated
}

func (t *DefaultTranslator) translateItemUpdatedParsed(entry *Entry) (updated *time.Time) {
	return entry.UpdatedParsed
}

func (t *DefaultTranslator) translateItemPublished(entry *Entry) (updated string) {
	return entry.Published
}

func (t *DefaultTranslator) translateItemPublishedParsed(entry *Entry) (updated *time.Time) {
	return entry.PublishedParsed
}

func (t *DefaultTranslator) translateItemAuthor(entry *Entry) (author *feed.Person) {
	a := t.firstPerson(entry.Authors)
	if a != nil {
		author = &feed.Person{}
		author.Name = a.Name
		author.Email = a.Email
	}
	return
}

func (t *DefaultTranslator) translateItemGuid(entry *Entry) (guid string) {
	return entry.ID
}

func (t *DefaultTranslator) translateItemImage(entry *Entry) (image *feed.Image) {
	return nil
}

func (t *DefaultTranslator) translateItemCategories(entry *Entry) (categories []string) {
	if entry.Categories != nil {
		categories := []string{}
		for _, c := range entry.Categories {
			categories = append(categories, c.Term)
		}
	}
	return
}

func (t *DefaultTranslator) translateItemEnclosures(entry *Entry) (enclosures []*feed.Enclosure) {
	if entry.Links != nil {
		enclosures := []*feed.Enclosure{}
		for _, e := range entry.Links {
			if e.Rel == "enclosure" {
				enclosure := &feed.Enclosure{}
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

func (t *DefaultTranslator) firstLinkWithType(linkType string, links []*Link) *Link {
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

func (t *DefaultTranslator) firstPerson(persons []*Person) *Person {
	if persons == nil {
		return nil
	}

	if len(persons) == 0 {
		return nil
	}

	return persons[0]
}
