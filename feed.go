package feed

type Feed struct {
	Title        string
	Description  string
	Items        []FeedItem
	FeedType     string
	FeedVersion  string
	CustomFields map[string]string
}

type FeedItem struct {
	Title        string
	Description  string
	CustomFields map[string]string
}
