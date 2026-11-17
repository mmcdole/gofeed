package shared

import (
	"strings"
	"testing"
	"time"
)

// Named timezone abbreviations must resolve to the right offset, not a silent
// zero offset (issue #237). The result is checked in UTC.
func TestParseDateNamedZones(t *testing.T) {
	cases := []struct {
		in      string
		wantUTC string
	}{
		{"Mon, 02 Jan 2006 15:04:05 EST", "2006-01-02 20:04:05"},    // -5
		{"Mon,02 January 2006 15:04:05 EST", "2006-01-02 20:04:05"}, // comma without space, full month (issue #308)
		{"Mon, 02 Jan 2006 15:04:05 EDT", "2006-01-02 19:04:05"},    // -4
		{"Mon, 02 Jan 2006 15:04:05 CST", "2006-01-02 21:04:05"},    // -6
		{"Mon, 02 Jan 2006 15:04:05 PST", "2006-01-02 23:04:05"},    // -8
		{"Mon, 02 Jan 2006 15:04:05 PDT", "2006-01-02 22:04:05"},    // -7
		{"Mon, 02 Jan 2006 15:04:05 CEST", "2006-01-02 13:04:05"},   // +2
		{"Mon, 02 Jan 2006 15:04:05 GMT", "2006-01-02 15:04:05"},    // 0
		{"Mon, 02 Jan 2006 15:04:05 UTC", "2006-01-02 15:04:05"},    // 0
		{"Mon, 02 Jan 2006 15:04:05 -0700", "2006-01-02 22:04:05"},  // numeric still works
		{"2006-01-02T15:04:05Z", "2006-01-02 15:04:05"},             // RFC3339 still works
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

// RFC 3339 5.6 allows the "T" separator and "Z" designator in lower case
// (issue #371).
func TestParseDateRFC3339LowerCase(t *testing.T) {
	for _, in := range []string{
		"2006-01-02t15:04:05Z",
		"2006-01-02T15:04:05z",
		"2006-01-02t15:04:05z",
		"2006-01-02t15:04:05.500z",
		"2006-01-02t17:04:05+02:00",
	} {
		got, err := ParseDate(in)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", in, err)
			continue
		}
		if g := got.UTC().Format("2006-01-02 15:04:05"); g != "2006-01-02 15:04:05" {
			t.Errorf("%q -> %s UTC, want 2006-01-02 15:04:05", in, g)
		}
	}
}

// Every layout must be able to parse its own rendering of the reference time.
// A malformed layout (like an hour written as 14) consumes fields twice and
// can never match anything (issue #308).
func TestDateLayoutsRoundTrip(t *testing.T) {
	ref := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*3600))
	all := make([]string, 0, len(dateFormats)+len(dateFormatsWithNamedZone))
	all = append(all, dateFormats...)
	all = append(all, dateFormatsWithNamedZone...)
	for _, layout := range all {
		rendered := ref.Format(layout)
		if _, err := time.Parse(layout, rendered); err != nil {
			t.Errorf("layout %q cannot parse its own rendering %q: %v", layout, rendered, err)
		}
	}
}

// Guards ParseDate's dispatch. UTC covers the layouts whose zone token is a
// literal "Z" or an RFC3339 "Z07:00"; a named zone covers the ones that end in
// a real abbreviation.
func TestNumericLayoutsNeverEndInZone(t *testing.T) {
	ref := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
	for _, zone := range []*time.Location{time.UTC, time.FixedZone("MST", -7*3600)} {
		for _, layout := range dateFormats {
			if s := ref.In(zone).Format(layout); endsInZoneAbbreviation(s) {
				t.Errorf("layout %q renders %q, which ParseDate would match against the named-zone layouts first", layout, s)
			}
		}
	}
}

// A feed writing midnight at the end of a day as hour 24 means the start of
// the next day, not a date gofeed cannot read at all.
func TestParseDateEndOfDayHour(t *testing.T) {
	cases := []struct {
		in      string
		wantUTC string
	}{
		{"Sun, 07 Dec 2025 24:00:00 GMT", "2025-12-08 00:00:00"},
		{"Sun, 07 Dec 2025 24:00:00 +0200", "2025-12-07 22:00:00"},
		{"Sun, 07 Dec 2025 24:00:00 EST", "2025-12-08 05:00:00"},
		{"Sun, 07 Dec 2025 24:00:00", "2025-12-08 00:00:00"},
		{"2025-12-07T24:00:00Z", "2025-12-08 00:00:00"},
		{"2025-12-07t24:00:00z", "2025-12-08 00:00:00"}, // with the lower-case form of issue #371
		{"2025-12-07 24:00:00", "2025-12-08 00:00:00"},
		{"Wed, 31 Dec 2025 24:00:00 GMT", "2026-01-01 00:00:00"}, // over a year
		{"Sat, 28 Feb 2026 24:00:00 GMT", "2026-03-01 00:00:00"}, // over a month
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

// The end-of-day rewrite must never fire on a date that already parses, so no
// feed that works today can change meaning. The reference time puts a 24 in
// every field one can legitimately appear in - the year, the day, the minute
// and the second - leaving only the hour.
func TestNormalizeEndOfDayHourLeavesSupportedLayoutsAlone(t *testing.T) {
	ref := time.Date(2024, time.January, 24, 15, 24, 24, 0, time.FixedZone("MST", -7*3600))
	all := make([]string, 0, len(dateFormats)+len(dateFormatsWithNamedZone))
	all = append(all, dateFormats...)
	all = append(all, dateFormatsWithNamedZone...)
	for _, layout := range all {
		rendered := ref.Format(layout)
		if got, fired := normalizeEndOfDayHour(rendered); fired {
			t.Errorf("layout %q renders %q, which the rewrite changed to %q", layout, rendered, got)
		}
		if !strings.Contains(rendered, "24") {
			t.Errorf("layout %q renders %q, which carries no 24 to guard against", layout, rendered)
		}
	}
}

func TestNormalizeEndOfDayHourIgnoresOtherFields(t *testing.T) {
	for _, in := range []string{
		"2024:01:02",                    // a 24 that ends a longer number
		"Sat, 24 Jan 2026 12:24:24 GMT", // a day, minute and second of 24
		"2026-01-02T15:04:05Z",
		"",
		"24",
		"not a date at all",
	} {
		if got, fired := normalizeEndOfDayHour(in); fired {
			t.Errorf("normalizeEndOfDayHour(%q) = %q, want it left alone", in, got)
		}
	}
}

func BenchmarkParseDate(b *testing.B) {
	cases := []struct {
		name string
		in   string
	}{
		{"NamedZoneGMT", "Mon, 10 Aug 2026 10:00:00 GMT"},
		{"NumericOffset", "Mon, 10 Aug 2026 10:00:00 +0000"},
		{"RFC3339", "2026-08-10T10:00:00Z"},
		{"Unparseable", "not a date at all"},
		// Worst case for the dispatch: the predicate fires, so both lists are
		// exhausted in the reordered order.
		{"UnparseableZoneTail", "Mon, 99 Xyz 2026 99:99:99 QQQ"},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = ParseDate(c.in)
			}
		})
	}
}
