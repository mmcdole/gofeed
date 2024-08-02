package ext_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hartza-capital/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestITunes_Extensions(t *testing.T) {
	files, _ := filepath.Glob("../testdata/extensions/itunes/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/extensions/itunes/%s.xml", name)
		f, _ := os.ReadFile(ff)

		// Parse actual feed
		fp := gofeed.NewParser(gofeed.DefaultUserAgent)
		actual, _ := fp.Parse(bytes.NewReader(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/extensions/itunes/%s.json", name)
		e, _ := os.ReadFile(ef)

		// Unmarshal expected feed
		expected := &gofeed.Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

func TestMedia_Extensions(t *testing.T) {
	files, _ := filepath.Glob("../testdata/extensions/media/*.xml")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/extensions/media/%s.xml", name)
		f, _ := os.ReadFile(ff)

		// Parse actual feed
		fp := gofeed.NewParser(gofeed.DefaultUserAgent)
		actual, _ := fp.Parse(bytes.NewReader(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/extensions/media/%s.json", name)
		e, _ := os.ReadFile(ef)

		// Unmarshal expected feed
		expected := &gofeed.Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, expected, actual, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}
