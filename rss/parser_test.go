package rss_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/rss"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	files, _ := filepath.Glob("../testdata/parser/rss/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/parser/rss/%s.xml", name)
		f, _ := os.ReadFile(ff)

		// Parse actual feed
		fp := &rss.Parser{}
		actual, _ := fp.Parse(bytes.NewReader(f), gofeed.NewParser().BuildRSSExtParsers())

		// the `Parsed` part of extensions is not correctly unmarshalled from JSON
		// workaround: move the actual extensions through a round of json marshalling so that we get the same
		for _, i := range actual.Items {
			if len(i.Extensions) > 0 {
				b, _ := json.Marshal(i.Extensions)
				json.Unmarshal(b, &i.Extensions)
			}
		}

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/parser/rss/%s.json", name)
		e, _ := os.ReadFile(ef)

		// Unmarshal expected feed
		expected := &rss.Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

// TODO: Examples
