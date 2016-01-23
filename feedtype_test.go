package feed_test

import (
	"io/ioutil"
	"testing"

	"github.com/mmcdole/go-feed"
	"github.com/stretchr/testify/assert"
)

func TestDetectFeedType_RSS090(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_rss090.xml")

	expected := feed.FeedTypeRSS
	result := feed.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_RSS091(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_rss091.xml")

	expected := feed.FeedTypeRSS
	result := feed.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_RSS092(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_rss092.xml")

	expected := feed.FeedTypeRSS
	result := feed.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_RSS10(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_rss10.xml")

	expected := feed.FeedTypeRSS
	result := feed.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_RSS20(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_rss20.xml")

	expected := feed.FeedTypeRSS
	result := feed.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_Atom10(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_atom10.xml")

	expected := feed.FeedTypeAtom
	result := feed.DetectFeedType(string(f))
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}

func TestDetectFeedType_JunkData(t *testing.T) {
	f := `like tears in the rain`

	expected := feed.FeedTypeUnknown
	result := feed.DetectFeedType(f)
	assert.Equal(t, expected, result, "Expected FeedType %d, got %d", expected, result)
}
