package rss_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/mmcdole/gofeed/rss"
	"github.com/stretchr/testify/assert"
)

func TestRSSParser_ParseFeed_DetectVersion(t *testing.T) {
	var verTests = []struct {
		file    string
		version string
	}{
		{"simple_rss090.xml", "0.9"},
		{"simple_rss091.xml", "0.91"},
		{"simple_rss092.xml", "0.92"},
		{"simple_rss10.xml", "1.0"},
		{"simple_rss20.xml", "2.0"},
		{"extensions_rss20.xml", "0.0"},
	}

	for _, test := range verTests {
		file := fmt.Sprintf("testdata/%s", test.file)
		f, _ := ioutil.ReadFile(file)
		fp := &rss.Parser{}

		rss, err := fp.ParseFeed(string(f))

		rssJson, _ := json.Marshal(rss)
		fmt.Printf("\n\n%s\n", string(rssJson))

		assert.Nil(t, err, "Failed to parse feed: %s", file)
		assert.Equal(t, test.version, rss.Version, "Expected RSS version %s, got %s", test.version, rss.Version)
	}
}

func TestRSSParser_ParseFeed_ExpectedResults(t *testing.T) {
	var verTests = []struct {
		feedFile     string
		expectedFile string
	}{
		//		{"simple_rss090.xml", "0.9"},
		//		{"simple_rss091.xml", "0.91"},
		//		{"simple_rss092.xml", "0.92"},
		//		{"simple_rss10.xml", "1.0"},
		//		{"simple_rss20.xml", "2.0"},
		{"complete_rss090.xml", "complete_rss090.json"},
		{"complete_rss091.xml", "complete_rss091.json"},
	}

	for _, test := range verTests {
		// Get actual source feed
		ff := fmt.Sprintf("testdata/%s", test.feedFile)
		f, _ := ioutil.ReadFile(ff)

		// Parse actual feed
		fp := &rss.Parser{}
		actual, _ := fp.ParseFeed(string(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("testdata/%s", test.expectedFile)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &rss.Feed{}
		json.Unmarshal(e, &expected)

		assert.Equal(t, actual, expected, "Feed file %s did not match expected output %s", test.feedFile, test.expectedFile)
	}
}

func TestRSSParser_ParseFeed_Extensions(t *testing.T) {
	f, _ := ioutil.ReadFile("testdata/twit.xml")
	fp := &rss.Parser{}

	rss, err := fp.ParseFeed(string(f))

	assert.Nil(t, err, "Failed to parse feed")

	// Channel Extension
	expected := "weekly"
	actual := rss.Extensions["sy"]["updatePeriod"][0].Value
	assert.Equal(t, expected, actual, "Expected extension value %s, got %s", expected, actual)

	// Channel Extension - Attribute
	expected = "http://twit.cachefly.net/coverart/twit/twit1400audio.jpg"
	actual = rss.Extensions["itunes"]["image"][0].Attrs["href"]
	assert.Equal(t, expected, actual, "Expected extension attr value %s, got %s", expected, actual)

	// Channel Extension - Nested
	expected = "Tech News"
	actual = rss.Extensions["itunes"]["category"][0].Children["category"][0].Attrs["text"]
	assert.Equal(t, expected, actual, "Expected nested extension value %s, got %s", expected, actual)

	item := rss.Items[0]

	// Item Extension
	expected = "TWiT"
	actual = item.Extensions["itunes"]["author"][0].Value
	assert.Equal(t, expected, actual, "Expected item extension value %s, got %s", expected, actual)

}
