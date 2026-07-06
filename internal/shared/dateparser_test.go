package shared

import "testing"

// Named timezone abbreviations must resolve to the right offset, not a silent
// zero offset (issue #237). The result is checked in UTC.
func TestParseDateNamedZones(t *testing.T) {
	cases := []struct {
		in      string
		wantUTC string
	}{
		{"Mon, 02 Jan 2006 15:04:05 EST", "2006-01-02 20:04:05"},   // -5
		{"Mon, 02 Jan 2006 15:04:05 EDT", "2006-01-02 19:04:05"},   // -4
		{"Mon, 02 Jan 2006 15:04:05 CST", "2006-01-02 21:04:05"},   // -6
		{"Mon, 02 Jan 2006 15:04:05 PST", "2006-01-02 23:04:05"},   // -8
		{"Mon, 02 Jan 2006 15:04:05 PDT", "2006-01-02 22:04:05"},   // -7
		{"Mon, 02 Jan 2006 15:04:05 CEST", "2006-01-02 13:04:05"},  // +2
		{"Mon, 02 Jan 2006 15:04:05 GMT", "2006-01-02 15:04:05"},   // 0
		{"Mon, 02 Jan 2006 15:04:05 UTC", "2006-01-02 15:04:05"},   // 0
		{"Mon, 02 Jan 2006 15:04:05 -0700", "2006-01-02 22:04:05"}, // numeric still works
		{"2006-01-02T15:04:05Z", "2006-01-02 15:04:05"},            // RFC3339 still works
	}
	for _, c := range cases {
		got, err := ParseDate(c.in)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", c.in, err)
			continue
		}
		if g := got.UTC().Format("2006-01-02 15:04:05"); g != c.wantUTC {
			t.Errorf("%q -> %s UTC, want %s", c.in, g, c.wantUTC)
		}
	}
}

// A representative sample of the supported layouts must parse to the correct
// instant, across weekday/no-weekday, offset/no-zone, fractional seconds and
// date-only forms. Result checked in UTC (issue #119).
func TestParseDateFormats(t *testing.T) {
	cases := []struct {
		in      string
		wantUTC string
	}{
		{"Mon, 02 Jan 2006 15:04:05 -0700", "2006-01-02 22:04:05"}, // RFC1123Z
		{"Mon, 2 Jan 2006 15:04:05 -0700", "2006-01-02 22:04:05"},  // single-digit day
		{"02 Jan 2006 15:04:05 -0700", "2006-01-02 22:04:05"},      // RFC822Z, no weekday
		{"2006-01-02T15:04:05-07:00", "2006-01-02 22:04:05"},       // RFC3339 offset
		{"2006-01-02T15:04:05.500Z", "2006-01-02 15:04:05"},        // fractional seconds
		{"2006-01-02 15:04:05", "2006-01-02 15:04:05"},             // space separator, no zone
		{"2006-01-02", "2006-01-02 00:00:00"},                      // date only
	}
	for _, c := range cases {
		got, err := ParseDate(c.in)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", c.in, err)
			continue
		}
		if g := got.UTC().Format("2006-01-02 15:04:05"); g != c.wantUTC {
			t.Errorf("%q -> %s UTC, want %s", c.in, g, c.wantUTC)
		}
	}
}
