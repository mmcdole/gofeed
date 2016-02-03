package gofeed

import (
	"strings"

	"github.com/mmcdole/goxpp"
)

// Namespaces taken from github.com/kurtmckee/feedparser
// These are used for determining canonical name space prefixes
// for many of the popular RSS/Atom extensions.
//
// These prefixes override any different prefixes used in the feed itself.
var globalNamespaces = map[string]string{
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

type BaseParser struct {
	// Map of all namespaces (url / prefix)
	// that have been defined in the feed.
	feedSpaces map[string]string
}

func (rp *BaseParser) parseExtension(fe FeedExtensions, p *xpp.XMLPullParser) (FeedExtensions, error) {
	prefix := rp.prefixForNamespace(p.Space)

	result, err := rp.parseExtensionElement(p)
	if err != nil {
		return nil, err
	}

	// Ensure the extension prefix map exists
	if _, ok := fe[prefix]; !ok {
		fe[prefix] = map[string][]Extension{}
	}
	// Ensure the extension element slice exists
	if _, ok := fe[prefix][p.Name]; !ok {
		fe[prefix][p.Name] = []Extension{}
	}

	fe[prefix][p.Name] = append(fe[prefix][p.Name], result)
	return fe, nil
}

func (bp *BaseParser) parseExtensionElement(p *xpp.XMLPullParser) (ext Extension, err error) {
	if err = p.Expect(xpp.StartTag, "*"); err != nil {
		return ext, err
	}

	ext.Name = p.Name
	ext.Children = map[string][]Extension{}
	ext.Attrs = map[string]string{}

	for _, attr := range p.Attrs {
		// TODO: Alright that we are stripping
		// namespace information from attributes ?
		ext.Attrs[attr.Name.Local] = attr.Value
	}

	for {
		tok, err := p.Next()
		if err != nil {
			return ext, err
		}

		if tok == xpp.EndTag {
			break
		}

		if tok == xpp.StartTag {
			child, err := bp.parseExtensionElement(p)
			if err != nil {
				return ext, err
			}

			if _, ok := ext.Children[child.Name]; !ok {
				ext.Children[child.Name] = []Extension{}
			}

			ext.Children[child.Name] = append(ext.Children[child.Name], child)
		} else if tok == xpp.Text {
			ext.Value = strings.TrimSpace(p.Text)
		}
	}

	if err = p.Expect(xpp.EndTag, ext.Name); err != nil {
		return ext, err
	}

	return ext, nil
}

func (bp *BaseParser) prefixForNamespace(space string) string {
	// First we check if the global namespace map
	// contains an entry for this namespace/prefix.
	// This way we can use the canonical prefix for this
	// ns instead of the one defined in the feed.
	if prefix, ok := globalNamespaces[space]; ok {
		return prefix
	}

	// Next we check if the feed itself defined this
	// this namespace and return it if we have a result.
	if prefix, ok := bp.feedSpaces[space]; ok {
		return prefix
	}

	// Lastly, any namespace which is not defined in the
	// the feed will be the prefix itself when using Go's
	// xml.Decoder.Token() method.
	return space
}

func (bp *BaseParser) isExtension(p *xpp.XMLPullParser) bool {
	space := strings.TrimSpace(p.Space)
	if prefix, ok := bp.feedSpaces[space]; ok {
		return prefix != ""
	}

	return p.Space != ""
}

func (bp *BaseParser) parseNamespaces(p *xpp.XMLPullParser) {
	for _, attr := range p.Attrs {
		if attr.Name.Space == "xmlns" {
			space := strings.TrimSpace(attr.Value)
			spacePrefix := strings.TrimSpace(strings.ToLower(attr.Name.Local))
			bp.feedSpaces[space] = spacePrefix
		} else if attr.Name.Local == "xmlns" {
			space := strings.TrimSpace(attr.Value)
			bp.feedSpaces[space] = ""
		}
	}
}

func (bp *BaseParser) parseText(p *xpp.XMLPullParser) (text string, err error) {
	text, err = p.NextText()
	if err != nil {
		return text, err
	}

	text = strings.TrimSpace(text)
	return text, nil
}
