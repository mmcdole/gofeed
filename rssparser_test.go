package feed_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/mmcdole/go-feed"
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
	}

	for _, test := range verTests {
		file := fmt.Sprintf("test/%s", test.file)
		f, _ := ioutil.ReadFile(file)
		fp := &feed.RSSParser{}

		rss, err := fp.ParseFeed(string(f))

		assert.Nil(t, err, "Failed to parse feed: %s", file)
		assert.Equal(t, test.version, rss.Version, "Expected RSS version %s, got %s", test.version, rss.Version)
	}
}

func TestRSSParser_ParseFeed_Extensions(t *testing.T) {
	f, _ := ioutil.ReadFile("test/twit.xml")
	fp := &feed.RSSParser{}

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
