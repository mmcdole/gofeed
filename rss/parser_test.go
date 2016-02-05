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
		{"complete_rss090.xml", "0.9"},
		{"complete_rss091.xml", "0.91"},
		{"complete_rss092.xml", "0.92"},
		{"complete_rss10.xml", "1.0"},
		{"complete_rss20.xml", "2.0"},
	}

	for _, test := range verTests {
		file := fmt.Sprintf("../testdata/%s", test.file)
		f, _ := ioutil.ReadFile(file)
		fp := &rss.Parser{}

		rss, err := fp.ParseFeed(string(f))

		//rssJson, _ := json.Marshal(rss)
		//fmt.Printf("\n\n%s\n", string(rssJson))

		assert.Nil(t, err, "Failed to parse feed: %s", file)
		assert.Equal(t, test.version, rss.Version, "Expected RSS version %s, got %s", test.version, rss.Version)
	}
}

func TestRSSParser_ParseFeed_ExpectedResults(t *testing.T) {

	var verTests = []struct {
		file string
	}{
		{"complete_rss090"},
		{"complete_rss091"},
		{"complete_rss092"},
		{"complete_rss10"},
		{"complete_rss20"},
	}

	for _, test := range verTests {
		// Get actual source feed
		ff := fmt.Sprintf("../testdata/%s.xml", test.file)
		f, _ := ioutil.ReadFile(ff)

		// Parse actual feed
		fp := &rss.Parser{}
		actual, _ := fp.ParseFeed(string(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/%s.json", test.file)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &rss.Feed{}
		json.Unmarshal(e, &expected)

		assert.Equal(t, actual, expected, "Feed file %s.xml did not match expected output %s.json", test.file, test.file)
	}
}
