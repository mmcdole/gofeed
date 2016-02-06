package atom

import (
	"time"
)

type Feed struct {
	Title         string
	ID            string
	Updated       string
	UpdatedParsed *time.Time
	Subtitle      string
	Link          string
	Generator     Generator
	Icon          string
	Logo          string
	Rights        string
	Contributors  []Person
	Authors       []Person
	Categories    []Category
	Source        string
	Version       string
}

type Entry struct {
	Title           string
	ID              string
	Link            Link
	Published       string
	PublishedParsed *time.Time
	Updated         string
	UpdatedParsed   *time.Time
	Content         string
}

type Category struct {
	Term   string
	Scheme string
	Label  string
}

type Person struct {
	Name  string
	Email string
	URI   string
}

type Link struct {
	Rel      string
	Type     string
	Href     string
	Hreflang string
	Title    string
	Length   string
}

type Content struct {
	Src  string
	Type string
}

type Generator struct {
	Value   string
	URI     string
	Version string
}
