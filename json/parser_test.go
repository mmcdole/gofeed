package json_test

import (
	"bytes"
	jsonEncoding "encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/json"
	"github.com/stretchr/testify/assert"
)

// Tests

func TestParser_Parse(t *testing.T) {
	files, _ := filepath.Glob("../testdata/parser/json/*.json")
	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Printf("Testing %s... ", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/parser/json/%s.json", name)
		f, _ := ioutil.ReadFile(ff)

		// Parse actual feed
		fp := &json.Parser{}
		actual, _ := fp.Parse(bytes.NewReader(f))

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/parser/json/%s_result.json", name)
		e, _ := ioutil.ReadFile(ef)

		// Unmarshal expected feed
		expected := &json.Feed{}
		jsonEncoding.Unmarshal(e, expected)

		if assert.Equal(t, expected, actual, "Feed file %s.json did not match expected output %s_expected.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

// TODO: Examples
