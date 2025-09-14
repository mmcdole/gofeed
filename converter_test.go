package gofeed_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/atom"
	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/mmcdole/gofeed/json"
	"github.com/mmcdole/gofeed/rss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRSSConverter_Render(t *testing.T) {
	files, _ := filepath.Glob("testdata/translator/rss/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing RSS converter round-trip for %s... ", name)

		// Parse original RSS feed
		ff := fmt.Sprintf("testdata/translator/rss/%s.xml", name)
		file, err := os.Open(ff)
		require.NoError(t, err)
		defer file.Close()

		rssParser := &rss.Parser{}
		originalRSSFeed, err := rssParser.Parse(file)
		require.NoError(t, err)

		// Translate to universal format
		translator := &gofeed.DefaultRSSTranslator{}
		universalFeed, err := translator.Translate(originalRSSFeed)
		require.NoError(t, err)

		// Render back to RSS
		converter := &gofeed.DefaultRSSConverter{}
		renderedRSSFeed, err := converter.Convert(universalFeed)
		require.NoError(t, err)

		// Compare key fields
		assert.Equal(t, originalRSSFeed.Title, renderedRSSFeed.Title, "Title should match")
		assert.Equal(t, originalRSSFeed.Description, renderedRSSFeed.Description, "Description should match")
		assert.Equal(t, originalRSSFeed.Link, renderedRSSFeed.Link, "Link should match")
		assert.Equal(t, originalRSSFeed.Language, renderedRSSFeed.Language, "Language should match")
		assert.Equal(t, originalRSSFeed.Copyright, renderedRSSFeed.Copyright, "Copyright should match")
		assert.Equal(t, originalRSSFeed.Generator, renderedRSSFeed.Generator, "Generator should match")
		assert.Equal(t, originalRSSFeed.PubDate, renderedRSSFeed.PubDate, "PubDate should match")
		assert.Equal(t, originalRSSFeed.LastBuildDate, renderedRSSFeed.LastBuildDate, "LastBuildDate should match")

		// Compare number of items
		assert.Equal(t, len(originalRSSFeed.Items), len(renderedRSSFeed.Items), "Number of items should match")

		// Compare first item if it exists
		if len(originalRSSFeed.Items) > 0 && len(renderedRSSFeed.Items) > 0 {
			origItem := originalRSSFeed.Items[0]
			rendItem := renderedRSSFeed.Items[0]

			assert.Equal(t, origItem.Title, rendItem.Title, "Item title should match")
			assert.Equal(t, origItem.Description, rendItem.Description, "Item description should match")
			assert.Equal(t, origItem.Link, rendItem.Link, "Item link should match")
			assert.Equal(t, origItem.PubDate, rendItem.PubDate, "Item pubDate should match")

			// Author formatting might change slightly due to round-trip conversion
			// Just check that both are present if original had author
			if origItem.Author != "" {
				assert.NotEmpty(t, rendItem.Author, "Item should have author if original had author")
			}
		}

		fmt.Printf("OK\n")
	}
}

func TestRSSConverter_Render_NilFeed(t *testing.T) {
	converter := &gofeed.DefaultRSSConverter{}
	result, err := converter.Convert(nil)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "feed cannot be nil")
}

func TestAtomConverter_Render(t *testing.T) {
	files, _ := filepath.Glob("testdata/translator/atom/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing Atom converter round-trip for %s... ", name)

		// Parse original Atom feed
		ff := fmt.Sprintf("testdata/translator/atom/%s.xml", name)
		file, err := os.Open(ff)
		require.NoError(t, err)
		defer file.Close()

		atomParser := &atom.Parser{}
		originalAtomFeed, err := atomParser.Parse(file)
		require.NoError(t, err)

		// Translate to universal format
		translator := &gofeed.DefaultAtomTranslator{}
		universalFeed, err := translator.Translate(originalAtomFeed)
		require.NoError(t, err)

		// Render back to Atom
		converter := &gofeed.DefaultAtomConverter{}
		renderedAtomFeed, err := converter.Convert(universalFeed)
		require.NoError(t, err)

		// Compare key fields
		assert.Equal(t, originalAtomFeed.Title, renderedAtomFeed.Title, "Title should match")
		assert.Equal(t, originalAtomFeed.Subtitle, renderedAtomFeed.Subtitle, "Subtitle should match")
		assert.Equal(t, originalAtomFeed.Language, renderedAtomFeed.Language, "Language should match")
		assert.Equal(t, originalAtomFeed.Rights, renderedAtomFeed.Rights, "Rights should match")
		assert.Equal(t, originalAtomFeed.Updated, renderedAtomFeed.Updated, "Updated should match")

		// Logo comparison - allow for Icon->Logo conversion during round-trip
		if originalAtomFeed.Logo != "" {
			assert.Equal(t, originalAtomFeed.Logo, renderedAtomFeed.Logo, "Logo should match")
		} else if originalAtomFeed.Icon != "" {
			assert.Equal(t, originalAtomFeed.Icon, renderedAtomFeed.Logo, "Icon should become Logo")
		}

		// Compare number of entries
		assert.Equal(t, len(originalAtomFeed.Entries), len(renderedAtomFeed.Entries), "Number of entries should match")

		// Compare first entry if it exists
		if len(originalAtomFeed.Entries) > 0 && len(renderedAtomFeed.Entries) > 0 {
			origEntry := originalAtomFeed.Entries[0]
			rendEntry := renderedAtomFeed.Entries[0]

			assert.Equal(t, origEntry.Title, rendEntry.Title, "Entry title should match")
			assert.Equal(t, origEntry.Summary, rendEntry.Summary, "Entry summary should match")
			assert.Equal(t, origEntry.Updated, rendEntry.Updated, "Entry updated should match")

			// Published date handling can differ due to fallbacks in translation
			// Just verify both are present if either is present
			if origEntry.Published != "" || rendEntry.Published != "" {
				// Allow for some differences in published date handling
				if origEntry.Published != "" {
					assert.NotEmpty(t, rendEntry.Published, "Entry should have published if original had published")
				}
			}

			// Compare ID (GUID handling)
			if origEntry.ID != "" {
				assert.Equal(t, origEntry.ID, rendEntry.ID, "Entry ID should match")
			}
		}

		fmt.Printf("OK\n")
	}
}

func TestAtomConverter_Render_NilFeed(t *testing.T) {
	converter := &gofeed.DefaultAtomConverter{}
	result, err := converter.Convert(nil)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "feed cannot be nil")
}

func TestJSONConverter_Render(t *testing.T) {
	files, _ := filepath.Glob("testdata/translator/json/*.json")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		if strings.HasSuffix(name, "expected") {
			continue
		}

		fmt.Printf("Testing JSON converter round-trip for %s... ", name)

		// Parse original JSON feed
		ff := fmt.Sprintf("testdata/translator/json/%s.json", name)
		file, err := os.Open(ff)
		require.NoError(t, err)
		defer file.Close()

		jsonParser := json.Parser{}
		originalJSONFeed, err := jsonParser.Parse(file)
		require.NoError(t, err)

		// Translate to universal format
		translator := &gofeed.DefaultJSONTranslator{}
		universalFeed, err := translator.Translate(originalJSONFeed)
		require.NoError(t, err)

		// Render back to JSON
		converter := &gofeed.DefaultJSONConverter{}
		renderedJSONFeed, err := converter.Convert(universalFeed)
		require.NoError(t, err)

		// Compare key fields
		assert.Equal(t, originalJSONFeed.Title, renderedJSONFeed.Title, "Title should match")
		assert.Equal(t, originalJSONFeed.Description, renderedJSONFeed.Description, "Description should match")
		assert.Equal(t, originalJSONFeed.HomePageURL, renderedJSONFeed.HomePageURL, "HomePageURL should match")
		assert.Equal(t, originalJSONFeed.FeedURL, renderedJSONFeed.FeedURL, "FeedURL should match")
		assert.Equal(t, originalJSONFeed.Language, renderedJSONFeed.Language, "Language should match")
		assert.Equal(t, originalJSONFeed.Icon, renderedJSONFeed.Icon, "Icon should match")

		// Compare number of items
		assert.Equal(t, len(originalJSONFeed.Items), len(renderedJSONFeed.Items), "Number of items should match")

		// Compare first item if it exists
		if len(originalJSONFeed.Items) > 0 && len(renderedJSONFeed.Items) > 0 {
			origItem := originalJSONFeed.Items[0]
			rendItem := renderedJSONFeed.Items[0]

			assert.Equal(t, origItem.Title, rendItem.Title, "Item title should match")
			assert.Equal(t, origItem.Summary, rendItem.Summary, "Item summary should match")
			assert.Equal(t, origItem.URL, rendItem.URL, "Item URL should match")
			assert.Equal(t, origItem.DatePublished, rendItem.DatePublished, "Item DatePublished should match")
			assert.Equal(t, origItem.DateModified, rendItem.DateModified, "Item DateModified should match")

			// Compare ID handling
			if origItem.ID != "" {
				assert.Equal(t, origItem.ID, rendItem.ID, "Item ID should match")
			}
		}

		fmt.Printf("OK\n")
	}
}

func TestJSONConverter_Render_NilFeed(t *testing.T) {
	converter := &gofeed.DefaultJSONConverter{}
	result, err := converter.Convert(nil)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "feed cannot be nil")
}

// Test specific field mapping scenarios
func TestRSSConverter_AuthorFormatting(t *testing.T) {
	converter := &gofeed.DefaultRSSConverter{}

	tests := []struct {
		name     string
		author   *gofeed.Person
		expected string
	}{
		{
			name:     "Email and name",
			author:   &gofeed.Person{Name: "John Doe", Email: "john@example.com"},
			expected: "john@example.com (John Doe)",
		},
		{
			name:     "Email only",
			author:   &gofeed.Person{Email: "john@example.com"},
			expected: "john@example.com",
		},
		{
			name:     "Name only",
			author:   &gofeed.Person{Name: "John Doe"},
			expected: "John Doe",
		},
		{
			name:     "Empty",
			author:   &gofeed.Person{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.formatPersonForRSS(tt.author)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtomConverter_LinkGeneration(t *testing.T) {
	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Link:        "https://example.com",
		FeedLink:    "https://example.com/feed.xml",
		Items:       []*gofeed.Item{},
	}

	converter := &gofeed.DefaultAtomConverter{}
	renderedAtomFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	// Should have both alternate and self links
	require.Len(t, renderedAtomFeed.Links, 2)

	var alternateLink, selfLink *atom.Link
	for _, link := range renderedAtomFeed.Links {
		if link.Rel == "alternate" {
			alternateLink = link
		} else if link.Rel == "self" {
			selfLink = link
		}
	}

	require.NotNil(t, alternateLink, "Should have alternate link")
	require.NotNil(t, selfLink, "Should have self link")

	assert.Equal(t, "https://example.com", alternateLink.Href)
	assert.Equal(t, "text/html", alternateLink.Type)
	assert.Equal(t, "https://example.com/feed.xml", selfLink.Href)
	assert.Equal(t, "application/atom+xml", selfLink.Type)
}

func TestJSONConverter_AuthorFormatting(t *testing.T) {
	converter := &gofeed.DefaultJSONConverter{}

	tests := []struct {
		name     string
		author   *gofeed.Person
		expected string
	}{
		{
			name:     "Email and name",
			author:   &gofeed.Person{Name: "John Doe", Email: "john@example.com"},
			expected: "John Doe <john@example.com>",
		},
		{
			name:     "Email only",
			author:   &gofeed.Person{Email: "john@example.com"},
			expected: "john@example.com",
		},
		{
			name:     "Name only",
			author:   &gofeed.Person{Name: "John Doe"},
			expected: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed := &gofeed.Feed{
				Title:       "Test Feed",
				Description: "Test Description",
				Author:      tt.author,
				Items:       []*gofeed.Item{},
			}

			jsonFeed, err := converter.Convert(feed)
			require.NoError(t, err)
			require.NotNil(t, jsonFeed.Author)
			assert.Equal(t, tt.expected, jsonFeed.Author.Name)
		})
	}
}

func TestJSONConverter_EnclosureConversion(t *testing.T) {
	converter := &gofeed.DefaultJSONConverter{}

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Items: []*gofeed.Item{
			{
				Title: "Test Item",
				Enclosures: []*gofeed.Enclosure{
					{
						URL:    "https://example.com/audio.mp3",
						Type:   "audio/mpeg",
						Length: "12345",
					},
				},
			},
		},
	}

	jsonFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	require.Len(t, jsonFeed.Items, 1)
	item := jsonFeed.Items[0]
	require.NotNil(t, item.Attachments)
	require.Len(t, *item.Attachments, 1)

	attachment := (*item.Attachments)[0]
	assert.Equal(t, "https://example.com/audio.mp3", attachment.URL)
	assert.Equal(t, "audio/mpeg", attachment.MimeType)
	assert.Equal(t, int64(12345), attachment.DurationInSeconds)
}

// Test known round-trip limitations
func TestRoundTripLimitations(t *testing.T) {
	t.Run("RSS Updated field preservation via DublinCore", func(t *testing.T) {
		// RSS doesn't have native item-level updated field
		// It should be stored in DublinCore extension Date field
		universalFeed := &gofeed.Feed{
			Title: "Test Feed",
			Items: []*gofeed.Item{
				{
					Title:   "Test Item",
					Updated: "2023-01-01T12:00:00Z",
				},
			},
		}

		converter := &gofeed.DefaultRSSConverter{}
		rssFeed, err := converter.Convert(universalFeed)
		require.NoError(t, err)

		// RSS item doesn't have Updated field in the struct, but should be in DublinCore
		assert.Equal(t, "Test Item", rssFeed.Items[0].Title)
		require.NotNil(t, rssFeed.Items[0].DublinCoreExt)
		require.Len(t, rssFeed.Items[0].DublinCoreExt.Date, 1)
		assert.Equal(t, "2023-01-01T12:00:00Z", rssFeed.Items[0].DublinCoreExt.Date[0])
	})

	t.Run("JSON feed-level date handling", func(t *testing.T) {
		// JSON translator derives feed dates from first item
		// But converter doesn't use feed-level dates
		universalFeed := &gofeed.Feed{
			Title:     "Test Feed",
			Updated:   "2023-01-01T12:00:00Z",
			Published: "2023-01-01T10:00:00Z",
			Items:     []*gofeed.Item{},
		}

		converter := &gofeed.DefaultJSONConverter{}
		jsonFeed, err := converter.Convert(universalFeed)
		require.NoError(t, err)

		// JSON Feed spec doesn't have feed-level dates
		// This is expected behavior
		assert.Equal(t, "Test Feed", jsonFeed.Title)
	})

	t.Run("Atom generator complexity", func(t *testing.T) {
		// Atom translator builds complex generator strings
		// Renderer just uses simple Value mapping
		universalFeed := &gofeed.Feed{
			Title:     "Test Feed",
			Generator: "Jekyll v4.2.0 https://jekyllrb.com/",
			Items:     []*gofeed.Item{},
		}

		converter := &gofeed.DefaultAtomConverter{}
		atomFeed, err := converter.Convert(universalFeed)
		require.NoError(t, err)

		// Generator is simplified in reverse mapping
		require.NotNil(t, atomFeed.Generator)
		assert.Equal(t, "Jekyll v4.2.0 https://jekyllrb.com/", atomFeed.Generator.Value)
		// Version and URI fields are not parsed out - this is acceptable
	})
}

// Test Custom field handling
func TestCustomFieldHandling(t *testing.T) {
	// Note: RSS feed-level custom fields are not supported
	// Custom fields are only supported at the item level in RSS

	t.Run("RSS item-level custom fields", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title:       "Test Feed",
			Description: "Test Description",
			Items: []*gofeed.Item{
				{
					Title: "Test Item",
					Custom: map[string]string{
						"itemCustom": "itemValue",
					},
				},
			},
		}

		converter := &gofeed.DefaultRSSConverter{}
		rssFeed, err := converter.Convert(feed)
		require.NoError(t, err)

		require.Len(t, rssFeed.Items, 1)
		require.NotNil(t, rssFeed.Items[0].Custom)
		assert.Equal(t, "itemValue", rssFeed.Items[0].Custom["itemCustom"])
	})

	t.Run("Atom and JSON don't support custom fields", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title: "Test Feed",
			Custom: map[string]string{
				"custom": "value",
			},
			Items: []*gofeed.Item{
				{
					Custom: map[string]string{
						"itemCustom": "itemValue",
					},
				},
			},
		}

		// Atom doesn't have Custom field support
		atomRenderer := &gofeed.DefaultAtomConverter{}
		atomFeed, err := atomRenderer.Convert(feed)
		require.NoError(t, err)
		// Can't test for absence since Atom struct doesn't have Custom field

		// JSON doesn't have Custom field support
		jsonRenderer := &gofeed.DefaultJSONConverter{}
		jsonFeed, err := jsonRenderer.Convert(feed)
		require.NoError(t, err)
		// Can't test for absence since JSON struct doesn't have Custom field

		// This is expected behavior - Custom is RSS-specific
		assert.NotNil(t, atomFeed)
		assert.NotNil(t, jsonFeed)
	})
}

// Test Item image handling across formats
func TestItemImageHandling(t *testing.T) {
	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Items: []*gofeed.Item{
			{
				Title: "Test Item",
				Image: &gofeed.Image{
					URL:   "https://example.com/item-image.jpg",
					Title: "Item Image",
				},
			},
		},
	}

	t.Run("RSS item image handling via enclosures", func(t *testing.T) {
		converter := &gofeed.DefaultRSSConverter{}
		rssFeed, err := converter.Convert(feed)
		require.NoError(t, err)

		require.Len(t, rssFeed.Items, 1)
		item := rssFeed.Items[0]

		// Should have image as enclosure
		require.NotNil(t, item.Enclosures)
		require.Len(t, item.Enclosures, 1)

		enclosure := item.Enclosures[0]
		assert.Equal(t, "https://example.com/item-image.jpg", enclosure.URL)
		assert.Equal(t, "", enclosure.Type)   // Type not guessed for arbitrary URLs
		assert.Equal(t, "", enclosure.Length) // Length not guessed

		// Primary enclosure should be set
		require.NotNil(t, item.Enclosure)
		assert.Equal(t, "https://example.com/item-image.jpg", item.Enclosure.URL)
	})

	t.Run("Atom item image as enclosure link", func(t *testing.T) {
		converter := &gofeed.DefaultAtomConverter{}
		atomFeed, err := converter.Convert(feed)
		require.NoError(t, err)

		require.Len(t, atomFeed.Entries, 1)
		entry := atomFeed.Entries[0]

		// Should have image as enclosure link
		foundImageLink := false
		for _, link := range entry.Links {
			if link.Href == "https://example.com/item-image.jpg" && link.Rel == "enclosure" {
				foundImageLink = true
				// Type not set for arbitrary URLs
				break
			}
		}
		assert.True(t, foundImageLink, "Item image should be added as enclosure link")
	})

	t.Run("JSON item image native support", func(t *testing.T) {
		converter := &gofeed.DefaultJSONConverter{}
		jsonFeed, err := converter.Convert(feed)
		require.NoError(t, err)

		require.Len(t, jsonFeed.Items, 1)
		assert.Equal(t, "https://example.com/item-image.jpg", jsonFeed.Items[0].Image)
	})
}

// Test edge cases and error conditions
func TestRendererEdgeCases(t *testing.T) {
	t.Run("Empty feed with all converters", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title: "Empty Feed",
			Items: []*gofeed.Item{},
		}

		// RSS converter
		rssRenderer := &gofeed.DefaultRSSConverter{}
		rssFeed, err := rssRenderer.Convert(feed)
		require.NoError(t, err)
		assert.Equal(t, "Empty Feed", rssFeed.Title)
		assert.Len(t, rssFeed.Items, 0)

		// Atom converter
		atomRenderer := &gofeed.DefaultAtomConverter{}
		atomFeed, err := atomRenderer.Convert(feed)
		require.NoError(t, err)
		assert.Equal(t, "Empty Feed", atomFeed.Title)
		assert.Len(t, atomFeed.Entries, 0)

		// JSON converter
		jsonRenderer := &gofeed.DefaultJSONConverter{}
		jsonFeed, err := jsonRenderer.Convert(feed)
		require.NoError(t, err)
		assert.Equal(t, "Empty Feed", jsonFeed.Title)
		assert.Len(t, jsonFeed.Items, 0)
	})

	t.Run("Feed with items but no content", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title: "Minimal Feed",
			Items: []*gofeed.Item{
				{},                    // Completely empty item
				{Title: "Just Title"}, // Item with only title
			},
		}

		// Test RSS converter
		rssRenderer := &gofeed.DefaultRSSConverter{}
		rssResult, err := rssRenderer.Convert(feed)
		require.NoError(t, err)
		assert.NotNil(t, rssResult)

		// Test Atom converter
		atomRenderer := &gofeed.DefaultAtomConverter{}
		atomResult, err := atomRenderer.Convert(feed)
		require.NoError(t, err)
		assert.NotNil(t, atomResult)

		// Test JSON converter
		jsonRenderer := &gofeed.DefaultJSONConverter{}
		jsonResult, err := jsonRenderer.Convert(feed)
		require.NoError(t, err)
		assert.NotNil(t, jsonResult)
	})

	t.Run("Mixed enclosures and images", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title: "Mixed Media Feed",
			Items: []*gofeed.Item{
				{
					Title: "Item with both enclosure and image",
					Enclosures: []*gofeed.Enclosure{
						{URL: "https://example.com/audio.mp3", Type: "audio/mpeg", Length: "12345"},
					},
					Image: &gofeed.Image{URL: "https://example.com/image.jpg"},
				},
			},
		}

		rssRenderer := &gofeed.DefaultRSSConverter{}
		rssFeed, err := rssRenderer.Convert(feed)
		require.NoError(t, err)

		item := rssFeed.Items[0]
		require.Len(t, item.Enclosures, 2) // Original enclosure + image enclosure

		// Original enclosure should be first
		assert.Equal(t, "https://example.com/audio.mp3", item.Enclosures[0].URL)
		assert.Equal(t, "audio/mpeg", item.Enclosures[0].Type)

		// Image enclosure should be second
		assert.Equal(t, "https://example.com/image.jpg", item.Enclosures[1].URL)

		// Primary enclosure should be the audio (first one)
		assert.Equal(t, "https://example.com/audio.mp3", item.Enclosure.URL)
	})

	t.Run("Don't duplicate image if already in enclosures", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title: "Duplicate Prevention Feed",
			Items: []*gofeed.Item{
				{
					Title: "Item with image already in enclosures",
					Enclosures: []*gofeed.Enclosure{
						{URL: "https://example.com/image.jpg", Type: "image/jpeg", Length: "54321"},
					},
					Image: &gofeed.Image{URL: "https://example.com/image.jpg"},
				},
			},
		}

		// RSS converter should not duplicate
		rssRenderer := &gofeed.DefaultRSSConverter{}
		rssFeed, err := rssRenderer.Convert(feed)
		require.NoError(t, err)

		item := rssFeed.Items[0]
		require.Len(t, item.Enclosures, 1) // Should not duplicate
		assert.Equal(t, "https://example.com/image.jpg", item.Enclosures[0].URL)
		assert.Equal(t, "image/jpeg", item.Enclosures[0].Type) // Original type preserved

		// Atom converter should not duplicate
		atomRenderer := &gofeed.DefaultAtomConverter{}
		atomFeed, err := atomRenderer.Convert(feed)
		require.NoError(t, err)

		entry := atomFeed.Entries[0]
		imageEnclosureCount := 0
		for _, link := range entry.Links {
			if link.Href == "https://example.com/image.jpg" && link.Rel == "enclosure" {
				imageEnclosureCount++
			}
		}
		assert.Equal(t, 1, imageEnclosureCount, "Image should only appear once as enclosure")
	})

	t.Run("Authors array vs single author", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title:  "Author Test Feed",
			Author: &gofeed.Person{Name: "Single Author", Email: "single@example.com"},
			Authors: []*gofeed.Person{
				{Name: "First Author", Email: "first@example.com"},
				{Name: "Second Author", Email: "second@example.com"},
			},
			Items: []*gofeed.Item{
				{
					Title:  "Test Item",
					Author: &gofeed.Person{Name: "Item Author", Email: "item@example.com"},
					Authors: []*gofeed.Person{
						{Name: "Item First", Email: "itemfirst@example.com"},
					},
				},
			},
		}

		// RSS should prefer Authors over Author
		rssRenderer := &gofeed.DefaultRSSConverter{}
		rssFeed, err := rssRenderer.Convert(feed)
		require.NoError(t, err)

		// Should use first author from Authors array
		assert.Contains(t, rssFeed.ManagingEditor, "First Author")
		assert.Contains(t, rssFeed.Items[0].Author, "Item First")

		// Atom should preserve all authors
		atomRenderer := &gofeed.DefaultAtomConverter{}
		atomFeed, err := atomRenderer.Convert(feed)
		require.NoError(t, err)

		require.Len(t, atomFeed.Authors, 2)
		assert.Equal(t, "First Author", atomFeed.Authors[0].Name)
		assert.Equal(t, "Second Author", atomFeed.Authors[1].Name)
	})
}

func TestRSSConverter_VersionDefaults(t *testing.T) {
	converter := &gofeed.DefaultRSSConverter{}

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Items:       []*gofeed.Item{},
	}

	rssFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	assert.Equal(t, "2.0", rssFeed.Version, "Should default to RSS 2.0")
}

func TestAtomConverter_VersionDefaults(t *testing.T) {
	converter := &gofeed.DefaultAtomConverter{}

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Items:       []*gofeed.Item{},
	}

	atomFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	assert.Equal(t, "1.0", atomFeed.Version, "Should default to Atom 1.0")
}

func TestJSONConverter_VersionDefaults(t *testing.T) {
	converter := &gofeed.DefaultJSONConverter{}

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Items:       []*gofeed.Item{},
	}

	jsonFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	assert.Equal(t, "https://jsonfeed.org/version/1.1", jsonFeed.Version, "Should default to JSON Feed v1.1")
}

// Test multiple authors handling
func TestAtomConverter_MultipleAuthors(t *testing.T) {
	converter := &gofeed.DefaultAtomConverter{}

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Authors: []*gofeed.Person{
			{Name: "John Doe", Email: "john@example.com"},
			{Name: "Jane Doe", Email: "jane@example.com"},
		},
		Items: []*gofeed.Item{},
	}

	atomFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	require.Len(t, atomFeed.Authors, 2)

	assert.Equal(t, "John Doe", atomFeed.Authors[0].Name)
	assert.Equal(t, "john@example.com", atomFeed.Authors[0].Email)
	assert.Equal(t, "Jane Doe", atomFeed.Authors[1].Name)
	assert.Equal(t, "jane@example.com", atomFeed.Authors[1].Email)
}

func TestJSONConverter_MultipleAuthors(t *testing.T) {
	converter := &gofeed.DefaultJSONConverter{}

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Authors: []*gofeed.Person{
			{Name: "John Doe", Email: "john@example.com"},
			{Name: "Jane Doe", Email: "jane@example.com"},
		},
		Items: []*gofeed.Item{},
	}

	jsonFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	// Should have feed-level author (first one)
	require.NotNil(t, jsonFeed.Author)
	assert.Equal(t, "John Doe <john@example.com>", jsonFeed.Author.Name)

	// Should have all authors in v1.1 format
	require.Len(t, jsonFeed.Authors, 2)
	assert.Equal(t, "John Doe <john@example.com>", jsonFeed.Authors[0].Name)
	assert.Equal(t, "Jane Doe <jane@example.com>", jsonFeed.Authors[1].Name)
}

// Test category handling
func TestRSSConverter_Categories(t *testing.T) {
	converter := &gofeed.DefaultRSSConverter{}

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Categories:  []string{"Tech", "News", "Programming"},
		Items: []*gofeed.Item{
			{
				Title:      "Test Item",
				Categories: []string{"Go", "Testing"},
			},
		},
	}

	rssFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	// Check feed categories
	require.Len(t, rssFeed.Categories, 3)
	assert.Equal(t, "Tech", rssFeed.Categories[0].Value)
	assert.Equal(t, "News", rssFeed.Categories[1].Value)
	assert.Equal(t, "Programming", rssFeed.Categories[2].Value)

	// Check item categories
	require.Len(t, rssFeed.Items, 1)
	require.Len(t, rssFeed.Items[0].Categories, 2)
	assert.Equal(t, "Go", rssFeed.Items[0].Categories[0].Value)
	assert.Equal(t, "Testing", rssFeed.Items[0].Categories[1].Value)
}

func TestJSONConverter_TagsHandling(t *testing.T) {
	converter := &gofeed.DefaultJSONConverter{}

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Description: "Test Description",
		Items: []*gofeed.Item{
			{
				Title:      "Test Item",
				Categories: []string{"golang", "testing", "feeds"},
			},
		},
	}

	jsonFeed, err := converter.Convert(feed)
	require.NoError(t, err)

	require.Len(t, jsonFeed.Items, 1)
	item := jsonFeed.Items[0]
	require.Len(t, item.Tags, 3)
	assert.Equal(t, []string{"golang", "testing", "feeds"}, item.Tags)
}

// Test multi-field author preservation
func TestRSSConverter_MultiFieldAuthorPreservation(t *testing.T) {
	t.Run("Feed-level author populates both ManagingEditor and DublinCore Creator", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title:       "Test Feed",
			Description: "Test Description",
			Authors: []*gofeed.Person{
				{Name: "John Doe", Email: "john@example.com"},
			},
			Items: []*gofeed.Item{},
		}

		converter := &gofeed.DefaultRSSConverter{}
		rssFeed, err := converter.Convert(feed)
		require.NoError(t, err)

		// Should populate ManagingEditor
		assert.Equal(t, "john@example.com (John Doe)", rssFeed.ManagingEditor)

		// Should also populate DublinCore Creator
		require.NotNil(t, rssFeed.DublinCoreExt)
		require.Len(t, rssFeed.DublinCoreExt.Creator, 1)
		assert.Equal(t, "john@example.com (John Doe)", rssFeed.DublinCoreExt.Creator[0])
	})

	t.Run("Item-level author populates Author, DublinCore Creator, and iTunes Author", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title:       "Test Feed",
			Description: "Test Description",
			Items: []*gofeed.Item{
				{
					Title: "Test Item",
					Authors: []*gofeed.Person{
						{Name: "Jane Smith", Email: "jane@example.com"},
					},
					ITunesExt: &ext.ITunesItemExtension{
						Duration: "30:00", // Pre-existing iTunes extension
					},
				},
			},
		}

		converter := &gofeed.DefaultRSSConverter{}
		rssFeed, err := converter.Convert(feed)
		require.NoError(t, err)

		require.Len(t, rssFeed.Items, 1)
		item := rssFeed.Items[0]

		// Should populate native Author field
		assert.Equal(t, "jane@example.com (Jane Smith)", item.Author)

		// Should also populate DublinCore Creator
		require.NotNil(t, item.DublinCoreExt)
		require.Len(t, item.DublinCoreExt.Creator, 1)
		assert.Equal(t, "jane@example.com (Jane Smith)", item.DublinCoreExt.Creator[0])

		// Should also populate iTunes Author
		require.NotNil(t, item.ITunesExt)
		assert.Equal(t, "jane@example.com (Jane Smith)", item.ITunesExt.Author)
		assert.Equal(t, "30:00", item.ITunesExt.Duration) // Should preserve existing iTunes fields
	})

	t.Run("Don't duplicate existing DublinCore Creator", func(t *testing.T) {
		feed := &gofeed.Feed{
			Title:       "Test Feed",
			Description: "Test Description",
			Authors: []*gofeed.Person{
				{Name: "Bob Wilson", Email: "bob@example.com"},
			},
			DublinCoreExt: &ext.DublinCoreExtension{
				Creator: []string{"bob@example.com (Bob Wilson)", "Other Creator"},
			},
			Items: []*gofeed.Item{},
		}

		converter := &gofeed.DefaultRSSConverter{}
		rssFeed, err := converter.Convert(feed)
		require.NoError(t, err)

		// Should not duplicate existing Creator
		require.NotNil(t, rssFeed.DublinCoreExt)
		require.Len(t, rssFeed.DublinCoreExt.Creator, 2)
		assert.Equal(t, "bob@example.com (Bob Wilson)", rssFeed.DublinCoreExt.Creator[0])
		assert.Equal(t, "Other Creator", rssFeed.DublinCoreExt.Creator[1])
	})
}
