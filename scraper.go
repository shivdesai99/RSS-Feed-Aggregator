package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/RSS-Feed-Aggregator/internal/database"
	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scrapping on %v goroutines with %v time between requests", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C { // Start when ran, then wait ticker
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("Error getting feeds to fetch: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
	
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

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

	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := robustParseDate(item.PubDate)
		if err != nil {
			log.Printf("Error parsing time: %v at time %v", err, item.PubDate)
			continue
		}
		db.CreatePost(context.Background(), database.CreatePostParams {
				ID: uuid.New(),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Title: item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url: item.Link,
				FeedID: feed.ID,
				})
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key") {
					continue
				}
				log.Printf("Error creating post: %v", err)
			}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Items))
}

func robustParseDate(dateStr string) (time.Time, error) {
    pubAt, err := dateparse.ParseAny(dateStr)
    if err != nil {
        return time.Time{}, err
    }
    return pubAt, nil
}