package gofeed_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	var feedTests = []struct {
		file      string
		feedType  string
		feedTitle string
		hasError  bool
	}{
		{"atom03_feed.xml", "atom", "Feed Title", false},
		{"atom10_feed.xml", "atom", "Feed Title", false},
		{"rss_feed.xml", "rss", "Feed Title", false},
		{"rss_feed_bom.xml", "rss", "Feed Title", false},
		{"rss_feed_leading_spaces.xml", "rss", "Feed Title", false},
		{"rdf_feed.xml", "rss", "Feed Title", false},
		{"sample.json", "json", "title", false},
		{"json10_feed.json", "json", "title", false},
		{"json11_feed.json", "json", "title", false},
		{"unknown_feed.xml", "", "", true},
		{"empty_feed.xml", "", "", true},
		{"invalid.json", "", "", true},
	}

	for _, test := range feedTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/universal/%s", test.file)
		f, _ := os.ReadFile(path)

		// Get actual value
		fp := gofeed.NewParser()
		feed, err := fp.Parse(bytes.NewReader(f))

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

func TestParser_ParseString(t *testing.T) {
	var feedTests = []struct {
		file      string
		feedType  string
		feedTitle string
		hasError  bool
	}{
		{"atom03_feed.xml", "atom", "Feed Title", false},
		{"atom10_feed.xml", "atom", "Feed Title", false},
		{"rss_feed.xml", "rss", "Feed Title", false},
		{"rss_feed_bom.xml", "rss", "Feed Title", false},
		{"rss_feed_leading_spaces.xml", "rss", "Feed Title", false},
		{"rdf_feed.xml", "rss", "Feed Title", false},
		{"sample.json", "json", "title", false},
		{"unknown_feed.xml", "", "", true},
		{"empty_feed.xml", "", "", true},
		{"invalid.json", "", "", true},
	}

	for _, test := range feedTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/universal/%s", test.file)
		f, _ := os.ReadFile(path)

		// Get actual value
		fp := gofeed.NewParser()
		feed, err := fp.ParseString(string(f))

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

func TestParser_ParseURL_Success(t *testing.T) {
	var feedTests = []struct {
		file      string
		feedType  string
		feedTitle string
		hasError  bool
	}{
		{"atom03_feed.xml", "atom", "Feed Title", false},
		{"atom10_feed.xml", "atom", "Feed Title", false},
		{"rss_feed.xml", "rss", "Feed Title", false},
		{"rss_feed_bom.xml", "rss", "Feed Title", false},
		{"rss_feed_leading_spaces.xml", "rss", "Feed Title", false},
		{"rdf_feed.xml", "rss", "Feed Title", false},
		{"json10_feed.json", "json", "title", false},
		{"json11_feed.json", "json", "title", false},
		{"unknown_feed.xml", "", "", true},
		{"invalid.json", "", "", true},
	}

	for _, test := range feedTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/universal/%s", test.file)
		f, _ := os.ReadFile(path)

		// Get actual value
		server, client := mockServerResponse(200, string(f), 0)
		fp := gofeed.NewParser()
		fp.Client = client
		feed, err := fp.ParseURL(server.URL)

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

func TestParser_ParseURLWithContext(t *testing.T) {
	server, client := mockServerResponse(404, "", 1*time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	fp.Client = client
	_, err := fp.ParseURLWithContext(server.URL, ctx)
	assert.True(t, strings.Contains(err.Error(), ctx.Err().Error()))
}

func TestParser_ParseURL_Failure(t *testing.T) {
	server, client := mockServerResponse(404, "", 0)
	fp := gofeed.NewParser()
	fp.Client = client
	feed, err := fp.ParseURL(server.URL)

	assert.NotNil(t, err)
	assert.IsType(t, gofeed.HTTPError{}, err)
	assert.Nil(t, feed)
}

func TestParser_ParseURLWithContextAndBasicAuth(t *testing.T) {
	server, client := mockServerResponse(404, "", 1*time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	fp.AuthConfig = &gofeed.Auth{
		Username: "foo",
		Password: "bar",
	}
	fp.Client = client
	_, err := fp.ParseURLWithContext(server.URL, ctx)
	assert.True(t, strings.Contains(err.Error(), ctx.Err().Error()))
}

// to detect race conditions, run with go test -race
func TestParser_Concurrent(t *testing.T) {

	var feedTests = []string{"atom03_feed.xml", "atom10_feed.xml", "rss_feed.xml", "rss_feed_bom.xml",
		"rss_feed_leading_spaces.xml", "rdf_feed.xml", "json10_feed.json",
		"json11_feed.json"}

	fp := gofeed.NewParser()
	fp.AtomTranslator = &gofeed.DefaultAtomTranslator{}
	fp.RSSTranslator = &gofeed.DefaultRSSTranslator{}
	fp.JSONTranslator = &gofeed.DefaultJSONTranslator{}
	wg := sync.WaitGroup{}
	for _, test := range feedTests {
		fmt.Printf("\nTesting concurrently %s... ", test)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/universal/%s", test)
		f, _ := os.ReadFile(path)

		wg.Add(1)
		go func() {
			defer wg.Done()
			fp.ParseString(string(f))
		}()
	}
	wg.Wait()
}

// Test Helpers

func mockServerResponse(code int, body string, delay time.Duration) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
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

// Examples

func ExampleParser_Parse() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	fp := gofeed.NewParser()
	feed, err := fp.Parse(strings.NewReader(feedData))
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}

func ExampleParser_ParseURL() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("http://feeds.twit.tv/twit.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}

func ExampleParser_ParseString() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(feedData)
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}

func ExampleParserWithBasicAuth_ParseURL() {
	fp := gofeed.NewParser()
	fp.AuthConfig = &gofeed.Auth{
		Username: "foo",
		Password: "bar",
	}
	feed, err := fp.ParseURL("http://feeds.twit.tv/twit.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(feed.Title)
}
