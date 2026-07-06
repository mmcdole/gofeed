package shared

import (
	"net/url"
	"testing"
)

// XmlBaseResolveUrl must not mutate the base it is given: the base is the live
// value held on the parser's xml:base stack, shared by every sibling in scope.
func TestXmlBaseResolveUrlDoesNotMutateBase(t *testing.T) {
	base, _ := url.Parse("http://example.com/a/b")
	before := base.String()

	if _, err := XmlBaseResolveUrl(base, "x"); err != nil {
		t.Fatal(err)
	}
	if after := base.String(); after != before {
		t.Errorf("base mutated: %q -> %q", before, after)
	}
}

// A directory-style resolution should still work (the path is treated as a
// directory), just without mutating the input.
func TestXmlBaseResolveUrlResolves(t *testing.T) {
	base, _ := url.Parse("http://example.com/a/b")
	got, err := XmlBaseResolveUrl(base, "x")
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != "http://example.com/a/b/x" {
		t.Errorf("resolved = %q, want %q", got.String(), "http://example.com/a/b/x")
	}
}

// Every URI attribute on an element must be resolved, not just the first
// one encountered.
func TestResolveHTMLResolvesAllURIAttributes(t *testing.T) {
	base, _ := url.Parse("http://example.com/dir/")

	got, err := ResolveHTML(base, `<video poster="p.jpg" src="v.mp4"></video>`)
	if err != nil {
		t.Fatal(err)
	}
	want := `<video poster="http://example.com/dir/p.jpg" src="http://example.com/dir/v.mp4"></video>`
	if got != want {
		t.Errorf("resolved = %q, want %q", got, want)
	}
}

func TestResolveHTMLResolvesAcrossElements(t *testing.T) {
	base, _ := url.Parse("http://example.com/dir/")

	got, err := ResolveHTML(base, `<a href="page.html"><img src="i.png"/></a>`)
	if err != nil {
		t.Fatal(err)
	}
	want := `<a href="http://example.com/dir/page.html"><img src="http://example.com/dir/i.png"/></a>`
	if got != want {
		t.Errorf("resolved = %q, want %q", got, want)
	}
}

func TestResolveHTMLLeavesAbsoluteURLs(t *testing.T) {
	base, _ := url.Parse("http://example.com/dir/")

	got, err := ResolveHTML(base, `<a href="https://other.org/page.html">x</a>`)
	if err != nil {
		t.Fatal(err)
	}
	want := `<a href="https://other.org/page.html">x</a>`
	if got != want {
		t.Errorf("resolved = %q, want %q", got, want)
	}
}

// A malformed URI in one attribute is kept as-is; it must not error out or
// stop other attributes from resolving.
func TestResolveHTMLKeepsMalformedURIs(t *testing.T) {
	base, _ := url.Parse("http://example.com/dir/")

	got, err := ResolveHTML(base, `<video poster="http://%zz" src="v.mp4"></video>`)
	if err != nil {
		t.Fatal(err)
	}
	want := `<video poster="http://%zz" src="http://example.com/dir/v.mp4"></video>`
	if got != want {
		t.Errorf("resolved = %q, want %q", got, want)
	}
}

func TestResolveHTMLNilBase(t *testing.T) {
	got, err := ResolveHTML(nil, `<a href="page.html">x</a>`)
	if err != nil {
		t.Fatal(err)
	}
	if got != `<a href="page.html">x</a>` {
		t.Errorf("resolved = %q, want input unchanged", got)
	}
}

// Plain text must survive the html.Render round trip: the document wrapper
// the renderer adds has to be stripped back off.
func TestResolveHTMLPlainTextRoundTrip(t *testing.T) {
	base, _ := url.Parse("http://example.com/dir/")

	got, err := ResolveHTML(base, "hello world")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello world" {
		t.Errorf("resolved = %q, want %q", got, "hello world")
	}
}
