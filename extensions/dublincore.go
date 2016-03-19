package ext

type DublinCoreExtension struct {
	Title       string `json:"title,omitempty"`
	Creator     string `json:"creator,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
	Publisher   string `json:"publisher,omitempty"`
	Contributor string `json:"contributor,omitempty"`
	Date        string `json:"date,omitempty"`
	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Identifier  string `json:"identifier,omitempty"`
	Source      string `json:"source,omitempty"`
	Language    string `json:"language,omitempty"`
	Relation    string `json:"relation,omitempty"`
	Coverage    string `json:"coverage,omitempty"`
	Rights      string `json:"rights,omitempty"`
}

func ParseDublinCoreExtension(extensions map[string][]Extension) *DublinCoreExtension {
	dc := &DublinCoreExtension{}
	dc.Title = ParseTextExtension("title", extensions)
	dc.Creator = ParseTextExtension("creator", extensions)
	dc.Subject = ParseTextExtension("subject", extensions)
	dc.Description = ParseTextExtension("description", extensions)
	dc.Publisher = ParseTextExtension("publisher", extensions)
	dc.Contributor = ParseTextExtension("contributor", extensions)
	dc.Date = ParseTextExtension("date", extensions)
	dc.Type = ParseTextExtension("type", extensions)
	dc.Format = ParseTextExtension("format", extensions)
	dc.Identifier = ParseTextExtension("identifier", extensions)
	dc.Source = ParseTextExtension("source", extensions)
	dc.Language = ParseTextExtension("language", extensions)
	dc.Relation = ParseTextExtension("relation", extensions)
	dc.Coverage = ParseTextExtension("coverage", extensions)
	dc.Rights = ParseTextExtension("rights", extensions)
	return dc
}
