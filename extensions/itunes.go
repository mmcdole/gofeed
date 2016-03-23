package ext

type ITunesFeedExtension struct {
	Author     string            `json:"author,omitempty"`
	Block      string            `json:"block,omitempty"`
	Categories []*ITunesCategory `json:"categories,omitempty"`
	Explicit   string            `json:"explicit,omitempty"`
	Keywords   string            `json:"keywords,omitempty"`
	Owner      *ITunesOwner      `json:"owner,omitempty"`
	Subtitle   string            `json:"subtitle,omitempty"`
	Summary    string            `json:"summary,omitempty"`
	Image      string            `json:"image,omitempty"`
	Complete   string            `json:"complete,omitempty"`
	NewFeedURL string            `json:"newFeedUrl,omitempty"`
}

type ITunesEntryExtension struct {
	Author            string `json:"author,omitempty"`
	Block             string `json:"block,omitempty"`
	Duration          string `json:"duration,omitempty"`
	Explicit          string `json:"explicit,omitempty"`
	Keywords          string `json:"keywords,omitempty"`
	Subtitle          string `json:"subtitle,omitempty"`
	Summary           string `json:"summary,omitempty"`
	Image             string `json:"image,omitempty"`
	IsClosedCaptioned string `json:"isClosedCaptioned,omitempty"`
	Order             string `json:"order,omitempty"`
}

type ITunesCategory struct {
	Text        string          `json:"text,omitempty"`
	Subcategory *ITunesCategory `json:"subcategory,omitempty"`
}

type ITunesOwner struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// Parse an iTunes Feed Extension from the "itunes" entry in the
// extension map.
func NewITunesFeedExtension(extensions map[string][]Extension) *ITunesFeedExtension {
	feed := &ITunesFeedExtension{}
	feed.Author = ParseTextExtension("author", extensions)
	feed.Block = ParseTextExtension("block", extensions)
	feed.Explicit = ParseTextExtension("explicit", extensions)
	feed.Keywords = ParseTextExtension("keywords", extensions)
	feed.Subtitle = ParseTextExtension("subtitle", extensions)
	feed.Summary = ParseTextExtension("summary", extensions)
	feed.Image = ParseTextExtension("image", extensions)
	feed.Complete = ParseTextExtension("complete", extensions)
	feed.NewFeedURL = ParseTextExtension("new-feed-url", extensions)
	feed.Categories = parseCategories(extensions)
	feed.Owner = parseOwner(extensions)
	return feed
}

// Parse an iTunes Entry Extension from the "itunes" entry in the
// extension map.
func NewITunesEntryExtension(extensions map[string][]Extension) *ITunesEntryExtension {
	entry := &ITunesEntryExtension{}
	entry.Author = ParseTextExtension("author", extensions)
	entry.Block = ParseTextExtension("block", extensions)
	entry.Duration = ParseTextExtension("duration", extensions)
	entry.Explicit = ParseTextExtension("explicit", extensions)
	entry.Subtitle = ParseTextExtension("subtitle", extensions)
	entry.Summary = ParseTextExtension("summary", extensions)
	entry.Image = ParseTextExtension("image", extensions)
	entry.IsClosedCaptioned = ParseTextExtension("isClosedCaptioned", extensions)
	entry.Order = ParseTextExtension("order", extensions)
	return entry
}

func parseOwner(extensions map[string][]Extension) (owner *ITunesOwner) {
	if extensions == nil {
		return
	}

	matches, ok := extensions["owner"]
	if !ok || len(matches) == 0 {
		return
	}

	owner = &ITunesOwner{}
	if name, ok := matches[0].Children["name"]; ok {
		owner.Name = name[0].Value
	}
	if email, ok := matches[0].Children["email"]; ok {
		owner.Email = email[0].Value
	}
	return
}

func parseCategories(extensions map[string][]Extension) (categories []*ITunesCategory) {
	if extensions == nil {
		return
	}

	matches, ok := extensions["category"]
	if !ok || len(matches) == 0 {
		return
	}

	categories = []*ITunesCategory{}
	for _, cat := range matches {
		c := &ITunesCategory{}
		if text, ok := cat.Attrs["text"]; ok {
			c.Text = text
		}

		if subs, ok := cat.Children["category"]; ok {
			s := &ITunesCategory{}
			if text, ok := subs[0].Attrs["text"]; ok {
				s.Text = text
			}
			c.Subcategory = s
		}
		categories = append(categories, c)
	}
	return
}
