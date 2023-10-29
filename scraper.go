package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/FilipLusnia/rssagg/internal/database"
	"github.com/google/uuid"
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

		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing \"%v\" date: %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.UUID{},
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Failed to create post: %v", err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
