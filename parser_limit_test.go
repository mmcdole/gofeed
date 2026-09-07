package gofeed_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"
)

type limitTestBody struct {
	io.Reader
	chunk  int
	read   int64
	closed bool
}

func (b *limitTestBody) Read(p []byte) (int, error) {
	if len(p) > b.chunk {
		p = p[:b.chunk]
	}
	n, err := b.Reader.Read(p)
	b.read += int64(n)
	return n, err
}

func (b *limitTestBody) Close() error { b.closed = true; return nil }

type limitTestTransport struct{ body *limitTestBody }

func (t limitTestTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: t.body}, nil
}

// Return the error with the last bytes, then EOF on subsequent reads.
type finalErrorReader struct {
	*strings.Reader
	err error
}

func (r finalErrorReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if n > 0 && r.Len() == 0 {
		return n, r.err
	}
	return n, err
}

func TestParseURLMaxByteSizeBoundaries(t *testing.T) {
	for _, format := range []struct{ name, prefix, suffix string }{
		{"rss", `<rss><channel><title>`, `</title></channel></rss>`},
		{"atom", `<feed xmlns="http://www.w3.org/2005/Atom"><title>`, `</title></feed>`},
		{"json", `{"version":"https://jsonfeed.org/version/1.1","title":"`, `","items":[]}`},
	} {
		for _, length := range []int{10, 5900} {
			doc := format.prefix + strings.Repeat("x", length) + format.suffix
			for _, trailing := range []int{0, 9000} {
				input := doc + strings.Repeat(" ", trailing)
				for _, limit := range []int64{0, 4096, int64(len(doc) - 1), int64(len(doc)), int64(len(input)), int64(len(input) + 1), math.MaxInt64} {
					for _, chunk := range []int{1, 4096, 16384} {
						for _, dataEOF := range []bool{false, true} {
							name := fmt.Sprintf("%s/text=%d/trailing=%d/limit=%d/chunk=%d/dataEOF=%t", format.name, length, trailing, limit, chunk, dataEOF)
							t.Run(name, func(t *testing.T) {
								var reader io.Reader = strings.NewReader(input)
								if dataEOF {
									reader = iotest.DataErrReader(reader)
								}
								body := &limitTestBody{Reader: reader, chunk: chunk}
								p := gofeed.NewParser()
								p.MaxByteSize = limit
								p.Client = &http.Client{Transport: limitTestTransport{body}}
								feed, err := p.ParseURL("https://example.com/feed")
								if limit > 0 && int64(len(input)) > limit {
									require.ErrorIs(t, err, gofeed.ErrResponseTooLarge)
									require.Nil(t, feed)
								} else {
									require.NoError(t, err)
									require.NotNil(t, feed)
									require.Equal(t, strings.Repeat("x", length), feed.Title)
								}
								if limit > 0 && limit < math.MaxInt64 {
									require.LessOrEqual(t, body.read, limit+1, "read at most one extra byte to detect overflow")
								}
								require.True(t, body.closed)
							})
						}
					}
				}
			}
		}
	}
}

func TestParseURLMaxByteSizeTrailingReaderError(t *testing.T) {
	boom := errors.New("response read failed")
	for _, prefix := range []string{"<rss><channel><title>", "<feed><title>"} {
		t.Run(prefix, func(t *testing.T) {
			suffix := "</title></channel></rss>"
			if prefix == "<feed><title>" {
				suffix = "</title></feed>"
			}
			doc := prefix + strings.Repeat("x", 5900) + suffix
			for _, withData := range []bool{false, true} {
				t.Run(fmt.Sprintf("withData=%t", withData), func(t *testing.T) {
					var reader io.Reader = io.MultiReader(strings.NewReader(doc), iotest.ErrReader(boom))
					if withData {
						reader = finalErrorReader{strings.NewReader(doc), boom}
					}
					body := &limitTestBody{Reader: reader, chunk: 4096}
					p := gofeed.NewParser()
					p.MaxByteSize = 10000
					p.Client = &http.Client{Transport: limitTestTransport{body}}
					feed, err := p.ParseURLWithContext("https://example.com/feed", context.Background())
					require.ErrorIs(t, err, boom)
					require.Nil(t, feed)
					require.True(t, body.closed)
				})
			}
		})
	}
}
