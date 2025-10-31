package rss

import (
	"context"
	"fmt"
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

	fmt.Println("RSS feed")
	fmt.Printf("Channel title: %s\n", feed.Channel.Title)
	fmt.Printf("Channel link: %s\n", feed.Channel.Link)
	fmt.Printf("Channel Description: %s\n", feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		fmt.Printf("Channel item %d\n", i)
		fmt.Printf("Item Title: %s\n", item.Title)
		fmt.Printf("Item Link: %s\n", item.Link)
		//fmt.Printf("Item Description: %s\n", item.Description)
		fmt.Printf("Item PubDate: %s\n", item.PubDate)
	}
}
