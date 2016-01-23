package feed

type AtomFeed struct {
	Title string
}

type AtomEntry struct {
	Title string
}

func ParseAtomFeed(feed string) (*AtomFeed, error) {
	return &AtomFeed{}, nil
}
