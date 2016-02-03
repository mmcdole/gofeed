package parseutil

import (
	"strings"

	. "github.com/mmcdole/gofeed/feed"
	"github.com/mmcdole/goxpp"
)

func ParseExtension(fe FeedExtensions, p *xpp.XMLPullParser) (FeedExtensions, error) {
	prefix := PrefixForNamespace(p.Space, p)

	result, err := parseExtensionElement(p)
	if err != nil {
		return nil, err
	}

	// Ensure the extension prefix map exists
	if _, ok := fe[prefix]; !ok {
		fe[prefix] = map[string][]Extension{}
	}
	// Ensure the extension element slice exists
	if _, ok := fe[prefix][p.Name]; !ok {
		fe[prefix][p.Name] = []Extension{}
	}

	fe[prefix][p.Name] = append(fe[prefix][p.Name], result)
	return fe, nil
}

func parseExtensionElement(p *xpp.XMLPullParser) (ext Extension, err error) {
	if err = p.Expect(xpp.StartTag, "*"); err != nil {
		return ext, err
	}

	ext.Name = p.Name
	ext.Children = map[string][]Extension{}
	ext.Attrs = map[string]string{}

	for _, attr := range p.Attrs {
		// TODO: Alright that we are stripping
		// namespace information from attributes ?
		ext.Attrs[attr.Name.Local] = attr.Value
	}

	for {
		tok, err := p.Next()
		if err != nil {
			return ext, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			child, err := parseExtensionElement(p)
			if err != nil {
				return ext, err
			}

			if _, ok := ext.Children[child.Name]; !ok {
				ext.Children[child.Name] = []Extension{}
			}

			ext.Children[child.Name] = append(ext.Children[child.Name], child)
		} else if tok == xpp.Text {
			ext.Value = strings.TrimSpace(p.Text)
		}
	}

	if err = p.Expect(xpp.EndTag, ext.Name); err != nil {
		return ext, err
	}

	return ext, nil
}
