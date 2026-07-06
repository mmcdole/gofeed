package gofeed

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed/atom"
	"github.com/mmcdole/gofeed/json"
	"github.com/mmcdole/gofeed/rss"
)

// defaultRequestTimeout bounds ParseURL, which would otherwise use a client
// with no timeout and hang forever on a slow or unresponsive server.
// ParseURLWithContext callers manage their own context and are not subject to
// this.
const defaultRequestTimeout = 30 * time.Second

// ErrResponseTooLarge is returned when a response body exceeds Parser.MaxByteSize.
var ErrResponseTooLarge = errors.New("gofeed: response body exceeds MaxByteSize")

// ErrFeedTypeNotDetected is returned when the detection system can not figure
// out the Feed format
var ErrFeedTypeNotDetected = errors.New("failed to detect feed type")

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
	UserAgent      string
	AuthConfig     *Auth
	Client         *http.Client
	// MaxByteSize limits how many bytes ParseURL/ParseURLWithContext will read
	// from a response body. Zero means no limit. Exceeding it returns
	// ErrResponseTooLarge rather than silently truncating.
	MaxByteSize int64
	// KeepOriginalFeed retains the source rss/atom/json feed on the result,
	// accessible via Feed.OriginalFeed(). Off by default: keeping it holds a
	// second copy of the feed in memory for the lifetime of the result.
	KeepOriginalFeed bool
	rp               *rss.Parser
	ap               *atom.Parser
	jp               *json.Parser
}

// Auth is a structure allowing to
// use the BasicAuth during the HTTP request
// It must be instantiated with your new Parser
type Auth struct {
	Username string
	Password string
}

// Shared defaults used when the corresponding Parser field is unset. The
// default translators are stateless and http.Client is safe for concurrent
// use, so single shared instances are fine and avoid per-parse allocation.
var (
	defaultAtomTranslator = &DefaultAtomTranslator{}
	defaultRSSTranslator  = &DefaultRSSTranslator{}
	defaultJSONTranslator = &DefaultJSONTranslator{}
	defaultClient         = &http.Client{}
)

// NewParser creates a universal feed parser.
func NewParser() *Parser {
	fp := Parser{
		rp:        &rss.Parser{},
		ap:        &atom.Parser{},
		jp:        &json.Parser{},
		UserAgent: "Gofeed/1.0",
	}
	return &fp
}

// Parse parses a RSS or Atom or JSON feed into
// the universal gofeed.Feed.  It takes an
// io.Reader which should return the xml/json content.
func (f *Parser) Parse(feed io.Reader) (*Feed, error) {
	// Read the input once up front: detection inspects the bytes and the
	// format parser reads them again, and reading here surfaces an I/O
	// error as itself rather than as a failed type detection.
	data, err := io.ReadAll(feed)
	if err != nil {
		return nil, err
	}

	switch DetectFeedType(bytes.NewReader(data)) {
	case FeedTypeAtom:
		return f.parseAtomFeed(bytes.NewReader(data))
	case FeedTypeRSS:
		return f.parseRSSFeed(bytes.NewReader(data))
	case FeedTypeJSON:
		return f.parseJSONFeed(bytes.NewReader(data))
	}

	return nil, ErrFeedTypeNotDetected
}

// ParseURL fetches the contents of a given url and
// attempts to parse the response into the universal feed type. It applies a
// default request timeout; use ParseURLWithContext to control cancellation.
func (f *Parser) ParseURL(feedURL string) (feed *Feed, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	return f.ParseURLWithContext(feedURL, ctx)
}

// ParseURLWithContext fetches contents of a given url and
// attempts to parse the response into the universal feed type.
// You can instantiate the Auth structure with your Username and Password
// to use the BasicAuth during the HTTP call.
// It will be automatically added to the header of the request
// Request could be canceled or timeout via given context
func (f *Parser) ParseURLWithContext(feedURL string, ctx context.Context) (feed *Feed, err error) {
	client := f.httpClient()

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", f.UserAgent)

	if f.AuthConfig != nil && f.AuthConfig.Username != "" && f.AuthConfig.Password != "" {
		req.SetBasicAuth(f.AuthConfig.Username, f.AuthConfig.Password)
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

	var body io.Reader = resp.Body
	if f.MaxByteSize > 0 {
		body = &limitedReader{r: resp.Body, left: f.MaxByteSize}
	}
	return f.Parse(body)
}

// limitedReader returns ErrResponseTooLarge once more than the configured
// number of bytes has been read, rather than silently truncating (which would
// produce a corrupt partial parse).
type limitedReader struct {
	r    io.Reader
	left int64
}

func (l *limitedReader) Read(p []byte) (int, error) {
	n, err := l.r.Read(p)
	l.left -= int64(n)
	if l.left < 0 {
		return n, ErrResponseTooLarge
	}
	return n, err
}

// ParseString parses a feed XML string and into the
// universal feed type.
func (f *Parser) ParseString(feed string) (*Feed, error) {
	return f.Parse(strings.NewReader(feed))
}

func (f *Parser) parseAtomFeed(feed io.Reader) (*Feed, error) {
	af, err := f.ap.Parse(feed)
	if err != nil {
		return nil, err
	}
	result, err := f.atomTrans().Translate(af)
	f.keepOriginal(result, af)
	return result, err
}

func (f *Parser) parseRSSFeed(feed io.Reader) (*Feed, error) {
	rf, err := f.rp.Parse(feed)
	if err != nil {
		return nil, err
	}

	result, err := f.rssTrans().Translate(rf)
	f.keepOriginal(result, rf)
	return result, err
}

func (f *Parser) parseJSONFeed(feed io.Reader) (*Feed, error) {
	jf, err := f.jp.Parse(feed)
	if err != nil {
		return nil, err
	}
	result, err := f.jsonTrans().Translate(jf)
	f.keepOriginal(result, jf)
	return result, err
}

// keepOriginal stashes the source feed on the result when KeepOriginalFeed is
// set. Gating here keeps the Translator interface free of parse options.
func (f *Parser) keepOriginal(result *Feed, original interface{}) {
	if f.KeepOriginalFeed && result != nil {
		result.originalFeed = original
	}
}

// These accessors return a shared default when the corresponding field is
// unset. They must not write back to the Parser: doing so races when one Parser
// is shared across goroutines (a common pattern for crawlers).

func (f *Parser) atomTrans() Translator {
	if f.AtomTranslator != nil {
		return f.AtomTranslator
	}
	return defaultAtomTranslator
}

func (f *Parser) rssTrans() Translator {
	if f.RSSTranslator != nil {
		return f.RSSTranslator
	}
	return defaultRSSTranslator
}

func (f *Parser) jsonTrans() Translator {
	if f.JSONTranslator != nil {
		return f.JSONTranslator
	}
	return defaultJSONTranslator
}

func (f *Parser) httpClient() *http.Client {
	if f.Client != nil {
		return f.Client
	}
	return defaultClient
}
