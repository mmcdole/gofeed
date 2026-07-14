package shared

import (
	"strings"

	"github.com/mmcdole/gofeed/extensions"
	xpp "github.com/mmcdole/goxpp/v2"
)

// IsExtension returns whether or not the current
// XML element is an extension element (if it has a
// non empty prefix)
func IsExtension(p *xpp.Parser) bool {
	prefix := PrefixForNamespace(p.Space(), p)
	return !(prefix == "" || prefix == "rss" || prefix == "rdf" || prefix == "content")
}

// ParseExtension parses the current element of the
// XMLPullParser as an extension element and updates
// the extension map
func ParseExtension(fe ext.Extensions, p *xpp.Parser) (ext.Extensions, error) {
	prefix := PrefixForNamespace(p.Space(), p)

	result, err := parseExtensionElement(p)
	if err != nil {
		return nil, err
	}

	// Ensure the extension prefix map exists
	if _, ok := fe[prefix]; !ok {
		fe[prefix] = map[string][]ext.Extension{}
	}
	// Ensure the extension element slice exists
	if _, ok := fe[prefix][p.Name()]; !ok {
		fe[prefix][p.Name()] = []ext.Extension{}
	}

	fe[prefix][p.Name()] = append(fe[prefix][p.Name()], result)
	return fe, nil
}

func parseExtensionElement(p *xpp.Parser) (e ext.Extension, err error) {
	if err = p.Expect(xpp.StartTag, "*"); err != nil {
		return e, err
	}

	e.Name = p.Name()
	e.Children = map[string][]ext.Extension{}
	e.Attrs = map[string]string{}

	for _, attr := range p.Attrs() {
		// TODO: Alright that we are stripping
		// namespace information from attributes ?
		e.Attrs[attr.Name.Local] = attr.Value
	}

	for {
		tok, err := p.Next()
		if err != nil {
			return e, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			child, err := parseExtensionElement(p)
			if err != nil {
				return e, err
			}

			if _, ok := e.Children[child.Name]; !ok {
				e.Children[child.Name] = []ext.Extension{}
			}

			e.Children[child.Name] = append(e.Children[child.Name], child)
		} else if tok == xpp.Text {
			e.Value += p.Text()
		}
	}

	e.Value = strings.TrimSpace(e.Value)

	if err = p.Expect(xpp.EndTag, e.Name); err != nil {
		return e, err
	}

	return e, nil
}

func PrefixForNamespace(space string, p *xpp.Parser) string {
	// Namespace attribute values may legally carry surrounding whitespace.
	// Trim here, once, so every lookup below (and every caller) agrees on
	// the key.
	space = strings.TrimSpace(space)

	// First we check if the global namespace map
	// contains an entry for this namespace/prefix.
	// This way we can use the canonical prefix for this
	// ns instead of the one defined in the feed.
	if prefix, ok := canonicalNamespaces[space]; ok {
		return prefix
	}

	// Next we check if the feed itself declared a prefix for this
	// namespace and return it if we have a result.
	if prefix, ok := p.PrefixForURI(space); ok {
		return prefix
	}

	// Lastly, any namespace which is not defined in the
	// the feed will be the prefix itself when using Go's
	// xml.Decoder.Token() method.
	return space
}

// Namespaces taken from github.com/kurtmckee/feedparser
// These are used for determining canonical name space prefixes
// for many of the popular RSS/Atom extensions.
//
// These canonical prefixes override any prefixes used in the feed itself.
var canonicalNamespaces = map[string]string{
	"http://webns.net/mvcb/":                                         "admin",
	"http://purl.org/rss/1.0/modules/aggregation/":                   "ag",
	"http://purl.org/rss/1.0/modules/annotate/":                      "annotate",
	"http://media.tangent.org/rss/1.0/":                              "audio",
	"http://backend.userland.com/blogChannelModule":                  "blogChannel",
	"http://creativecommons.org/ns#license":                          "cc",
	"http://web.resource.org/cc/":                                    "cc",
	"http://cyber.law.harvard.edu/rss/creativeCommonsRssModule.html": "creativeCommons",
	"http://backend.userland.com/creativeCommonsRssModule":           "creativeCommons",
	"http://purl.org/rss/1.0/modules/company":                        "co",
	"http://purl.org/rss/1.0/modules/content/":                       "content",
	"http://my.theinfo.org/changed/1.0/rss/":                         "cp",
	"http://purl.org/dc/elements/1.1/":                               "dc",
	"http://purl.org/dc/terms/":                                      "dcterms",
	"http://purl.org/rss/1.0/modules/email/":                         "email",
	"http://purl.org/rss/1.0/modules/event/":                         "ev",
	"http://rssnamespace.org/feedburner/ext/1.0":                     "feedburner",
	"http://freshmeat.net/rss/fm/":                                   "fm",
	"http://xmlns.com/foaf/0.1/":                                     "foaf",
	"http://www.w3.org/2003/01/geo/wgs84_pos#":                       "geo",
	"http://www.georss.org/georss":                                   "georss",
	"http://www.opengis.net/gml":                                     "gml",
	"http://postneo.com/icbm/":                                       "icbm",
	"http://purl.org/rss/1.0/modules/image/":                         "image",
	"http://www.itunes.com/DTDs/PodCast-1.0.dtd":                     "itunes",
	"http://example.com/DTDs/PodCast-1.0.dtd":                        "itunes",
	"http://purl.org/rss/1.0/modules/link/":                          "l",
	"http://search.yahoo.com/mrss":                                   "media",
	"http://search.yahoo.com/mrss/":                                  "media",
	"http://madskills.com/public/xml/rss/module/pingback/":           "pingback",
	"http://prismstandard.org/namespaces/1.2/basic/":                 "prism",
	"http://www.w3.org/1999/02/22-rdf-syntax-ns#":                    "rdf",
	"http://www.w3.org/2000/01/rdf-schema#":                          "rdfs",
	"http://purl.org/rss/1.0/modules/reference/":                     "ref",
	"http://purl.org/rss/1.0/modules/richequiv/":                     "reqv",
	"http://purl.org/rss/1.0/modules/search/":                        "search",
	"http://purl.org/rss/1.0/modules/slash/":                         "slash",
	"http://schemas.xmlsoap.org/soap/envelope/":                      "soap",
	"http://purl.org/rss/1.0/modules/servicestatus/":                 "ss",
	"http://hacks.benhammersley.com/rss/streaming/":                  "str",
	"http://purl.org/rss/1.0/modules/subscription/":                  "sub",
	"http://purl.org/rss/1.0/modules/syndication/":                   "sy",
	"http://schemas.pocketsoap.com/rss/myDescModule/":                "szf",
	"http://purl.org/rss/1.0/modules/taxonomy/":                      "taxo",
	"http://purl.org/rss/1.0/modules/threading/":                     "thr",
	"http://purl.org/rss/1.0/modules/textinput/":                     "ti",
	"http://madskills.com/public/xml/rss/module/trackback/":          "trackback",
	"http://wellformedweb.org/commentAPI/":                           "wfw",
	"http://purl.org/rss/1.0/modules/wiki/":                          "wiki",
	"http://www.w3.org/1999/xhtml":                                   "xhtml",
	"http://www.w3.org/1999/xlink":                                   "xlink",
	"http://www.w3.org/XML/1998/namespace":                           "xml",
	"http://podlove.org/simple-chapters":                             "psc",
}

// CustomPrefix is the pseudo namespace prefix that files non-namespaced
// unknown elements in the extension map.
const CustomPrefix = "_custom"

// ParseCustom parses the current element as an arbitrary non-namespaced
// element and files it in the extension map under CustomPrefix, preserving
// nesting, attributes and repetition the same way namespaced extensions are
// kept. The parsed element is also returned so callers can maintain
// compatibility shims from it.
func ParseCustom(fe ext.Extensions, p *xpp.Parser) (ext.Extensions, ext.Extension, error) {
	result, err := parseExtensionElement(p)
	if err != nil {
		return nil, ext.Extension{}, err
	}

	if _, ok := fe[CustomPrefix]; !ok {
		fe[CustomPrefix] = map[string][]ext.Extension{}
	}
	fe[CustomPrefix][result.Name] = append(fe[CustomPrefix][result.Name], result)
	return fe, result, nil
}
