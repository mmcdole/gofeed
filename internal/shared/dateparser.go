package shared

import (
	"fmt"
	"strings"
	"time"
)

// DateFormats taken from github.com/mjibson/goread
var dateFormats = []string{
	// Numeric-offset and zoneless formats only. Named-zone formats (those
	// ending in "MST") must NOT be here: time.Parse resolves an unknown zone
	// abbreviation to a zero offset, so a named-zone date matched here returns
	// silently shifted. Those layouts live in dateFormatsWithNamedZone instead.
	time.RFC822Z, // RSS
	time.RFC3339, // Atom
	time.RubyDate,
	time.RFC1123Z,
	time.ANSIC,
	"Mon, January 2 2006 15:04:05 -0700",
	"Mon, Jan 2 2006 15:04:05 -700",
	"Mon, Jan 2 2006 15:04:05 -0700",
	"Mon Jan 2 15:04 2006",
	"Mon Jan 02, 2006 3:04 pm",
	"Mon Jan 02 2006 15:04:05 -0700",
	"Mon Jan 02 2006 15:04:05 GMT-0700 (MST)",
	"Monday, January 2, 2006 03:04 PM",
	"Monday, January 2, 2006",
	"Monday, January 02, 2006",
	"Monday, 2 January 2006 15:04:05 -0700",
	"Monday, 2 Jan 2006 15:04:05 -0700",
	"Monday, 02 January 2006 15:04:05 -0700",
	"Monday, 02 January 2006 15:04:05",
	"Mon, 2 January 2006, 15:04 -0700",
	"Mon, 2 January 2006 15:04:05 -0700",
	"Mon, 2 January 2006",
	"Mon, 2 Jan 2006 3:04:05 PM -0700",
	"Mon, 2 Jan 2006 15:4:5 -0700 GMT",
	"Mon, 2, Jan 2006 15:4",
	"Mon, 2 Jan 2006, 15:04 -0700",
	"Mon, 2 Jan 2006 15:04 -0700",
	"Mon, 2 Jan 2006 15:04:05 UT",
	"Mon, 2 Jan 2006 15:04:05 -0700 MST",
	"Mon, 2 Jan 2006 15:04:05-0700",
	"Mon, 2 Jan 2006 15:04:05-07:00",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05",
	"Mon, 2 Jan 2006 15:04",
	"Mon,2 Jan 2006",
	"Mon, 2 Jan 2006",
	"Mon, 2 Jan 06 15:04:05 -0700",
	"Mon, 2006-01-02 15:04",
	"Mon, 02 January 2006",
	"Mon, 02 Jan 2006 15 -0700",
	"Mon, 02 Jan 2006 15:04 -0700",
	"Mon, 02 Jan 2006 15:04:05 Z",
	"Mon, 02 Jan 2006 15:04:05 UT",
	"Mon, 02 Jan 2006 15:04:05 MST-07:00",
	"Mon, 02 Jan 2006 15:04:05 MST -0700",
	"Mon, 02 Jan 2006 15:04:05 GMT-0700",
	"Mon,02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -07:00",
	"Mon, 02 Jan 2006 15:04:05 --0700",
	"Mon 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -07",
	"Mon, 02 Jan 2006 15:04:05 00",
	"Mon, 02 Jan 2006 15:04:05",
	"Mon, 02 Jan 2006 15:4:5 Z",
	"Mon, 02 Jan 2006",
	"January 2, 2006 3:04 PM",
	"January 2, 2006, 3:04 p.m.",
	"January 2, 2006 15:04:05",
	"January 2, 2006 03:04 PM",
	"January 2, 2006",
	"January 02, 2006 15:04",
	"January 02, 2006 03:04 PM",
	"January 02, 2006",
	"Jan 2, 2006 3:04:05 PM",
	"Jan 2, 2006",
	"Jan 02 2006 03:04:05PM",
	"Jan 02, 2006",
	"Jan 2 2006 15:04:05",
	"6/1/2 15:04",
	"6-1-2 15:04",
	"2 January 2006 15:04:05 -0700",
	"2 January 2006",
	"2 Jan 2006 15:04:05 Z",
	"2 Jan 2006 15:04:05 -0700",
	"2 Jan 2006",
	"2.1.2006 15:04:05",
	"2/1/2006",
	"2-1-2006",
	"2006 January 02",
	"2006-1-2T15:04:05Z",
	"2006-1-2 15:04:05",
	"2006-1-2",
	"2006-1-02T15:04:05Z",
	"2006-01-02T15:04Z",
	"2006-01-02T15:04-07:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05-07:00:00",
	"2006-01-02T15:04:05:-0700",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05 -0700",
	"2006-01-02T15:04:05:00",
	"2006-01-02T15:04:05",
	"2006-01-02 at 15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05-0700",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04",
	"2006-01-02 00:00:00.0 15:04:05.0 -0700",
	"2006/01/02",
	"2006-01-02",
	"15:04 02.01.2006 -0700",
	"1/2/2006 3:04:05 PM",
	"1/2/2006",
	"06/1/2 15:04",
	"06-1-2 15:04",
	"02 Monday, Jan 2006 15:04",
	"02 Jan 2006 15:04:05 UT",
	"02 Jan 2006 15:04:05 -0700",
	"02 Jan 2006 15:04:05",
	"02 Jan 2006",
	"02.01.2006 15:04:05",
	"02/01/2006 15:04:05",
	"02.01.2006 15:04",
	"02/01/2006 - 15:04",
	"02.01.2006 -0700",
	"02/01/2006",
	"02-01-2006",
	"01/02/2006 3:04 PM",
	"01/02/2006 - 15:04",
	"01/02/2006",
	"01-02-2006",
}

// Named zone cannot be consistently loaded, so handle separately
var dateFormatsWithNamedZone = []string{
	time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC850,   // "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC822,   // "02 Jan 06 15:04 MST"
	time.UnixDate, // "Mon Jan _2 15:04:05 MST 2006"
	"Mon, January 02, 2006, 15:04:05 MST",
	"Mon, January 02, 2006 15:04:05 MST",
	"Mon, Jan 2, 2006 15:04 MST",
	"Mon, Jan 2 2006 15:04 MST",
	"Mon, Jan 2, 2006 15:04:05 MST",
	"Mon Jan 2 15:04:05 2006 MST",
	"Mon, Jan 02,2006 15:04:05 MST",
	"Monday, January 2, 2006 15:04:05 MST",
	"Monday, 2 January 2006 15:04:05 MST",
	"Monday, 2 Jan 2006 15:04:05 MST",
	"Monday, 02 January 2006 15:04:05 MST",
	"Mon, 2 January 2006 15:04 MST",
	"Mon, 2 January 2006, 15:04:05 MST",
	"Mon, 2 January 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:4:5 MST",
	"Mon, 2 Jan 2006 15:04 MST",
	"Mon, 2 Jan 2006 15:04:05MST",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon 2 Jan 2006 15:04:05 MST",
	"mon,2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 15:04:05 MST",
	"Mon, 2 Jan 06 15:04:05 MST",
	"Mon,02 January 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 3:04:05 PM MST",
	"Mon,02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 2006, 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05MST",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon , 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 06 15:04:05 MST",
	"January 2, 2006 15:04:05 MST",
	"January 02, 2006 15:04:05 MST",
	"Jan 2, 2006 3:04:05 PM MST",
	"Jan 2, 2006 15:04:05 MST",
	"2 January 2006 15:04:05 MST",
	"2 Jan 2006 15:04:05 MST",
	"2006-01-02 15:04:05 MST",
	"1/2/2006 3:04:05 PM MST",
	"1/2/2006 15:04:05 MST",
	"02 Jan 2006 15:04 MST",
	"02 Jan 2006 15:04:05 MST",
	"02/01/2006 15:04 MST",
	"02-01-2006 15:04:05 MST",
	"01/02/2006 15:04:05 MST",
}

// timezoneAbbreviations maps the zone abbreviations seen in feed dates to their
// offset in seconds. Go's time package only knows UTC and the host's local
// zone, so it parses any other abbreviation as a zero offset (and its handling
// of the ones it does know depends on where the code runs). We resolve them
// ourselves for a deterministic result. A few abbreviations are genuinely
// ambiguous (CST, BST); we use the North American / Western European reading
// that dominates real-world feeds.
var timezoneAbbreviations = map[string]int{
	"UT":   0,
	"GMT":  0,
	"UTC":  0,
	"EST":  -5 * 3600,
	"EDT":  -4 * 3600,
	"CST":  -6 * 3600,
	"CDT":  -5 * 3600,
	"MST":  -7 * 3600,
	"MDT":  -6 * 3600,
	"PST":  -8 * 3600,
	"PDT":  -7 * 3600,
	"WET":  0,
	"WEST": 1 * 3600,
	"BST":  1 * 3600,
	"CET":  1 * 3600,
	"CEST": 2 * 3600,
	"EET":  2 * 3600,
	"EEST": 3 * 3600,
}

// endsInZoneAbbreviation reports whether d ends in a bare uppercase timezone
// abbreviation, in a shape that no layout in dateFormats can match.
//
// Three to five characters follows time.parseTimeZone's internal rule, and it also
// excludes every trailing alphabetic token dateFormats layouts end in ("UT"/"Z"/"PM").
// The only overlap left is a zone name after a numeric offset ("-0700 MST" and its
// "GMT" sibling),hence the check on the field in front.
func endsInZoneAbbreviation(d string) bool {
	i := len(d)
	for i > 0 && d[i-1] >= 'A' && d[i-1] <= 'Z' {
		i--
	}
	// i == 0 means the whole string is letters, so it carries no date at all.
	if n := len(d) - i; i == 0 || n < 3 || n > 5 {
		return false
	}

	head := strings.TrimRight(d[:i], " ")
	if j := strings.LastIndexByte(head, ' '); j >= 0 {
		head = head[j+1:]
	}
	return !strings.HasPrefix(head, "+") && !strings.HasPrefix(head, "-")
}

// ParseDate parses a given date string using a large
// list of commonly found feed date formats.
func ParseDate(ds string) (time.Time, error) {
	d := strings.TrimSpace(ds)
	if d == "" {
		return time.Time{}, fmt.Errorf("date string is empty")
	}
	// A date ending in a bare zone abbreviation cannot match dateFormats, so
	// look at the named-zone layouts first. Only the order changes: both lists
	// are still tried, so a wrong guess costs time and nothing else.
	namedZoneFirst := endsInZoneAbbreviation(d)
	if namedZoneFirst {
		if t, ok := parseNamedZone(d); ok {
			return t, nil
		}
	}
	for _, f := range dateFormats {
		if t, err := time.Parse(f, d); err == nil {
			return t, nil
		}
	}
	if !namedZoneFirst {
		if t, ok := parseNamedZone(d); ok {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse date: %s", ds)
}

func parseNamedZone(d string) (time.Time, bool) {
	for _, f := range dateFormatsWithNamedZone {
		t, err := time.Parse(f, d)
		if err != nil {
			continue
		}

		// time.Parse gives an unknown zone abbreviation a zero offset, and its
		// handling of known ones is host-dependent. If we recognise the
		// abbreviation, rebuild the time with our own offset so the result is
		// correct and deterministic. An unrecognised abbreviation is left as
		// parsed (best effort).
		if name, offset := t.Zone(); name != "" {
			if want, ok := timezoneAbbreviations[name]; ok && want != offset {
				t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
					t.Minute(), t.Second(), t.Nanosecond(), time.FixedZone(name, want))
			}
		}
		return t, true
	}

	return time.Time{}, false
}
