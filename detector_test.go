package gofeed_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/v2"
	"github.com/stretchr/testify/assert"
)

func TestDetectFeedType(t *testing.T) {
	var feedTypeTests = []struct {
		file     string
		expected gofeed.FeedType
	}{
		{"atom03_feed.xml", gofeed.FeedTypeAtom},
		{"atom10_feed.xml", gofeed.FeedTypeAtom},
		{"rss_feed.xml", gofeed.FeedTypeRSS},
		{"rss_feed_bom.xml", gofeed.FeedTypeRSS},
		{"rss_feed_leading_spaces.xml", gofeed.FeedTypeRSS},
		{"rdf_feed.xml", gofeed.FeedTypeRSS},
		{"unknown_feed.xml", gofeed.FeedTypeUnknown},
		{"empty_feed.xml", gofeed.FeedTypeUnknown},
		{"json10_feed.json", gofeed.FeedTypeJSON},
	}

	for _, test := range feedTypeTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/universal/%s", test.file)
		f, _ := os.ReadFile(path)

		// Get actual value
		actual := gofeed.DetectFeedType(bytes.NewReader(f))

		if assert.Equal(t, actual, test.expected, "Feed file %s did not match expected type %d", test.file, test.expected) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

// Examples

func ExampleDetectFeedType() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	feedType := gofeed.DetectFeedType(strings.NewReader(feedData))
	if feedType == gofeed.FeedTypeRSS {
		fmt.Println("Wow! This is an RSS feed!")
	}
}
