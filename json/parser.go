package json

import (
	"bytes"
	"encoding/json"
	"io"
)

// Parser is an JSON Feed Parser
type Parser struct{}

// Parse parses an json feed into an json.Feed
func (ap *Parser) Parse(feed io.Reader) (*Feed, error) {
	jsonFeed := &Feed{}

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(feed); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buffer.Bytes(), jsonFeed); err != nil {
		return nil, err
	}
	return jsonFeed, nil
}
