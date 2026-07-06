package json

import (
	"strings"
	"testing"
)

func parseFeed(t *testing.T, s string) *Feed {
	t.Helper()
	f, err := (&Parser{}).Parse(strings.NewReader(s))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	return f
}

// A numeric id must be coerced to a string rather than failing the feed.
func TestNumericIDCoerced(t *testing.T) {
	f := parseFeed(t, `{"version":"https://jsonfeed.org/version/1.1","title":"t","items":[{"id":123,"content_text":"x"}]}`)
	if len(f.Items) != 1 || f.Items[0].ID != "123" {
		t.Fatalf("id = %q, want \"123\"", f.Items[0].ID)
	}
}

func TestStringIDPreserved(t *testing.T) {
	f := parseFeed(t, `{"version":"https://jsonfeed.org/version/1","title":"t","items":[{"id":"https://x/1"}]}`)
	if f.Items[0].ID != "https://x/1" {
		t.Fatalf("id = %q", f.Items[0].ID)
	}
}

func TestStringExpiredCoerced(t *testing.T) {
	f := parseFeed(t, `{"version":"https://jsonfeed.org/version/1","title":"t","expired":"true","items":[{"id":"a"}]}`)
	if !f.Expired {
		t.Fatal("expired string \"true\" should coerce to true")
	}
}

func TestFloatSizeCoerced(t *testing.T) {
	f := parseFeed(t, `{"version":"https://jsonfeed.org/version/1","title":"t","items":[{"id":"a","attachments":[{"url":"u","size_in_bytes":5000000.0,"duration_in_seconds":3600}]}]}`)
	att := (*f.Items[0].Attachments)[0]
	if att.SizeInBytes != 5000000 {
		t.Fatalf("size = %d, want 5000000", att.SizeInBytes)
	}
	if att.DurationInSeconds != 3600 {
		t.Fatalf("duration = %d, want 3600", att.DurationInSeconds)
	}
}

// Values outside the int64 range (and NaN) must coerce to 0, the same
// fallback as other unusable values, not implementation-defined garbage.
func TestOutOfRangeSizeCoercedToZero(t *testing.T) {
	f := parseFeed(t, `{"version":"https://jsonfeed.org/version/1","title":"t","items":[{"id":"a","attachments":[{"url":"u","size_in_bytes":1e300,"duration_in_seconds":-1e300}]}]}`)
	att := (*f.Items[0].Attachments)[0]
	if att.SizeInBytes != 0 {
		t.Fatalf("size = %d, want 0", att.SizeInBytes)
	}
	if att.DurationInSeconds != 0 {
		t.Fatalf("duration = %d, want 0", att.DurationInSeconds)
	}
}

func TestOutOfRangeStringSizeCoercedToZero(t *testing.T) {
	f := parseFeed(t, `{"version":"https://jsonfeed.org/version/1","title":"t","items":[{"id":"a","attachments":[{"url":"u","size_in_bytes":"1e300","duration_in_seconds":"NaN"}]}]}`)
	att := (*f.Items[0].Attachments)[0]
	if att.SizeInBytes != 0 {
		t.Fatalf("size = %d, want 0", att.SizeInBytes)
	}
	if att.DurationInSeconds != 0 {
		t.Fatalf("duration = %d, want 0", att.DurationInSeconds)
	}
}

// The int64 boundaries themselves stay exact through the integer path.
func TestInt64BoundarySizes(t *testing.T) {
	f := parseFeed(t, `{"version":"https://jsonfeed.org/version/1","title":"t","items":[{"id":"a","attachments":[{"url":"u","size_in_bytes":9223372036854775807,"duration_in_seconds":-9223372036854775808}]}]}`)
	att := (*f.Items[0].Attachments)[0]
	if att.SizeInBytes != 9223372036854775807 {
		t.Fatalf("size = %d, want MaxInt64", att.SizeInBytes)
	}
	if att.DurationInSeconds != -9223372036854775808 {
		t.Fatalf("duration = %d, want MinInt64", att.DurationInSeconds)
	}
}
