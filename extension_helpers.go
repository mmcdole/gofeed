package gofeed

import (
	ext "github.com/mmcdole/gofeed/extensions"
)

// CustomNamespace is the pseudo namespace under which non-namespaced unknown
// elements are filed in Extensions, preserving nesting, attributes and
// repetition.
const CustomNamespace = "_custom"

// GetExtension returns the extension elements for the given namespace prefix
// and element name, or nil when absent. Non-namespaced custom elements live
// under CustomNamespace.
func (f *Feed) GetExtension(namespace, element string) []ext.Extension {
	return extensionsFor(f.Extensions, namespace, element)
}

// GetExtension returns the extension elements for the given namespace prefix
// and element name, or nil when absent. Non-namespaced custom elements live
// under CustomNamespace.
func (i *Item) GetExtension(namespace, element string) []ext.Extension {
	return extensionsFor(i.Extensions, namespace, element)
}

// GetExtensionValue returns the text of the first matching extension
// element, or "" when absent.
func (f *Feed) GetExtensionValue(namespace, element string) string {
	return firstExtensionValue(f.Extensions, namespace, element)
}

// GetExtensionValue returns the text of the first matching extension
// element, or "" when absent.
func (i *Item) GetExtensionValue(namespace, element string) string {
	return firstExtensionValue(i.Extensions, namespace, element)
}

// GetCustomValue returns the text of the first non-namespaced custom element
// with the given name, or "" when absent. Unlike the flat Custom map, the
// backing extension tree also holds nested elements, attributes and repeated
// values; use GetExtension(CustomNamespace, element) for those.
func (f *Feed) GetCustomValue(element string) string {
	return firstExtensionValue(f.Extensions, CustomNamespace, element)
}

// GetCustomValue returns the text of the first non-namespaced custom element
// with the given name, or "" when absent. Unlike the flat Custom map, the
// backing extension tree also holds nested elements, attributes and repeated
// values; use GetExtension(CustomNamespace, element) for those.
func (i *Item) GetCustomValue(element string) string {
	return firstExtensionValue(i.Extensions, CustomNamespace, element)
}

func extensionsFor(exts ext.Extensions, namespace, element string) []ext.Extension {
	if exts == nil {
		return nil
	}
	nsMap, ok := exts[namespace]
	if !ok {
		return nil
	}
	return nsMap[element]
}

func firstExtensionValue(exts ext.Extensions, namespace, element string) string {
	matches := extensionsFor(exts, namespace, element)
	if len(matches) == 0 {
		return ""
	}
	return matches[0].Value
}
