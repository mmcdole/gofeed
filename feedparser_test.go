package gofeed_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestDetectFeedType(t *testing.T) {
	var feedTypeTests = []struct {
		file     string
		expected gofeed.FeedType
	}{
		{"feedtype_atom03.xml", gofeed.FeedTypeAtom},
		{"feedtype_atom10.xml", gofeed.FeedTypeAtom},
		{"feedtype_rss.xml", gofeed.FeedTypeRSS},
		{"feedtype_rdf.xml", gofeed.FeedTypeRSS},
		{"feedtype_unknown.xml", gofeed.FeedTypeUnknown},
	}

	for _, test := range feedTypeTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/feed/%s", test.file)
		f, _ := ioutil.ReadFile(path)

		// Get actual value
		actual := gofeed.DetectFeedType(string(f))

		if assert.Equal(t, actual, test.expected, "Feed file %s did not match expected type %d", test.file, test.expected) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

func ExampleDetectFeedType() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	feedType := gofeed.DetectFeedType(feedata)
}
