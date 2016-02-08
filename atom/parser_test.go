package atom_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/mmcdole/gofeed/atom"
	"github.com/stretchr/testify/assert"
)

func TestRSSParser_ParseFeed_DetectVersion(t *testing.T) {
	var verTests = []struct {
		file    string
		version string
	}{
		{"complete_atom10.xml", "1.0"},
	}
	for _, test := range verTests {
		file := fmt.Sprintf("../testdata/%s", test.file)
		f, _ := ioutil.ReadFile(file)
		fp := &atom.Parser{}

		atom, err := fp.ParseFeed(string(f))

		atomJson, _ := json.Marshal(atom)
		fmt.Printf("\n\n%s\n", string(atomJson))

		assert.Nil(t, err, "Failed to parse feed: %s", file)
		assert.Equal(t, test.version, atom.Version, "Expected RSS version %s, got %s", test.version, atom.Version)
	}
}
