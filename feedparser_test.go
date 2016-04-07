package gofeed_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mmcdole/gofeed"
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
		{"rdf_feed.xml", gofeed.FeedTypeRSS},
		{"unknown_feed.xml", gofeed.FeedTypeUnknown},
		{"empty_feed.xml", gofeed.FeedTypeUnknown},
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

func TestFeedParser_ParseFeed(t *testing.T) {
	var feedTests = []struct {
		file      string
		feedType  string
		feedTitle string
		hasError  bool
	}{
		{"atom03_feed.xml", "atom", "Atom Title", false},
		{"atom10_feed.xml", "atom", "Atom Title", false},
		{"rss_feed.xml", "rss", "RSS Title", false},
		{"rdf_feed.xml", "rss", "RDF Title", false},
		{"unknown_feed.xml", "", "", true},
		{"empty_feed.xml", "", "", true},
	}

	for _, test := range feedTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/feed/%s", test.file)
		f, _ := ioutil.ReadFile(path)

		// Get actual value
		fp := gofeed.NewFeedParser()
		feed, err := fp.ParseFeed(string(f))

		if test.hasError {
			assert.NotNil(t, err)
			assert.Nil(t, feed)
		} else {
			assert.NotNil(t, feed)
			assert.Nil(t, err)
			assert.Equal(t, feed.FeedType, test.feedType)
			assert.Equal(t, feed.Title, test.feedTitle)
		}
	}
}

func TestFeedParser_ParseFeedURL_Success(t *testing.T) {
	var feedTests = []struct {
		file      string
		feedType  string
		feedTitle string
		hasError  bool
	}{
		{"atom03_feed.xml", "atom", "Atom Title", false},
		{"atom10_feed.xml", "atom", "Atom Title", false},
		{"rss_feed.xml", "rss", "RSS Title", false},
		{"rdf_feed.xml", "rss", "RDF Title", false},
		{"unknown_feed.xml", "", "", true},
	}

	for _, test := range feedTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/feed/%s", test.file)
		f, _ := ioutil.ReadFile(path)

		// Get actual value
		server, client := mockServerResponse(200, string(f))
		fp := gofeed.NewFeedParser()
		fp.Client = client
		feed, err := fp.ParseFeedURL(server.URL)

		if test.hasError {
			assert.NotNil(t, err)
			assert.Nil(t, feed)
		} else {
			assert.NotNil(t, feed)
			assert.Nil(t, err)
			assert.Equal(t, feed.FeedType, test.feedType)
			assert.Equal(t, feed.Title, test.feedTitle)
		}
	}
}

func TestFeedParser_ParseFeedURL_Failure(t *testing.T) {
	server, client := mockServerResponse(404, "")
	fp := gofeed.NewFeedParser()
	fp.Client = client
	feed, err := fp.ParseFeedURL(server.URL)

	assert.NotNil(t, err)
	assert.Nil(t, feed)
}

func ExampleDetectFeedType() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	feedType := gofeed.DetectFeedType(feedData)
	if feedType == gofeed.FeedTypeRSS {
		fmt.Println("Wow! This is an RSS feed!")
	}
}

func ExampleFeedParser_ParseFeedURL() {
	fp := gofeed.NewFeedParser()
	feed, err := fp.ParseFeedURL("http://feeds.twit.tv/twit.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}

func ExampleFeedParser_ParseFeed() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	fp := gofeed.NewFeedParser()
	feed, err := fp.ParseFeed(feedData)
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}

func mockServerResponse(code int, body string) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/xml")
		io.WriteString(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	client := &http.Client{Transport: transport}
	return server, client
}
