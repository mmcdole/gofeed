package xpp

import (
	"encoding/xml"
	"io"
)

type XMLEventType int

const (
	StartDocument XMLEventType = iota
	EndDocument
	StartTag
	EndTag
	Text
	//IgnorableWhitespace TODO
	//CDSECT TODO
)

type XMLPullParser struct {
	Depth int

	// Current Token State
	Event       EventType
	Attrs       xml.Attr
	Name        string
	SpacePrefix string
	Space       string
	Text        string

	decoder *xml.Decoder
	token   interface{}
}

func NewXMLPullParser(r io.Reader) *XMLPullParser {
	d := xml.NewDecoder(r)
	d.Strict = false
	return &XMLPullParser{decoder: d, TokenType: StartDocument, Depth: 0}
}

func (p *XMLPullParser) Next() (XMLEventType, error) {

}

func (p *XMLPullParser) NextTag() (XMLEventType, error) {

}

func (p *XMLPullParser) NextToken() (event XMLEventType, err error) {
	tok, err := p.decoder.Token()
	if err != nil {
		return
	}

	p.token = tok

	switch tt := tok.(type) {
	case xml.StartElement:
		p.processStartToken(tt)
	case xml.EndElement:
		p.processEndToken(tt)
	case xml.CharData:
		p.processTextToken(tt)
	}
}

func (p *XMLPullParser) NextText() (string, error) {

}

func (p *XMLPullParser) Skip() error {
	for {
		tok, err := p.NextToken()
		if err != nil {
			return err
		}
		if tok == StartTag {
			if err := p.Skip(); err != nil {
				return err
			}
		} else if tok == EndTag {
			return nil
		}
	}
}

func (p *XMLPullParser) Attribute(name string) *xml.Attr {
}

func (p *XMLPullParser) Matches(event XMLEventType, namespace *string, name *string) bool {
	return p.Event == event && (namespace == nil || p.Namespace == namespace) && (name == nil || p.Name == name)
}

func (p *XMLPullParser) processStartToken(t *xml.StartElement) {
	p.Depth++
	p.Event = StartTag
	p.Attrs = t.Attr
	p.Name = t.Name
	p.Space = t.Space
}

func (p *XMLPullParser) processEndToken(t *xml.EndElement) {
	p.Depth--
	p.Event = EndTag
	p.Name = t.Name
}

func (p *XMLPullParser) processTextToken(t *xml.CharData) {
	p.Event = Text
	p.Text = string([]byte(t))
}
