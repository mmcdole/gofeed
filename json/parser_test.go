package json_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	jsonParser "github.com/mmcdole/gofeed/json"
	"github.com/stretchr/testify/assert"
)

// Tests

// TODO: add tests for invalid

func TestParser_Parse(t *testing.T) {
	files, _ := filepath.Glob("../testdata/parser/json/*.json")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		if strings.HasSuffix(name, "expected") {
			continue
		}

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/parser/json/%s.json", name)
		f, _ := os.ReadFile(ff)

		// Parse actual feed
		fp := &jsonParser.Parser{}
		actual, _ := fp.Parse(bytes.NewReader(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/parser/json/%s_expected.json", name)
		e, _ := os.ReadFile(ef)

		// Unmarshal expected feed
		expected := &jsonParser.Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.json did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

// TODO: Remove redundant tests
func TestParser_ParseInvalidAndStruct(t *testing.T) {
	name := "invalid"
	fmt.Printf("Testing %s... ", name)

	// Get actual source feed
	ff := fmt.Sprintf("../testdata/parser/json/invalid/%s.json", name)
	fmt.Println(ff)
	f, _ := os.ReadFile(ff)

	// Parse actual feed
	fp := &jsonParser.Parser{}
	_, err := fp.Parse(bytes.NewReader(f))
	assert.Contains(t, err.Error(), "expect }")

	name = "version_json_10"
	fmt.Printf("Testing %s... ", name)

	// Get actual source feed
	ff = fmt.Sprintf("../testdata/parser/json/%s.json", name)
	fmt.Println(ff)
	f, _ = os.ReadFile(ff)

	// Parse actual feed
	actual, _ := fp.Parse(bytes.NewReader(f))

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

// TODO: Examples
