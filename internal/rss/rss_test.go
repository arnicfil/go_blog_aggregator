package rss

import (
	"context"
	"testing"
)

func TestFetchFeed(t *testing.T) {
	ctx := context.Background()
	feed, err := FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		t.Fatalf("Error occuded during fetchfeed: %v", err)
	}

	if feed == nil {
		t.Fatalf("Feed is nil")
	}
}
