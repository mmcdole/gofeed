package feed_test

import (
	"io/ioutil"
	"testing"

	"github.com/mmcdole/go-feed"
)

func TestDetectFeedType_RSS10(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_rss10.xml")

	expectedType := feed.FeedTypeRSS
	result := feed.DetectFeedType(string(f))
	if result != expectedType {
		t.Fatalf("Expected FeedType %d, got %d", expectedType, result)
	}
}

func TestDetectFeedType_RSS20(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_rss20.xml")

	expectedType := feed.FeedTypeRSS
	result := feed.DetectFeedType(string(f))
	if result != expectedType {
		t.Fatalf("Expected FeedType %d, got %d", expectedType, result)
	}
}

func TestDetectFeedType_Atom10(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/simple_atom10.xml")

	expectedType := feed.FeedTypeAtom
	result := feed.DetectFeedType(string(f))
	if result != expectedType {
		t.Fatalf("Expected FeedType %d, got %d", expectedType, result)
	}
}

func TestDetectFeedType_JunkData(t *testing.T) {
	f := `like tears in the rain`

	expectedType := feed.FeedTypeUnknown
	result := feed.DetectFeedType(f)
	if result != expectedType {
		t.Fatalf("Expected FeedType %d, got %d", expectedType, result)
	}
}
