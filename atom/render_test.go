package atom

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/extensions"
)

// readTestFile reads a test file from the testdata directory
func readTestFile(t *testing.T, filename string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "render", filename))
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", filename, err)
	}
	return string(data)
}

// normalizeXML normalizes XML for comparison by removing extra whitespace
func normalizeXML(xml string) string {
	// Remove extra whitespace between XML elements
	re := regexp.MustCompile(`>\s+<`)
	normalized := re.ReplaceAllString(xml, "><")
	// Remove leading/trailing whitespace
	return strings.TrimSpace(normalized)
}

func TestAtomFeedRender_Comprehensive(t *testing.T) {
	tests := []struct {
		name         string
		feed         *Feed
		expectedFile string
		checkStrings []string
	}{
		{
			name: "minimal Atom feed",
			feed: &Feed{
				Title:   "Test Feed",
				ID:      "http://example.com/feed",
				Updated: "2024-01-01T12:00:00Z",
			},
			expectedFile: "minimal_feed.xml",
		},
		{
			name: "comprehensive Atom feed with all fields",
			feed: &Feed{
				Title:    "Comprehensive Atom Feed",
				ID:       "http://example.com/comprehensive-feed",
				Updated:  "2024-01-01T12:00:00Z",
				Subtitle: "A comprehensive test feed showcasing all Atom 1.0 features",
				Language: "en-US",
				Icon:     "http://example.com/icon.png",
				Logo:     "http://example.com/logo.png",
				Rights:   "Copyright 2024 Example Corp. All rights reserved.",
				Links: []*Link{
					{
						Href: "http://example.com",
						Rel:  "alternate",
						Type: "text/html",
					},
					{
						Href: "http://example.com/feed.xml",
						Rel:  "self",
						Type: "application/atom+xml",
					},
					{
						Href:     "http://example.com/related",
						Rel:      "related",
						Type:     "text/html",
						Title:    "Related Site",
						Hreflang: "en",
						Length:   "1024",
					},
				},
				Generator: &Generator{
					URI:     "http://atomgenerator.example.com",
					Version: "1.0",
					Value:   "Atom Generator",
				},
				Contributors: []*Person{
					{
						Name:  "Bob Contributor",
						Email: "bob@example.com",
						URI:   "http://bob.example.com",
					},
				},
				Authors: []*Person{
					{
						Name:  "John Doe",
						Email: "john@example.com",
						URI:   "http://johndoe.example.com",
					},
					{
						Name:  "Jane Smith",
						Email: "jane@example.com",
					},
				},
				Categories: []*Category{
					{
						Term:   "technology",
						Scheme: "http://example.com/categories",
						Label:  "Technology",
					},
					{
						Term:  "programming",
						Label: "Programming",
					},
				},
			},
			expectedFile: "comprehensive_feed.xml",
		},
		{
			name: "Atom feed with comprehensive entries",
			feed: &Feed{
				Title:   "Test Feed",
				ID:      "http://example.com/feed",
				Updated: "2024-01-01T12:00:00Z",
				Entries: []*Entry{
					{
						Title:     "Comprehensive Entry",
						ID:        "http://example.com/entry1",
						Updated:   "2024-01-01T11:00:00Z",
						Published: "2024-01-01T10:00:00Z",
						Summary:   "A comprehensive Atom entry with all fields",
						Rights:    "Entry-specific rights",
						Content: &Content{
							Type:  "html",
							Value: "<p>Rich HTML content for the entry</p>",
						},
						Authors: []*Person{
							{
								Name:  "Entry Author",
								Email: "entryauthor@example.com",
							},
						},
						Contributors: []*Person{
							{
								Name: "Entry Contributor",
								URI:  "http://contributor.example.com",
							},
						},
						Categories: []*Category{
							{
								Term:   "news",
								Scheme: "http://example.com/entry-categories",
								Label:  "News",
							},
						},
						Links: []*Link{
							{
								Href: "http://example.com/entry1",
								Rel:  "alternate",
								Type: "text/html",
							},
							{
								Href:   "http://example.com/entry1.mp3",
								Rel:    "enclosure",
								Type:   "audio/mpeg",
								Length: "1048576",
							},
						},
						Source: &Source{
							Title:   "Original Source Feed",
							ID:      "http://originalsource.com/feed",
							Updated: "2024-01-01T09:00:00Z",
						},
					},
					{
						Title:   "Simple Entry",
						ID:      "http://example.com/entry2",
						Updated: "2024-01-01T11:30:00Z",
						Summary: "A simple Atom entry",
						Content: &Content{
							Type:  "text",
							Value: "Plain text content",
						},
					},
				},
			},
			expectedFile: "feed_with_comprehensive_entries.xml",
		},
		{
			name: "Atom feed with content variations",
			feed: &Feed{
				Title:   "Content Test Feed",
				ID:      "http://example.com/content-feed",
				Updated: "2024-01-01T12:00:00Z",
				Entries: []*Entry{
					{
						Title:   "Entry with external content",
						ID:      "http://example.com/external-content",
						Updated: "2024-01-01T11:00:00Z",
						Content: &Content{
							Type: "html",
							Src:  "http://example.com/external-content.html",
						},
					},
					{
						Title:   "Entry with XHTML content",
						ID:      "http://example.com/xhtml-content",
						Updated: "2024-01-01T11:15:00Z",
						Content: &Content{
							Type:  "xhtml",
							Value: "<div xmlns=\"http://www.w3.org/1999/xhtml\"><p>XHTML content</p></div>",
						},
					},
				},
			},
			expectedFile: "feed_with_content_variations.xml",
		},
		{
			name: "Atom feed with extensions",
			feed: &Feed{
				Title:   "Feed with Extensions",
				ID:      "http://example.com/feed",
				Updated: "2024-01-01T12:00:00Z",
				Extensions: ext.Extensions{
					"custom": map[string][]ext.Extension{
						"field": {
							{Name: "field", Value: "value"},
						},
					},
				},
			},
			checkStrings: []string{
				"<title>Feed with Extensions</title>",
				// Extensions should not appear in rendered XML
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.feed.Render(&buf)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			output := buf.String()

			// Verify it's valid XML
			var parsed interface{}
			if err := xml.Unmarshal(buf.Bytes(), &parsed); err != nil {
				t.Errorf("Rendered XML is not valid: %v\nOutput:\n%s", err, output)
			}

			// Compare with expected output if file exists
			if tt.expectedFile != "" {
				expected := readTestFile(t, tt.expectedFile)
				if normalizeXML(output) != normalizeXML(expected) {
					t.Errorf("Output doesn't match expected.\nExpected:\n%s\n\nActual:\n%s",
						normalizeXML(expected), normalizeXML(output))
				}
			}

			// Check for specific strings if provided
			for _, checkStr := range tt.checkStrings {
				if !strings.Contains(output, checkStr) {
					t.Errorf("Expected output to contain %q, but it didn't.\nActual output:\n%s", checkStr, output)
				}
			}

			// Verify extensions don't appear in output
			if strings.Contains(output, "custom") {
				t.Error("Extensions should not appear in rendered XML")
			}
		})
	}
}

func TestRender_NilFeed(t *testing.T) {
	var feed *Feed
	var buf bytes.Buffer
	err := feed.Render(&buf)
	if err != nil {
		t.Errorf("Render() error = %v, want nil", err)
	}
	if buf.Len() != 0 {
		t.Errorf("Expected empty output for nil feed, got: %s", buf.String())
	}
}
