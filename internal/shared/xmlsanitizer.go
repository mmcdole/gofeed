package shared

import (
	"io"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// NewXMLSanitizerReader creates an io.Reader that
// wraps another io.Reader and removes illegal xml
// characters from the io stream.
func NewXMLSanitizerReader(xml io.Reader) io.Reader {
	isIllegal := runes.Predicate(func() func(rune) bool {
		return func(r rune) bool {
			return !(r == 0x09 ||
				r == 0x0A ||
				r == 0x0D ||
				r >= 0x20 && r <= 0xDF77 ||
				r >= 0xE000 && r <= 0xFFFD ||
				r >= 0x10000 && r <= 0x10FFFF)
		}
	}())
	t := transform.Chain(runes.Remove(isIllegal))
	return transform.NewReader(xml, t)
}
