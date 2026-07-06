package gofeed_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/rss"
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

func ExampleParser_ParseURL_basicAuth() {
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

const concurrencyFeed = `<rss version="2.0"><channel><title>t</title><item><title>i</title></item></channel></rss>`

// TestParserConcurrentParseString shares one Parser across goroutines. Before
// the lazy-init fix this races on the AtomTranslator/RSSTranslator/JSONTranslator
// fields under -race.
func TestParserConcurrentParseString(t *testing.T) {
	p := gofeed.NewParser()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := p.ParseString(concurrencyFeed); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

// TestParserConcurrentParseURL exercises the httpClient() lazy init the same way.
func TestParserConcurrentParseURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, concurrencyFeed)
	}))
	defer srv.Close()

	p := gofeed.NewParser()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := p.ParseURL(srv.URL); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}

func TestParseURLMaxByteSize(t *testing.T) {
	big := `<rss version="2.0"><channel><title>` + strings.Repeat("x", 200000) + `</title></channel></rss>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, big)
	}))
	defer srv.Close()

	p := gofeed.NewParser()

	p.MaxByteSize = 1000
	if _, err := p.ParseURL(srv.URL); !errors.Is(err, gofeed.ErrResponseTooLarge) {
		t.Errorf("small limit: got %v, want ErrResponseTooLarge", err)
	}

	p.MaxByteSize = 10_000_000
	if _, err := p.ParseURL(srv.URL); err != nil {
		t.Errorf("large limit: unexpected error %v", err)
	}

	p.MaxByteSize = 0 // unlimited
	if _, err := p.ParseURL(srv.URL); err != nil {
		t.Errorf("unlimited: unexpected error %v", err)
	}
}

// Confirms request-context cancellation works, which the default ParseURL
// timeout relies on.
func TestParseURLContextTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	if _, err := gofeed.NewParser().ParseURLWithContext(srv.URL, ctx); err == nil {
		t.Fatal("expected a timeout error")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("did not time out promptly: %v", elapsed)
	}
}

func TestParserKeepOriginalFeed(t *testing.T) {
	const feed = `<rss version="2.0"><channel><title>t</title><item><title>i</title></item></channel></rss>`

	// Off by default: OriginalFeed() is nil.
	p := gofeed.NewParser()
	f, err := p.ParseString(feed)
	if err != nil {
		t.Fatal(err)
	}
	if f.OriginalFeed() != nil {
		t.Errorf("OriginalFeed() = %T, want nil when KeepOriginalFeed is off", f.OriginalFeed())
	}

	// On: OriginalFeed() returns the source *rss.Feed.
	p.KeepOriginalFeed = true
	f, err = p.ParseString(feed)
	if err != nil {
		t.Fatal(err)
	}
	orig, ok := f.OriginalFeed().(*rss.Feed)
	if !ok {
		t.Fatalf("OriginalFeed() = %T, want *rss.Feed", f.OriginalFeed())
	}
	if orig.Title != "t" {
		t.Errorf("original feed title = %q, want %q", orig.Title, "t")
	}
}

// An I/O error from the reader must surface as itself, not be masked as a
// failed type detection (issue #311).
func TestParser_Parse_ReaderError(t *testing.T) {
	boom := errors.New("boom")
	r := io.MultiReader(strings.NewReader(`<rss version="2.0"><channel>`), iotest.ErrReader(boom))

	_, err := gofeed.NewParser().Parse(r)
	assert.ErrorIs(t, err, boom)
}
