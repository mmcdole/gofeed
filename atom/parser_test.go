package atom_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/v2/atom"
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
		f, _ := os.ReadFile(ff)

		// Parse actual feed
		fp := &atom.Parser{}
		actual, _ := fp.Parse(bytes.NewReader(f), nil)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/parser/atom/%s.json", name)
		e, _ := os.ReadFile(ef)

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

// TODO: Examples
