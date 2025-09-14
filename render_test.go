package gofeed

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
)

func TestRenderJSON(t *testing.T) {
	// Create a universal feed
	feed := &Feed{
		Title:       "Test JSON Feed",
		Description: "A test feed for JSON rendering",
		Link:        "http://example.com",
		FeedLink:    "http://example.com/feed.json",
		Updated:     "2024-01-01T12:00:00Z",
		Authors: []*Person{
			{
				Name:  "Test Author",
				Email: "test@example.com",
			},
		},
		Items: []*Item{
			{
				Title:       "Test Item",
				Description: "A test item",
				Link:        "http://example.com/item1",
				Published:   "2024-01-01T10:00:00Z",
				GUID:        "item-1",
			},
		},
	}

	var buf bytes.Buffer
	err := feed.RenderJSON(&buf, nil)
	if err != nil {
		t.Fatalf("RenderJSON() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("RenderJSON() produced empty output")
	}

	// Verify it's valid JSON
	var parsed interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Errorf("RenderJSON() produced invalid JSON: %v", err)
	}

	// Verify it contains expected content
	if !strings.Contains(output, "Test JSON Feed") {
		t.Error("RenderJSON() output doesn't contain expected title")
	}
	if !strings.Contains(output, "jsonfeed.org/version/1.1") {
		t.Error("RenderJSON() output doesn't contain JSON Feed version")
	}
}

func TestRenderAtom(t *testing.T) {
	// Create a universal feed
	feed := &Feed{
		Title:       "Test Atom Feed",
		Description: "A test feed for Atom rendering",
		Link:        "http://example.com",
		FeedLink:    "http://example.com/feed.xml",
		Updated:     "2024-01-01T12:00:00Z",
		Authors: []*Person{
			{
				Name:  "Test Author",
				Email: "test@example.com",
			},
		},
		Items: []*Item{
			{
				Title:       "Test Item",
				Description: "A test item",
				Link:        "http://example.com/item1",
				Published:   "2024-01-01T10:00:00Z",
				GUID:        "item-1",
			},
		},
	}

	var buf bytes.Buffer
	err := feed.RenderAtom(&buf, nil)
	if err != nil {
		t.Fatalf("RenderAtom() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("RenderAtom() produced empty output")
	}

	// Verify it's valid XML
	var parsed interface{}
	if err := xml.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Errorf("RenderAtom() produced invalid XML: %v", err)
	}

	// Verify it contains expected content
	if !strings.Contains(output, "Test Atom Feed") {
		t.Error("RenderAtom() output doesn't contain expected title")
	}
	if !strings.Contains(output, "http://www.w3.org/2005/Atom") {
		t.Error("RenderAtom() output doesn't contain Atom namespace")
	}
}

func TestRenderRSS(t *testing.T) {
	// Create a universal feed
	feed := &Feed{
		Title:       "Test RSS Feed",
		Description: "A test feed for RSS rendering",
		Link:        "http://example.com",
		FeedLink:    "http://example.com/feed.xml",
		Updated:     "2024-01-01T12:00:00Z",
		Authors: []*Person{
			{
				Name:  "Test Author",
				Email: "test@example.com",
			},
		},
		Items: []*Item{
			{
				Title:       "Test Item",
				Description: "A test item",
				Link:        "http://example.com/item1",
				Published:   "2024-01-01T10:00:00Z",
				GUID:        "item-1",
			},
		},
	}

	var buf bytes.Buffer
	err := feed.RenderRSS(&buf, nil)
	if err != nil {
		t.Fatalf("RenderRSS() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("RenderRSS() produced empty output")
	}

	// Verify it's valid XML
	var parsed interface{}
	if err := xml.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Errorf("RenderRSS() produced invalid XML: %v", err)
	}

	// Verify it contains expected content
	if !strings.Contains(output, "Test RSS Feed") {
		t.Error("RenderRSS() output doesn't contain expected title")
	}
	if !strings.Contains(output, "<rss version=\"2.0\">") {
		t.Error("RenderRSS() output doesn't contain RSS version")
	}
}

func TestRender_WithCustomConverter(t *testing.T) {
	feed := &Feed{
		Title:       "Test Feed",
		Description: "A test feed",
		Link:        "http://example.com",
	}

	// Test with custom RSS converter
	customRSSConverter := &DefaultRSSConverter{}
	var buf bytes.Buffer
	err := feed.RenderRSS(&buf, customRSSConverter)
	if err != nil {
		t.Fatalf("RenderRSS() with custom converter error = %v", err)
	}

	if buf.Len() == 0 {
		t.Error("RenderRSS() with custom converter produced empty output")
	}
}

func TestRender_NilFeed(t *testing.T) {
	var feed *Feed

	t.Run("RenderJSON with nil feed", func(t *testing.T) {
		var buf bytes.Buffer
		err := feed.RenderJSON(&buf, nil)
		if err == nil {
			t.Error("Expected error for nil feed, got nil")
		}
	})

	t.Run("RenderAtom with nil feed", func(t *testing.T) {
		var buf bytes.Buffer
		err := feed.RenderAtom(&buf, nil)
		if err == nil {
			t.Error("Expected error for nil feed, got nil")
		}
	})

	t.Run("RenderRSS with nil feed", func(t *testing.T) {
		var buf bytes.Buffer
		err := feed.RenderRSS(&buf, nil)
		if err == nil {
			t.Error("Expected error for nil feed, got nil")
		}
	})
}

func TestRender_RoundTrip(t *testing.T) {
	// Test round-trip: parse -> render -> parse
	originalRSSData := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Test RSS Feed</title>
		<link>http://example.com</link>
		<description>A test RSS feed</description>
		<item>
			<title>Test Item</title>
			<link>http://example.com/item1</link>
			<description>A test item</description>
		</item>
	</channel>
</rss>`

	t.Run("RSS round-trip", func(t *testing.T) {
		// Parse RSS
		fp := NewParser()
		feed, err := fp.ParseString(originalRSSData)
		if err != nil {
			t.Fatalf("Failed to parse RSS: %v", err)
		}

		// Render back to RSS
		var buf bytes.Buffer
		err = feed.RenderRSS(&buf, nil)
		if err != nil {
			t.Fatalf("Failed to render RSS: %v", err)
		}

		// Parse the rendered RSS
		renderedFeed, err := fp.ParseString(buf.String())
		if err != nil {
			t.Fatalf("Failed to parse rendered RSS: %v", err)
		}

		// Verify key fields match
		if renderedFeed.Title != feed.Title {
			t.Errorf("Title mismatch: got %q, want %q", renderedFeed.Title, feed.Title)
		}
		if renderedFeed.Link != feed.Link {
			t.Errorf("Link mismatch: got %q, want %q", renderedFeed.Link, feed.Link)
		}
	})

	t.Run("Cross-format rendering", func(t *testing.T) {
		// Parse RSS and render as Atom
		fp := NewParser()
		feed, err := fp.ParseString(originalRSSData)
		if err != nil {
			t.Fatalf("Failed to parse RSS: %v", err)
		}

		// Render as Atom
		var atomBuf bytes.Buffer
		err = feed.RenderAtom(&atomBuf, nil)
		if err != nil {
			t.Fatalf("Failed to render as Atom: %v", err)
		}

		atomOutput := atomBuf.String()
		if !strings.Contains(atomOutput, "Test RSS Feed") {
			t.Error("Atom output doesn't contain original title")
		}
		if !strings.Contains(atomOutput, "http://www.w3.org/2005/Atom") {
			t.Error("Atom output doesn't contain Atom namespace")
		}

		// Render as JSON
		var jsonBuf bytes.Buffer
		err = feed.RenderJSON(&jsonBuf, nil)
		if err != nil {
			t.Fatalf("Failed to render as JSON: %v", err)
		}

		jsonOutput := jsonBuf.String()
		if !strings.Contains(jsonOutput, "Test RSS Feed") {
			t.Error("JSON output doesn't contain original title")
		}
		if !strings.Contains(jsonOutput, "jsonfeed.org/version/1.1") {
			t.Error("JSON output doesn't contain JSON Feed version")
		}
	})
}
