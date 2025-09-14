package gofeed

import (
	"fmt"
	"strconv"

	"github.com/mmcdole/gofeed/atom"
	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/mmcdole/gofeed/json"
	"github.com/mmcdole/gofeed/rss"
)

// Renderers convert a universal Feed struct into format-specific feed structures.
// These are the reverse of the Translator interfaces.
//
// IMPORTANT LIMITATIONS:
// Perfect round-trip conversion (parse -> translate -> render -> parse) is not
// always possible.
//
// RSS limitations:
//    - Complex extension-based author fallbacks may not preserve exact original location
//
// JSON Feed limitations:
//    - Feed-level dates are derived from first item in translator, not used in renderer
//    - Some JSON Feed specific fields (UserComment, NextURL, etc.) are not supported

// RSSRenderer converts a Feed struct into an RSS feed structure
type RSSRenderer struct{}

// Render converts the universal Feed into an RSS feed
func (r *RSSRenderer) Render(feed *Feed) (*rss.Feed, error) {
	if feed == nil {
		return nil, fmt.Errorf("feed cannot be nil")
	}

	rssFeed := &rss.Feed{}
	rssFeed.Title = feed.Title
	rssFeed.Description = feed.Description
	rssFeed.Link = feed.Link
	rssFeed.Links = make([]string, len(feed.Links))
	copy(rssFeed.Links, feed.Links)
	rssFeed.Language = feed.Language
	rssFeed.Copyright = feed.Copyright
	rssFeed.Generator = feed.Generator
	rssFeed.Version = "2.0"

	// Handle dates
	rssFeed.LastBuildDate = feed.Updated
	rssFeed.LastBuildDateParsed = feed.UpdatedParsed
	rssFeed.PubDate = feed.Published
	rssFeed.PubDateParsed = feed.PublishedParsed

	// Handle extensions first
	rssFeed.DublinCoreExt = feed.DublinCoreExt
	rssFeed.ITunesExt = feed.ITunesExt
	rssFeed.Extensions = feed.Extensions

	// Handle author information with extension fallbacks
	var feedAuthor *Person
	if len(feed.Authors) > 0 {
		feedAuthor = feed.Authors[0]
	} else if feed.Author != nil {
		feedAuthor = feed.Author
	}

	if feedAuthor != nil {
		rssFeed.ManagingEditor = r.FormatPersonForRSS(feedAuthor)

		// Also populate DublinCore Creator to improve round-trip fidelity
		if rssFeed.DublinCoreExt == nil {
			rssFeed.DublinCoreExt = &ext.DublinCoreExtension{}
		}
		// Only add if not already present
		dcAuthor := r.FormatPersonForRSS(feedAuthor)
		creatorExists := false
		for _, creator := range rssFeed.DublinCoreExt.Creator {
			if creator == dcAuthor {
				creatorExists = true
				break
			}
		}
		if !creatorExists {
			rssFeed.DublinCoreExt.Creator = append(rssFeed.DublinCoreExt.Creator, dcAuthor)
		}
	}

	// Handle image
	if feed.Image != nil {
		rssFeed.Image = &rss.Image{
			URL:   feed.Image.URL,
			Title: feed.Image.Title,
			Link:  feed.Link, // RSS spec requires link in image
		}
	}

	// Handle categories
	if len(feed.Categories) > 0 {
		rssFeed.Categories = make([]*rss.Category, len(feed.Categories))
		for i, cat := range feed.Categories {
			rssFeed.Categories[i] = &rss.Category{Value: cat}
		}
	}

	// Handle items
	rssFeed.Items = make([]*rss.Item, len(feed.Items))
	for i, item := range feed.Items {
		rssFeed.Items[i] = r.renderItem(item)
	}

	return rssFeed, nil
}

func (r *RSSRenderer) renderItem(item *Item) *rss.Item {
	rssItem := &rss.Item{}
	rssItem.Title = item.Title
	rssItem.Description = item.Description
	rssItem.Content = item.Content
	rssItem.Link = item.Link
	rssItem.Links = make([]string, len(item.Links))
	copy(rssItem.Links, item.Links)
	rssItem.PubDate = item.Published
	rssItem.PubDateParsed = item.PublishedParsed

	// Handle extensions first
	rssItem.DublinCoreExt = item.DublinCoreExt
	rssItem.ITunesExt = item.ITunesExt
	rssItem.Extensions = item.Extensions
	rssItem.Custom = item.Custom

	// Handle author with extension fallbacks
	var itemAuthor *Person
	if len(item.Authors) > 0 {
		itemAuthor = item.Authors[0]
	} else if item.Author != nil {
		itemAuthor = item.Author
	}

	if itemAuthor != nil {
		rssItem.Author = r.FormatPersonForRSS(itemAuthor)

		// Also populate DublinCore Creator to improve round-trip fidelity
		if rssItem.DublinCoreExt == nil {
			rssItem.DublinCoreExt = &ext.DublinCoreExtension{}
		}
		// Only add if not already present
		dcAuthor := r.FormatPersonForRSS(itemAuthor)
		creatorExists := false
		for _, creator := range rssItem.DublinCoreExt.Creator {
			if creator == dcAuthor {
				creatorExists = true
				break
			}
		}
		if !creatorExists {
			rssItem.DublinCoreExt.Creator = append(rssItem.DublinCoreExt.Creator, dcAuthor)
		}

		// Also populate iTunes Author if iTunes extension exists
		if rssItem.ITunesExt != nil && rssItem.ITunesExt.Author == "" {
			rssItem.ITunesExt.Author = r.FormatPersonForRSS(itemAuthor)
		}
	}

	// Handle GUID
	if item.GUID != "" {
		rssItem.GUID = &rss.GUID{Value: item.GUID}
	}

	// Handle categories
	if len(item.Categories) > 0 {
		rssItem.Categories = make([]*rss.Category, len(item.Categories))
		for i, cat := range item.Categories {
			rssItem.Categories[i] = &rss.Category{Value: cat}
		}
	}

	// Handle enclosures and item image
	enclosures := make([]*rss.Enclosure, 0, len(item.Enclosures)+1)

	// Add existing enclosures
	for _, enc := range item.Enclosures {
		enclosures = append(enclosures, &rss.Enclosure{
			URL:    enc.URL,
			Type:   enc.Type,
			Length: enc.Length,
		})
	}

	// Add item image as enclosure if not already present
	// (this is how RSS translator extracts item images)
	if item.Image != nil {
		imageAlreadyExists := false
		for _, enc := range enclosures {
			if enc.URL == item.Image.URL {
				imageAlreadyExists = true
				break
			}
		}

		if !imageAlreadyExists {
			enclosures = append(enclosures, &rss.Enclosure{
				URL: item.Image.URL,
				// Type and Length omitted - don't guess unknown values
			})
		}
	}

	if len(enclosures) > 0 {
		rssItem.Enclosures = enclosures
		// Set first enclosure as the primary one for RSS compatibility
		rssItem.Enclosure = enclosures[0]
	}

	// Handle item Updated field via DublinCore extension (RSS doesn't have native updated field)
	if item.Updated != "" {
		if rssItem.DublinCoreExt == nil {
			rssItem.DublinCoreExt = &ext.DublinCoreExtension{}
		}
		// Only add if not already present in DublinCore Date
		dateExists := false
		for _, date := range rssItem.DublinCoreExt.Date {
			if date == item.Updated {
				dateExists = true
				break
			}
		}
		if !dateExists {
			rssItem.DublinCoreExt.Date = append(rssItem.DublinCoreExt.Date, item.Updated)
		}
	}

	return rssItem
}

func (r *RSSRenderer) FormatPersonForRSS(person *Person) string {
	if person.Email != "" && person.Name != "" {
		return person.Email + " (" + person.Name + ")"
	} else if person.Email != "" {
		return person.Email
	} else if person.Name != "" {
		return person.Name
	}
	return ""
}

// AtomRenderer converts a Feed struct into an Atom feed structure
type AtomRenderer struct{}

// Render converts the universal Feed into an Atom feed
func (r *AtomRenderer) Render(feed *Feed) (*atom.Feed, error) {
	if feed == nil {
		return nil, fmt.Errorf("feed cannot be nil")
	}

	atomFeed := &atom.Feed{}
	atomFeed.Title = feed.Title
	atomFeed.Subtitle = feed.Description
	atomFeed.Language = feed.Language
	atomFeed.Rights = feed.Copyright
	atomFeed.Updated = feed.Updated
	atomFeed.UpdatedParsed = feed.UpdatedParsed
	atomFeed.Version = "1.0"

	// Generate an ID if not present
	atomFeed.ID = feed.FeedLink
	if atomFeed.ID == "" {
		atomFeed.ID = feed.Link
	}

	// Handle links
	if feed.Link != "" {
		atomFeed.Links = append(atomFeed.Links, &atom.Link{
			Href: feed.Link,
			Rel:  "alternate",
			Type: "text/html",
		})
	}
	if feed.FeedLink != "" {
		atomFeed.Links = append(atomFeed.Links, &atom.Link{
			Href: feed.FeedLink,
			Rel:  "self",
			Type: "application/atom+xml",
		})
	}

	// Handle additional links
	for _, link := range feed.Links {
		if link != feed.Link && link != feed.FeedLink {
			atomFeed.Links = append(atomFeed.Links, &atom.Link{
				Href: link,
				Rel:  "alternate",
			})
		}
	}

	// Handle authors
	if len(feed.Authors) > 0 {
		atomFeed.Authors = make([]*atom.Person, len(feed.Authors))
		for i, author := range feed.Authors {
			atomFeed.Authors[i] = &atom.Person{
				Name:  author.Name,
				Email: author.Email,
			}
		}
	} else if feed.Author != nil {
		atomFeed.Authors = []*atom.Person{
			{
				Name:  feed.Author.Name,
				Email: feed.Author.Email,
			},
		}
	}

	// Handle images - support both Icon and Logo for full round-trip fidelity
	if feed.Image != nil {
		// Check if we have icon data in Custom field (indicates Logo is primary)
		if feed.Custom != nil && feed.Custom[CustomAtomIcon] != "" {
			// Primary image is Logo, secondary is Icon
			atomFeed.Logo = feed.Image.URL
			atomFeed.Icon = feed.Custom[CustomAtomIcon]
		} else {
			// Only one image exists, default to Logo
			atomFeed.Logo = feed.Image.URL
		}
	}

	if feed.GeneratorDetail != nil {
		atomFeed.Generator = &atom.Generator{
			Value:   feed.GeneratorDetail.Value,
			URI:     feed.GeneratorDetail.URI,
			Version: feed.GeneratorDetail.Version,
		}
	} else if feed.Generator != "" {
		atomFeed.Generator = &atom.Generator{Value: feed.Generator}
	}

	// Handle categories
	if len(feed.Categories) > 0 {
		atomFeed.Categories = make([]*atom.Category, len(feed.Categories))
		for i, cat := range feed.Categories {
			atomFeed.Categories[i] = &atom.Category{
				Term:  cat,
				Label: cat,
			}
		}
	}

	// Handle entries
	atomFeed.Entries = make([]*atom.Entry, len(feed.Items))
	for i, item := range feed.Items {
		atomFeed.Entries[i] = r.renderEntry(item)
	}

	// Handle extensions
	atomFeed.Extensions = feed.Extensions

	return atomFeed, nil
}

func (r *AtomRenderer) renderEntry(item *Item) *atom.Entry {
	entry := &atom.Entry{}
	entry.Title = item.Title
	entry.Summary = item.Description
	entry.Updated = item.Updated
	entry.UpdatedParsed = item.UpdatedParsed
	entry.Published = item.Published
	entry.PublishedParsed = item.PublishedParsed

	// Use GUID as ID, or generate from link
	entry.ID = item.GUID
	if entry.ID == "" {
		entry.ID = item.Link
	}

	// Handle content
	if item.Content != "" {
		entry.Content = &atom.Content{
			Value: item.Content,
			Type:  "html",
		}
	}

	// Handle links
	if item.Link != "" {
		entry.Links = append(entry.Links, &atom.Link{
			Href: item.Link,
			Rel:  "alternate",
			Type: "text/html",
		})
	}

	// Handle additional links
	for _, link := range item.Links {
		if link != item.Link {
			entry.Links = append(entry.Links, &atom.Link{
				Href: link,
				Rel:  "alternate",
			})
		}
	}

	// Handle enclosures
	for _, enc := range item.Enclosures {
		entry.Links = append(entry.Links, &atom.Link{
			Href:   enc.URL,
			Rel:    "enclosure",
			Type:   enc.Type,
			Length: enc.Length,
		})
	}

	// Handle item image as enclosure link if not already present
	// (Atom doesn't have native item image support)
	if item.Image != nil {
		imageAlreadyExists := false
		for _, link := range entry.Links {
			if link.Href == item.Image.URL {
				imageAlreadyExists = true
				break
			}
		}

		if !imageAlreadyExists {
			entry.Links = append(entry.Links, &atom.Link{
				Href: item.Image.URL,
				Rel:  "enclosure",
				// Type: omitted - don't guess MIME type for arbitrary URLs
			})
		}
	}

	// Handle authors
	if len(item.Authors) > 0 {
		entry.Authors = make([]*atom.Person, len(item.Authors))
		for i, author := range item.Authors {
			entry.Authors[i] = &atom.Person{
				Name:  author.Name,
				Email: author.Email,
			}
		}
	} else if item.Author != nil {
		entry.Authors = []*atom.Person{
			{
				Name:  item.Author.Name,
				Email: item.Author.Email,
			},
		}
	}

	// Handle categories
	if len(item.Categories) > 0 {
		entry.Categories = make([]*atom.Category, len(item.Categories))
		for i, cat := range item.Categories {
			entry.Categories[i] = &atom.Category{
				Term:  cat,
				Label: cat,
			}
		}
	}

	// Handle extensions
	entry.Extensions = item.Extensions

	return entry
}

// JSONRenderer converts a Feed struct into a JSON feed structure
type JSONRenderer struct{}

// Render converts the universal Feed into a JSON feed
func (r *JSONRenderer) Render(feed *Feed) (*json.Feed, error) {
	if feed == nil {
		return nil, fmt.Errorf("feed cannot be nil")
	}

	jsonFeed := &json.Feed{}
	jsonFeed.Version = "https://jsonfeed.org/version/1.1"
	jsonFeed.Title = feed.Title
	jsonFeed.Description = feed.Description
	jsonFeed.HomePageURL = feed.Link
	jsonFeed.FeedURL = feed.FeedLink
	jsonFeed.Language = feed.Language

	// Handle image - prefer Icon over Logo for JSON feeds
	if feed.Image != nil {
		jsonFeed.Icon = feed.Image.URL
	}

	// Handle author
	if len(feed.Authors) > 0 {
		// Use the first author for the feed-level author
		jsonFeed.Author = &json.Author{
			Name: feed.Authors[0].Name,
		}
		if feed.Authors[0].Email != "" && feed.Authors[0].Name != "" {
			jsonFeed.Author.Name = feed.Authors[0].Name + " <" + feed.Authors[0].Email + ">"
		} else if feed.Authors[0].Email != "" {
			jsonFeed.Author.Name = feed.Authors[0].Email
		}

		// Handle multiple authors (JSON Feed v1.1)
		jsonFeed.Authors = make([]*json.Author, len(feed.Authors))
		for i, author := range feed.Authors {
			jsonFeed.Authors[i] = &json.Author{
				Name: author.Name,
			}
			if author.Email != "" && author.Name != "" {
				jsonFeed.Authors[i].Name = author.Name + " <" + author.Email + ">"
			} else if author.Email != "" {
				jsonFeed.Authors[i].Name = author.Email
			}
		}
	} else if feed.Author != nil {
		jsonFeed.Author = &json.Author{
			Name: feed.Author.Name,
		}
		// JSON Feed doesn't support email in author, but we could put it in name
		if feed.Author.Email != "" && feed.Author.Name != "" {
			jsonFeed.Author.Name = feed.Author.Name + " <" + feed.Author.Email + ">"
		} else if feed.Author.Email != "" {
			jsonFeed.Author.Name = feed.Author.Email
		}
	}

	// Handle items
	jsonFeed.Items = make([]*json.Item, len(feed.Items))
	for i, item := range feed.Items {
		jsonFeed.Items[i] = r.renderItem(item)
	}

	return jsonFeed, nil
}

func (r *JSONRenderer) renderItem(item *Item) *json.Item {
	jsonItem := &json.Item{}
	jsonItem.ID = item.GUID
	if jsonItem.ID == "" {
		jsonItem.ID = item.Link
	}
	jsonItem.URL = item.Link
	jsonItem.Title = item.Title
	jsonItem.Summary = item.Description
	jsonItem.DatePublished = item.Published
	jsonItem.DateModified = item.Updated
	jsonItem.Language = "" // Item-level language not available in universal Feed

	// Handle content - prefer HTML over text
	if item.Content != "" {
		jsonItem.ContentHTML = item.Content
	} else if item.Description != "" {
		jsonItem.ContentText = item.Description
	}

	// Handle image
	if item.Image != nil {
		jsonItem.Image = item.Image.URL
	}

	// Handle author
	if len(item.Authors) > 0 {
		// Use the first author for the item-level author
		jsonItem.Author = &json.Author{
			Name: item.Authors[0].Name,
		}
		if item.Authors[0].Email != "" && item.Authors[0].Name != "" {
			jsonItem.Author.Name = item.Authors[0].Name + " <" + item.Authors[0].Email + ">"
		} else if item.Authors[0].Email != "" {
			jsonItem.Author.Name = item.Authors[0].Email
		}

		// Handle multiple authors (JSON Feed v1.1)
		jsonItem.Authors = make([]*json.Author, len(item.Authors))
		for i, author := range item.Authors {
			jsonItem.Authors[i] = &json.Author{
				Name: author.Name,
			}
			if author.Email != "" && author.Name != "" {
				jsonItem.Authors[i].Name = author.Name + " <" + author.Email + ">"
			} else if author.Email != "" {
				jsonItem.Authors[i].Name = author.Email
			}
		}
	} else if item.Author != nil {
		jsonItem.Author = &json.Author{
			Name: item.Author.Name,
		}
		if item.Author.Email != "" && item.Author.Name != "" {
			jsonItem.Author.Name = item.Author.Name + " <" + item.Author.Email + ">"
		} else if item.Author.Email != "" {
			jsonItem.Author.Name = item.Author.Email
		}
	}

	// Handle categories as tags
	if len(item.Categories) > 0 {
		jsonItem.Tags = make([]string, len(item.Categories))
		copy(jsonItem.Tags, item.Categories)
	}

	// Handle enclosures as attachments
	if len(item.Enclosures) > 0 {
		attachments := make([]json.Attachments, len(item.Enclosures))
		for i, enc := range item.Enclosures {
			attachments[i] = json.Attachments{
				URL:      enc.URL,
				MimeType: enc.Type,
			}
			// Convert length string to duration if it's numeric
			if enc.Length != "" {
				if duration, err := strconv.ParseInt(enc.Length, 10, 64); err == nil {
					attachments[i].DurationInSeconds = duration
				}
			}
		}
		jsonItem.Attachments = &attachments
	}

	return jsonItem
}
