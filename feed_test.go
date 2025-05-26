package gofeed_test

import (
	"sort"
	"testing"
	"time"

	"github.com/mmcdole/gofeed/v2"
	ext "github.com/mmcdole/gofeed/v2/extensions"
)

func TestFeedSort(t *testing.T) {
	oldestItem := &gofeed.Item{
		PublishedParsed: &[]time.Time{time.Unix(0, 0)}[0],
	}
	inbetweenItem := &gofeed.Item{
		PublishedParsed: &[]time.Time{time.Unix(1, 0)}[0],
	}
	newestItem := &gofeed.Item{
		PublishedParsed: &[]time.Time{time.Unix(2, 0)}[0],
	}

	feed := gofeed.Feed{
		Items: []*gofeed.Item{
			newestItem,
			oldestItem,
			inbetweenItem,
		},
	}
	expected := gofeed.Feed{
		Items: []*gofeed.Item{
			oldestItem,
			inbetweenItem,
			newestItem,
		},
	}

	sort.Sort(feed)

	for i, item := range feed.Items {
		if !item.PublishedParsed.Equal(
			*expected.Items[i].PublishedParsed,
		) {
			t.Errorf(
				"Item PublishedParsed = %s; want %s",
				item.PublishedParsed,
				expected.Items[i].PublishedParsed,
			)
		}
	}
}

func TestItemGetExtension(t *testing.T) {
	item := &gofeed.Item{
		Extensions: ext.Extensions{
			"dc": {
				"creator": []ext.Extension{
					{Name: "creator", Value: "John Doe"},
				},
			},
			"_custom": {
				"customField": []ext.Extension{
					{Name: "customField", Value: "Custom Value", Attrs: map[string]string{"id": "123"}},
				},
			},
		},
	}

	// Test getting existing extension
	creators := item.GetExtension("dc", "creator")
	if len(creators) != 1 {
		t.Errorf("Expected 1 creator, got %d", len(creators))
	}
	if creators[0].Value != "John Doe" {
		t.Errorf("Expected creator value 'John Doe', got '%s'", creators[0].Value)
	}

	// Test getting custom element
	custom := item.GetExtension("_custom", "customField")
	if len(custom) != 1 {
		t.Errorf("Expected 1 custom field, got %d", len(custom))
	}
	if custom[0].Value != "Custom Value" {
		t.Errorf("Expected custom value 'Custom Value', got '%s'", custom[0].Value)
	}
	if custom[0].Attrs["id"] != "123" {
		t.Errorf("Expected custom attr id '123', got '%s'", custom[0].Attrs["id"])
	}

	// Test getting non-existent extension
	missing := item.GetExtension("missing", "field")
	if missing != nil {
		t.Errorf("Expected nil for missing extension, got %v", missing)
	}
}

func TestItemGetExtensionValue(t *testing.T) {
	item := &gofeed.Item{
		Extensions: ext.Extensions{
			"dc": {
				"creator": []ext.Extension{
					{Name: "creator", Value: "John Doe"},
				},
			},
			"_custom": {
				"customField": []ext.Extension{
					{Name: "customField", Value: "Custom Value"},
				},
			},
		},
	}

	// Test getting existing value
	if v := item.GetExtensionValue("dc", "creator"); v != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", v)
	}
	if v := item.GetExtensionValue("_custom", "customField"); v != "Custom Value" {
		t.Errorf("Expected 'Custom Value', got '%s'", v)
	}

	// Test getting non-existent value
	if v := item.GetExtensionValue("missing", "field"); v != "" {
		t.Errorf("Expected empty string, got '%s'", v)
	}
}

func TestItemGetCustomValue(t *testing.T) {
	item := &gofeed.Item{
		Extensions: ext.Extensions{
			"_custom": {
				"customField": []ext.Extension{
					{Name: "customField", Value: "Custom Value"},
				},
				"anotherField": []ext.Extension{
					{Name: "anotherField", Value: "Another Value"},
				},
			},
		},
	}

	// Test getting custom values
	if v := item.GetCustomValue("customField"); v != "Custom Value" {
		t.Errorf("Expected 'Custom Value', got '%s'", v)
	}
	if v := item.GetCustomValue("anotherField"); v != "Another Value" {
		t.Errorf("Expected 'Another Value', got '%s'", v)
	}

	// Test getting non-existent custom value
	if v := item.GetCustomValue("missing"); v != "" {
		t.Errorf("Expected empty string, got '%s'", v)
	}
}

func TestFeedGetExtension(t *testing.T) {
	feed := &gofeed.Feed{
		Extensions: ext.Extensions{
			"sy": {
				"updatePeriod": []ext.Extension{
					{Name: "updatePeriod", Value: "hourly"},
				},
			},
			"_custom": {
				"customFeedData": []ext.Extension{
					{Name: "customFeedData", Value: "Feed Custom Value"},
				},
			},
		},
	}

	// Test getting existing extension
	period := feed.GetExtension("sy", "updatePeriod")
	if len(period) != 1 {
		t.Errorf("Expected 1 updatePeriod, got %d", len(period))
	}
	if period[0].Value != "hourly" {
		t.Errorf("Expected 'hourly', got '%s'", period[0].Value)
	}

	// Test getting custom element
	custom := feed.GetExtension("_custom", "customFeedData")
	if len(custom) != 1 {
		t.Errorf("Expected 1 custom feed data, got %d", len(custom))
	}
	if custom[0].Value != "Feed Custom Value" {
		t.Errorf("Expected 'Feed Custom Value', got '%s'", custom[0].Value)
	}
}

func TestFeedGetCustomValue(t *testing.T) {
	feed := &gofeed.Feed{
		Extensions: ext.Extensions{
			"_custom": {
				"customFeedId": []ext.Extension{
					{Name: "customFeedId", Value: "feed-123"},
				},
				"updateFrequency": []ext.Extension{
					{Name: "updateFrequency", Value: "hourly"},
				},
			},
		},
	}

	// Test getting custom values at feed level
	if v := feed.GetCustomValue("customFeedId"); v != "feed-123" {
		t.Errorf("Expected 'feed-123', got '%s'", v)
	}
	if v := feed.GetCustomValue("updateFrequency"); v != "hourly" {
		t.Errorf("Expected 'hourly', got '%s'", v)
	}

	// Test getting non-existent custom value
	if v := feed.GetCustomValue("missing"); v != "" {
		t.Errorf("Expected empty string, got '%s'", v)
	}
}

func TestMultipleExtensionsWithSameName(t *testing.T) {
	item := &gofeed.Item{
		Extensions: ext.Extensions{
			"_custom": {
				"tag": []ext.Extension{
					{Name: "tag", Value: "First"},
					{Name: "tag", Value: "Second"},
					{Name: "tag", Value: "Third"},
				},
			},
		},
	}

	// Test getting all tags
	tags := item.GetExtension("_custom", "tag")
	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}
	if tags[0].Value != "First" {
		t.Errorf("Expected 'First', got '%s'", tags[0].Value)
	}
	if tags[1].Value != "Second" {
		t.Errorf("Expected 'Second', got '%s'", tags[1].Value)
	}
	if tags[2].Value != "Third" {
		t.Errorf("Expected 'Third', got '%s'", tags[2].Value)
	}

	// GetExtensionValue returns the first one
	if v := item.GetExtensionValue("_custom", "tag"); v != "First" {
		t.Errorf("Expected 'First' (first value), got '%s'", v)
	}
}