package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	const rssFeed = `<rss version="2.0"><channel><title>Example</title></channel></rss>`
	const atomFeed = `<feed xmlns="http://www.w3.org/2005/Atom"><title>Example</title></feed>`
	for _, tt := range []struct {
		name     string
		args     []string
		input    string
		feedType string
	}{
		{"default", nil, rssFeed, "rss"},
		{"long flag", []string{"--type", "rss"}, rssFeed, ""},
		{"short flag", []string{"-t", "atom"}, atomFeed, ""},
		{"equals", []string{"--type=atom"}, atomFeed, ""},
		{"rss alias", []string{"-t", "R"}, rssFeed, ""},
		{"atom alias", []string{"-t", "A"}, atomFeed, ""},
		{"universal", []string{"-t", "universal"}, atomFeed, "atom"},
		{"unknown type falls back", []string{"-t", "other"}, rssFeed, "rss"},
		{"end of flags", []string{"--"}, rssFeed, "rss"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "feed.xml")
			if err := os.WriteFile(path, []byte(tt.input), 0600); err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			if err := run(append(tt.args, path), &out); err != nil {
				t.Fatal(err)
			}
			var got struct {
				Title    string `json:"title"`
				FeedType string `json:"feedType"`
			}
			if err := json.Unmarshal(out.Bytes(), &got); err != nil {
				t.Fatalf("invalid output: %v\n%s", err, out.String())
			}
			if got.Title != "Example" || got.FeedType != tt.feedType {
				t.Fatalf("got title=%q feedType=%q, want Example and %q", got.Title, got.FeedType, tt.feedType)
			}
		})
	}
}

func TestRunHelp(t *testing.T) {
	for _, arg := range []string{"-h", "--help", "help"} {
		t.Run(arg, func(t *testing.T) {
			var out bytes.Buffer
			if err := run([]string{arg}, &out); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), "Usage: ftest") {
				t.Fatalf("missing help: %q", out.String())
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
	}{
		{"missing path", nil},
		{"unknown flag", []string{"--unknown"}},
		{"missing flag value", []string{"--type"}},
		{"missing file", []string{filepath.Join(t.TempDir(), "missing.xml")}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args, io.Discard); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestRunURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `<rss version="2.0"><channel><title>Remote</title></channel></rss>`)
	}))
	defer server.Close()
	var out bytes.Buffer
	if err := run([]string{server.URL}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"Remote"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
