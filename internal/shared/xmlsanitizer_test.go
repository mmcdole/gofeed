package shared

import (
	"io"
	"strings"
	"testing"
)

func TestControlCharFilterReader(t *testing.T) {
	// Illegal C0 controls (0x00, 0x08) are dropped; tab, LF, CR and normal
	// text are kept.
	in := "ab\x08cd\x00ef\tgh\nij\rk"
	got, err := io.ReadAll(NewControlCharFilterReader(strings.NewReader(in)))
	if err != nil {
		t.Fatal(err)
	}
	want := "abcdef\tgh\nij\rk"
	if string(got) != want {
		t.Errorf("got %q, want %q", string(got), want)
	}
}
