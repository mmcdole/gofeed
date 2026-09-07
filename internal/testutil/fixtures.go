// Package testutil contains helpers used only by gofeed's tests.
package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ReadFile reads a test input, failing the current test if it is unavailable.
func ReadFile(t testing.TB, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	return data
}

// RunFixtures compares parsed inputs with their typed JSON expectations.
// XML inputs use a sibling .json file; JSON inputs use _expected.json.
func RunFixtures[T any](t *testing.T, pattern string, parse func(io.Reader) (*T, error)) {
	t.Helper()
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, path := range files {
		ext := filepath.Ext(path)
		stem := strings.TrimSuffix(path, ext)
		if ext == ".json" && strings.HasSuffix(stem, "_expected") {
			continue
		}
		count++
		t.Run(filepath.Base(stem), func(t *testing.T) {
			expectedPath := stem + ".json"
			if ext == ".json" {
				expectedPath = stem + "_expected.json"
			}
			expected := new(T)
			if err := json.Unmarshal(ReadFile(t, expectedPath), expected); err != nil {
				t.Fatalf("decode fixture %s: %v", expectedPath, err)
			}
			actual, err := parse(bytes.NewReader(ReadFile(t, path)))
			if err != nil {
				t.Fatalf("parse fixture %s: %v", path, err)
			}
			assert.Equal(t, expected, actual, "expected output: %s", expectedPath)
		})
	}
	if count == 0 {
		t.Fatalf("no input fixtures matched %q", pattern)
	}
}
