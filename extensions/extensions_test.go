package ext_test

import (
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/internal/testutil"
)

func TestITunes_Extensions(t *testing.T) {
	testutil.RunFixtures(t, "../testdata/extensions/itunes/*.xml", gofeed.NewParser().Parse)
}

func TestMedia_Extensions(t *testing.T) {
	testutil.RunFixtures(t, "../testdata/extensions/media/*.xml", gofeed.NewParser().Parse)
}

func TestDublinCore_Extensions(t *testing.T) {
	testutil.RunFixtures(t, "../testdata/extensions/dublincore/*.xml", gofeed.NewParser().Parse)
}
