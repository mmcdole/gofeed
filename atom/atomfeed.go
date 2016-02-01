package atom

type AtomFeed struct {
	Title string
}

type AtomEntry struct {
	Title string
}

func ParseAtomFeed(feed string) (*AtomFeed, error) {
	return &AtomFeed{}, nil
}
