package gofeed_test

import (
	"sort"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

func TestFeedSort(t *testing.T) {
	oldestItem := &gofeed.Item{
		PublishedParsed: &[]time.Time{time.Unix(0, 0)}[0],
	}
	inbetweenItem := &gofeed.Item{
		PublishedParsed: &[]time.Time{time.Unix(1, 0)}[0],
	}
	newestItem := &gofeed.Item{
		PublishedParsed: &[]time.Time{time.Unix(2, 0)}[0],
	}

	feed := gofeed.Feed{
		Items: []*gofeed.Item{
			newestItem,
			oldestItem,
			inbetweenItem,
		},
	}
	expected := gofeed.Feed{
		Items: []*gofeed.Item{
			oldestItem,
			inbetweenItem,
			newestItem,
		},
	}

	sort.Sort(feed)

	for i, item := range feed.Items {
		if !item.PublishedParsed.Equal(
			*expected.Items[i].PublishedParsed,
		) {
			t.Errorf(
				"Item PublishedParsed = %s; want %s",
				item.PublishedParsed,
				expected.Items[i].PublishedParsed,
			)
		}
	}
}
