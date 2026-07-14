package gofeed_test

import (
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestExtensionHelpers(t *testing.T) {
	feed := `<rss version="2.0" xmlns:dc="http://purl.org/dc/elements/1.1/">
		<channel>
			<customFeedId>feed-1</customFeedId>
			<item>
				<dc:creator>Jane</dc:creator>
				<event><venue city="Austin">Hall</venue></event>
				<simple>plain</simple>
				<simple>again</simple>
			</item>
		</channel>
	</rss>`

	f, err := gofeed.NewParser().Parse(strings.NewReader(feed))
	assert.NoError(t, err)
	item := f.Items[0]

	// Namespaced extensions through the same accessors.
	assert.Equal(t, "Jane", item.GetExtensionValue("dc", "creator"))
	assert.Len(t, item.GetExtension("dc", "creator"), 1)

	// Feed level custom element, previously dropped entirely.
	assert.Equal(t, "feed-1", f.GetCustomValue("customFeedId"))

	// Item level custom values, including repetition the flat map loses.
	assert.Equal(t, "plain", item.GetCustomValue("simple"))
	assert.Len(t, item.GetExtension(gofeed.CustomNamespace, "simple"), 2)

	// Nested custom elements keep children and attributes in the tree.
	events := item.GetExtension(gofeed.CustomNamespace, "event")
	if assert.Len(t, events, 1) {
		venue := events[0].Children["venue"][0]
		assert.Equal(t, "Hall", venue.Value)
		assert.Equal(t, "Austin", venue.Attrs["city"])
	}

	// The flat Custom map still works for childless elements (last wins),
	// and nested elements no longer corrupt it.
	assert.Equal(t, "again", item.Custom["simple"])
	_, hasEvent := item.Custom["event"]
	assert.False(t, hasEvent)

	// Absent lookups are empty, not panics.
	assert.Nil(t, f.GetExtension("nope", "x"))
	assert.Equal(t, "", item.GetCustomValue("absent"))
}
