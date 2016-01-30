package feed

import ()

type Extension struct {
	Name     string
	Value    string
	Attrs    map[string]string
	Children map[string][]Extension
}

type Feed struct {
	Title       string
	Description string
	Items       []FeedItem
	FeedType    string
	FeedVersion string
	Custom      map[string]string
	Extensions  map[string]map[string][]Extension
}

type FeedItem struct {
	Title       string
	Description string
	Custom      map[string]string
	Extensions  map[string]map[string][]Extension
}
