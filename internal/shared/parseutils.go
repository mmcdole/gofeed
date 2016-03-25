package shared

import (
	"regexp"
	"strings"

	"github.com/mmcdole/goxpp"
)

var (
	emailNameRgx = regexp.MustCompile(`^([^@]+@[^\s]+)\s+\(([^@]+)\)$`)
	nameEmailRgx = regexp.MustCompile(`^([^@]+)\s+\(([^@]+@[^)]+)\)$`)
	nameOnlyRgx  = regexp.MustCompile(`^([^@()]+)$`)
	emailOnlyRgx = regexp.MustCompile(`^([^@()]+@[^@()]+)$`)
)

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

func ParseNameAddress(nameAddressText string) (name string, address string) {
	if nameAddressText == "" {
		return
	}

	if emailNameRgx.MatchString(nameAddressText) {
		result := emailNameRgx.FindStringSubmatch(nameAddressText)
		address = result[1]
		name = result[1]
	} else if nameEmailRgx.MatchString(nameAddressText) {
		result := nameEmailRgx.FindStringSubmatch(nameAddressText)
		name = result[1]
		address = result[2]
	} else if nameOnlyRgx.MatchString(nameAddressText) {
		result := nameOnlyRgx.FindStringSubmatch(nameAddressText)
		name = result[1]
	} else if emailOnlyRgx.MatchString(nameAddressText) {
		result := emailOnlyRgx.FindStringSubmatch(nameAddressText)
		address = result[1]
	}
	return
}
