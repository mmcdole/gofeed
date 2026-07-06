package shared

import (
	"strings"
	"testing"

	xpp "github.com/mmcdole/goxpp/v2"

	ext "github.com/mmcdole/gofeed/extensions"
)

// parserOn returns a NewXMLParser positioned on the first StartTag with the
// given local name.
func parserOn(t *testing.T, doc, name string) *xpp.Parser {
	t.Helper()
	p := NewXMLParser(strings.NewReader(doc))
	for {
		tok, err := p.NextToken()
		if err != nil {
			t.Fatalf("positioning on <%s>: %v", name, err)
		}
		if tok == xpp.StartTag && p.Name() == name {
			return p
		}
		if tok == xpp.EndDocument {
			t.Fatalf("no <%s> element found", name)
		}
	}
}

func TestIsExtension(t *testing.T) {
	doc := `<rss xmlns:itunes="http://www.itunes.com/DTDs/PodCast-1.0.dtd" xmlns:content="http://purl.org/rss/1.0/modules/content/">
		<channel>
			<title>t</title>
			<itunes:author>a</itunes:author>
			<content:encoded>c</content:encoded>
		</channel>
	</rss>`

	cases := []struct {
		element string
		want    bool
	}{
		{"title", false},   // no prefix
		{"author", true},   // itunes prefix
		{"encoded", false}, // content prefix is exempt
	}
	for _, c := range cases {
		p := parserOn(t, doc, c.element)
		if got := IsExtension(p); got != c.want {
			t.Errorf("IsExtension(<%s>) = %v, want %v", c.element, got, c.want)
		}
	}
}

func TestParseExtension(t *testing.T) {
	doc := `<rss xmlns:itunes="http://www.itunes.com/DTDs/PodCast-1.0.dtd">
		<channel>
			<itunes:owner>
				<itunes:name>Alice</itunes:name>
				<itunes:email>a@example.org</itunes:email>
			</itunes:owner>
			<itunes:category text="Tech"/>
			<itunes:category text="News"/>
		</channel>
	</rss>`

	p := NewXMLParser(strings.NewReader(doc))
	extensions := ext.Extensions{}
	for {
		tok, err := p.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		if tok == xpp.EndDocument {
			break
		}
		if tok == xpp.StartTag && IsExtension(p) {
			extensions, err = ParseExtension(extensions, p)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	owner := extensions["itunes"]["owner"]
	if len(owner) != 1 {
		t.Fatalf("owner entries = %d, want 1", len(owner))
	}
	if got := owner[0].Children["name"][0].Value; got != "Alice" {
		t.Errorf("owner name = %q, want Alice", got)
	}
	if got := owner[0].Children["email"][0].Value; got != "a@example.org" {
		t.Errorf("owner email = %q, want a@example.org", got)
	}

	cats := extensions["itunes"]["category"]
	if len(cats) != 2 {
		t.Fatalf("category entries = %d, want 2 (same-name elements must append)", len(cats))
	}
	if cats[0].Attrs["text"] != "Tech" || cats[1].Attrs["text"] != "News" {
		t.Errorf("category attrs = %v, %v", cats[0].Attrs, cats[1].Attrs)
	}
	if cats[0].Name != "category" {
		t.Errorf("extension Name = %q, want category", cats[0].Name)
	}
}

func TestPrefixForNamespace(t *testing.T) {
	doc := `<rss xmlns:pod="http://www.itunes.com/DTDs/PodCast-1.0.dtd" xmlns:custom="http://example.org/ns">
		<channel><pod:author>a</pod:author></channel>
	</rss>`
	p := parserOn(t, doc, "author")

	// A canonical namespace URI maps to its canonical prefix, not the
	// prefix the feed used.
	if got := PrefixForNamespace("http://www.itunes.com/DTDs/PodCast-1.0.dtd", p); got != "itunes" {
		t.Errorf("canonical = %q, want itunes", got)
	}
	// Whitespace-padded URIs map the same way.
	if got := PrefixForNamespace(" http://www.itunes.com/DTDs/PodCast-1.0.dtd ", p); got != "itunes" {
		t.Errorf("padded canonical = %q, want itunes", got)
	}
	// A feed-declared namespace maps to the feed's prefix.
	if got := PrefixForNamespace("http://example.org/ns", p); got != "custom" {
		t.Errorf("feed-declared = %q, want custom", got)
	}
	// An undeclared prefix comes through encoding/xml as the raw prefix in
	// Space; it maps to itself.
	if got := PrefixForNamespace("media", p); got != "media" {
		t.Errorf("undeclared = %q, want media", got)
	}
}

// NewXMLParser must be non-strict (real feeds carry bare ampersands) and must
// convert declared non-UTF-8 encodings.
func TestNewXMLParserLeniencyAndCharset(t *testing.T) {
	p := parserOn(t, `<root>Fish & Chips</root>`, "root")
	text, err := p.NextText()
	if err != nil {
		t.Fatalf("bare ampersand should tokenize in non-strict mode: %v", err)
	}
	if !strings.Contains(text, "&") {
		t.Errorf("text = %q, want the ampersand preserved", text)
	}

	// 0xE9 is é in ISO-8859-1 and invalid UTF-8; parsing it proves the
	// charset reader is wired up.
	latin1 := `<?xml version="1.0" encoding="ISO-8859-1"?><root>caf` + "\xe9" + `</root>`
	p2 := parserOn(t, latin1, "root")
	text, err = p2.NextText()
	if err != nil {
		t.Fatal(err)
	}
	if text != "café" {
		t.Errorf("text = %q, want café", text)
	}
}

func TestParseCustom(t *testing.T) {
	doc := `<rss version="2.0">
		<channel>
			<item>
				<event><venue city="Austin">Hall</venue></event>
				<simple>plain</simple>
				<simple>again</simple>
			</item>
		</channel>
	</rss>`

	p := NewXMLParser(strings.NewReader(doc))
	extensions := ext.Extensions{}
	var parsed []ext.Extension
	for {
		tok, err := p.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		if tok == xpp.EndDocument {
			break
		}
		if tok == xpp.StartTag {
			switch p.Name() {
			case "event", "simple":
				var e ext.Extension
				var err error
				extensions, e, err = ParseCustom(extensions, p)
				if err != nil {
					t.Fatal(err)
				}
				parsed = append(parsed, e)
			}
		}
	}

	if len(parsed) != 3 {
		t.Fatalf("parsed %d elements, want 3", len(parsed))
	}

	// Nesting and attributes survive in the tree.
	events := extensions[CustomPrefix]["event"]
	if len(events) != 1 {
		t.Fatalf("event entries = %d, want 1", len(events))
	}
	venue := events[0].Children["venue"][0]
	if venue.Value != "Hall" || venue.Attrs["city"] != "Austin" {
		t.Fatalf("venue = %+v", venue)
	}

	// Repetition appends rather than overwrites.
	if n := len(extensions[CustomPrefix]["simple"]); n != 2 {
		t.Fatalf("simple entries = %d, want 2", n)
	}

	// The returned element mirrors what was filed.
	if parsed[0].Name != "event" || len(parsed[0].Children) != 1 {
		t.Fatalf("returned element = %+v", parsed[0])
	}
	if parsed[1].Value != "plain" || len(parsed[1].Children) != 0 {
		t.Fatalf("returned simple = %+v", parsed[1])
	}
}
