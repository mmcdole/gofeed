package rss_test

import (
	"strings"
	"testing"

	"github.com/mmcdole/gofeed/internal/testutil"
	"github.com/mmcdole/gofeed/rss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	testutil.RunFixtures(t, "../testdata/parser/rss/*.xml", (&rss.Parser{}).Parse)
}

func TestParser_Parse_CloudTruncated(t *testing.T) {
	// A stream that ends inside <cloud> must surface the underlying XML
	// error from the skip to the matching end tag.
	feed := `<rss version="2.0"><channel><cloud domain="x"><child>`
	fp := &rss.Parser{}
	_, err := fp.Parse(strings.NewReader(feed))
	assert.Error(t, err)
}

func TestParser_Parse_KnownForeignRootNamespace(t *testing.T) {
	// Only an unrecognized root default namespace is tolerated as RSS
	// core; a recognized extension namespace keeps its meaning, so this
	// channel stays foreign and is skipped.
	feed := `<rss version="2.0" xmlns="http://www.w3.org/2005/Atom"><channel><title>x</title></channel></rss>`
	f, err := (&rss.Parser{}).Parse(strings.NewReader(feed))
	require.NoError(t, err)
	require.NotNil(t, f)
	assert.Empty(t, f.Title)
}

func TestParser_Parse_UnknownRoot(t *testing.T) {
	// Neither <rss> nor <rdf>: the parse must fail rather than return an
	// empty feed.
	_, err := (&rss.Parser{}).Parse(strings.NewReader(`<foo><channel/></foo>`))
	assert.Error(t, err)
}
