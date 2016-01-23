package feed

type RSSFeed struct {
	Title      string
	Items      []RSSItem
	Extensions map[string]interface{}
}

type RSSItem struct {
	Title      string
	Extensions map[string]interface{}
}

func ParseRSSFeed(feed string) (*RSSFeed, error) {
	return &RSSFeed{}, nil
}
