package atom

import (
	"time"

	"github.com/mmcdole/gofeed/feed"
)

type Feed struct {
	Title         string              `json:"title"`
	ID            string              `json:"id"`
	Updated       string              `json:"updated"`
	UpdatedParsed *time.Time          `json:"updatedParsed,omitempty"`
	Subtitle      string              `json:"subtitle,omitempty"`
	Link          *Link               `json:"link,omitempty"`
	Generator     *Generator          `json:"generator,omitempty"`
	Icon          string              `json:"icon,omitempty"`
	Logo          string              `json:"logo,omitempty"`
	Rights        string              `json:"rights,omitempty"`
	Contributors  []*Person           `json:"contributors,omitempty"`
	Authors       []*Person           `json:"authors,omitempty"`
	Categories    []*Category         `json:"categories,omitempty"`
	Source        *Source             `json:"source,omitempty"`
	Entries       []*Entry            `json:"entries"`
	Extensions    feed.FeedExtensions `json:"extensions"`
	Version       string              `json:"version"`
}

type Entry struct {
	Title           string              `json:"title"`
	ID              string              `json:"id"`
	Updated         string              `json:"updated"`
	UpdatedParsed   *time.Time          `json:"updatedParsed,omitempty"`
	Authors         []*Person           `json:"authors,omitempty"`
	Contributors    []*Person           `json:"contributors,omitempty"`
	Categories      []*Category         `json:"categories,omitempty"`
	Link            *Link               `json:"link,omitempty"`
	Published       string              `json:"published,omitempty"`
	PublishedParsed *time.Time          `json:"publishedParsed,omitempty"`
	Content         *Content            `json:"content,omitempty"`
	Extensions      feed.FeedExtensions `json:"extensions"`
}

type Category struct {
	Term   string `json:"term"`
	Scheme string `json:"scheme,omitempty"`
	Label  string `json:"label,omitempty"`
}

type Person struct {
	Name  string `json:'name'`
	Email string `json:"email,omitempty"`
	URI   string `json:"uri,omitempty"`
}

type Link struct {
	Href     string `json:"href"`
	Hreflang string `json:"hreflang,omitempty"`
	Rel      string `json:"rel,omitempty"`
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Length   string `json:"length,omitempty"`
}

type Content struct {
	Src   string `json:"src,omitempty"`
	Value string `json:"value"`
}

type Generator struct {
	Value   string `json:"value"`
	URI     string `json:"uri,omitempty"`
	Version string `json:"version,omitempty"`
}

type Source struct {
	Title         string              `json:"title"`
	ID            string              `json:"id"`
	Updated       string              `json:"updated"`
	UpdatedParsed *time.Time          `json:"updatedParsed,omitempty"`
	Subtitle      string              `json:"subtitle,omitempty"`
	Link          *Link               `json:"link,omitempty"`
	Generator     *Generator          `json:"generator,omitempty"`
	Icon          string              `json:"icon,omitempty"`
	Logo          string              `json:"logo,omitempty"`
	Rights        string              `json:"rights,omitempty"`
	Contributors  []*Person           `json:"contributors,omitempty"`
	Authors       []*Person           `json:"authors,omitempty"`
	Categories    []*Category         `json:"categories,omitempty"`
	Extensions    feed.FeedExtensions `json:"extensions"`
}
