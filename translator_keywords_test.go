package gofeed

import (
	"reflect"
	"testing"
)

// Regression test for #367: iTunes keywords must be trimmed and empty entries
// dropped, matching how a plain <category> contributes a trimmed value.
func TestSplitKeywords(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"news, politics, tech", []string{"news", "politics", "tech"}},
		{"a,,c", []string{"a", "c"}},
		{"  spaced  ,  out  ", []string{"spaced", "out"}},
		{"single", []string{"single"}},
		{",,", []string{}},
	}
	for _, c := range cases {
		got := splitKeywords(c.in)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("splitKeywords(%q) = %#v, want %#v", c.in, got, c.want)
		}
	}
}
