package shared

import (
	"strings"

	"github.com/mmcdole/goxpp"
)

//func Expect(p *xpp.XMLPullParser, event xpp.XMLEventType, name string) (err error) {
//	if !(p.Event == event && strings.ToLower(p.Name) == strings.ToLower(name)) {
//		err = fmt.Errorf("Expected Name:%s Event:%s but got Name:%s Event:%s", name, p.EventName(event), p.Name, p.EventName(p.Event))
//	}
//	return
//}

func ParseText(p *xpp.XMLPullParser) (string, error) {
	var text struct {
		Type     string `xml:"type,attr"`
		Body     string `xml:",chardata"`
		InnerXML string `xml:",innerxml"`
	}

	err := p.DecodeElement(&text)
	if err != nil {
		return "", err
	}

	result := ""
	if len(text.InnerXML) > 0 {
		result = text.InnerXML
	} else if len(text.Body) > 0 {
		result = text.Body
	}

	result = strings.TrimSpace(result)
	result = DecodeEntities(result)
	return result, nil
}

func DecodeEntities(str string) string {
	str = strings.Replace(str, "&lt;", "<", -1)
	str = strings.Replace(str, "&gt;", ">", -1)
	str = strings.Replace(str, "&quot;", "\"", -1)
	str = strings.Replace(str, "&apos;", "'", -1)
	str = strings.Replace(str, "&amp;", "&", -1)
	return str
}
