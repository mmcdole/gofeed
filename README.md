# gofeed 

[![Build Status](https://travis-ci.org/mmcdole/gofeed.svg?branch=master)](https://travis-ci.org/mmcdole/gofeed) [![Coverage Status](https://coveralls.io/repos/github/mmcdole/gofeed/badge.svg?branch=master)](https://coveralls.io/github/mmcdole/gofeed?branch=master)

gofeed is a robust feed parser that supports parsing RSS 0.90, Netscape RSS 0.91, Userland RSS 0.91, RSS 0.92, RSS 0.93, RSS 0.94, RSS 1.0, RSS 2.0, Atom 0.3, Atom 1.0 feeds.  It also provides support for parsing several popular extension modules, including Dublin Core and Apple’s iTunes extensions.

gofeed is currently considered **Alpha Quality**. While this package is backed by a [large number of unit tests](https://github.com/mmcdole/gofeed/tree/master/testdata) it still has not achieved a public 1.0 release.  There are a few features still pending such as Atom relative URL resolution and the extension parsers still lack unit tests.

## Table of Contents
- [Design](#design)
- [Basic Usage](#basic-usage)
- [Advanced Usage](#advanced-usage)
- [Default Mappings](#default-mappings)
- [Dependencies](#dependencies)
- [License](#license)

## Design

The design of gofeed is unique from many feed parser libraries in that it performs it's feed parsing in 2 stages.  It first parses the feed into its true representation (RSS or Atom specific models).  These models cover every field possible for their respective specification.  They are then *translated* into a generic feed model that is a hybrid of the RSS and Atom specification.  This keeps the parsing code for Atom feeds completely seperate from RSS feeds which I think makes the codebase easier to understand and maintain.

The default translators that convert from ```atom.Feed``` or ```rss.Feed``` to the generic ```Feed``` struct can be swapped out with your own.  This allows you to change the precedence for field mappings as well as make a translator that is aware of particular extension module that isn't supported by default.  See the [Default Mappings](#default-mappings) section for information on what the default translator's field precedence is as well as the [Advanced Usage](#advanced-usage) section to see how to provide your own translators.

## Basic Usage

Parse a feed from a URL:

```go
fp := gofeed.NewFeedParser()
feed := fp.ParseFeedURL("http://feeds.twit.tv/twit.xml")
fmt.Println(feed.Title)
```

Parse a feed from a string:

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

## Advanced Usage

TODO

## Default Mappings

The ```DefaultRSSTranslator``` and the ```DefaultAtomTranslator``` map the following ```rss.Feed``` and ```atom.Feed``` fields to their respective ```Feed``` fields.  They are listed in order of precedence (highest to lowest):


/atom03:feed/atom03:modified
/atom10:feed/atom10:updated


/rss/channel/dc:date
/rss/channel/lastBuildDate

Feed | RSS | Atom
--- | --- | ---
Title | /rss/channel/title<br>/rdf:RDF/channel/title<br>/rss/channel/dc:title<br>/rdf:RDF/channel/dc:title | /feed/title
Description | /rss/channel/description<br>/rdf:RDF/channel/description<br>/rss/channel/itunes:subtitle<br>/rdf:RDF/channel/itunes:subtitle | /feed/subtitle<br>/feed/tagline
Link | /rss/channel/link<br>/rdf:RDF/channel/link | /feed/link[@rel=”alternate”]/@href<br>/feed/link[not(@rel)]/@href
FeedLink | /rss/channel/atom:link[@rel="self]/@href<br>/rdf:RDF/channel/atom:link[@rel="self] | /feed/link[@rel="self"]/@href
Updated | /rss/channel/lastBuildDate<br>/rss/channel/dc:date<br>/rdf:RDF/channel/dc:date<br>/rdf:RDF/channel/dcterms:modified | /feed/updated<br>/feed/modified
