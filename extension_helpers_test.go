package gofeed

import (
	"testing"

	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/stretchr/testify/assert"
)

func TestItemGetExtension(t *testing.T) {
	item := &Item{
		Extensions: ext.Extensions{
			"dc": {
				"creator": []ext.Extension{
					{Name: "creator", Value: "John Doe"},
				},
			},
			"rss": {
				"customField": []ext.Extension{
					{Name: "customField", Value: "Custom Value", Attrs: map[string]string{"id": "123"}},
				},
			},
		},
	}

	// Test getting existing extension
	creators := item.GetExtension("dc", "creator")
	assert.Len(t, creators, 1)
	assert.Equal(t, "John Doe", creators[0].Value)

	// Test getting custom RSS element
	custom := item.GetExtension("rss", "customField")
	assert.Len(t, custom, 1)
	assert.Equal(t, "Custom Value", custom[0].Value)
	assert.Equal(t, "123", custom[0].Attrs["id"])

	// Test getting non-existent extension
	missing := item.GetExtension("missing", "field")
	assert.Nil(t, missing)
}

func TestItemGetExtensionValue(t *testing.T) {
	item := &Item{
		Extensions: ext.Extensions{
			"dc": {
				"creator": []ext.Extension{
					{Name: "creator", Value: "John Doe"},
				},
			},
			"rss": {
				"customField": []ext.Extension{
					{Name: "customField", Value: "Custom Value"},
				},
			},
		},
	}

	// Test getting existing value
	assert.Equal(t, "John Doe", item.GetExtensionValue("dc", "creator"))
	assert.Equal(t, "Custom Value", item.GetExtensionValue("rss", "customField"))

	// Test getting non-existent value
	assert.Equal(t, "", item.GetExtensionValue("missing", "field"))
}

func TestItemGetCustomValue(t *testing.T) {
	item := &Item{
		Extensions: ext.Extensions{
			"rss": {
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
	assert.Equal(t, "Custom Value", item.GetCustomValue("customField"))
	assert.Equal(t, "Another Value", item.GetCustomValue("anotherField"))

	// Test getting non-existent custom value
	assert.Equal(t, "", item.GetCustomValue("missing"))
}

func TestFeedGetExtension(t *testing.T) {
	feed := &Feed{
		Extensions: ext.Extensions{
			"sy": {
				"updatePeriod": []ext.Extension{
					{Name: "updatePeriod", Value: "hourly"},
				},
			},
			"rss": {
				"customFeedData": []ext.Extension{
					{Name: "customFeedData", Value: "Feed Custom Value"},
				},
			},
		},
	}

	// Test getting existing extension
	period := feed.GetExtension("sy", "updatePeriod")
	assert.Len(t, period, 1)
	assert.Equal(t, "hourly", period[0].Value)

	// Test getting custom RSS element
	custom := feed.GetExtension("rss", "customFeedData")
	assert.Len(t, custom, 1)
	assert.Equal(t, "Feed Custom Value", custom[0].Value)
}

func TestMultipleExtensionsWithSameName(t *testing.T) {
	item := &Item{
		Extensions: ext.Extensions{
			"rss": {
				"tag": []ext.Extension{
					{Name: "tag", Value: "First"},
					{Name: "tag", Value: "Second"},
					{Name: "tag", Value: "Third"},
				},
			},
		},
	}

	// Test getting all tags
	tags := item.GetExtension("rss", "tag")
	assert.Len(t, tags, 3)
	assert.Equal(t, "First", tags[0].Value)
	assert.Equal(t, "Second", tags[1].Value)
	assert.Equal(t, "Third", tags[2].Value)

	// GetExtensionValue returns the first one
	assert.Equal(t, "First", item.GetExtensionValue("rss", "tag"))
}