package parseutil

// Namespaces taken from github.com/kurtmckee/feedparser
// These are used for determining canonical name space prefixes
// for many of the popular RSS/Atom extensions.
//
// These prefixes override any different prefixes used in the feed itself.
var globalNamespaces = map[string]string{
	"": "",
	"http://backend.userland.com/rss":                                "",
	"http://blogs.law.harvard.edu/tech/rss":                          "",
	"http://purl.org/rss/1.0/":                                       "",
	"http://my.netscape.com/rdf/simple/0.9/":                         "",
	"http://example.com/newformat#":                                  "",
	"http://example.com/necho":                                       "",
	"http://purl.org/echo/":                                          "",
	"uri/of/echo/namespace#":                                         "",
	"http://purl.org/pie/":                                           "",
	"http://purl.org/atom/ns#":                                       "",
	"http://www.w3.org/2005/Atom":                                    "",
	"http://purl.org/rss/1.0/modules/rss091#":                        "",
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

func PrefixForNamespace(namespace string) (prefix string) {
	if p, ok := globalNamespaces[namespace]; ok {
		prefix = p
	}
	return prefix
}
