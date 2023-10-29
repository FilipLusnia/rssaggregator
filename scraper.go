package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/FilipLusnia/rssagg/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {

	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
			continue
		}

		waitGroup := &sync.WaitGroup{}
		for _, feed := range feeds {
			waitGroup.Add(1)
			go scrapeFeed(db, waitGroup, feed)
		}
		waitGroup.Wait()
	}
}

func scrapeFeed(db *database.Queries, waitGroup *sync.WaitGroup, feed database.Feed) {
	defer waitGroup.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Printf("Found post %s on feed %s", item.Title, feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
