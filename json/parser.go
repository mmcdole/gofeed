package json

import (
	"bytes"
	"encoding/json"
	"io"
)

var utf8BOM = []byte("\xef\xbb\xbf")

// Parser is an JSON Feed Parser
type Parser struct{}

// Parse parses an json feed into an json.Feed
func (ap *Parser) Parse(feed io.Reader) (*Feed, error) {
	jsonFeed := &Feed{}

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(feed); err != nil {
		return nil, err
	}

	b := bytes.TrimPrefix(buffer.Bytes(), utf8BOM)
	if err := json.Unmarshal(b, jsonFeed); err != nil {
		return nil, err
	}
	return jsonFeed, nil
}
