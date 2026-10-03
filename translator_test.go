package gofeed_test

import (
	jsonEncoding "encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/internal/testutil"
	"github.com/mmcdole/gofeed/json"
	"github.com/mmcdole/gofeed/rss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultRSSTranslator_Translate(t *testing.T) {
	testutil.RunFixtures(t, "testdata/translator/rss/*.xml", func(r io.Reader) (*gofeed.Feed, error) {
		source, err := (&rss.Parser{}).Parse(r)
		if err != nil {
			return nil, err
		}
		return (&gofeed.DefaultRSSTranslator{}).Translate(source)
	})
}

func TestDefaultRSSTranslator_Translate_WrongType(t *testing.T) {
	translator := &gofeed.DefaultRSSTranslator{}
	af, err := translator.Translate("wrong type")
	assert.Nil(t, af)
	assert.NotNil(t, err)
}

func TestDefaultAtomTranslator_Translate(t *testing.T) {
	testutil.RunFixtures(t, "testdata/translator/atom/*.xml", func(r io.Reader) (*gofeed.Feed, error) {
		source, err := (&atom.Parser{}).Parse(r)
		if err != nil {
			return nil, err
		}
		return (&gofeed.DefaultAtomTranslator{}).Translate(source)
	})
}

func TestDefaultAtomTranslator_Translate_WrongType(t *testing.T) {
	translator := &gofeed.DefaultAtomTranslator{}
	af, err := translator.Translate("wrong type")
	assert.Nil(t, af)
	assert.NotNil(t, err)
}

func TestDefaultJSONTranslator_Translate(t *testing.T) {
	testutil.RunFixtures(t, "testdata/translator/json/*.json", func(r io.Reader) (*gofeed.Feed, error) {
		source, err := (&json.Parser{}).Parse(r)
		if err != nil {
			return nil, err
		}
		return (&gofeed.DefaultJSONTranslator{}).Translate(source)
	})
}

func TestDefaultJSONTranslator_Translate_WrongType(t *testing.T) {
	translator := &gofeed.DefaultJSONTranslator{}
	af, err := translator.Translate("wrong type")
	assert.Nil(t, af)
	assert.NotNil(t, err)
}

// A JSON Feed attachment's size_in_bytes maps to the universal Enclosure.Length
// (bytes), not its duration.
func TestJSONAttachmentEnclosureLength(t *testing.T) {
	in := `{"version":"https://jsonfeed.org/version/1","title":"t","items":[{"id":"a","attachments":[{"url":"u","mime_type":"audio/mpeg","size_in_bytes":5000000,"duration_in_seconds":3600}]}]}`
	feed, err := gofeed.NewParser().ParseString(in)
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, feed)
	require.Len(t, feed.Items, 1)
	require.NotNil(t, feed.Items[0])
	require.Len(t, feed.Items[0].Enclosures, 1)
	require.NotNil(t, feed.Items[0].Enclosures[0])
	enc := feed.Items[0].Enclosures[0]
	if enc.Length != "5000000" {
		t.Fatalf("enclosure length = %q, want \"5000000\" (bytes, not duration)", enc.Length)
	}
}

// DisableContentImageScan turns off the HTML-parsing fallback that finds a
// first <img> in feed and item content; explicit images are unaffected.
func TestDisableContentImageScan(t *testing.T) {
	feed := `<rss version="2.0"><channel>
		<description><![CDATA[<p><img src="http://example.org/feed.png"/></p>]]></description>
		<item><description><![CDATA[<img src="http://example.org/item.png">]]></description></item>
	</channel></rss>`

	fp := &rss.Parser{}
	rssFeed, err := fp.Parse(strings.NewReader(feed))
	require.NoError(t, err)

	def := &gofeed.DefaultRSSTranslator{}
	out, err := def.Translate(rssFeed)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Len(t, out.Items, 1)
	require.NotNil(t, out.Items[0])
	if assert.NotNil(t, out.Image) {
		assert.Equal(t, "http://example.org/feed.png", out.Image.URL)
	}
	if assert.NotNil(t, out.Items[0].Image) {
		assert.Equal(t, "http://example.org/item.png", out.Items[0].Image.URL)
	}

	off := &gofeed.DefaultRSSTranslator{DisableContentImageScan: true}
	out, err = off.Translate(rssFeed)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Len(t, out.Items, 1)
	require.NotNil(t, out.Items[0])
	assert.Nil(t, out.Image)
	assert.Nil(t, out.Items[0].Image)
}

func TestPersonURLJSON(t *testing.T) {
	for _, tt := range []struct {
		name   string
		person gofeed.Person
		want   string
	}{
		{name: "present", person: gofeed.Person{Name: "Author", URL: "https://example.org/author"}, want: `{"name":"Author","url":"https://example.org/author"}`},
		{name: "absent", person: gofeed.Person{Name: "Author"}, want: `{"name":"Author"}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonEncoding.Marshal(tt.person)
			if err != nil {
				t.Fatal(err)
			}
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

// JSON Feed items without their own authors inherit the feed's top-level
// authors (JSON Feed 1.1) or author (JSON Feed 1.0).
func TestJSONFeedItemAuthorFallback(t *testing.T) {
	t.Run("JSON Feed 1.1 top-level authors fallback", func(t *testing.T) {
		in := `{"version":"https://jsonfeed.org/version/1.1","title":"t",
 "authors":[{"name":"Ann","url":"https://example.com/ann"}],
 "items":[{"id":"1","content_text":"x","authors":[{"name":"Bob"}]},
          {"id":"2","content_text":"x"},
          {"id":"3","content_text":"x","authors":[]}]}`

		feed, err := gofeed.NewParser().ParseString(in)
		require.NoError(t, err)
		require.NotNil(t, feed)
		require.Len(t, feed.Authors, 1)
		assert.Equal(t, "Ann", feed.Authors[0].Name)

		require.Len(t, feed.Items, 3)

		// Item 1 specifies authors, should not inherit
		require.Len(t, feed.Items[0].Authors, 1)
		assert.Equal(t, "Bob", feed.Items[0].Authors[0].Name)

		// Item 2 specifies no authors, should inherit feed authors
		require.Len(t, feed.Items[1].Authors, 1)
		assert.Equal(t, "Ann", feed.Items[1].Authors[0].Name)
		assert.Equal(t, "https://example.com/ann", feed.Items[1].Authors[0].URL)

		// Item 3 explicitly specifies empty authors, should not inherit
		assert.Empty(t, feed.Items[2].Authors)
	})

	t.Run("JSON Feed 1.0 top-level author fallback", func(t *testing.T) {
		in := `{"version":"https://jsonfeed.org/version/1","title":"t",
 "author":{"name":"Ann","url":"https://example.com/ann"},
 "items":[{"id":"1","content_text":"x","author":{"name":"Bob"}},
          {"id":"2","content_text":"x"}]}`

		feed, err := gofeed.NewParser().ParseString(in)
		require.NoError(t, err)
		require.NotNil(t, feed)
		require.NotNil(t, feed.Author)
		assert.Equal(t, "Ann", feed.Author.Name)
		require.Len(t, feed.Authors, 1)
		assert.Equal(t, "Ann", feed.Authors[0].Name)

		require.Len(t, feed.Items, 2)

		// Item 1 specifies author, should not inherit
		require.NotNil(t, feed.Items[0].Author)
		assert.Equal(t, "Bob", feed.Items[0].Author.Name)
		require.Len(t, feed.Items[0].Authors, 1)
		assert.Equal(t, "Bob", feed.Items[0].Authors[0].Name)

		// Item 2 specifies no author, should inherit feed author
		require.NotNil(t, feed.Items[1].Author)
		assert.Equal(t, "Ann", feed.Items[1].Author.Name)
		assert.Equal(t, "https://example.com/ann", feed.Items[1].Author.URL)
		require.Len(t, feed.Items[1].Authors, 1)
		assert.Equal(t, "Ann", feed.Items[1].Authors[0].Name)
	})

	t.Run("Feed without authors leaves item authors nil", func(t *testing.T) {
		in := `{"version":"https://jsonfeed.org/version/1.1","title":"t",
 "items":[{"id":"1","content_text":"x"}]}`

		feed, err := gofeed.NewParser().ParseString(in)
		require.NoError(t, err)
		require.NotNil(t, feed)
		assert.Nil(t, feed.Author)
		assert.Nil(t, feed.Authors)

		require.Len(t, feed.Items, 1)
		assert.Nil(t, feed.Items[0].Author)
		assert.Nil(t, feed.Items[0].Authors)
	})
}
