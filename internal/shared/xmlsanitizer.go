package shared

import "io"

// NewControlCharFilterReader wraps r and drops the C0 control bytes that are
// illegal in XML (bytes below 0x20 other than tab, LF, and CR). Those byte
// values are identical across UTF-8 and the single-byte encodings feeds use,
// and never occur inside a multi-byte UTF-8 sequence, so filtering them on the
// raw stream is safe regardless of the declared encoding. This lets feeds with
// stray control characters parse instead of failing outright.
func NewControlCharFilterReader(r io.Reader) io.Reader {
	return &controlCharFilter{r: r}
}

type controlCharFilter struct {
	r io.Reader
}

func (c *controlCharFilter) Read(p []byte) (int, error) {
	for {
		n, err := c.r.Read(p)
		w := 0
		for i := 0; i < n; i++ {
			if b := p[i]; b < 0x20 && b != 0x09 && b != 0x0A && b != 0x0D {
				continue
			}
			p[w] = p[i]
			w++
		}
		// Avoid returning (0, nil), which is discouraged for io.Reader: if this
		// read was entirely illegal bytes, read again.
		if w > 0 || err != nil {
			return w, err
		}
	}
}
