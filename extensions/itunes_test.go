package ext_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
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
		f, _ := ioutil.ReadFile(ff)

		// Parse actual feed
		fp := gofeed.NewParser()
		actual, _ := fp.Parse(bytes.NewReader(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/extensions/itunes/%s.json", name)
		e, _ := ioutil.ReadFile(ef)

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
