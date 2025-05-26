package gofeed

import (
	"github.com/mmcdole/gofeed/v2/internal/shared"
)

// ParseOptions re-exports the shared type
type ParseOptions = shared.ParseOptions

// StrictnessOptions re-exports the shared type  
type StrictnessOptions = shared.StrictnessOptions

// RequestOptions re-exports the shared type
type RequestOptions = shared.RequestOptions

// DefaultParseOptions returns sensible defaults
func DefaultParseOptions() *ParseOptions {
	opts := shared.DefaultParseOptions()
	// Fix the AuthConfig type to use our Auth type
	if opts.RequestOptions.AuthConfig == nil {
		opts.RequestOptions.AuthConfig = (*Auth)(nil)
	}
	return opts
}