package json

import (
	"encoding/json"
	"io"
	"time"

	validator "gopkg.in/go-playground/validator.v5"
)

//ParseV1 will parse a feed in V1 specification, https://jsonfeed.org/version/1
func ParseV1(feed io.Reader) (*FeedV1, error) {
	f := &FeedV1{}
	if err := json.NewDecoder(feed).Decode(f); err != nil {
		return nil, err
	}

	validate := validator.New("validate", validator.BakedInValidators)

	if errs := validate.Struct(f); errs != nil {
		return nil, errs
	}

	return f, nil
}

//FeedV1 is designed according to specification at https://jsonfeed.org/version/1
type FeedV1 struct {
	Version string `json:"version" validate:"required"`
	Title   string `json:"title" validate:"required"`

	HomePageURL string `json:"home_page_url"` // this should be considered required for feeds on the public web
	FeedURL     string `json:"feed_url"`      // this should be considered required for feeds on the public web

	Description string `json:"description"`
	UserComment string `json:"user_comment"`
	NextURL     string `json:"next_url"`

	Icon    string `json:"icon"`
	Favicon string `json:"favicon"`

	Author *V1Author `json:"author"`

	Expired bool `json:"expired"`

	Hubs []*V1Hub `json:"hubs"`

	Items []*V1Item `json:"items" validate:"required,dive"`

	Attachments []*V1Attachment `json:"attachments" validate:"dive"`
}

//V1Author is the Author object for V1 JSON feeds
type V1Author struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Avatar string `json:"avatar"`
}

//V1Item is the Item object for V1 JSON feeds
type V1Item struct {
	ID          string `json:"id" validate:"required"`
	URL         string `json:"url"`
	ExternalURL string `json:"external_url"`
	Title       string `json:"title"`

	ContentHTML string `json:"content_html"` //either content_html and content_text is required
	ContentText string `json:"content_text"` //either content_html and content_text is required

	Summary       string     `json:"summary"`
	Image         string     `json:"image"`
	BannerImage   string     `json:"banner_image"`
	DatePublished *time.Time `json:"date_published"` // date in RFC 3339 format. (Example: 2010-02-07T14:04:00-05:00.)
	DateModified  *time.Time `json:"date_modified"`  // date in RFC 3339 format. (Example: 2010-02-07T14:04:00-05:00.)

	Author *V1Author `json:"author"`

	Tags []string `json:"tags"`
}

//V1Attachment is the Attachment object for V1 JSON feeds
type V1Attachment struct {
	URL               string  `json:"url" validate:"required"`
	MimeType          string  `json:"mime_type" validate:"required"`
	Title             string  `json:"title"`
	SizeInBytes       float64 `json:"size_in_bytes"`
	DurationInSeconds float64 `json:"duration_in_seconds"`
}

//V1Hub is the Hub object for V1 JSON feeds
type V1Hub struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
