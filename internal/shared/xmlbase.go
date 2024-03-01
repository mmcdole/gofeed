package shared

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	xpp "github.com/mmcdole/goxpp"
	"golang.org/x/net/html"
)

var (
	// HTML attributes which contain URIs
	// https://pythonhosted.org/feedparser/resolving-relative-links.html
	// To catch every possible URI attribute is non-trivial:
	// https://stackoverflow.com/questions/2725156/complete-list-of-html-tag-attributes-which-have-a-url-value
	htmlURIAttrs = map[string]bool{
		"action":     true,
		"background": true,
		"cite":       true,
		"codebase":   true,
		"data":       true,
		"href":       true,
		"poster":     true,
		"profile":    true,
		"scheme":     true,
		"src":        true,
		"uri":        true,
		"usemap":     true,
	}

	// List of xml attributes that contain URIs to be resolved relative to
	// xml:base
	// From the Atom spec https://tools.ietf.org/html/rfc4287
	uriAttrs = map[string]bool{
		"href":   true,
		"scheme": true,
		"src":    true,
		"uri":    true,
	}
)

// XMLBase.NextTag iterates through the tokens until it reaches a StartTag or
// EndTag. It resolves urls in tag attributes relative to the current xml:base.
//
// NextTag is similar to goxpp's NextTag method except it wont throw an error
// if the next immediate token isnt a Start/EndTag.  Instead, it will continue
// to consume tokens until it hits a Start/EndTag or EndDocument.
func NextTag(p *xpp.XMLPullParser) (event xpp.XMLEventType, err error) {
	for {
		event, err = p.Next()
		if err != nil {
			return event, err
		}

		if event == xpp.EndTag {
			break
		}

		if event == xpp.StartTag {
			if err != nil {
				return
			}

			err = resolveAttrs(p)
			if err != nil {
				return
			}

			break
		}

		if event == xpp.EndDocument {
			return event, fmt.Errorf("Failed to find NextTag before reaching the end of the document.")
		}

	}
	return
}

// resolve relative URI attributes according to xml:base
func resolveAttrs(p *xpp.XMLPullParser) error {
	for i, attr := range p.Attrs {
		lowerName := strings.ToLower(attr.Name.Local)
		if uriAttrs[lowerName] {
			absURL, err := XmlBaseResolveUrl(p.BaseStack.Top(), attr.Value)
			if err != nil {
				return err
			}
			if absURL != nil {
				p.Attrs[i].Value = absURL.String()
			}
		}
	}
	return nil
}

// resolve u relative to b
func XmlBaseResolveUrl(b *url.URL, u string) (*url.URL, error) {
	relURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return relURL, nil
	}

	if b.Path != "" && u != "" && b.Path[len(b.Path)-1] != '/' {
		// There's no reason someone would use a path in xml:base if they
		// didn't mean for it to be a directory
		b.Path = b.Path + "/"
	}
	absURL := b.ResolveReference(relURL)
	return absURL, nil
}

// Transforms html by resolving any relative URIs in attributes
// if an error occurs during parsing or serialization, then the original string
// is returned along with the error.
func ResolveHTML(base *url.URL, relHTML string) (string, error) {
	if base == nil {
		return relHTML, nil
	}

	htmlReader := strings.NewReader(relHTML)

	doc, err := html.Parse(htmlReader)
	if err != nil {
		return relHTML, err
	}

	var visit func(*html.Node)

	// recursively traverse HTML resolving any relative URIs in attributes
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, a := range n.Attr {
				if htmlURIAttrs[a.Key] {
					absVal, err := XmlBaseResolveUrl(base, a.Val)
					if absVal != nil && err == nil {
						n.Attr[i].Val = absVal.String()
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}

	visit(doc)
	var w bytes.Buffer
	err = html.Render(&w, doc)
	if err != nil {
		return relHTML, err
	}

	// html.Render() always writes a complete html5 document, so strip the html
	// and body tags
	absHTML := w.String()
	absHTML = strings.TrimPrefix(absHTML, "<html><head></head><body>")
	absHTML = strings.TrimSuffix(absHTML, "</body></html>")

	return absHTML, err
}
