package gofeed

import (
	"encoding/json"
	"time"

	ext "github.com/mmcdole/gofeed/v2/extensions"
)

// Feed is the universal Feed type that atom.Feed
// and rss.Feed gets translated to. It represents
// a web feed.
// Sorting with sort.Sort will order the Items by
// oldest to newest publish time.
type Feed struct {
	Title           string                   `json:"title,omitempty"`
	Description     string                   `json:"description,omitempty"`
	Link            string                   `json:"link,omitempty"`
	FeedLink        string                   `json:"feedLink,omitempty"`
	Links           []string                 `json:"links,omitempty"`
	Updated         string                   `json:"updated,omitempty"`
	UpdatedParsed   *time.Time               `json:"updatedParsed,omitempty"`
	Published       string                   `json:"published,omitempty"`
	PublishedParsed *time.Time               `json:"publishedParsed,omitempty"`
	Author          *Person                  `json:"author,omitempty"` // Deprecated: Use feed.Authors instead
	Authors         []*Person                `json:"authors,omitempty"`
	Language        string                   `json:"language,omitempty"`
	Image           *Image                   `json:"image,omitempty"`
	Copyright       string                   `json:"copyright,omitempty"`
	Generator       string                   `json:"generator,omitempty"`
	Categories      []string                 `json:"categories,omitempty"`
	DublinCoreExt   *ext.DublinCoreExtension `json:"dcExt,omitempty"`
	ITunesExt       *ext.ITunesFeedExtension `json:"itunesExt,omitempty"`
	Extensions      ext.Extensions           `json:"extensions,omitempty"`
	Items           []*Item                  `json:"items"`
	FeedType        string                   `json:"feedType"`
	FeedVersion     string                   `json:"feedVersion"`
	
	// Original format-specific feed data (only populated if KeepOriginalFeed is true)
	OriginalFeed interface{} `json:"-"`
}

// String returns a JSON representation of the Feed for debugging purposes.
func (f Feed) String() string {
	json, _ := json.MarshalIndent(f, "", "    ")
	return string(json)
}

// GetExtension retrieves extension values by namespace and element name.
// Returns a slice of Extension structs for the given namespace and element.
// For non-namespaced RSS elements, use "rss" as the namespace.
func (f *Feed) GetExtension(namespace, element string) []ext.Extension {
	if f.Extensions == nil {
		return nil
	}
	
	nsMap, ok := f.Extensions[namespace]
	if !ok {
		return nil
	}
	
	return nsMap[element]
}

// GetExtensionValue is a convenience method that returns the text value
// of the first matching extension element, or empty string if not found.
func (f *Feed) GetExtensionValue(namespace, element string) string {
	exts := f.GetExtension(namespace, element)
	if len(exts) == 0 {
		return ""
	}
	return exts[0].Value
}

// GetCustomValue retrieves the text value of a non-namespaced custom element at the feed level.
// Returns empty string if the element is not found.
func (f *Feed) GetCustomValue(element string) string {
	return f.GetExtensionValue("_custom", element)
}

// Item is the universal Item type that atom.Entry
// and rss.Item gets translated to.  It represents
// a single entry in a given feed.
type Item struct {
	Title           string                   `json:"title,omitempty"`
	Description     string                   `json:"description,omitempty"`
	Content         string                   `json:"content,omitempty"`
	Link            string                   `json:"link,omitempty"`
	Links           []string                 `json:"links,omitempty"`
	Updated         string                   `json:"updated,omitempty"`
	UpdatedParsed   *time.Time               `json:"updatedParsed,omitempty"`
	Published       string                   `json:"published,omitempty"`
	PublishedParsed *time.Time               `json:"publishedParsed,omitempty"`
	Author          *Person                  `json:"author,omitempty"` // Deprecated: Use item.Authors instead
	Authors         []*Person                `json:"authors,omitempty"`
	GUID            string                   `json:"guid,omitempty"`
	Image           *Image                   `json:"image,omitempty"`
	Categories      []string                 `json:"categories,omitempty"`
	Enclosures      []*Enclosure             `json:"enclosures,omitempty"`
	DublinCoreExt   *ext.DublinCoreExtension `json:"dcExt,omitempty"`
	ITunesExt       *ext.ITunesItemExtension `json:"itunesExt,omitempty"`
	Extensions      ext.Extensions           `json:"extensions,omitempty"`
}

// GetExtension retrieves extension values by namespace and element name.
// Returns a slice of Extension structs for the given namespace and element.
// For non-namespaced RSS elements, use "rss" as the namespace.
func (i *Item) GetExtension(namespace, element string) []ext.Extension {
	if i.Extensions == nil {
		return nil
	}
	
	nsMap, ok := i.Extensions[namespace]
	if !ok {
		return nil
	}
	
	return nsMap[element]
}

// GetExtensionValue is a convenience method that returns the text value
// of the first matching extension element, or empty string if not found.
func (i *Item) GetExtensionValue(namespace, element string) string {
	exts := i.GetExtension(namespace, element)
	if len(exts) == 0 {
		return ""
	}
	return exts[0].Value
}

// GetCustomValue retrieves the text value of a non-namespaced custom element.
// This is a convenience method that replaces the previous Item.Custom[key] access pattern.
// Returns empty string if the element is not found.
func (i *Item) GetCustomValue(element string) string {
	return i.GetExtensionValue("_custom", element)
}

// Person is an individual specified in a feed
// (e.g. an author)
type Person struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// Image is an image that is the artwork for a given
// feed or item.
type Image struct {
	URL   string `json:"url,omitempty"`
	Title string `json:"title,omitempty"`
}

// Enclosure is a file associated with a given Item.
type Enclosure struct {
	URL    string `json:"url,omitempty"`
	Length string `json:"length,omitempty"`
	Type   string `json:"type,omitempty"`
}

// Len returns the length of Items.
func (f Feed) Len() int {
	return len(f.Items)
}

// Less compares PublishedParsed of Items[i], Items[k]
// and returns true if Items[i] is less than Items[k].
func (f Feed) Less(i, k int) bool {
	iParsed := f.Items[i].PublishedParsed
	kParsed := f.Items[k].PublishedParsed
	
	if iParsed == nil && kParsed == nil {
		return false
	}
	if iParsed == nil {
		return true
	}
	if kParsed == nil {
		return false
	}
	return iParsed.Before(*kParsed)
}

// Swap swaps Items[i] and Items[k].
func (f Feed) Swap(i, k int) {
	f.Items[i], f.Items[k] = f.Items[k], f.Items[i]
}
