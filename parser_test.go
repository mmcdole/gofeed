package gofeed_test

import (
	"bytes"
	"context"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/internal/testutil"
	jsonfeed "github.com/mmcdole/gofeed/json"
	"github.com/mmcdole/gofeed/rss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_NullJSONElements(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
	}{
		{"first_item", "items[0]"},
		{"later_item", "items[1]"},
		{"first_author", "authors[0]"},
		{"later_author", "authors[1]"},
		{"first_item_author", "items[0].authors[0]"},
		{"later_item_author", "items[1].authors[1]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := testutil.ReadFile(t, "testdata/parser/json/invalid/issue_355_"+tc.name+".json")
			want := "gofeed: JSON feed " + tc.path + " must not be null"
			t.Run("json", func(t *testing.T) {
				feed, err := (&jsonfeed.Parser{}).Parse(bytes.NewReader(data))
				require.EqualError(t, err, want)
				require.Nil(t, feed)
			})
			t.Run("universal", func(t *testing.T) {
				feed, err := gofeed.NewParser().Parse(bytes.NewReader(data))
				require.EqualError(t, err, want)
				require.Nil(t, feed)
			})
		})
	}
}

func TestParser_NullOptionalJSONFields(t *testing.T) {
	for _, items := range []string{"null", `[{"id":"a","content_text":"text","author":null,"authors":null}]`} {
		t.Run(items, func(t *testing.T) {
			data := `{"version":"https://jsonfeed.org/version/1.1","title":"Null fields","author":null,"authors":null,"items":` + items + `}`
			feed, err := gofeed.NewParser().ParseString(data)
			require.NoError(t, err)
			require.NotNil(t, feed)
			require.Nil(t, feed.Author)
			require.Nil(t, feed.Authors)
			if items == "null" {
				require.Empty(t, feed.Items)
			} else {
				require.Len(t, feed.Items, 1)
				require.Equal(t, "a", feed.Items[0].GUID)
				require.Equal(t, "text", feed.Items[0].Content)
				require.Nil(t, feed.Items[0].Author)
				require.Nil(t, feed.Items[0].Authors)
			}
		})
	}
}

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
		t.Run(test.file, func(t *testing.T) {
			data := testutil.ReadFile(t, "testdata/parser/universal/"+test.file)
			fp := gofeed.NewParser()
			feed, err := fp.Parse(bytes.NewReader(data))
			if test.hasError {
				assert.Error(t, err)
				assert.Nil(t, feed)
			} else {
				require.NoError(t, err)
				require.NotNil(t, feed)
				assert.Equal(t, test.feedType, feed.FeedType)
				assert.Equal(t, test.feedTitle, feed.Title)
			}
		})
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
		t.Run(test.file, func(t *testing.T) {
			data := testutil.ReadFile(t, "testdata/parser/universal/"+test.file)
			fp := gofeed.NewParser()
			feed, err := fp.ParseString(string(data))
			if test.hasError {
				assert.Error(t, err)
				assert.Nil(t, feed)
			} else {
				require.NoError(t, err)
				require.NotNil(t, feed)
				assert.Equal(t, test.feedType, feed.FeedType)
				assert.Equal(t, test.feedTitle, feed.Title)
			}
		})
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
		t.Run(test.file, func(t *testing.T) {
			data := testutil.ReadFile(t, "testdata/parser/universal/"+test.file)
			fp := gofeed.NewParser()
			server, client := mockServerResponse(t, 200, string(data), 0)
			fp.Client = client
			feed, err := fp.ParseURL(server.URL)
			if test.hasError {
				assert.Error(t, err)
				assert.Nil(t, feed)
			} else {
				require.NoError(t, err)
				require.NotNil(t, feed)
				assert.Equal(t, test.feedType, feed.FeedType)
				assert.Equal(t, test.feedTitle, feed.Title)
			}
		})
	}
}

func TestParser_ParseURLWithContext(t *testing.T) {
	server, client := mockServerResponse(t, 404, "", 1*time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	fp.Client = client
	_, err := fp.ParseURLWithContext(server.URL, ctx)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestParser_ParseURL_Failure(t *testing.T) {
	server, client := mockServerResponse(t, 404, "", 0)
	fp := gofeed.NewParser()
	fp.Client = client
	feed, err := fp.ParseURL(server.URL)

	assert.NotNil(t, err)
	assert.IsType(t, gofeed.HTTPError{}, err)
	assert.Nil(t, feed)
}

func TestParser_ParseURLWithContextAndBasicAuth(t *testing.T) {
	server, client := mockServerResponse(t, 404, "", 1*time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	fp.AuthConfig = &gofeed.Auth{
		Username: "foo",
		Password: "bar",
	}
	fp.Client = client
	_, err := fp.ParseURLWithContext(server.URL, ctx)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

// to detect race conditions, run with go test -race
func TestParser_Concurrent(t *testing.T) {
	files := []string{"atom03_feed.xml", "atom10_feed.xml", "rss_feed.xml", "rss_feed_bom.xml",
		"rss_feed_leading_spaces.xml", "rdf_feed.xml", "json10_feed.json", "json11_feed.json"}
	fp := gofeed.NewParser()
	fp.AtomTranslator = &gofeed.DefaultAtomTranslator{}
	fp.RSSTranslator = &gofeed.DefaultRSSTranslator{}
	fp.JSONTranslator = &gofeed.DefaultJSONTranslator{}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
			data := testutil.ReadFile(t, "testdata/parser/universal/"+file)
			_, err := fp.ParseString(string(data))
			require.NoError(t, err)
		})
	}
}

// Test Helpers

func mockServerResponse(t *testing.T, code int, body string, delay time.Duration) (*httptest.Server, *http.Client) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-r.Context().Done():
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(code)
		io.WriteString(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	client := &http.Client{Transport: transport}
	t.Cleanup(func() {
		transport.CloseIdleConnections()
		server.Close()
	})
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

// A feed much larger than the detection window must parse completely: the
// format parser reads from the start of the stream, not from after the peek.
func TestParser_Parse_LargeFeed(t *testing.T) {
	var sb strings.Builder
	sb.WriteString(`<rss version="2.0"><channel><title>big</title>`)
	for i := 0; i < 2000; i++ {
		fmt.Fprintf(&sb, `<item><title>item %d</title><guid>g%d</guid></item>`, i, i)
	}
	sb.WriteString(`</channel></rss>`)

	feed, err := gofeed.NewParser().Parse(strings.NewReader(sb.String()))
	require.NoError(t, err)
	require.NotNil(t, feed)
	require.Len(t, feed.Items, 2000)
	require.NotNil(t, feed.Items[1999])
	assert.Equal(t, "item 1999", feed.Items[1999].Title)
}

// JSON feeds larger than the detection window must be classified from their
// prefix and then validated and decoded from the complete stream (issue #344).
func TestParser_Parse_LargeJSONFeed(t *testing.T) {
	content := strings.Repeat("x", 8192)
	body := fmt.Sprintf(
		`{"version":"https://jsonfeed.org/version/1.1","title":"big","items":[{"id":"1","content_text":%q}]}`,
		content,
	)

	feed, err := gofeed.NewParser().Parse(strings.NewReader(body))
	require.NoError(t, err)
	if assert.NotNil(t, feed) && assert.Len(t, feed.Items, 1) {
		require.NotNil(t, feed.Items[0])
		assert.Equal(t, "json", feed.FeedType)
		assert.Equal(t, "big", feed.Title)
		assert.Equal(t, content, feed.Items[0].Content)
	}
}

// Once a large JSON document is detected, syntax errors beyond the detection
// window must come from the JSON parser rather than being masked as a type
// detection failure.
func TestParser_Parse_LargeMalformedJSON(t *testing.T) {
	body := `{"version":"https://jsonfeed.org/version/1.1","items":[{"id":"1","content_text":"` +
		strings.Repeat("x", 8192) + `"}`

	feed, err := gofeed.NewParser().Parse(strings.NewReader(body))
	assert.Nil(t, feed)
	var syntaxError *stdjson.SyntaxError
	assert.ErrorAs(t, err, &syntaxError)
	assert.NotErrorIs(t, err, gofeed.ErrFeedTypeNotDetected)
}

// Detection only inspects the first few KB; a root element pushed beyond
// that window by leading junk is reported as undetected, not misparsed.
func TestParser_Parse_RootBeyondDetectionWindow(t *testing.T) {
	pad := "<!-- " + strings.Repeat("x", 8192) + " -->"
	_, err := gofeed.NewParser().Parse(strings.NewReader(pad + `<rss version="2.0"><channel></channel></rss>`))
	assert.ErrorIs(t, err, gofeed.ErrFeedTypeNotDetected)
}

// An HTML page served in place of a feed — a login wall, error page, or
// bot-protection challenge — should report that it wasn't a feed, while still
// matching ErrFeedTypeNotDetected for callers that check for it (issue #294).
func TestParser_Parse_HTMLNotFeed(t *testing.T) {
	pages := []string{
		"<!DOCTYPE html>\n<html><head><title>Log in</title></head><body>Please sign in</body></html>",
		`<html lang="en"><body>403 Forbidden</body></html>`,
		"\n  \t<HTML>\n<BODY>Checking your browser...</BODY>\n</HTML>",
		"\uFEFF<!doctype html><html></html>",
		"<head><meta charset=\"utf-8\"></head>",
	}
	for _, page := range pages {
		_, err := gofeed.NewParser().Parse(strings.NewReader(page))
		assert.ErrorIs(t, err, gofeed.ErrFeedTypeNotDetected)
		assert.Contains(t, err.Error(), "HTML", "page %q should be reported as HTML", page)
	}
}

// A non-feed payload that isn't recognizably HTML keeps the plain
// ErrFeedTypeNotDetected, with no invented HTML detail. In particular a
// document that merely begins with an XML comment, or an element whose name
// only starts with "html", must not be mistaken for an HTML page (issue #294).
func TestParser_Parse_NonHTMLNotFeed(t *testing.T) {
	inputs := []string{
		"<note><to>Tove</to><from>Jani</from></note>",
		"<!-- just a comment --><data>value</data>",
		"<htmlish>not actually html</htmlish>",
		"plain text, not markup at all",
	}
	for _, in := range inputs {
		_, err := gofeed.NewParser().Parse(strings.NewReader(in))
		assert.ErrorIs(t, err, gofeed.ErrFeedTypeNotDetected)
		assert.NotContains(t, err.Error(), "HTML", "input %q should not be reported as HTML", in)
	}
}

// The real-world case from the issue: a server answers with 2xx but hands back
// an HTML challenge page instead of the feed. ParseURL should surface the
// clearer "not a feed" error rather than a bare detection failure (issue #294).
func TestParser_ParseURL_HTMLBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, "<!DOCTYPE html><html><body>Just a moment...</body></html>")
	}))
	defer srv.Close()

	_, err := gofeed.NewParser().ParseURL(srv.URL)
	assert.ErrorIs(t, err, gofeed.ErrFeedTypeNotDetected)
	assert.Contains(t, err.Error(), "HTML")
}
