package atom_test

import (
	"testing"

	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/internal/testutil"
)

func TestParser_Parse(t *testing.T) {
	testutil.RunFixtures(t, "../testdata/parser/atom/*.xml", (&atom.Parser{}).Parse)
}
