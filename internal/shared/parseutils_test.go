package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeEntities(t *testing.T) {
	tests := []struct {
		str string
		res string
	}{
		{"", ""},
		{"foo", "foo"},

		{"&lt;foo&gt;", "<foo>"},
		{"a &quot;b&quot; &apos;c&apos;", "a \"b\" 'c'"},
		{"foo &amp;&amp; bar", "foo && bar"},

		{"&#34;foo&#34;", "\"foo\""},
		{"&#x61;&#x062;&#x0063;", "abc"},
		{"r&#xe9;sum&#x00E9;", "résumé"},
	}

	for _, test := range tests {
		res, err := DecodeEntities(test.str)
		assert.Nil(t, err, "cannot decode %q", test.str)
		assert.Equal(t, res, test.res,
			"%q was decoded to %q instead of %q",
			test.str, res, test.res)
	}
}

func TestDecodeEntitiesInvalid(t *testing.T) {
	tests := []string{
		// Predefined entities
		"&",     // truncated
		"&foo",  // truncated
		"&foo;", // unknown
		"&lt",   // known but truncated

		// Numerical character references
		"&#",      // truncated
		"&#;",     // missing number
		"&#x;",    // missing hexadecimal number
		"&#12a;",  // invalid decimal number
		"&#xfoo;", // invalid hexadecimal number
	}

	for _, test := range tests {
		res, err := DecodeEntities(test)
		assert.NotNil(t, err, "%q was decoded to %q", test, res)
	}
}
