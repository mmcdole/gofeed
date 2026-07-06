package shared

import (
	"strings"

	xpp "github.com/mmcdole/goxpp/v2"
)

// ForEachChild advances through the children of the element the parser is
// currently on, calling handle for each child start tag with its lowercased
// local name. handle must fully consume the child element, by parsing or
// skipping it. ForEachChild returns once it reaches the enclosing element's
// end tag, leaving the parser on it, so callers can assert it with Expect.
func ForEachChild(p *xpp.Parser, handle func(name string) error) error {
	for {
		tok, err := NextTag(p)
		if err != nil {
			return err
		}
		if tok == xpp.EndTag {
			return nil
		}
		if tok == xpp.StartTag {
			if err := handle(strings.ToLower(p.Name())); err != nil {
				return err
			}
		}
	}
}
