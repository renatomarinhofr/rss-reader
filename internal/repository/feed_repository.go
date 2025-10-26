package repository

import (
	"context"

	"rssreader/internal/domain/feed"
)

// FeedFetcher abstracts fetching raw feed data from an external source.
type FeedFetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// FeedStore persists feed snapshots for later retrieval.
type FeedStore interface {
	// Save stores or updates the feed snapshot for the given source URL.
	Save(ctx context.Context, entry *feed.Feed) error
	// ListRecent returns the most recent feeds ordered by fetched at descending.
	ListRecent(ctx context.Context, limit int) ([]feed.Summary, error)
	// FindByURL returns the latest stored feed for a given URL, if any.
	FindByURL(ctx context.Context, url string) (*feed.Feed, error)
	// Clear removes all stored feed snapshots.
	Clear(ctx context.Context) error
}
