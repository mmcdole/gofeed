package gofeed

import (
	ext "github.com/mmcdole/gofeed/extensions"
)

// GetExtension retrieves extension values by namespace and element name.
// Returns a slice of Extension structs for the given namespace and element.
// For non-namespaced RSS elements, use "rss" as the namespace.
func (f *Feed) GetExtension(namespace, element string) []ext.Extension {
	if f.Extensions == nil {
		return nil
	}
	
	nsMap, ok := f.Extensions[namespace]
	if !ok {
		return nil
	}
	
	return nsMap[element]
}

// GetExtension retrieves extension values by namespace and element name.
// Returns a slice of Extension structs for the given namespace and element.
// For non-namespaced RSS elements, use "rss" as the namespace.
func (i *Item) GetExtension(namespace, element string) []ext.Extension {
	if i.Extensions == nil {
		return nil
	}
	
	nsMap, ok := i.Extensions[namespace]
	if !ok {
		return nil
	}
	
	return nsMap[element]
}

// GetExtensionValue is a convenience method that returns the text value
// of the first matching extension element, or empty string if not found.
func (f *Feed) GetExtensionValue(namespace, element string) string {
	exts := f.GetExtension(namespace, element)
	if len(exts) == 0 {
		return ""
	}
	return exts[0].Value
}

// GetExtensionValue is a convenience method that returns the text value
// of the first matching extension element, or empty string if not found.
func (i *Item) GetExtensionValue(namespace, element string) string {
	exts := i.GetExtension(namespace, element)
	if len(exts) == 0 {
		return ""
	}
	return exts[0].Value
}

// GetCustomValue retrieves the text value of a non-namespaced custom RSS element.
// This is a convenience method that replaces the previous Item.Custom[key] access pattern.
// Returns empty string if the element is not found.
func (i *Item) GetCustomValue(element string) string {
	return i.GetExtensionValue("rss", element)
}