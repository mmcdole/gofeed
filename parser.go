package gofeed

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mmcdole/gofeed/v2/atom"
	"github.com/mmcdole/gofeed/v2/json"
	"github.com/mmcdole/gofeed/v2/rss"
)

// ErrFeedTypeNotDetected is returned when the detection system can not figure
// out the Feed format
var ErrFeedTypeNotDetected = errors.New("Failed to detect feed type")

// HTTPError represents an HTTP error returned by a server.
type HTTPError struct {
	StatusCode int
	Status     string
}

// Error returns the string representation of the HTTP error.
func (err HTTPError) Error() string {
	return fmt.Sprintf("http error: %s", err.Status)
}

// Parser is a universal feed parser that detects
// a given feed type, parsers it, and translates it
// to the universal feed type.
type Parser struct {
	AtomTranslator Translator
	RSSTranslator  Translator
	JSONTranslator Translator
	rp             *rss.Parser
	ap             *atom.Parser
	jp             *json.Parser
}

// Auth is a structure allowing to
// use the BasicAuth during the HTTP request
// It must be instantiated with your new Parser
type Auth struct {
	Username string
	Password string
}

// NewParser creates a universal feed parser.
func NewParser() *Parser {
	fp := Parser{
		rp:        &rss.Parser{},
		ap:        &atom.Parser{},
		jp:        &json.Parser{},
	}
	return &fp
}

// Parse parses a RSS or Atom or JSON feed into
// the universal gofeed.Feed.  It takes an
// io.Reader which should return the xml/json content.
func (f *Parser) Parse(feed io.Reader, opts *ParseOptions) (*Feed, error) {
	if opts == nil {
		opts = DefaultParseOptions()
	}
	
	// Wrap the feed io.Reader in a io.TeeReader
	// so we can capture all the bytes read by the
	// DetectFeedType function and construct a new
	// reader with those bytes intact for when we
	// attempt to parse the feeds.
	var buf bytes.Buffer
	tee := io.TeeReader(feed, &buf)
	feedType := DetectFeedType(tee)

	// Glue the read bytes from the detect function
	// back into a new reader
	r := io.MultiReader(&buf, feed)

	switch feedType {
	case FeedTypeAtom:
		return f.parseAtomFeed(r, opts)
	case FeedTypeRSS:
		return f.parseRSSFeed(r, opts)
	case FeedTypeJSON:
		return f.parseJSONFeed(r, opts)
	}

	return nil, ErrFeedTypeNotDetected
}

// ParseURL fetches the contents of a given url and attempts to parse
// the response into the universal feed type. Context can be used to control
// timeout and cancellation.
func (f *Parser) ParseURL(ctx context.Context, feedURL string, opts *ParseOptions) (feed *Feed, err error) {
	if opts == nil {
		opts = DefaultParseOptions()
	}

	client := opts.RequestOptions.Client
	if client == nil {
		client = &http.Client{
			Timeout: opts.RequestOptions.Timeout,
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	
	// Set user agent
	if opts.RequestOptions.UserAgent != "" {
		req.Header.Set("User-Agent", opts.RequestOptions.UserAgent)
	}

	// TODO: Implement conditional request support (IfNoneMatch, IfModifiedSince)
	// This will be implemented as part of issue #247

	// Set auth if provided
	if auth, ok := opts.RequestOptions.AuthConfig.(*Auth); ok && auth != nil && auth.Username != "" && auth.Password != "" {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			if ce := resp.Body.Close(); ce != nil && err == nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	// Parse the feed
	return f.Parse(resp.Body, opts)
}

// ParseString parses a feed XML string and into the
// universal feed type.
func (f *Parser) ParseString(feed string, opts *ParseOptions) (*Feed, error) {
	return f.Parse(strings.NewReader(feed), opts)
}

func (f *Parser) parseAtomFeed(feed io.Reader, opts *ParseOptions) (*Feed, error) {
	af, err := f.ap.Parse(feed, opts)
	if err != nil {
		return nil, err
	}
	
	result, err := f.atomTrans().Translate(af, opts)
	if err != nil {
		return nil, err
	}
	
	if opts != nil && opts.KeepOriginalFeed {
		result.OriginalFeed = af
	}
	
	return result, nil
}

func (f *Parser) parseRSSFeed(feed io.Reader, opts *ParseOptions) (*Feed, error) {
	rf, err := f.rp.Parse(feed, opts)
	if err != nil {
		return nil, err
	}

	result, err := f.rssTrans().Translate(rf, opts)
	if err != nil {
		return nil, err
	}
	
	if opts != nil && opts.KeepOriginalFeed {
		result.OriginalFeed = rf
	}
	
	return result, nil
}

func (f *Parser) parseJSONFeed(feed io.Reader, opts *ParseOptions) (*Feed, error) {
	jf, err := f.jp.Parse(feed, opts)
	if err != nil {
		return nil, err
	}
	
	result, err := f.jsonTrans().Translate(jf, opts)
	if err != nil {
		return nil, err
	}
	
	if opts != nil && opts.KeepOriginalFeed {
		result.OriginalFeed = jf
	}
	
	return result, nil
}

func (f *Parser) atomTrans() Translator {
	if f.AtomTranslator != nil {
		return f.AtomTranslator
	}
	f.AtomTranslator = &DefaultAtomTranslator{}
	return f.AtomTranslator
}

func (f *Parser) rssTrans() Translator {
	if f.RSSTranslator != nil {
		return f.RSSTranslator
	}
	f.RSSTranslator = &DefaultRSSTranslator{}
	return f.RSSTranslator
}

func (f *Parser) jsonTrans() Translator {
	if f.JSONTranslator != nil {
		return f.JSONTranslator
	}
	f.JSONTranslator = &DefaultJSONTranslator{}
	return f.JSONTranslator
}

