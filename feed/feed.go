package feed

type FeedExtensions map[string]map[string][]Extension

type Extension struct {
	Name     string                 `json:"name"`
	Value    string                 `json:"value,omitempty"`
	Attrs    map[string]string      `json:"attrs,omitempty"`
	Children map[string][]Extension `json:"children,omitempty"`
}

type Feed struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Items       []Item            `json:"items"`
	FeedType    string            `json:"feedType"`
	FeedVersion string            `json:"feedVersion"`
	Custom      map[string]string `json:"custom"`
	Extensions  FeedExtensions    `json:"extensions"`
}

type Item struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Custom      map[string]string `json:"custom"`
	Extensions  FeedExtensions    `json:"extensions"`
}
