package json_test

import (
	"fmt"
	"strings"
	"time"

	json_feed "github.com/francoishill/gofeed/json"
	validator "gopkg.in/go-playground/validator.v5"
)

func ExamplePodCast() {
	content := `
		{
			"version": "https://jsonfeed.org/version/1",
			"user_comment": "This is a podcast feed. You can add this feed to your podcast client using the following URL: http://therecord.co/feed.json",
			"title": "The Record",
			"home_page_url": "http://therecord.co/",
			"feed_url": "http://therecord.co/feed.json",
			"items": [
				{
					"id": "http://therecord.co/chris-parrish",
					"title": "Special #1 - Chris Parrish",
					"url": "http://therecord.co/chris-parrish",
					"content_text": "Chris has worked at Adobe and as a founder of Rogue Sheep, which won an Apple Design Award for Postage. Chris’s new company is Aged & Distilled with Guy English — which shipped Napkin, a Mac app for visual collaboration. Chris is also the co-host of The Record. He lives on Bainbridge Island, a quick ferry ride from Seattle.",
					"content_html": "Chris has worked at <a href=\"http://adobe.com/\">Adobe</a> and as a founder of Rogue Sheep, which won an Apple Design Award for Postage. Chris’s new company is Aged & Distilled with Guy English — which shipped <a href=\"http://aged-and-distilled.com/napkin/\">Napkin</a>, a Mac app for visual collaboration. Chris is also the co-host of The Record. He lives on <a href=\"http://www.ci.bainbridge-isl.wa.us/\">Bainbridge Island</a>, a quick ferry ride from Seattle.",
					"summary": "Brent interviews Chris Parrish, co-host of The Record and one-half of Aged & Distilled.",
					"date_published": "2014-05-09T14:04:00-07:00",
					"attachments": [
						{
							"url": "http://therecord.co/downloads/The-Record-sp1e1-ChrisParrish.m4a",
							"mime_type": "audio/x-m4a",
							"size_in_bytes": 89970236,
							"duration_in_seconds": 6629
						}
					]
				}
			]
		}
	`

	feed, err := json_feed.ParseV1(strings.NewReader(content))
	if err != nil {
		panic(err)
	}

	fmt.Println(feed.Version)
	fmt.Println(feed.UserComment)
	fmt.Println(feed.Title)
	fmt.Println(feed.HomePageURL)
	fmt.Println(feed.FeedURL)
	fmt.Println(len(feed.Items))
	fmt.Println(feed.Items[0].ID)
	fmt.Println(feed.Items[0].Title)
	fmt.Println(feed.Items[0].URL)
	fmt.Println(feed.Items[0].ContentText)
	fmt.Println(feed.Items[0].ContentHTML)
	fmt.Println(feed.Items[0].Summary)
	fmt.Println(feed.Items[0].DatePublished.Format(time.RFC3339))

	// Output: https://jsonfeed.org/version/1
	// This is a podcast feed. You can add this feed to your podcast client using the following URL: http://therecord.co/feed.json
	// The Record
	// http://therecord.co/
	// http://therecord.co/feed.json
	// 1
	// http://therecord.co/chris-parrish
	// Special #1 - Chris Parrish
	// http://therecord.co/chris-parrish
	// Chris has worked at Adobe and as a founder of Rogue Sheep, which won an Apple Design Award for Postage. Chris’s new company is Aged & Distilled with Guy English — which shipped Napkin, a Mac app for visual collaboration. Chris is also the co-host of The Record. He lives on Bainbridge Island, a quick ferry ride from Seattle.
	// Chris has worked at <a href="http://adobe.com/">Adobe</a> and as a founder of Rogue Sheep, which won an Apple Design Award for Postage. Chris’s new company is Aged & Distilled with Guy English — which shipped <a href="http://aged-and-distilled.com/napkin/">Napkin</a>, a Mac app for visual collaboration. Chris is also the co-host of The Record. He lives on <a href="http://www.ci.bainbridge-isl.wa.us/">Bainbridge Island</a>, a quick ferry ride from Seattle.
	// Brent interviews Chris Parrish, co-host of The Record and one-half of Aged & Distilled.
	// 2014-05-09T14:04:00-07:00
}

func ExampleMicroblog() {
	content := `
		{
			"version": "https://jsonfeed.org/version/1",
			"user_comment": "This is a microblog feed. You can add this to your feed reader using the following URL: https://example.org/feed.json",
			"title": "Brent Simmons’s Microblog",
			"home_page_url": "https://example.org/",
			"feed_url": "https://example.org/feed.json",
			"author": {
				"name": "Brent Simmons",
				"url": "http://example.org/",
				"avatar": "https://example.org/avatar.png"
			},
			"items": [
				{
					"id": "2347259",
					"url": "https://example.org/2347259",
					"content_text": "Cats are neat. \n\nhttps://example.org/cats",
					"date_published": "2016-02-09T14:22:00-07:00"
				}
			]
		}
	`

	feed, err := json_feed.ParseV1(strings.NewReader(content))
	if err != nil {
		panic(err)
	}

	fmt.Println(feed.Version)
	fmt.Println(feed.UserComment)
	fmt.Println(feed.Title)
	fmt.Println(feed.HomePageURL)
	fmt.Println(feed.FeedURL)
	fmt.Println(feed.Author.Name)
	fmt.Println(feed.Author.URL)
	fmt.Println(feed.Author.Avatar)
	fmt.Println(len(feed.Items))
	fmt.Println(feed.Items[0].ID)
	fmt.Println(feed.Items[0].URL)
	fmt.Println(strings.Replace(feed.Items[0].ContentText, "\n", "\\n", -1)) //escape new lines in text
	fmt.Println(feed.Items[0].DatePublished.Format(time.RFC3339))

	// Output: https://jsonfeed.org/version/1
	// This is a microblog feed. You can add this to your feed reader using the following URL: https://example.org/feed.json
	// Brent Simmons’s Microblog
	// https://example.org/
	// https://example.org/feed.json
	// Brent Simmons
	// http://example.org/
	// https://example.org/avatar.png
	// 1
	// 2347259
	// https://example.org/2347259
	// Cats are neat. \n\nhttps://example.org/cats
	// 2016-02-09T14:22:00-07:00
}

func ExampleMissingRequiredInMainFeed() {
	content := `{}`

	_, err := json_feed.ParseV1(strings.NewReader(content))
	if err != nil {
		validationErrs, ok := err.(*validator.StructErrors)
		if !ok {
			panic(err)
		}

		fmt.Println("ERRORS:", len(validationErrs.Flatten()))

		fmt.Println(validationErrs.Errors["Version"].Error())
		fmt.Println(validationErrs.Errors["Title"].Error())
		fmt.Println(validationErrs.Errors["Items"].Error())
	}

	// Output: ERRORS: 3
	// Field validation for "Version" failed on the "required" tag
	// Field validation for "Title" failed on the "required" tag
	// Field validation for "Items" failed on the "required" tag
}

func ExampleMissingRequiredInItem() {
	content := `{
		"version": "https://jsonfeed.org/version/1",
		"title": "The Record",
		"items": [
			{}
		]
	}`

	_, err := json_feed.ParseV1(strings.NewReader(content))
	if err != nil {
		validationErrs, ok := err.(*validator.StructErrors)
		if !ok {
			panic(err)
		}

		fmt.Println("ERRORS:", len(validationErrs.Flatten()))

		fmt.Println("Items ERRORS:", len(validationErrs.Errors["Items"].Flatten()))
		fmt.Println(validationErrs.Errors["Items"].Flatten()["[0].V1Item.ID"])
	}

	// Output: ERRORS: 1
	// Items ERRORS: 1
	// Field validation for "ID" failed on the "required" tag
}

func ExampleMissingRequiredInAttachment() {
	content := `{
		"version": "https://jsonfeed.org/version/1",
		"title": "The Record",
		"items": [
			{
				"id": "http://therecord.co/chris-parrish"
			}
		],
		"attachments": [
			{}
		]
	}`

	_, err := json_feed.ParseV1(strings.NewReader(content))
	if err != nil {
		validationErrs, ok := err.(*validator.StructErrors)
		if !ok {
			panic(err)
		}

		fmt.Println("ERRORS:", len(validationErrs.Flatten()))

		fmt.Println("Attachments ERRORS:", len(validationErrs.Errors["Attachments"].Flatten()))
		fmt.Println(validationErrs.Errors["Attachments"].Flatten()["[0].V1Attachment.URL"])
		fmt.Println(validationErrs.Errors["Attachments"].Flatten()["[0].V1Attachment.MimeType"])
	}

	// Output: ERRORS: 2
	// Attachments ERRORS: 2
	// Field validation for "URL" failed on the "required" tag
	// Field validation for "MimeType" failed on the "required" tag
}
