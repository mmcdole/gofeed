package atom

import (
	"encoding/json"
	"encoding/xml"
	"time"

	"github.com/mmcdole/gofeed/extensions"
)

// Feed is an Atom Feed
type Feed struct {
	XMLName       xml.Name       `xml:"http://www.w3.org/2005/Atom feed" json:"-"`
	Title         string         `json:"title,omitempty" xml:"title"`
	ID            string         `json:"id,omitempty" xml:"id"`
	Updated       string         `json:"updated,omitempty" xml:"updated"`
	UpdatedParsed *time.Time     `json:"updatedParsed,omitempty" xml:"-"`
	Subtitle      string         `json:"subtitle,omitempty" xml:"subtitle,omitempty"`
	Links         []*Link        `json:"links,omitempty" xml:"link"`
	Language      string         `json:"language,omitempty" xml:"lang,attr,omitempty"`
	Generator     *Generator     `json:"generator,omitempty" xml:"generator"`
	Icon          string         `json:"icon,omitempty" xml:"icon,omitempty"`
	Logo          string         `json:"logo,omitempty" xml:"logo,omitempty"`
	Rights        string         `json:"rights,omitempty" xml:"rights,omitempty"`
	Contributors  []*Person      `json:"contributors,omitempty" xml:"contributor"`
	Authors       []*Person      `json:"authors,omitempty" xml:"author"`
	Categories    []*Category    `json:"categories,omitempty" xml:"category"`
	Entries       []*Entry       `json:"entries" xml:"entry"`
	Extensions    ext.Extensions `json:"extensions,omitempty" xml:"-"`
	Version       string         `json:"version" xml:"-"`
}

func (f Feed) String() string {
	json, _ := json.MarshalIndent(f, "", "    ")
	return string(json)
}

// Entry is an Atom Entry
type Entry struct {
	Title           string         `json:"title,omitempty" xml:"title"`
	ID              string         `json:"id,omitempty" xml:"id"`
	Updated         string         `json:"updated,omitempty" xml:"updated"`
	UpdatedParsed   *time.Time     `json:"updatedParsed,omitempty" xml:"-"`
	Summary         string         `json:"summary,omitempty" xml:"summary"`
	Authors         []*Person      `json:"authors,omitempty" xml:"author"`
	Contributors    []*Person      `json:"contributors,omitempty" xml:"contributor"`
	Categories      []*Category    `json:"categories,omitempty" xml:"category"`
	Links           []*Link        `json:"links,omitempty" xml:"link"`
	Rights          string         `json:"rights,omitempty" xml:"rights,omitempty"`
	Published       string         `json:"published,omitempty" xml:"published"`
	PublishedParsed *time.Time     `json:"publishedParsed,omitempty" xml:"-"`
	Source          *Source        `json:"source,omitempty" xml:"source"`
	Content         *Content       `json:"content,omitempty" xml:"content"`
	Extensions      ext.Extensions `json:"extensions,omitempty" xml:"-"`
}

// Category is category metadata for Feeds and Entries
type Category struct {
	Term   string `json:"term,omitempty" xml:"term,attr"`
	Scheme string `json:"scheme,omitempty" xml:"scheme,attr"`
	Label  string `json:"label,omitempty" xml:"label,attr"`
}

// Person represents a person in an Atom feed
// for things like Authors, Contributors, etc
type Person struct {
	Name  string `json:"name,omitempty" xml:"name"`
	Email string `json:"email,omitempty" xml:"email"`
	URI   string `json:"uri,omitempty" xml:"uri"`
}

// Link is an Atom link that defines a reference
// from an entry or feed to a Web resource
type Link struct {
	Href     string `json:"href,omitempty" xml:"href,attr"`
	Hreflang string `json:"hreflang,omitempty" xml:"hreflang,attr,omitempty"`
	Rel      string `json:"rel,omitempty" xml:"rel,attr"`
	Type     string `json:"type,omitempty" xml:"type,attr"`
	Title    string `json:"title,omitempty" xml:"title,attr,omitempty"`
	Length   string `json:"length,omitempty" xml:"length,attr,omitempty"`
}

// Content either contains or links to the content of
// the entry
type Content struct {
	Src   string `json:"src,omitempty" xml:"src,attr"`
	Type  string `json:"type,omitempty" xml:"type,attr"`
	Value string `json:"value,omitempty" xml:",chardata"`
}

// Generator identifies the agent used to generate a
// feed, for debugging and other purposes.
type Generator struct {
	Value   string `json:"value,omitempty" xml:",chardata"`
	URI     string `json:"uri,omitempty" xml:"uri,attr"`
	Version string `json:"version,omitempty" xml:"version,attr"`
}

// Source contains the feed information for another
// feed if a given entry came from that feed.
type Source struct {
	Title         string         `json:"title,omitempty" xml:"title"`
	ID            string         `json:"id,omitempty" xml:"id"`
	Updated       string         `json:"updated,omitempty" xml:"updated"`
	UpdatedParsed *time.Time     `json:"updatedParsed,omitempty" xml:"-"`
	Subtitle      string         `json:"subtitle,omitempty" xml:"subtitle,omitempty"`
	Links         []*Link        `json:"links,omitempty" xml:"link"`
	Generator     *Generator     `json:"generator,omitempty" xml:"generator"`
	Icon          string         `json:"icon,omitempty" xml:"icon,omitempty"`
	Logo          string         `json:"logo,omitempty" xml:"logo,omitempty"`
	Rights        string         `json:"rights,omitempty" xml:"rights,omitempty"`
	Contributors  []*Person      `json:"contributors,omitempty" xml:"contributor"`
	Authors       []*Person      `json:"authors,omitempty" xml:"author"`
	Categories    []*Category    `json:"categories,omitempty" xml:"category"`
	Extensions    ext.Extensions `json:"extensions,omitempty" xml:"-"`
}
