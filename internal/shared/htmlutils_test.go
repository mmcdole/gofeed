package shared

import (
	"testing"
)

func TestFindFirstImgSrc(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "simple img tag",
			html:     `<img src="test.jpg">`,
			expected: "test.jpg",
		},
		{
			name:     "img tag with other attributes",
			html:     `<img alt="test" src="image.png" width="100">`,
			expected: "image.png",
		},
		{
			name:     "nested img tag",
			html:     `<div><p><img src="nested.gif"></p></div>`,
			expected: "nested.gif",
		},
		{
			name:     "multiple img tags",
			html:     `<img src="first.jpg"><img src="second.jpg">`,
			expected: "first.jpg",
		},
		{
			name:     "img without src",
			html:     `<img alt="no source">`,
			expected: "",
		},
		{
			name:     "no img tag",
			html:     `<div><p>No images here</p></div>`,
			expected: "",
		},
		{
			name:     "empty html",
			html:     ``,
			expected: "",
		},
		{
			name:     "invalid html",
			html:     `<div><img src="test.jpg"`,
			expected: "",
		},
		{
			name:     "img in html with body",
			html:     `<html><body><img src="body-image.jpg"></body></html>`,
			expected: "body-image.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindFirstImgSrc(tt.html)
			if result != tt.expected {
				t.Errorf("FindFirstImgSrc() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStripWrappingDiv(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "simple wrapping div",
			html:     `<div>content</div>`,
			expected: `content`,
		},
		{
			name:     "div with nested content",
			html:     `<div><p>paragraph</p><span>span</span></div>`,
			expected: `<p>paragraph</p><span>span</span>`,
		},
		{
			name:     "div with attributes",
			html:     `<div class="wrapper">content</div>`,
			expected: `content`,
		},
		{
			name:     "multiple root elements",
			html:     `<div>first</div><div>second</div>`,
			expected: `<div>first</div><div>second</div>`,
		},
		{
			name:     "non-div root element",
			html:     `<p>paragraph</p>`,
			expected: `<p>paragraph</p>`,
		},
		{
			name:     "div with text sibling",
			html:     `<div>content</div>extra text`,
			expected: `<div>content</div>extra text`,
		},
		{
			name:     "empty content",
			html:     ``,
			expected: ``,
		},
		{
			name:     "only whitespace",
			html:     `   `,
			expected: `   `,
		},
		{
			name:     "div with only whitespace siblings",
			html:     `  <div>content</div>  `,
			expected: `content`,
		},
		{
			name:     "full html document with single div in body",
			html:     `<html><body><div>content</div></body></html>`,
			expected: `content`,
		},
		{
			name:     "full html document with multiple elements",
			html:     `<html><body><div>first</div><p>second</p></body></html>`,
			expected: `<html><body><div>first</div><p>second</p></body></html>`,
		},
		{
			name:     "nested divs",
			html:     `<div><div>nested</div></div>`,
			expected: `<div>nested</div>`,
		},
		{
			name:     "invalid html - unclosed div gets auto-closed",
			html:     `<div>unclosed`,
			expected: `unclosed`, // HTML parser auto-closes the div
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripWrappingDiv(tt.html)
			if result != tt.expected {
				t.Errorf("StripWrappingDiv() = %q, want %q", result, tt.expected)
			}
		})
	}
}