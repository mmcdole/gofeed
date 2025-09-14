package rss

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	ext "github.com/mmcdole/gofeed/extensions"
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

func TestRSSFeedRender_Comprehensive(t *testing.T) {
	tests := []struct {
		name         string
		feed         *Feed
		expectedFile string
		checkStrings []string // Additional specific strings to verify
	}{
		{
			name: "minimal RSS feed",
			feed: &Feed{
				Title:       "Test Feed",
				Link:        "http://example.com",
				Description: "Test Description",
				Version:     "2.0",
			},
			expectedFile: "minimal_feed.xml",
		},
		{
			name: "comprehensive RSS feed with all standard fields",
			feed: &Feed{
				Title:          "Comprehensive Test Feed",
				Link:           "http://example.com",
				Description:    "A comprehensive test feed with all RSS 2.0 fields",
				Language:       "en-us",
				Copyright:      "Copyright 2024 Example Corp",
				ManagingEditor: "editor@example.com (John Editor)",
				WebMaster:      "webmaster@example.com (Jane Webmaster)",
				PubDate:        "Mon, 01 Jan 2024 12:00:00 GMT",
				LastBuildDate:  "Mon, 01 Jan 2024 13:00:00 GMT",
				Generator:      "Test RSS Generator v1.0",
				Docs:           "http://www.rssboard.org/rss-specification",
				TTL:            "60",
				Rating:         "PICS-1.1 'http://www.rsac.org/ratingsv01.html' l gen true comment 'RSACi North America Server'",
				SkipHours:      []string{"0", "1", "2", "3", "4", "5"},
				SkipDays:       []string{"Saturday", "Sunday"},
				Version:        "2.0",
				Categories: []*Category{
					{Domain: "http://example.com/categories", Value: "Technology"},
					{Value: "Software"},
				},
				Image: &Image{
					URL:         "http://example.com/logo.png",
					Title:       "Test Feed Logo",
					Link:        "http://example.com",
					Width:       "144",
					Height:      "144",
					Description: "The official logo of Test Feed",
				},
				Cloud: &Cloud{
					Domain:            "rpc.example.com",
					Port:              "80",
					Path:              "/RPC2",
					RegisterProcedure: "myCloud.rssPleaseNotify",
					Protocol:          "xml-rpc",
				},
				TextInput: &TextInput{
					Title:       "Search",
					Description: "Search our feed",
					Name:        "query",
					Link:        "http://example.com/search",
				},
			},
			expectedFile: "comprehensive_feed.xml",
		},
		{
			name: "RSS feed with comprehensive items",
			feed: &Feed{
				Title:       "Test Feed",
				Link:        "http://example.com",
				Description: "Test Description",
				Version:     "2.0",
				Items: []*Item{
					{
						Title:       "Comprehensive Item",
						Link:        "http://example.com/item1",
						Description: "A comprehensive RSS item with all fields",
						Content:     "<p>Rich HTML content for the item</p>",
						Author:      "author@example.com (Item Author)",
						Comments:    "http://example.com/item1/comments",
						PubDate:     "Mon, 01 Jan 2024 10:00:00 GMT",
						Categories: []*Category{
							{Domain: "http://example.com/itemcats", Value: "News"},
							{Value: "Breaking"},
						},
						Enclosure: &Enclosure{
							URL:    "http://example.com/item1.mp3",
							Type:   "audio/mpeg",
							Length: "1048576",
						},
						Enclosures: []*Enclosure{
							{
								URL:    "http://example.com/item1.mp3",
								Type:   "audio/mpeg",
								Length: "1048576",
							},
							{
								URL:    "http://example.com/item1.pdf",
								Type:   "application/pdf",
								Length: "2097152",
							},
						},
						GUID: &GUID{
							Value:       "http://example.com/item1/guid",
							IsPermalink: "true",
						},
						Source: &Source{
							Title: "Original Source",
							URL:   "http://originalsource.com/rss",
						},
					},
					{
						Title:       "Simple Item",
						Link:        "http://example.com/item2",
						Description: "A simple RSS item",
						PubDate:     "Tue, 02 Jan 2024 10:00:00 GMT",
						GUID: &GUID{
							Value:       "unique-item-2",
							IsPermalink: "false",
						},
					},
				},
			},
			expectedFile: "feed_with_comprehensive_items.xml",
		},
		{
			name: "RSS feed with iTunes extensions",
			feed: &Feed{
				Title:       "iTunes Podcast Feed",
				Link:        "http://example.com",
				Description: "A podcast feed with iTunes extensions",
				Version:     "2.0",
				ITunesExt: &ext.ITunesFeedExtension{
					Author:     "Podcast Author",
					Block:      "yes",
					Categories: []*ext.ITunesCategory{{Text: "Technology"}},
					Image:      "http://example.com/podcast-image.jpg",
					Explicit:   "no",
					Complete:   "no",
					NewFeedURL: "",
					Owner:      &ext.ITunesOwner{Name: "Owner Name", Email: "owner@example.com"},
					Subtitle:   "A great podcast about technology",
					Summary:    "This is a comprehensive summary of our technology podcast.",
					Keywords:   "technology,programming,software",
				},
				Items: []*Item{
					{
						Title:       "Podcast Episode 1",
						Link:        "http://example.com/episode1",
						Description: "First episode of our podcast",
						PubDate:     "Mon, 01 Jan 2024 10:00:00 GMT",
						GUID:        &GUID{Value: "episode-1"},
						Enclosure: &Enclosure{
							URL:    "http://example.com/episode1.mp3",
							Type:   "audio/mpeg",
							Length: "10485760",
						},
						ITunesExt: &ext.ITunesItemExtension{
							Author:            "Episode Author",
							Block:             "no",
							Image:             "http://example.com/episode1.jpg",
							Duration:          "45:30",
							Explicit:          "no",
							IsClosedCaptioned: "no",
							Order:             "1",
							Subtitle:          "Episode 1 Subtitle",
							Summary:           "Episode 1 Summary",
							Keywords:          "intro,technology",
						},
					},
				},
			},
			expectedFile: "feed_with_itunes_extensions.xml",
			// Since iTunes extensions are excluded from rendering, we verify the basic structure
			checkStrings: []string{
				"<title>iTunes Podcast Feed</title>",
				"<title>Podcast Episode 1</title>",
			},
		},
		{
			name: "RSS feed with Dublin Core extensions",
			feed: &Feed{
				Title:       "Dublin Core Extended Feed",
				Link:        "http://example.com",
				Description: "A feed with Dublin Core metadata",
				Version:     "2.0",
				DublinCoreExt: &ext.DublinCoreExtension{
					Creator:     []string{"John Creator", "Jane Creator"},
					Subject:     []string{"Technology", "Science"},
					Description: []string{"Extended description"},
					Publisher:   []string{"Example Publisher"},
					Contributor: []string{"Contributor 1", "Contributor 2"},
					Date:        []string{"2024-01-01"},
					Type:        []string{"Text"},
					Format:      []string{"application/rss+xml"},
					Identifier:  []string{"http://example.com/feed"},
					Source:      []string{"Original Source"},
					Language:    []string{"en-US"},
					Relation:    []string{"Related Resource"},
					Coverage:    []string{"Global"},
					Rights:      []string{"All rights reserved"},
				},
				Items: []*Item{
					{
						Title:       "Item with DC metadata",
						Link:        "http://example.com/item1",
						Description: "An item with Dublin Core metadata",
						DublinCoreExt: &ext.DublinCoreExtension{
							Creator: []string{"Item Creator"},
							Subject: []string{"Item Subject"},
							Date:    []string{"2024-01-01T10:00:00Z"},
						},
					},
				},
			},
			// Dublin Core extensions are excluded from rendering
			checkStrings: []string{
				"<title>Dublin Core Extended Feed</title>",
				"<title>Item with DC metadata</title>",
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
