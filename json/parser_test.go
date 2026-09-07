package json_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/mmcdole/gofeed/internal/testutil"
	jsonParser "github.com/mmcdole/gofeed/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	testutil.RunFixtures(t, "../testdata/parser/json/*.json", (&jsonParser.Parser{}).Parse)
}

func TestParser_ParseInvalid(t *testing.T) {
	data := testutil.ReadFile(t, "../testdata/parser/json/invalid/invalid.json")
	_, err := (&jsonParser.Parser{}).Parse(bytes.NewReader(data))
	assert.Error(t, err)
}

func TestParser_ParseStruct(t *testing.T) {
	data := testutil.ReadFile(t, "../testdata/parser/json/version_json_10.json")
	actual, err := (&jsonParser.Parser{}).Parse(bytes.NewReader(data))
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.NotNil(t, actual.Author)
	require.Len(t, actual.Items, 1)
	require.NotNil(t, actual.Items[0])
	require.NotNil(t, actual.Items[0].Author)
	require.Len(t, actual.Items[0].Tags, 2)
	require.NotNil(t, actual.Items[0].Attachments)
	require.Len(t, *actual.Items[0].Attachments, 1)

	assert.Equal(t, "1.0", actual.Version)
	assert.Equal(t, "title", actual.Title)
	assert.Equal(t, "https://sample-json-feed.com", actual.HomePageURL)
	assert.Equal(t, "https://sample-json-feed.com/feed.json", actual.FeedURL)
	assert.Equal(t, "description", actual.Description)
	assert.Equal(t, "user_comment", actual.UserComment)
	assert.Equal(t, "https://sample-json-feed.com/feed.json?next=500", actual.NextURL)
	assert.Equal(t, "https://sample-json-feed.com/icon.png", actual.Icon)
	assert.Equal(t, "https://sample-json-feed.com/favicon.png", actual.Favicon)
	assert.Equal(t, "author_name", actual.Author.Name)
	assert.Equal(t, "https://sample-feed-author.com", actual.Author.URL)
	assert.Equal(t, "https://sample-feed-author.com/me.png", actual.Author.Avatar)
	assert.Equal(t, false, actual.Expired)
	assert.Equal(t, "id", actual.Items[0].ID)
	assert.Equal(t, "https://sample-json-feed.com/id", actual.Items[0].URL)
	assert.Equal(t, "https://sample-json-feed.com/external", actual.Items[0].ExternalURL)
	assert.Equal(t, "title", actual.Items[0].Title)
	assert.Contains(t, actual.Items[0].ContentHTML, "content_html")
	assert.Equal(t, "content_text", actual.Items[0].ContentText)
	assert.Equal(t, "summary", actual.Items[0].Summary)
	assert.Equal(t, "https://sample-json-feed.com/image.png", actual.Items[0].Image)
	assert.Equal(t, "https://sample-json-feed.com/banner_image.png", actual.Items[0].BannerImage)
	assert.Equal(t, "2019-10-12T07:20:50.52Z", actual.Items[0].DatePublished)
	assert.Equal(t, "2019-10-12T07:20:50.52Z", actual.Items[0].DateModified)
	assert.Equal(t, "author_name", actual.Items[0].Author.Name)
	assert.Equal(t, "https://sample-feed-author.com", actual.Items[0].Author.URL)
	assert.Equal(t, "https://sample-feed-author.com/me.png", actual.Items[0].Author.Avatar)
	assert.Equal(t, "tag1", actual.Items[0].Tags[0])
	assert.Equal(t, "tag2", actual.Items[0].Tags[1])
	assert.Equal(t, "https://sample-json-feed.com/attachment", (*actual.Items[0].Attachments)[0].URL)
	assert.Equal(t, "audio/mpeg", (*actual.Items[0].Attachments)[0].MimeType)
	assert.Equal(t, "title", (*actual.Items[0].Attachments)[0].Title)
	assert.Equal(t, int64(100), (*actual.Items[0].Attachments)[0].SizeInBytes)
	assert.Equal(t, int64(100), (*actual.Items[0].Attachments)[0].DurationInSeconds)

	assert.Contains(t, actual.String(), "https://sample-json-feed.com/attachment")
}

// An I/O error from the reader must surface as itself, not as a misleading
// JSON syntax error from a truncated buffer (issue #311).
func TestParser_Parse_ReaderError(t *testing.T) {
	boom := errors.New("boom")
	r := io.MultiReader(strings.NewReader(`{"version":"https://jsonfeed.org/version/1"`), iotest.ErrReader(boom))

	_, err := (&jsonParser.Parser{}).Parse(r)
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}
}
