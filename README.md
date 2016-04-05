# gofeed 

[![Build Status](https://travis-ci.org/mmcdole/gofeed.svg?branch=master)](https://travis-ci.org/mmcdole/gofeed) [![Coverage Status](https://coveralls.io/repos/github/mmcdole/gofeed/badge.svg?branch=master)](https://coveralls.io/github/mmcdole/gofeed?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/mmcdole/gofeed)](https://goreportcard.com/report/github.com/mmcdole/gofeed) [![](https://godoc.org/github.com/mmcdole/gofeed?status.svg)](http://godoc.org/github.com/mmcdole/gofeed) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

The `gofeed` library is a robust feed parser that supports parsing both [RSS](https://en.wikipedia.org/wiki/RSS) and [Atom](https://en.wikipedia.org/wiki/Atom_(standard)) feeds.  These can be parsed faithfully into their respective ```atom.Feed``` and ```rss.Feed``` representations using the ```atom.Parser``` or ```rss.Parser```. You can also parse them with the universal ```gofeed.FeedParser``` that will detect the feed type, parse it and then normalize both types of feeds into a hybrid ```gofeed.Feed``` representation.

##### Supported feed types:
* RSS 0.90
* Netscape RSS 0.91
* Userland RSS 0.91
* RSS 0.92
* RSS 0.93
* RSS 0.94
* RSS 1.0
* RSS 2.0
* Atom 0.3
* Atom 1.0

It also provides support for parsing several popular extension modules, including [Dublin Core](http://dublincore.org/documents/dces/) and [Apple’s iTunes](https://help.apple.com/itc/podcasts_connect/#/itcb54353390) extensions.  See the [Extensions](#extensions) section for more details.

`gofeed` is currently considered **Alpha Quality**. While this package is backed by a [large number of unit tests](https://github.com/mmcdole/gofeed/tree/master/testdata) it still has not achieved a public 1.0 release.  There are a few features still pending such as Atom relative URL resolution and the extension parsers still lack unit tests.

## Table of Contents
- [Design](#design)
- [Basic Usage](#basic-usage)
- [Advanced Usage](#advanced-usage)
- [Extensions](#extensions)
- [Default Mappings](#default-mappings)
- [Dependencies](#dependencies)
- [License](#license)
- [Credits](#credits)

## Design

When using the `gofeed` library as a [universal feed parser](#universal-feed-parser), it performs feed parsing in 2 stages.  It first parses the feed into its true representation (RSS or Atom specific models).  These models cover every field possible for their respective feed types.  They are then *translated* into a more generic model that is a hybrid of the RSS and Atom specification.  Most feed parsing libraries will parse and translate to a universal model in a single pass.  However, by doing it in 2 passes it allows for more flexibility and keeps the code base more maintainable by seperating RSS and Atom parsing in to seperate packages.

![Diagram](https://raw.githubusercontent.com/mmcdole/gofeed/master/docs/sequence.png)

Default translators (`DefaultRSSTranslator` and `DefaultAtomTranslator`) have been provided for you and are used transparently behind the scenes when you use `gofeed.FeedParser` with its default settings.  You can see how they translate fields from ```atom.Feed``` or ```rss.Feed``` to the universal ```gofeed.Feed``` struct in the [Default Mappings](#default-mappings) section.  However, should you disagree with the way certain fields are translated you can easily supply your own `RSSTranslator` or `AtomTranslator` and override this behavior.  See the [Advanced Usage](#advanced-usage) section for an example how to do this.

## Basic Usage

#### Universal Feed Parser

The most common usage scenario will be to use ```gofeed.FeedParser``` to parse an arbitrary RSS or Atom feed into the hybrid ```gofeed.Feed```.  This is useful for when you don't know what feed type your feeds will be ahead of time.  This hybrid struct has a lot of the common properties between the two formats (but does not have all the properties).  See the [default mappings](#default-mappings) section for more details.

##### Parse a feed from an URL:

```go
fp := gofeed.NewFeedParser()
feed := fp.ParseFeedURL("http://feeds.twit.tv/twit.xml")
fmt.Println(feed.Title)
```

##### Parse a feed from a string:

```go
feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
fp := gofeed.NewFeedParser()
feed := fp.ParseFeed(feedData)
fmt.Println(feed.Title)
```

#### Feed Specific Parsers

If you know in advanced that you will be parsing an RSS or Atom feed it can sometimes be desirable to utilize the ```atom.Parser``` or the ```rss.Parser``` directly.  Not only will they parse the feed more efficiently but they also expose all fields of their respective feed formats (some of which will be missing from the universal ```gofeed.Feed```).

##### Parse a RSS feed into a `rss.Feed`

```go
feedData := `<rss version="2.0">
<channel>
<webMaster>example@site.com (Example Name)</webMaster>
</channel>
</rss>`
fp := rss.Parser{}
rssFeed := fp.ParseFeed(feedData)
fmt.Println(rssFeed.WebMaster)
```

##### Parse an Atom feed into a `atom.Feed`

```go
feedData := `<feed xmlns="http://www.w3.org/2005/Atom">
<subtitle>Example Atom</subtitle>
</feed>`
fp := atom.Parser{}
atomFeed := fp.ParseFeed(feedData)
fmt.Println(atomFeed.Subtitle)
```

## Advanced Usage

##### Parse a feed while using a custom translator

The mappings and precedence order that are outlined in the [Default Mappings](#default-mappings) section are provided by the following two structs: `DefaultRSSTranslator` and `DefaultAtomTranslator`.  If you have fields that you think should have a different precedence, or if you want to make a translator that is aware of an unsupported extension you can do this by specifying your own RSS or Atom translator when using the `gofeed.FeedParser`.

Here is a simple example of creating a custom `RSSTranslator` that makes the `/rss/channel/itunes:author` extension field have a higher precedence than the `/rss/channel/managingEditor` field.  We will wrap the existing `DefaultRSSTranslator` since we only want to change the behavior for a single field.

```go
type MyCustomTranslator struct {
    defaultTranslator *DefaultRSSTranslator
}

func NewMyCustomTranslator() *MyCustomTranslator {
  t := &MyCustomTranslator{}
  
  // We create a DefaultRSSTranslator internally so we can wrap its call
  // since we only want to modify the precedence for a single field.
  t.defaultTranslator = &DefaultRSSTranslator{}
  return t
}

func (ct* MyCustomTranslator) Translate(rss *rss.Feed) *Feed {
  f := ct.Translate(rss)
  
  if rss.ITunesExt != nil && rss.ITunesExt.Author != "" {
      f.Author = rss.ITunesExt.Author
  } else {
      f.Author = rss.ManagingEditor
  }
  return f
}

func main() {
    feedData := `<rss version="2.0">
    <channel>
    <managingEditor>Ender Wiggin</managingEditor>
    <itunes:author>Valentine Wiggin</itunes:author>
    </channel>
    </rss>`
    
    fp := gofeed.NewFeedParser()
    fp.RSSTrans = NewMyCustomTranslator()
    feed := fp.ParseFeed(feedData)
    fmt.Println(feed.Author) // Valentine Wiggin
}
```

## Extensions 

Every element which does not belong to the default namespace is considered an extension by `gofeed`.  These are parsed and stored in a tree-like structure located at `Feed.Extensions` and `Item.Extensions`.  These fields should allow you to access and read any custom extension elements.

In addition to the generic handling of extensions, `gofeed` also has built in support for parsing certain popular extensions into their own structs for convenience.  It currently supports the explicit parsing of the [Dublin Core](http://dublincore.org/documents/dces/) and [Apple iTunes](https://help.apple.com/itc/podcasts_connect/#/itcb54353390) extensions which you can access at `Feed.ItunesExt`, `feed.DublinCoreExt` and `Item.ITunesExt` and `Item.DublinCoreExt`

## Default Mappings

The ```DefaultRSSTranslator``` and the ```DefaultAtomTranslator``` map the following ```rss.Feed``` and ```atom.Feed``` fields to their respective ```gofeed.Feed``` fields.  They are listed in order of precedence (highest to lowest):


Feed | RSS | Atom
--- | --- | ---
Title | /rss/channel/title<br>/rdf:RDF/channel/title<br>/rss/channel/dc:title<br>/rdf:RDF/channel/dc:title | /feed/title
Description | /rss/channel/description<br>/rdf:RDF/channel/description<br>/rss/channel/itunes:subtitle | /feed/subtitle<br>/feed/tagline
Link | /rss/channel/link<br>/rdf:RDF/channel/link | /feed/link[@rel=”alternate”]/@href<br>/feed/link[not(@rel)]/@href
FeedLink | /rss/channel/atom:link[@rel="self"]/@href<br>/rdf:RDF/channel/atom:link[@rel="self"]/@href | /feed/link[@rel="self"]/@href
Updated | /rss/channel/lastBuildDate<br>/rss/channel/dc:date<br>/rdf:RDF/channel/dc:date | /feed/updated<br>/feed/modified
Published | /rss/channel/pubDate |
Author | /rss/channel/managingEditor<br>/rss/channel/webMaster<br>/rss/channel/dc:author<br>/rdf:RDF/channel/dc:author<br>/rss/channel/dc:creator<br>/rdf:RDF/channel/dc:creator<br>/rss/channel/itunes:author | /feed/author
Language | /rss/channel/language<br>/rss/channel/dc:language<br>/rdf:RDF/channel/dc:language | /feed/@xml:lang
Image | /rss/channel/image<br>/rdf:RDF/image<br>/rss/channel/itunes:image | /feed/logo
Copyright | /rss/channel/copyright<br>/rss/channel/dc:rights<br>/rdf:RDF/channel/dc:rights | /feed/rights<br>/feed/copyright
Generator | /rss/channel/generator | /feed/generator
Categories | /rss/channel/category<br>/rss/channel/itunes:category<br>/rss/channel/itunes:keywords<br>/rss/channel/dc:subject<br>/rdf:RDF/channel/dc:subject | /feed/category


Item | RSS | Atom
--- | --- | ---
Title | /rss/channel/item/title<br>/rdf:RDF/item/title<br>/rdf:RDF/item/dc:title<br>/rss/channel/item/dc:title | /feed/entry/title
Description | /rss/channel/item/description<br>/rdf:RDF/item/description<br>/rss/channel/item/dc:description<br>/rdf:RDF/item/dc:description | /feed/entry/summary
Content | | /feed/entry/content
Link | /rss/channel/item/link<br>/rdf:RDF/item/link | /feed/entry/link[@rel=”alternate”]/@href<br>/feed/entry/link[not(@rel)]/@href
Updated | /rss/channel/item/dc:date<br>/rdf:RDF/rdf:item/dc:date | /feed/entry/modified<br>/feed/entry/updated
Published | /rss/channel/item/pubDate | /feed/entry/published<br>/feed/entry/issued
Author | /rss/channel/item/author<br>/rss/channel/item/dc:author<br>/rdf:RDF/item/dc:author<br>/rss/channel/item/dc:creator<br>/rdf:RDF/item/dc:creator<br>/rss/channel/item/itunes:author | /feed/entry/author
Guid |  /rss/channel/item/guid | /feed/entry/id
Image | /rss/channel/item/itunes:image<br>/rss/channel/item/media:image |
Categories | /rss/channel/item/category<br>/rss/channel/item/dc:subject<br>/rss/channel/item/itunes:keywords<br>/rdf:RDF/channel/item/dc:subject | /feed/entry/category
Enclosures | /rss/channel/item/enclosure | /feed/entry/link[@rel=”enclosure”]

## Dependencies

* [goxpp](https://github.com/mmcdole/goxpp) - XML Pull Parser
* [goquery](https://github.com/PuerkitoBio/goquery) - Go jQuery-like interface
* [testify](https://github.com/stretchr/testify) - Unit test enhancements

## License

This project is licensed under the [MIT License](https://raw.githubusercontent.com/mmcdole/gofeed/master/LICENSE)

## Credits

* [Mark Pilgrim](https://en.wikipedia.org/wiki/Mark_Pilgrim) for his work on the excellent [Universal Feed Parser](https://github.com/kurtmckee/feedparser) Python library.  This library was referenced several times during the development of `gofeed`.  It's unit test cases were also ported to `gofeed` project as well.
* [Dan MacTough](http://blog.mact.me) for his work on [node-feedparser](https://github.com/danmactough/node-feedparser).  It provided inspiration for the `gofeed.Feed` properties.
* [Matt Jibson](https://mattjibson.com/) for his date parsing function in the [goread](https://github.com/mjibson/goread) project.
* [Jim Teeuwen](https://github.com/jteeuwen) for his method of representing arbitrary feed extensions in the [go-pkg-rss](https://github.com/jteeuwen/go-pkg-rss) library.
