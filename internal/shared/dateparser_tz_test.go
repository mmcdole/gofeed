package shared

import (
	"testing"
	"time"
)

func TestParseDateRFC822Zones(t *testing.T) {
	tests := []struct {
		input          string
		expectedOffset int // seconds east of UTC
		expectedUTCHr  int // expected hour after converting to UTC
	}{
		{"Mon, 21 Apr 2025 06:00:00 EDT", -4 * 3600, 10},
		{"Mon, 21 Apr 2025 06:00:00 CDT", -5 * 3600, 11},
		{"Mon, 21 Apr 2025 06:00:00 MDT", -6 * 3600, 12},
		{"Mon, 21 Apr 2025 06:00:00 PDT", -7 * 3600, 13},
		{"Mon, 21 Apr 2025 06:00:00 EST", -5 * 3600, 11},
		{"Mon, 21 Apr 2025 06:00:00 CST", -6 * 3600, 12},
		{"Mon, 21 Apr 2025 06:00:00 MST", -7 * 3600, 13},
		{"Mon, 21 Apr 2025 06:00:00 PST", -8 * 3600, 14},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parsed, err := ParseDate(tt.input)
			if err != nil {
				t.Fatalf("ParseDate(%q) returned error: %v", tt.input, err)
			}
			_, offset := parsed.Zone()
			if offset != tt.expectedOffset {
				t.Errorf("offset: got %d, want %d", offset, tt.expectedOffset)
			}
			if parsed.Hour() != 6 {
				t.Errorf("local hour: got %d, want 6", parsed.Hour())
			}
			if parsed.UTC().Hour() != tt.expectedUTCHr {
				t.Errorf("UTC hour: got %d, want %d", parsed.UTC().Hour(), tt.expectedUTCHr)
			}
		})
	}
}

func TestParseDateNumericOffsetUnchanged(t *testing.T) {
	// Numeric offsets must not be affected by the RFC 822 zone fix.
	input := "Mon, 21 Apr 2025 06:00:00 -0400"
	parsed, err := ParseDate(input)
	if err != nil {
		t.Fatalf("ParseDate(%q) returned error: %v", input, err)
	}
	_, offset := parsed.Zone()
	if offset != -4*3600 {
		t.Errorf("offset: got %d, want %d", offset, -4*3600)
	}
}

func TestParseDateGMTUnchanged(t *testing.T) {
	// GMT should still parse as offset 0.
	input := "Mon, 21 Apr 2025 06:00:00 GMT"
	parsed, err := ParseDate(input)
	if err != nil {
		t.Fatalf("ParseDate(%q) returned error: %v", input, err)
	}
	if !parsed.UTC().Equal(time.Date(2025, 4, 21, 6, 0, 0, 0, time.UTC)) {
		t.Errorf("expected 06:00 UTC, got %s", parsed.UTC())
	}
}
