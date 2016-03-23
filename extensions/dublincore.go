package ext

type DublinCoreExtension struct {
	Title       []string `json:"title,omitempty"`
	Creator     []string `json:"creator,omitempty"`
	Subject     []string `json:"subject,omitempty"`
	Description []string `json:"description,omitempty"`
	Publisher   []string `json:"publisher,omitempty"`
	Contributor []string `json:"contributor,omitempty"`
	Date        []string `json:"date,omitempty"`
	Type        []string `json:"type,omitempty"`
	Format      []string `json:"format,omitempty"`
	Identifier  []string `json:"identifier,omitempty"`
	Source      []string `json:"source,omitempty"`
	Language    []string `json:"language,omitempty"`
	Relation    []string `json:"relation,omitempty"`
	Coverage    []string `json:"coverage,omitempty"`
	Rights      []string `json:"rights,omitempty"`
}

func NewDublinCoreExtension(extensions map[string][]Extension) *DublinCoreExtension {
	dc := &DublinCoreExtension{}
	dc.Title = ParseTextArrayExtension("title", extensions)
	dc.Creator = ParseTextArrayExtension("creator", extensions)
	dc.Subject = ParseTextArrayExtension("subject", extensions)
	dc.Description = ParseTextArrayExtension("description", extensions)
	dc.Publisher = ParseTextArrayExtension("publisher", extensions)
	dc.Contributor = ParseTextArrayExtension("contributor", extensions)
	dc.Date = ParseTextArrayExtension("date", extensions)
	dc.Type = ParseTextArrayExtension("type", extensions)
	dc.Format = ParseTextArrayExtension("format", extensions)
	dc.Identifier = ParseTextArrayExtension("identifier", extensions)
	dc.Source = ParseTextArrayExtension("source", extensions)
	dc.Language = ParseTextArrayExtension("language", extensions)
	dc.Relation = ParseTextArrayExtension("relation", extensions)
	dc.Coverage = ParseTextArrayExtension("coverage", extensions)
	dc.Rights = ParseTextArrayExtension("rights", extensions)
	return dc
}
