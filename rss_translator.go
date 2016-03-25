package rss

import (
	"regexp"
	"time"

	"github.com/mmcdole/gofeed/feed"
)

var (
	emailNameRgx = regexp.MustCompile(`^([^@]+@[^\s]+)\s+\(([^@]+)\)$`)
	nameEmailRgx = regexp.MustCompile(`^([^@]+)\s+\(([^@]+@[^)]+)\)$`)
	nameOnlyRgx  = regexp.MustCompile(`^([^@()]+)$`)
	emailOnlyRgx = regexp.MustCompile(`^([^@()]+@[^@()]+)$`)
)

// Translator converts an rss.Feed struct
// into the generic feed.Feed struct
type Translator interface {
	Translate(rss *Feed) *feed.Feed
}

// DefaulTranslator converts an rss.Feed struct
// into the generic feed.Feed struct.
//
// This default implementation defines a set of
// mapping rules between rss.Feed -> feed.Feed
// for each of the fields in feed.Feed.
type DefaultTranslator struct{}

func (t *DefaultTranslator) Translate(rss *Feed) *feed.Feed {
	feed := &feed.Feed{}
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

func (t *DefaultTranslator) translateFeedItem(rssItem *Item) (item *feed.Item) {
	item = &feed.Item{}
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

func (t *DefaultTranslator) translateFeedTitle(rss *Feed) (title string) {
	return rss.Title
}

func (t *DefaultTranslator) translateFeedDescription(rss *Feed) (desc string) {
	return rss.Description
}

func (t *DefaultTranslator) translateFeedLink(rss *Feed) (link string) {
	return rss.Link
}

func (t *DefaultTranslator) translateFeedFeedLink(rss *Feed) (link string) {
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

func (t *DefaultTranslator) translateFeedPublished(rss *Feed) (updated string) {
	return rss.PubDate
}

func (t *DefaultTranslator) translateFeedPublishedParsed(rss *Feed) (updated *time.Time) {
	return rss.PubDateParsed
}

func (t *DefaultTranslator) translateFeedAuthor(rss *Feed) (author *feed.Person) {
	if rss.ManagingEditor != "" {

	} else if rss.WebMaster != "" {

	}
	return
}

func (t *DefaultTranslator) translateFeedLanguage(rss *Feed) (language string) {
	return rss.Language
}

func (t *DefaultTranslator) translateFeedImage(rss *Feed) (image *feed.Image) {
	if rss.Image != nil {
		image = &feed.Image{}
		image.Title = rss.Image.Title
		image.URL = rss.Image.URL
	}
	return
}

func (t *DefaultTranslator) translateFeedCopyright(rss *Feed) (rights string) {
	return rss.Copyright
}

func (t *DefaultTranslator) translateFeedGenerator(rss *Feed) (generator string) {
	return rss.Generator
}

func (t *DefaultTranslator) translateFeedCategories(rss *Feed) (categories []string) {
	if rss.Categories != nil {
		categories = []string{}
		for _, c := range rss.Categories {
			categories = append(categories, c.Value)
		}
	}
	return
}

func (t *DefaultTranslator) translateFeedItems(rss *Feed) (items []*feed.Item) {
	items = []*feed.Item{}
	for _, i := range rss.Items {
		items = append(items, t.translateFeedItem(i))
	}
	return
}

func (t *DefaultTranslator) translateItemTitle(rssItem *Item) (title string) {
	return rssItem.Title
}

func (t *DefaultTranslator) translateItemDescription(rssItem *Item) (desc string) {
	return rssItem.Description
}

func (t *DefaultTranslator) translateItemLink(rssItem *Item) (link string) {
	return rssItem.Link
}

func (t *DefaultTranslator) translateItemPublished(rssItem *Item) (updated string) {
	return rssItem.PubDate
}

func (t *DefaultTranslator) translateItemPublishedParsed(rssItem *Item) (updated *time.Time) {
	return rssItem.PubDateParsed
}

func (t *DefaultTranslator) translateItemAuthor(rssItem *Item) (author *feed.Person) {
	return // TODO: dc and itunes
}

func (t *DefaultTranslator) translateItemGuid(rssItem *Item) (guid string) {
	if rssItem.Guid != nil {
		guid = rssItem.Guid.Value
	}
	return
}

func (t *DefaultTranslator) translateItemImage(rssItem *Item) (image *feed.Image) {
	return // TODO: itunes
}

func (t *DefaultTranslator) translateItemCategories(rssItem *Item) (categories []string) {
	if rssItem.Categories != nil {
		categories = []string{}
		for _, c := range rssItem.Categories {
			categories = append(categories, c.Value)
		}
	}
	return
}

func (t *DefaultTranslator) translateItemEnclosures(rssItem *Item) (enclosures []*feed.Enclosure) {
	if rssItem.Enclosure != nil {
		e := &feed.Enclosure{}
		e.URL = rssItem.Enclosure.URL
		e.Type = rssItem.Enclosure.Type
		e.Length = rssItem.Enclosure.Length
		enclosures = []*feed.Enclosure{e}
	}
	return
}

func (t *DefaultTranslator) parsePersonText(personText string) (person *feed.Person) {
	if personText == "" {
		return
	}

	if emailNameRgx.MatchString(personText) {
		result := emailNameRgx.FindStringSubmatch(personText)
		person = &feed.Person{}
		person.Email = result[1]
		person.Name = result[2]
	} else if nameEmailRgx.MatchString(personText) {
		result := nameEmailRgx.FindStringSubmatch(personText)
		person = &feed.Person{}
		person.Name = result[1]
		person.Email = result[2]
	} else if nameOnlyRgx.MatchString(personText) {
		result := nameOnlyRgx.FindStringSubmatch(personText)
		person = &feed.Person{}
		person.Name = result[1]
	} else if emailOnlyRgx.MatchString(personText) {
		result := emailOnlyRgx.FindStringSubmatch(personText)
		person = &feed.Person{}
		person.Email = result[1]
	}
	return
}
