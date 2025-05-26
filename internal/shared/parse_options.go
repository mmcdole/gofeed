package shared

import (
	"net/http"
	"time"
)

// ParseOptions configures how feeds are parsed
type ParseOptions struct {
	// Keep reference to the original format-specific feed
	KeepOriginalFeed bool
	
	// Whether to parse dates (can be disabled for performance)
	ParseDates bool
	
	// Parsing behavior options
	StrictnessOptions StrictnessOptions
	
	// Limit the number of items parsed from the feed
	MaxItems int
	
	// HTTP request configuration for ParseURL
	RequestOptions RequestOptions
}

// StrictnessOptions controls parsing strictness
type StrictnessOptions struct {
	AllowInvalidDates     bool
	AllowMissingRequired  bool
	AllowUnescapedMarkup  bool
}

// RequestOptions configures HTTP requests for ParseURL
type RequestOptions struct {
	UserAgent        string
	Timeout          time.Duration
	IfNoneMatch      string     // ETag for conditional requests
	IfModifiedSince  time.Time  // For conditional requests
	Client           *http.Client
	AuthConfig       interface{} // Will be *Auth from parent package
}

// DefaultParseOptions returns sensible defaults
func DefaultParseOptions() *ParseOptions {
	return &ParseOptions{
		KeepOriginalFeed: false,
		ParseDates: true,
		StrictnessOptions: StrictnessOptions{
			AllowInvalidDates:     true,
			AllowMissingRequired:  true,
			AllowUnescapedMarkup:  true,
		},
		MaxItems: 0, // No limit
		RequestOptions: RequestOptions{
			UserAgent: "gofeed/2.0 (+https://github.com/mmcdole/gofeed)",
			Timeout:   60 * time.Second,
		},
	}
}