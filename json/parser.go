package json

import (
	"bytes"
	"io"

	jsoniter "github.com/json-iterator/go"
)


// Parser is an JSON Feed Parser
type Parser struct{}

// Parse parses an json feed into an json.Feed
func (ap *Parser) Parse(feed io.Reader) (*Feed, error) {
	jsonFeed := &Feed{}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(feed)

	j := jsoniter.ConfigCompatibleWithStandardLibrary
	err := j.Unmarshal(buffer.Bytes(), jsonFeed)
	if err != nil {
		return nil, err
	}
	return jsonFeed, err
}
