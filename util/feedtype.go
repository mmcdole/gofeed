package parsers

type FeedType int

const (
	FeedTypeUnknown FeedType = iota
	FeedTypeAtom
	FeedTypeRSS
)

func DetectFeedType(feed string) FeedType {
	p := xpp.NewXMLPullParser(strings.NewReader(feed))

	_, err := p.NextTag()
	if err != nil {
		return FeedTypeUnknown
	}

	name := strings.ToLower(p.Name)
	switch name {
	case "rdf":
		return FeedTypeRSS
	case "rss":
		return FeedTypeRSS
	case "feed":
		return FeedTypeAtom
	default:
		return FeedTypeUnknown
	}
}
