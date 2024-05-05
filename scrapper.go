package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Amaan-Khan14/RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C { //infinite loop
		feeds, err := db.GetNextFeedsTofetch(
			context.Background(),
			int32(concurrency))
		if err != nil {
			log.Print("Error fetching feeds", err)
			continue
		}
		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapFeed(db *database.Queries, wg sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Print("Error marking feed as fetched ", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Print("Error fetching feed ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		descripiton := sql.NullString{}
		if item.Description != "" {
			descripiton.String = item.Description
			descripiton.Valid = true
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Print("Error parsing time ", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: descripiton,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Print("Error creating post ", err)
		}
	}
	log.Printf("Feed %s collected ,%v posts found", feed.Name, len(rssFeed.Channel.Item))
}

// func startScrapping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
// 	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequest, concurrency)
// 	ticker := time.NewTicker(timeBetweenRequest)

// 	for ; ; <-ticker.C {
// 		feeds, err := db.GetNextFeedsTofetch(context.Background(), int32(concurrency))
// 		if err != nil {
// 			log.Println("Couldn't get next feeds to fetch", err)
// 			continue
// 		}
// 		log.Printf("Found %v feeds to fetch!", len(feeds))

// 		wg := &sync.WaitGroup{}
// 		for _, feed := range feeds {
// 			wg.Add(1)
// 			go scrapeFeed(db, wg, feed)
// 		}
// 		wg.Wait()
// 	}
// }

// func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
// 	defer wg.Done()
// 	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
// 	if err != nil {
// 		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
// 		return
// 	}

// 	feedData, err := urlToFeed(feed.Url)
// 	if err != nil {
// 		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
// 		return
// 	}
// 	for _, item := range feedData.Channel.Item {
// 		log.Println("Found post", item.Title)
// 	}
// 	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
// }
