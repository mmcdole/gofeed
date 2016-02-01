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

func (rp *RSSParser) parseExtension(fe FeedExtensions, p *xpp.XMLPullParser) (FeedExtensions, error) {
	prefix := rp.prefixForNamespace(p.Space)

	result, err := rp.parseExtensionElement(p)
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

func (bp *BaseParser) parseExtensionElement(p *xpp.XMLPullParser) (ext Extension, err error) {
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
			child, err := bp.parseExtensionElement(p)
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

func (bp *BaseParser) prefixForNamespace(space string) string {
	// First we check if the global namespace map
	// contains an entry for this namespace/prefix.
	// This way we can use the canonical prefix for this
	// ns instead of the one defined in the feed.
	if prefix, ok := globalNamespaces[space]; ok {
		return prefix
	}

	// Next we check if the feed itself defined this
	// this namespace and return it if we have a result.
	if prefix, ok := bp.feedSpaces[space]; ok {
		return prefix
	}

	// Lastly, any namespace which is not defined in the
	// the feed will be the prefix itself when using Go's
	// xml.Decoder.Token() method.
	return space
}

func (bp *BaseParser) isExtension(p *xpp.XMLPullParser) bool {
	space := strings.TrimSpace(p.Space)
	if prefix, ok := globalNamespaces[space]; ok {
		return prefix != ""
	}

	if prefix, ok := bp.feedSpaces[space]; ok {
		return prefix != ""
	}

	return p.Space != ""
}

func (bp *BaseParser) parseNamespaces(p *xpp.XMLPullParser) {
	for _, attr := range p.Attrs {
		if attr.Name.Space == "xmlns" {
			space := strings.TrimSpace(attr.Value)
			spacePrefix := strings.TrimSpace(strings.ToLower(attr.Name.Local))
			bp.feedSpaces[space] = spacePrefix
		} else if attr.Name.Local == "xmlns" {
			space := strings.TrimSpace(attr.Value)
			bp.feedSpaces[space] = ""
		}
	}
}

func (bp *BaseParser) parseText(p *xpp.XMLPullParser) (text string, err error) {
	text, err = p.NextText()
	if err != nil {
		return text, err
	}

	text = strings.TrimSpace(text)
	return text, nil
}
