package feed_test

import (
	"io/ioutil"
	"testing"

	"github.com/mmcdole/go-feed"
	"github.com/stretchr/testify/assert"
)

func TestDetectFeedType_RSS090(t *testing.T) {
	f, _ := ioutil.ReadFile("test/simple_rss090.xml")

	fp := feed.NewFeedParser()
	expected := feed.FeedTypeRSS
	result := fp.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_RSS091(t *testing.T) {
	f, _ := ioutil.ReadFile("test/simple_rss091.xml")

	fp := feed.NewFeedParser()
	expected := feed.FeedTypeRSS
	result := fp.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_RSS092(t *testing.T) {
	f, _ := ioutil.ReadFile("test/simple_rss092.xml")

	fp := feed.NewFeedParser()
	expected := feed.FeedTypeRSS
	result := fp.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_RSS10(t *testing.T) {
	f, _ := ioutil.ReadFile("test/simple_rss10.xml")

	fp := feed.NewFeedParser()
	expected := feed.FeedTypeRSS
	result := fp.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_RSS20(t *testing.T) {
	f, _ := ioutil.ReadFile("test/simple_rss20.xml")

	fp := feed.NewFeedParser()
	expected := feed.FeedTypeRSS
	result := fp.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_Atom10(t *testing.T) {
	f, _ := ioutil.ReadFile("test/simple_atom10.xml")

	fp := feed.NewFeedParser()
	expected := feed.FeedTypeAtom
	result := fp.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_JunkData(t *testing.T) {
	f := `like tears in the rain`

	fp := feed.NewFeedParser()
	expected := feed.FeedTypeUnknown
	result := fp.DetectFeedType(f)
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}
