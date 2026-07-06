package gofeed

import (
	"bytes"
	"fmt"
	"strings"
	"time"

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
type DefaultRSSTranslator struct {
	// DisableContentImageScan turns off the fallback that parses item and
	// feed HTML (content, description) to find a first <img> when no
	// explicit image is present. The scan runs a full HTML parse per item,
	// so large feeds may want it off. Default off (the scan runs), matching
	// historical behavior.
	DisableContentImageScan bool
}

// Translate converts an RSS feed into the universal
// feed type.
func (t *DefaultRSSTranslator) Translate(feed interface{}) (*Feed, error) {
	rss, found := feed.(*rss.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *rss.Feed")
	}

	result := &Feed{
		Published:       rss.PubDate,
		PublishedParsed: rss.PubDateParsed,
		Generator:       rss.Generator,
		ITunesExt:       rss.ITunesExt,
		DublinCoreExt:   rss.DublinCoreExt,
		Extensions:      rss.Extensions,
		FeedVersion:     rss.Version,
		FeedType:        "rss",
	}

	dc := rss.DublinCoreExt

	result.Title = rss.Title
	if result.Title == "" && dc != nil {
		result.Title = firstString(dc.Title)
	}

	result.Description = rss.Description
	if result.Description == "" && rss.ITunesExt != nil {
		result.Description = rss.ITunesExt.Summary
	}

	result.Language = rss.Language
	if result.Language == "" && dc != nil {
		result.Language = firstString(dc.Language)
	}

	result.Copyright = rss.Copyright
	if result.Copyright == "" && dc != nil {
		result.Copyright = firstString(dc.Rights)
	}

	result.Updated = rss.LastBuildDate
	if result.Updated == "" && dc != nil {
		result.Updated = firstString(dc.Date)
	}
	result.UpdatedParsed = rss.LastBuildDateParsed
	if result.UpdatedParsed == nil && dc != nil && dc.Date != nil {
		if date, err := shared.ParseDate(firstString(dc.Date)); err == nil {
			result.UpdatedParsed = &date
		}
	}

	result.Link = t.translateFeedLink(rss)
	result.FeedLink = t.translateFeedFeedLink(rss)
	result.Links = t.translateFeedLinks(rss)

	if author := t.translateFeedAuthor(rss); author != nil {
		result.Author = author
		result.Authors = []*Person{author}
	}

	result.Image = t.translateFeedImage(rss)
	result.Categories = t.translateFeedCategories(rss)

	result.Items = make([]*Item, 0, len(rss.Items))
	for _, i := range rss.Items {
		result.Items = append(result.Items, t.translateFeedItem(i))
	}

	return result, nil
}

func (t *DefaultRSSTranslator) translateFeedItem(rssItem *rss.Item) *Item {
	item := &Item{
		Link:          rssItem.Link,
		DublinCoreExt: rssItem.DublinCoreExt,
		ITunesExt:     rssItem.ITunesExt,
		Extensions:    rssItem.Extensions,
		Custom:        rssItem.Custom,
	}

	dc := rssItem.DublinCoreExt

	item.Title = rssItem.Title
	if item.Title == "" && dc != nil {
		item.Title = firstString(dc.Title)
	}

	item.Description = t.translateItemDescription(rssItem)

	item.Content = rssItem.Content
	if item.Content == "" {
		item.Content = t.atomExtValue(rssItem.Extensions, "content")
	}

	if len(rssItem.Links) > 0 {
		item.Links = append(item.Links, rssItem.Links...)
	}

	item.Published = t.translateItemPublished(rssItem)
	item.PublishedParsed = t.translateItemPublishedParsed(rssItem)
	item.Updated = t.translateItemUpdated(rssItem)
	item.UpdatedParsed = t.translateItemUpdatedParsed(rssItem)

	if author := t.translateItemAuthor(rssItem); author != nil {
		item.Author = author
		item.Authors = []*Person{author}
	}

	if rssItem.GUID != nil {
		item.GUID = rssItem.GUID.Value
	}

	item.Image = t.translateItemImage(rssItem)
	item.Categories = t.translateItemCategories(rssItem)
	item.Enclosures = t.translateItemEnclosures(rssItem)
	return item
}

func (t *DefaultRSSTranslator) translateFeedLink(rss *rss.Feed) (link string) {
	if rss.Link != "" {
		return rss.Link
	}
	// Fall back to an embedded atom:link with rel="alternate" (or no rel),
	// which points at the site the way <link> does.
	for _, ex := range t.extensionsForKeys([]string{"atom", "atom10", "atom03"}, rss.Extensions) {
		if links, ok := ex["link"]; ok {
			for _, l := range links {
				if l.Attrs["rel"] == "" || l.Attrs["rel"] == "alternate" {
					return l.Attrs["href"]
				}
			}
		}
	}
	return ""
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

// translateFeedAuthor picks the feed author from the first populated source:
// managingEditor, webMaster, dc:author, dc:creator, then itunes:author.
func (t *DefaultRSSTranslator) translateFeedAuthor(rss *rss.Feed) *Person {
	switch {
	case rss.ManagingEditor != "":
		return personFromText(rss.ManagingEditor)
	case rss.WebMaster != "":
		return personFromText(rss.WebMaster)
	case rss.DublinCoreExt != nil && rss.DublinCoreExt.Author != nil:
		return personFromText(firstString(rss.DublinCoreExt.Author))
	case rss.DublinCoreExt != nil && rss.DublinCoreExt.Creator != nil:
		return personFromText(firstString(rss.DublinCoreExt.Creator))
	case rss.ITunesExt != nil && rss.ITunesExt.Author != "":
		return personFromText(rss.ITunesExt.Author)
	}
	return nil
}

// translateFeedImage picks the feed image from the first populated source:
// the channel image, itunes:image, a media:content image, then a scan of the
// channel description HTML (unless disabled).
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
	if t.DisableContentImageScan {
		return nil
	}
	return firstImageFromHtmlDocument(rss.Description)
}

// translateFeedCategories merges plain channel categories with itunes
// keywords, itunes categories (and subcategories) and dc:subject values.
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

// translateItemDescription picks the item description from the first
// populated source: description, dc:description, itunes:summary, then an
// embedded atom summary.
func (t *DefaultRSSTranslator) translateItemDescription(rssItem *rss.Item) (desc string) {
	if rssItem.Description != "" {
		desc = rssItem.Description
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Description != nil {
		desc = firstString(rssItem.DublinCoreExt.Description)
	} else if rssItem.ITunesExt != nil && rssItem.ITunesExt.Summary != "" {
		desc = rssItem.ITunesExt.Summary
	}
	if desc == "" {
		desc = t.atomExtValue(rssItem.Extensions, "summary")
	}
	return
}

func (t *DefaultRSSTranslator) translateItemUpdated(rssItem *rss.Item) (updated string) {
	if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		updated = firstString(rssItem.DublinCoreExt.Date)
	}
	if updated == "" {
		updated = t.atomExtValue(rssItem.Extensions, "updated")
	}
	return updated
}

func (t *DefaultRSSTranslator) translateItemUpdatedParsed(rssItem *rss.Item) (updated *time.Time) {
	if updatedText := t.translateItemUpdated(rssItem); updatedText != "" {
		if updatedDate, err := shared.ParseDate(updatedText); err == nil {
			updated = &updatedDate
		}
	}
	return
}

func (t *DefaultRSSTranslator) translateItemPublished(rssItem *rss.Item) (pubDate string) {
	if rssItem.PubDate != "" {
		return rssItem.PubDate
	} else if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		return firstString(rssItem.DublinCoreExt.Date)
	}
	return t.atomExtValue(rssItem.Extensions, "published")
}

func (t *DefaultRSSTranslator) translateItemPublishedParsed(rssItem *rss.Item) (pubDate *time.Time) {
	if rssItem.PubDateParsed != nil {
		return rssItem.PubDateParsed
	}
	pubDateText := ""
	if rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Date != nil {
		pubDateText = firstString(rssItem.DublinCoreExt.Date)
	}
	if pubDateText == "" {
		pubDateText = t.atomExtValue(rssItem.Extensions, "published")
	}
	if pubDateText != "" {
		if pubDateParsed, err := shared.ParseDate(pubDateText); err == nil {
			pubDate = &pubDateParsed
		}
	}
	return
}

// translateItemAuthor picks the item author from the first populated source:
// author, dc:author, dc:creator, itunes:author, then an embedded atom
// author's name and email children.
func (t *DefaultRSSTranslator) translateItemAuthor(rssItem *rss.Item) *Person {
	switch {
	case rssItem.Author != "":
		return personFromText(rssItem.Author)
	case rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Author != nil:
		return personFromText(firstString(rssItem.DublinCoreExt.Author))
	case rssItem.DublinCoreExt != nil && rssItem.DublinCoreExt.Creator != nil:
		return personFromText(firstString(rssItem.DublinCoreExt.Creator))
	case rssItem.ITunesExt != nil && rssItem.ITunesExt.Author != "":
		return personFromText(rssItem.ITunesExt.Author)
	}
	if name, email := t.atomExtChild(rssItem.Extensions, "author", "name"), t.atomExtChild(rssItem.Extensions, "author", "email"); name != "" || email != "" {
		return &Person{Name: name, Email: email}
	}
	return nil
}

// translateItemImage picks the item image from the first populated source:
// itunes:image, a media:content image, an image enclosure, then a scan of the
// item content and description HTML (unless disabled).
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
	if t.DisableContentImageScan {
		return nil
	}
	if img := firstImageFromHtmlDocument(rssItem.Content); img != nil {
		return img
	}
	if img := firstImageFromHtmlDocument(rssItem.Description); img != nil {
		return img
	}
	return nil
}

// translateItemCategories merges plain item categories with itunes keywords,
// dc:subject values and embedded atom:category terms.
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

	for _, m := range t.extensionsForKeys([]string{"atom", "atom10", "atom03"}, rssItem.Extensions) {
		for _, c := range m["category"] {
			if term := c.Attrs["term"]; term != "" {
				cats = append(cats, term)
			}
		}
	}

	if len(cats) > 0 {
		categories = cats
	}

	return
}

func (t *DefaultRSSTranslator) translateItemEnclosures(rssItem *rss.Item) (enclosures []*Enclosure) {
	for _, enc := range rssItem.Enclosures {
		enclosures = append(enclosures, &Enclosure{
			URL:    enc.URL,
			Type:   enc.Type,
			Length: enc.Length,
		})
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

// atomExtValue returns the text of the first matching Atom element embedded in
// an RSS feed (across the atom/atom10/atom03 namespaces), or "" if absent. This
// promotes Atom tags in RSS to the universal fields, the same way dc:/itunes:
// values already are.
func (t *DefaultRSSTranslator) atomExtValue(exts ext.Extensions, name string) string {
	for _, m := range t.extensionsForKeys([]string{"atom", "atom10", "atom03"}, exts) {
		if es, ok := m[name]; ok && len(es) > 0 {
			return es[0].Value
		}
	}
	return ""
}

// atomExtChild returns the text of a child of the first matching Atom element
// (e.g. author > name).
func (t *DefaultRSSTranslator) atomExtChild(exts ext.Extensions, parent, child string) string {
	for _, m := range t.extensionsForKeys([]string{"atom", "atom10", "atom03"}, exts) {
		if ps, ok := m[parent]; ok && len(ps) > 0 {
			if cs, ok := ps[0].Children[child]; ok && len(cs) > 0 {
				return cs[0].Value
			}
		}
	}
	return ""
}

// firstString returns the first entry of a string slice, or "" when empty.
func firstString(entries []string) string {
	if len(entries) == 0 {
		return ""
	}
	return entries[0]
}

// personFromText builds a Person from a free-form author string like
// "Example Name (example@site.com)".
func personFromText(text string) *Person {
	name, address := shared.ParseNameAddress(text)
	return &Person{Name: name, Email: address}
}

func firstImageFromHtmlDocument(document string) *Image {
	doc, err := html.Parse(bytes.NewBufferString(document))
	if err != nil {
		return nil
	}
	if img := firstImgWithSrc(doc); img != nil {
		for _, attr := range img.Attr {
			if attr.Key == "src" {
				return &Image{URL: attr.Val}
			}
		}
	}
	return nil
}

// firstImgWithSrc returns the first <img> element in document order that
// carries a src attribute.
func firstImgWithSrc(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if img := firstImgWithSrc(c); img != nil {
			return img
		}
	}
	return nil
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
	atomFeed, found := feed.(*atom.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *atom.Feed")
	}

	result := &Feed{
		Title:         atomFeed.Title,
		Description:   atomFeed.Subtitle,
		Updated:       atomFeed.Updated,
		UpdatedParsed: atomFeed.UpdatedParsed,
		Language:      atomFeed.Language,
		Copyright:     atomFeed.Rights,
		Extensions:    atomFeed.Extensions,
		FeedVersion:   atomFeed.Version,
		FeedType:      "atom",
	}

	if l := firstLinkWithRel("alternate", atomFeed.Links); l != nil {
		result.Link = l.Href
	}
	if l := firstLinkWithRel("self", atomFeed.Links); l != nil {
		result.FeedLink = l.Href
	}
	for _, l := range atomFeed.Links {
		if l.Rel == "" || l.Rel == "alternate" || l.Rel == "self" {
			result.Links = append(result.Links, l.Href)
		}
	}

	if len(atomFeed.Authors) > 0 {
		result.Authors = atomPersons(atomFeed.Authors)
		result.Author = result.Authors[0]
	}

	if atomFeed.Logo != "" {
		result.Image = &Image{URL: atomFeed.Logo}
	} else if atomFeed.Icon != "" {
		result.Image = &Image{URL: atomFeed.Icon}
	}

	if atomFeed.Generator != nil {
		generator := atomFeed.Generator.Value
		if atomFeed.Generator.Version != "" {
			generator += " v" + atomFeed.Generator.Version
		}
		if atomFeed.Generator.URI != "" {
			generator += " " + atomFeed.Generator.URI
		}
		result.Generator = strings.TrimSpace(generator)
	}

	result.Categories = atomCategories(atomFeed.Categories)

	result.Items = make([]*Item, 0, len(atomFeed.Entries))
	for _, entry := range atomFeed.Entries {
		result.Items = append(result.Items, t.translateFeedItem(entry))
	}

	return result, nil
}

func (t *DefaultAtomTranslator) translateFeedItem(entry *atom.Entry) *Item {
	item := &Item{
		Title:         entry.Title,
		Description:   entry.Summary,
		Updated:       entry.Updated,
		UpdatedParsed: entry.UpdatedParsed,
		GUID:          entry.ID,
		Extensions:    entry.Extensions,
	}

	if entry.Content != nil {
		item.Content = entry.Content.Value
	}

	if l := firstLinkWithRel("alternate", entry.Links); l != nil {
		item.Link = l.Href
	}
	for _, l := range entry.Links {
		if l.Rel == "" || l.Rel == "alternate" || l.Rel == "self" {
			item.Links = append(item.Links, l.Href)
		}
	}

	// Published falls back to the update time when absent.
	item.Published = entry.Published
	if item.Published == "" {
		item.Published = entry.Updated
	}
	item.PublishedParsed = entry.PublishedParsed
	if item.PublishedParsed == nil {
		item.PublishedParsed = entry.UpdatedParsed
	}

	if len(entry.Authors) > 0 {
		item.Authors = atomPersons(entry.Authors)
		item.Author = item.Authors[0]
	}

	item.Categories = atomCategories(entry.Categories)

	for _, l := range entry.Links {
		if l.Rel == "enclosure" {
			item.Enclosures = append(item.Enclosures, &Enclosure{
				URL:    l.Href,
				Length: l.Length,
				Type:   l.Type,
			})
		}
	}

	return item
}

// firstLinkWithRel returns the first link carrying the given rel, or nil.
func firstLinkWithRel(rel string, links []*atom.Link) *atom.Link {
	for _, link := range links {
		if link.Rel == rel {
			return link
		}
	}
	return nil
}

// atomPersons converts atom persons to universal Persons.
func atomPersons(persons []*atom.Person) []*Person {
	if persons == nil {
		return nil
	}
	out := make([]*Person, 0, len(persons))
	for _, p := range persons {
		out = append(out, &Person{Name: p.Name, Email: p.Email})
	}
	return out
}

// atomCategories flattens atom categories, using the label when present and
// the term otherwise.
func atomCategories(categories []*atom.Category) []string {
	if categories == nil {
		return nil
	}
	out := make([]string, 0, len(categories))
	for _, c := range categories {
		if c.Label != "" {
			out = append(out, c.Label)
		} else {
			out = append(out, c.Term)
		}
	}
	return out
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
	jsonFeed, found := feed.(*json.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *json.Feed")
	}

	result := &Feed{
		FeedVersion: jsonFeed.Version,
		Title:       jsonFeed.Title,
		Link:        jsonFeed.HomePageURL,
		FeedLink:    jsonFeed.FeedURL,
		Description: jsonFeed.Description,
		Language:    jsonFeed.Language,
		FeedType:    "json",
	}

	if jsonFeed.HomePageURL != "" {
		result.Links = append(result.Links, jsonFeed.HomePageURL)
	}
	if jsonFeed.FeedURL != "" {
		result.Links = append(result.Links, jsonFeed.FeedURL)
	}

	// The icon is used over json.Feed.Image: the spec describes it as the
	// square image suitable for a timeline, which is what Feed.Image holds.
	if jsonFeed.Icon != "" {
		result.Image = &Image{URL: jsonFeed.Icon}
	}

	if jsonFeed.Author != nil {
		result.Author = personFromText(jsonFeed.Author.Name)
	}
	if jsonFeed.Authors != nil {
		result.Authors = jsonPersons(jsonFeed.Authors)
	} else if result.Author != nil {
		result.Authors = []*Person{result.Author}
	}

	// The feed-level times mirror the first (most recent) item's.
	if len(jsonFeed.Items) > 0 {
		result.Updated = jsonFeed.Items[0].DateModified
		if date, err := shared.ParseDate(result.Updated); err == nil {
			result.UpdatedParsed = &date
		}
		result.Published = jsonFeed.Items[0].DatePublished
		if date, err := shared.ParseDate(result.Published); err == nil {
			result.PublishedParsed = &date
		}
	}

	result.Items = make([]*Item, 0, len(jsonFeed.Items))
	for _, i := range jsonFeed.Items {
		result.Items = append(result.Items, t.translateFeedItem(i))
	}

	// TODO UserComment is missing in global Feed
	// TODO NextURL is missing in global Feed
	// TODO Favicon is missing in global Feed
	// TODO Exipred is missing in global Feed
	// TODO Hubs is not supported in json.Feed
	// TODO Extensions is not supported in json.Feed
	return result, nil
}

func (t *DefaultJSONTranslator) translateFeedItem(jsonItem *json.Item) *Item {
	item := &Item{
		GUID:        jsonItem.ID,
		Link:        jsonItem.URL,
		Title:       jsonItem.Title,
		Description: jsonItem.Summary,
		Published:   jsonItem.DatePublished,
		Updated:     jsonItem.DateModified,
	}

	if jsonItem.URL != "" {
		item.Links = append(item.Links, jsonItem.URL)
	}
	if jsonItem.ExternalURL != "" {
		item.Links = append(item.Links, jsonItem.ExternalURL)
	}

	item.Content = jsonItem.ContentHTML
	if item.Content == "" {
		item.Content = jsonItem.ContentText
	}

	if jsonItem.Image != "" {
		item.Image = &Image{URL: jsonItem.Image}
	} else if jsonItem.BannerImage != "" {
		item.Image = &Image{URL: jsonItem.BannerImage}
	}

	if jsonItem.DatePublished != "" {
		if date, err := shared.ParseDate(jsonItem.DatePublished); err == nil {
			item.PublishedParsed = &date
		}
	}
	if jsonItem.DateModified != "" {
		if date, err := shared.ParseDate(jsonItem.DateModified); err == nil {
			item.UpdatedParsed = &date
		}
	}

	if jsonItem.Author != nil {
		item.Author = personFromText(jsonItem.Author.Name)
	}
	if jsonItem.Authors != nil {
		item.Authors = jsonPersons(jsonItem.Authors)
	} else if item.Author != nil {
		item.Authors = []*Person{item.Author}
	}

	if len(jsonItem.Tags) > 0 {
		item.Categories = jsonItem.Tags
	}

	if jsonItem.Attachments != nil {
		for _, attachment := range *jsonItem.Attachments {
			e := &Enclosure{
				URL:  attachment.URL,
				Type: attachment.MimeType,
			}
			// RSS enclosure length is the size in bytes, not the duration.
			if attachment.SizeInBytes > 0 {
				e.Length = fmt.Sprintf("%d", attachment.SizeInBytes)
			}
			// Title is not defined in global enclosure
			item.Enclosures = append(item.Enclosures, e)
		}
	}

	// TODO ExternalURL is missing in global Feed
	// TODO BannerImage is missing in global Feed
	// Author.URL and Author.Avatar are missing in global feed
	return item
}

// jsonPersons converts json feed authors to universal Persons, splitting
// free-form "Name (email)" values.
func jsonPersons(authors []*json.Author) []*Person {
	out := make([]*Person, 0, len(authors))
	for _, a := range authors {
		out = append(out, personFromText(a.Name))
	}
	return out
}
