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
		feedType gofeed.FeedType
	}{
		{"simple_rss090.xml", gofeed.FeedTypeRSS},
		{"simple_rss091.xml", gofeed.FeedTypeRSS},
		{"simple_rss092.xml", gofeed.FeedTypeRSS},
		{"simple_rss10.xml", gofeed.FeedTypeRSS},
		{"simple_rss20.xml", gofeed.FeedTypeRSS},
		{"simple_atom10.xml", gofeed.FeedTypeAtom},
		{"invalid.xml", gofeed.FeedTypeUnknown},
	}

	fp := gofeed.NewFeedParser()
	for _, test := range verTests {
		file := fmt.Sprintf("testdata/%s", test.file)
		f, _ := ioutil.ReadFile(file)
		actual := fp.DetectFeedType(string(f))

		assert.Equal(t, test.feedType, actual, "Expected feed type %d, got %d in %s", test.feedType, actual, test.file)
	}
}
