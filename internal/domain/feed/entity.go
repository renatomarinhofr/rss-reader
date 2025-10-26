package feed

import "time"

// Feed represents the RSS feed metadata and entries.
type Feed struct {
	SourceURL   string
	Title       string
	Description string
	Link        string
	Items       []Item
	FetchedAt   time.Time
}

// Item represents a single entry in the RSS feed.
type Item struct {
	Title       string
	Link        string
	Description string
	PublishedAt time.Time
}

// Summary represents persisted metadata for a feed.
type Summary struct {
	SourceURL   string
	Title       string
	Description string
	Link        string
	FetchedAt   time.Time
}
