package json

import (
	"bytes"
	"encoding/json"
	"io"
	
	"github.com/mmcdole/gofeed/v2/internal/shared"
)


// Parser is an JSON Feed Parser
type Parser struct{}

// NewParser creates a new JSON Feed parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses an json feed into an json.Feed
func (ap *Parser) Parse(feed io.Reader, opts *shared.ParseOptions) (*Feed, error) {
	jsonFeed := &Feed{}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(feed)

	err := json.Unmarshal(buffer.Bytes(), jsonFeed)
	if err != nil {
		return nil, err
	}
	
	return jsonFeed, nil
}
