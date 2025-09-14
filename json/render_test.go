package json

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

// normalizeJSON normalizes JSON for comparison by parsing and re-encoding
func normalizeJSON(jsonStr string) string {
	var parsed interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return jsonStr // return original if parsing fails
	}
	normalized, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return jsonStr // return original if marshaling fails
	}
	return string(normalized)
}

func TestJSONFeedRender_Comprehensive(t *testing.T) {
	tests := []struct {
		name         string
		feed         *Feed
		expectedFile string
		checkStrings []string
	}{
		{
			name: "minimal JSON feed",
			feed: &Feed{
				Version: "https://jsonfeed.org/version/1.1",
				Title:   "Test Feed",
			},
			expectedFile: "minimal_feed.json",
		},
		{
			name: "comprehensive JSON feed with all fields",
			feed: &Feed{
				Version:     "https://jsonfeed.org/version/1.1",
				Title:       "Comprehensive JSON Feed",
				HomePageURL: "http://example.com",
				FeedURL:     "http://example.com/feed.json",
				Description: "A comprehensive test feed showcasing all JSON Feed 1.1 features",
				UserComment: "This is a test feed for demonstration purposes",
				NextURL:     "http://example.com/feed.json?page=2",
				Icon:        "http://example.com/icon.png",
				Favicon:     "http://example.com/favicon.ico",
				Author: &Author{
					Name:   "John Doe",
					URL:    "http://johndoe.example.com",
					Avatar: "http://johndoe.example.com/avatar.png",
				},
				Language: "en-US",
				Authors: []*Author{
					{
						Name:   "John Doe",
						URL:    "http://johndoe.example.com",
						Avatar: "http://johndoe.example.com/avatar.png",
					},
					{
						Name: "Jane Smith",
						URL:  "http://janesmith.example.com",
					},
				},
			},
			expectedFile: "comprehensive_feed.json",
		},
		{
			name: "JSON feed with comprehensive items",
			feed: &Feed{
				Version:     "https://jsonfeed.org/version/1.1",
				Title:       "Test Feed",
				HomePageURL: "http://example.com",
				Items: []*Item{
					{
						ID:            "http://example.com/item1",
						URL:           "http://example.com/item1",
						ExternalURL:   "http://external.example.com/article",
						Title:         "Comprehensive Item",
						ContentHTML:   "<p>Rich HTML content for the item</p>",
						ContentText:   "Plain text version of the content",
						Summary:       "A comprehensive JSON Feed item with all fields",
						Image:         "http://example.com/item1-image.jpg",
						BannerImage:   "http://example.com/item1-banner.jpg",
						DatePublished: "2024-01-01T10:00:00Z",
						DateModified:  "2024-01-01T11:00:00Z",
						Author: &Author{
							Name:   "Item Author",
							URL:    "http://itemauthor.example.com",
							Avatar: "http://itemauthor.example.com/avatar.png",
						},
						Tags: []string{"technology", "programming", "json"},
						Attachments: &[]Attachments{
							{
								URL:               "http://example.com/item1.mp3",
								MimeType:          "audio/mpeg",
								Title:             "Episode Audio",
								SizeInBytes:       1048576,
								DurationInSeconds: 3600,
							},
							{
								URL:         "http://example.com/item1.pdf",
								MimeType:    "application/pdf",
								Title:       "Episode Transcript",
								SizeInBytes: 2097152,
							},
						},
						Authors: []*Author{
							{
								Name: "Item Author",
								URL:  "http://itemauthor.example.com",
							},
							{
								Name: "Co-Author",
								URL:  "http://coauthor.example.com",
							},
						},
						Language: "en-US",
					},
					{
						ID:            "http://example.com/item2",
						URL:           "http://example.com/item2",
						Title:         "Simple Item",
						ContentText:   "Simple plain text content",
						DatePublished: "2024-01-02T10:00:00Z",
						Tags:          []string{"simple"},
					},
					{
						ID:           "http://example.com/item3",
						Title:        "Item with only HTML content",
						ContentHTML:  "<h1>HTML Only</h1><p>This item has only HTML content</p>",
						DateModified: "2024-01-03T10:00:00Z",
					},
				},
			},
			expectedFile: "feed_with_comprehensive_items.json",
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

			// Verify it's valid JSON
			var parsed interface{}
			if err := json.Unmarshal([]byte(output), &parsed); err != nil {
				t.Errorf("Rendered JSON is not valid: %v\nOutput:\n%s", err, output)
			}

			// Compare with expected output if file exists
			if tt.expectedFile != "" {
				expected := readTestFile(t, tt.expectedFile)
				if normalizeJSON(output) != normalizeJSON(expected) {
					t.Errorf("Output doesn't match expected.\nExpected:\n%s\n\nActual:\n%s",
						normalizeJSON(expected), normalizeJSON(output))
				}
			}

			// Check for specific strings if provided
			for _, checkStr := range tt.checkStrings {
				if !strings.Contains(output, checkStr) {
					t.Errorf("Expected output to contain %q, but it didn't.\nActual output:\n%s", checkStr, output)
				}
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
