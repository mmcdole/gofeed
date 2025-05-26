package shared

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

// FindFirstImgSrc finds the first <img> tag with a src attribute in the HTML document
// and returns the src value. Returns empty string if no img with src is found.
func FindFirstImgSrc(document string) string {
	doc, err := html.Parse(strings.NewReader(document))
	if err != nil {
		return ""
	}

	var findImg func(*html.Node) string
	findImg = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					return attr.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if src := findImg(c); src != "" {
				return src
			}
		}
		return ""
	}

	return findImg(doc)
}

// StripWrappingDiv removes a single wrapping <div> element from HTML content
// if the body contains only that div (and optional whitespace).
// Returns the inner content of the div, or the original content if conditions aren't met.
func StripWrappingDiv(content string) string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return content
	}

	// Find body node
	var findBody func(*html.Node) *html.Node
	findBody = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "body" {
			return n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if b := findBody(c); b != nil {
				return b
			}
		}
		return nil
	}

	body := findBody(doc)
	if body == nil {
		return content
	}

	// Find first non-whitespace child
	var firstElement *html.Node
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			firstElement = c
			break
		}
		if c.Type == html.TextNode && strings.TrimSpace(c.Data) != "" {
			// Non-whitespace text node before any element
			return content
		}
	}

	if firstElement == nil || firstElement.Data != "div" {
		return content
	}

	// Check no element siblings after the div
	for c := firstElement.NextSibling; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			return content // Has element siblings
		}
		if c.Type == html.TextNode && strings.TrimSpace(c.Data) != "" {
			return content // Has non-whitespace text siblings
		}
	}

	// Return inner HTML of the div
	var buf bytes.Buffer
	for c := firstElement.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&buf, c)
	}
	return buf.String()
}