package feed

import (
	"strings"

	"github.com/mmcdole/go-xpp"
)

type BaseParser struct {
	// Map of all namespaces (url / prefix)
	// that have been defined in the feed.
	feedSpaces map[string]string
}

func (bp *BaseParser) parseExtension(p *xpp.XMLPullParser) (ext Extension, err error) {
	if err = p.Expect(xpp.StartTag, "*"); err != nil {
		return ext, err
	}

	ext.Name = p.Name
	ext.Attrs = map[string]string{}
	ext.Children = map[string][]Extension{}

	for _, attr := range p.Attrs {
		// TODO: Alright that we are stripping
		// namespace information from attributes
		// for the usecase of feed parsing?
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
			child, err := bp.parseExtension(p)
			if err != nil {
				return ext, err
			}

			if _, ok := ext.Children[child.Name]; !ok {
				ext.Children[child.Name] = []Extension{}
			}

			ext.Children[child.Name] = append(ext.Children[child.Name], child)
		} else if tok == xpp.Text {
			ext.Value = p.Text
		}
	}

	if err = p.Expect(xpp.EndTag, ext.Name); err != nil {
		return ext, err
	}

	return ext, nil
}

func (bp *BaseParser) prefixForNamespace(space string) string {
	lspace := strings.ToLower(space)

	// First we check if the global namespace map
	// contains an entry for this namespace/prefix.
	// This way we can use the canonical prefix for this
	// ns instead of the one defined in the feed.
	if prefix, ok := globalNamespaces[lspace]; ok {
		return prefix
	}

	// Next we check if the feed itself defined this
	// this namespace and return it if we have a result.
	if prefix, ok := bp.feedSpaces[lspace]; ok {
		return prefix
	}

	// Lastly, any namespace which is not defined in the
	// the feed will be the prefix itself when using Go's
	// xml.Decoder.Token() method.
	return space
}

func (bp *BaseParser) parseNamespaces(p *xpp.XMLPullParser) {
	for _, attr := range p.Attrs {
		if attr.Name.Space == "xmlns" {
			spacePrefix := strings.ToLower(attr.Name.Local)
			bp.feedSpaces[attr.Value] = spacePrefix
		}
	}
}
