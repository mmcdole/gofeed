package json

import (
	"encoding/json"
	"io"
)

// Render writes the JSON feed as JSON Feed 1.1 to the provided io.Writer
func (f *Feed) Render(w io.Writer) error {
	if f == nil {
		return nil
	}

	// Create JSON encoder with proper formatting
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(f)
}
