package atom_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/atom"
	"github.com/stretchr/testify/assert"
)

// Tests

func TestParser_Parse(t *testing.T) {
	files, _ := filepath.Glob("../testdata/parser/atom/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/parser/atom/%s.xml", name)
		f, _ := ioutil.ReadFile(ff)

		// Parse actual feed
		fp := &atom.Parser{}
		actual, _ := fp.Parse(bytes.NewReader(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/parser/atom/%s.json", name)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &atom.Feed{}
		json.Unmarshal(e, expected)

		if assert.Equal(t, expected, actual, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

func TestFeed_ToXML(t *testing.T) {
	f, _ := ioutil.ReadFile("../testdata/parser/atom/atom10_feed_with_entry.xml")
	fp := atom.Parser{}
	feed, err := fp.Parse(bytes.NewReader(f))
	assert.NotNil(t, feed)
	assert.Nil(t, err)

	b, err := xml.MarshalIndent(feed, "", "    ")
	assert.Nil(t, err)

	// Verify that we can write out the same fields that we read in
	// assuming that it was in the same order and formatting
	gotXml := string(b)
	wantXml := string(f)
	assert.Equal(t, wantXml, gotXml, "Feed xml did not match expected output atom10_feed_with_entry.xml")
}

// TODO: Examples
