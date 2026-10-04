package gofeed

import (
	"strings"
	"testing"
)

// Regression test for #369: an RSS item whose only link is an embedded
// <atom:link rel="alternate"> must populate Item.Link and Item.Links, the same
// way the channel's atom:link populates Feed.Link.
func TestItemAtomLinkFallback(t *testing.T) {
	feed := `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel><title>t</title>
<atom:link rel="alternate" href="https://example.com/"/>
<item><title>i</title>
<atom:link rel="alternate" href="https://example.com/1"/>
</item></channel></rss>`
	f, err := NewParser().Parse(strings.NewReader(feed))
	if err != nil {
		t.Fatal(err)
	}
	it := f.Items[0]
	if it.Link != "https://example.com/1" {
		t.Fatalf("Item.Link = %q, want https://example.com/1", it.Link)
	}
	if len(it.Links) != 1 || it.Links[0] != "https://example.com/1" {
		t.Fatalf("Item.Links = %v, want [https://example.com/1]", it.Links)
	}
}
