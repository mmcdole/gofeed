package gofeed_test

import (
	jsonEncoding "encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/v2"
	"github.com/mmcdole/gofeed/v2/atom"
	"github.com/mmcdole/gofeed/v2/json"
	"github.com/mmcdole/gofeed/v2/rss"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRSSTranslator_Translate(t *testing.T) {
	files, _ := filepath.Glob("testdata/translator/rss/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("testdata/translator/rss/%s.xml", name)
		f, _ := os.Open(ff)
		defer f.Close()

		// Parse actual feed
		translator := &gofeed.DefaultRSSTranslator{}
		fp := rss.NewParser()
		rssFeed, _ := fp.Parse(f, nil)
		actual, _ := translator.Translate(rssFeed, nil)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("testdata/translator/rss/%s.json", name)
		e, _ := os.ReadFile(ef)

		// Unmarshal expected feed
		expected := &gofeed.Feed{}
		jsonEncoding.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

func TestDefaultRSSTranslator_Translate_WrongType(t *testing.T) {
	translator := &gofeed.DefaultRSSTranslator{}
	af, err := translator.Translate("wrong type", nil)
	assert.Nil(t, af)
	assert.NotNil(t, err)
}

func TestDefaultAtomTranslator_Translate(t *testing.T) {
	files, _ := filepath.Glob("testdata/translator/atom/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("testdata/translator/atom/%s.xml", name)
		f, _ := os.Open(ff)
		defer f.Close()

		// Parse actual feed
		translator := &gofeed.DefaultAtomTranslator{}
		fp := atom.NewParser()
		atomFeed, _ := fp.Parse(f, nil)
		actual, _ := translator.Translate(atomFeed, nil)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("testdata/translator/atom/%s.json", name)
		e, _ := os.ReadFile(ef)

		// Unmarshal expected feed
		expected := &gofeed.Feed{}
		jsonEncoding.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

func TestDefaultAtomTranslator_Translate_WrongType(t *testing.T) {
	translator := &gofeed.DefaultAtomTranslator{}
	af, err := translator.Translate("wrong type", nil)
	assert.Nil(t, af)
	assert.NotNil(t, err)
}

func TestDefaultJSONTranslator_Translate(t *testing.T) {
	files, _ := filepath.Glob("testdata/translator/json/*.json")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		if strings.HasSuffix(name, "expected") {
			continue
		}

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("testdata/translator/json/%s.json", name)
		f, _ := os.Open(ff)
		defer f.Close()

		// Parse actual feed
		translator := &gofeed.DefaultJSONTranslator{}
		fp := json.NewParser()
		jsonFeed, _ := fp.Parse(f, nil)
		actual, _ := translator.Translate(jsonFeed, nil)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("testdata/translator/json/%s_expected.json", name)
		e, _ := os.ReadFile(ef)

		// Unmarshal expected feed
		expected := &gofeed.Feed{}
		jsonEncoding.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.json did not match expected output %s_expected.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

/*

func TestDefaultJSONTranslator_Translate(t *testing.T) {
	name := "sample"
	fmt.Printf("Testing %s... ", name)

	// Get actual source feed
	ff := fmt.Sprintf("testdata/translator/json/%s.json", name)
	fmt.Println(ff)
	f, _ := ioutil.ReadFile(ff)

	// Parse actual feed
	translator := &gofeed.DefaultJSONTranslator{}
	fp := json.NewParser()
	feed, _ := fp.Parse(bytes.NewReader(f), nil)
	actual, _ := translator.Translate(feed, nil)

	assert.Equal(t, "title", actual.Title)
	assert.Equal(t, "description", actual.Description)
	assert.Equal(t, "https://sample-json-feed.com", actual.Link)
	assert.Equal(t, "https://sample-json-feed.com/feed.json", actual.FeedLink)
	assert.Equal(t, "2019-10-12T07:20:50.52Z", actual.Updated)
	assert.Equal(t, "2019-10-12T07:20:50Z", actual.UpdatedParsed.Format(time.RFC3339))
	assert.Equal(t, "2019-10-12T07:20:50.52Z", actual.Published)
	assert.Equal(t, "2019-10-12T07:20:50Z", actual.PublishedParsed.Format(time.RFC3339))
	assert.Equal(t, "author_name", actual.Author.Name)
	assert.Equal(t, "", actual.Author.Email)
	assert.Equal(t, "", actual.Language)
	assert.Equal(t, "https://sample-json-feed.com/icon.png", actual.Image.URL)
	assert.Equal(t, "", actual.Image.Title)
	assert.Equal(t, "", actual.Copyright)
	assert.Equal(t, "", actual.Generator)
	assert.Equal(t, 0, len(actual.Categories))
	assert.Equal(t, (*ext.DublinCoreExtension)(nil), actual.DublinCoreExt)
	assert.Equal(t, (*ext.ITunesFeedExtension)(nil), actual.ITunesExt)
	assert.Equal(t, ext.Extensions(nil), actual.Extensions)
	assert.Equal(t, "json", actual.FeedType)
	assert.Equal(t, "1.0", actual.FeedVersion)
	assert.Equal(t, "title", actual.Items[0].Title)
	assert.Equal(t, "summary", actual.Items[0].Description)
	assert.Equal(t, "<p>content_html</p>", actual.Items[0].Content)
	assert.Equal(t, "https://sample-json-feed.com/id", actual.Items[0].Link)
	assert.Equal(t, "2019-10-12T07:20:50.52Z", actual.Items[0].Updated)
	assert.Equal(t, "2019-10-12T07:20:50Z", actual.Items[0].UpdatedParsed.Format(time.RFC3339))
	assert.Equal(t, "2019-10-12T07:20:50.52Z", actual.Items[0].Published)
	assert.Equal(t, "2019-10-12T07:20:50Z", actual.Items[0].PublishedParsed.Format(time.RFC3339))
	assert.Equal(t, "author_name", actual.Items[0].Author.Name)
	assert.Equal(t, "", actual.Items[0].Author.Email)
	assert.Equal(t, "id", actual.Items[0].GUID)
	assert.Equal(t, "https://sample-json-feed.com/image.png", actual.Items[0].Image.URL)
	assert.Equal(t, "", actual.Items[0].Image.Title)
	assert.Equal(t, "tag1", actual.Items[0].Categories[0])
	assert.Equal(t, "tag2", actual.Items[0].Categories[1])
	assert.Equal(t, "https://sample-json-feed.com/attachment", (actual.Items[0].Enclosures)[0].URL)
	assert.Equal(t, "100", (actual.Items[0].Enclosures)[0].Length)
	assert.Equal(t, "audio/mpeg", (actual.Items[0].Enclosures)[0].Type)
	assert.Equal(t, (*ext.DublinCoreExtension)(nil), actual.Items[0].DublinCoreExt)
	assert.Equal(t, (*ext.ITunesItemExtension)(nil), actual.Items[0].ITunesExt)
	assert.Equal(t, ext.Extensions(nil), actual.Items[0].Extensions)

	name = "sample2"
	fmt.Printf("Testing %s... ", name)

	// Get actual source feed
	ff = fmt.Sprintf("testdata/translator/json/%s.json", name)
	fmt.Println(ff)
	f, _ = ioutil.ReadFile(ff)

	// Parse actual feed
	feed, _ = fp.Parse(bytes.NewReader(f), nil)
	actual, _ = translator.Translate(feed, nil)

	assert.Equal(t, "content_text", actual.Items[0].Content)
	assert.Equal(t, "https://sample-json-feed.com/banner_image.png", actual.Items[0].Image.URL)

}
*/

func TestDefaultJSONTranslator_Translate_WrongType(t *testing.T) {
	translator := &gofeed.DefaultJSONTranslator{}
	af, err := translator.Translate("wrong type", nil)
	assert.Nil(t, af)
	assert.NotNil(t, err)
}
