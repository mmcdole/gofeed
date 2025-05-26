package json

import (
	"bytes"
	"encoding/json"
	"io"
	
	"github.com/mmcdole/gofeed/v2/internal/shared"
)


// Parser is an JSON Feed Parser
type Parser struct{}

// Parse parses an json feed into an json.Feed
func (ap *Parser) Parse(feed io.Reader, opts *shared.ParseOptions) (*Feed, error) {
	jsonFeed := &Feed{}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(feed)

	err := json.Unmarshal(buffer.Bytes(), jsonFeed)
	if err != nil {
		return nil, err
	}
	
	// Apply MaxItems limit after unmarshaling
	if opts != nil && opts.MaxItems > 0 && len(jsonFeed.Items) > opts.MaxItems {
		jsonFeed.Items = jsonFeed.Items[:opts.MaxItems]
	}
	
	return jsonFeed, nil
}
