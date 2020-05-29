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
		{"skip & normal & amps", "skip & normal & amps"},
		{"not & entity;hello &ne xt;one", "not & entity;hello &ne xt;one"},

		{"&lt;foo&gt;", "<foo>"},
		{"a &quot;b&quot; &apos;c&apos;", "a \"b\" 'c'"},
		{"foo &amp;&amp; bar", "foo && bar"},

		{"&#34;foo&#34;", "\"foo\""},
		{"&#x61;&#x062;&#x0063;", "abc"},
		{"r&#xe9;sum&#x00E9;", "résumé"},
		{"r&eacute;sum&eacute;", "résumé"},
		{"&", "&"},
		{"&foo", "&foo"},
		{"&lt", "&lt"},
		{"&#", "&#"},
	}

	for _, test := range tests {
		res, err := DecodeEntities(test.str)
		assert.Nil(t, err, "cannot decode %q", test.str)
		assert.Equal(t, res, test.res,
			"%q was decoded to %q instead of %q",
			test.str, res, test.res)
	}
}

func TestStripCDATA(t *testing.T) {
	tests := []struct {
		str string
		res string
	}{
		{"<![CDATA[ test ]]>test", " test test"},
		{"<![CDATA[test &]]> &lt;", "test & <"},
		{"", ""},
		{"test", "test"},
		{"]]>", "]]>"},
		{"<![CDATA[", "<![CDATA["},
		{"<![CDATA[testtest", "<![CDATA[testtest"},
		{`<![CDATA[
    Since this is a CDATA section
    I can use all sorts of reserved characters
    like > < " and &
    or write things like
    <foo></bar>
    but my document is still well formed!
]]>`, `
    Since this is a CDATA section
    I can use all sorts of reserved characters
    like > < " and &
    or write things like
    <foo></bar>
    but my document is still well formed!
`},
		{`<![CDATA[
Within this Character Data block I can
use double dashes as much as I want (along with <, &, ', and ")
*and* %MyParamEntity; will be expanded to the text
"Has been expanded" ... however, I can't use
the CEND sequence. If I need to use CEND I must escape one of the
brackets or the greater-than sign using concatenated CDATA sections.
]]>`, `
Within this Character Data block I can
use double dashes as much as I want (along with <, &, ', and ")
*and* %MyParamEntity; will be expanded to the text
"Has been expanded" ... however, I can't use
the CEND sequence. If I need to use CEND I must escape one of the
brackets or the greater-than sign using concatenated CDATA sections.
`},
		// 		{`<![CDATA[ test ]]><!--
		// Within this comment I can use ]]>
		// and other reserved characters like <
		// &, ', and ", but %MyParamEntity; will not be expanded
		// (if I retrieve the text of this node it will contain
		// %MyParamEntity; and not "Has been expanded")
		// and I can't place two dashes next to each other.
		// -->`, ` test <!--
		// Within this comment I can use ]]>
		// and other reserved characters like <
		// &, ', and ", but %MyParamEntity; will not be expanded
		// (if I retrieve the text of this node it will contain
		// %MyParamEntity; and not "Has been expanded")
		// and I can't place two dashes next to each other.
		// -->`,
		// 		},
		{`<![CDATA[ test ]]><!-- test -->`, ` test <!-- test -->`}, // TODO: probably wrong
		{`An example of escaped CENDs`, `An example of escaped CENDs`},
		{`<![CDATA[This text contains a CEND ]]]]><![CDATA[>]]>`, `This text contains a CEND ]]>`},
		{`<![CDATA[This text contains a CEND ]]]><![CDATA[]>]]>`, `This text contains a CEND ]]>`},
	}

	for _, test := range tests {
		res := StripCDATA(test.str)
		assert.Equal(t, test.res, res)
	}
}
