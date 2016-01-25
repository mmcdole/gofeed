package xpp

import (
	"encoding/xml"
	"io"
)

type XMLEventType int

const (
	XML_START_DOCUMENT XMLEventType = iota
	XML_END_DOCUMENT
	XML_START_TAG
	XML_END_TAG
	XML_TEXT
	XML_IGNORABLE_WHITESPACE
)

type XMLPullParser struct {
	Depth int

	// Current Token State
	TokenType EventType
	Attrs     xml.Attr
	Name      string
	Space     string
	Text      string

	decoder []xml.Decoder
}

func NewXMLPullParser(r io.Reader) *XMLPullParser {
	d := xml.NewDecoder(r)
	d.Strict = false
	return &XMLPullParser{decoder: d, TokenType: XML_START_DOCUMENT}
}

func (p *XMLPullParser) Next() (XMLEventType, error) {

}

func (p *XMLPullParser) NextTag() (XMLEventType, error) {

}

func (p *XMLPullParser) NextToken() (XMLEventType, error) {

}

func (p *XMLPullParser) NextText() (string, error) {

}

func (p *XMLPullParser) TagText() (string, error) {

}

func (p *XMLPullParser) GetName() string {

}

func (p *XMLPullParser) GetNamespace() string {

}

func (p *XMLPullParser) GetPrefixNamespace(prefix string) string {

}

func (p *XMLPullParser) GetAttributes() []xml.Attr {

}

func (p *XMLPullParser) GetAttribute(name string) *xml.Attr {

}
