package gofeed_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestFeedParser_DetectFeedType(t *testing.T) {
	var verTests = []struct {
		file     string
		feedType feed.FeedType
	}{
		{"rss/simple_rss090.xml", feed.FeedTypeRSS},
		{"rss/simple_rss091.xml", feed.FeedTypeRSS},
		{"rss/simple_rss092.xml", feed.FeedTypeRSS},
		{"rss/simple_rss10.xml", feed.FeedTypeRSS},
		{"rss/simple_rss20.xml", feed.FeedTypeRSS},
		{"atom/simple_atom10.xml", feed.FeedTypeAtom},
		{"invalid.xml", feed.FeedTypeUnknown},
	}

	fp := feed.NewFeedParser()
	for _, test := range verTests {
		file := fmt.Sprintf("test/%s", test.file)
		f, _ := ioutil.ReadFile(file)

		actual := fp.DetectFeedType(string(f))

		assert.Equal(t, test.feedType, actual, "Expected feed type %d, got %d", test.feedType, actual)
	}
}
