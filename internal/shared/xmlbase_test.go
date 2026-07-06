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
