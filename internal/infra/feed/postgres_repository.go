package feed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rssreader/internal/domain/feed"
)

// PostgresStore persists feed snapshots in PostgreSQL.
type PostgresStore struct {
	pool *pgxpool.Pool
}

// NewPostgresStore creates a new Postgres-backed FeedStore and ensures schema exists.
func NewPostgresStore(ctx context.Context, pool *pgxpool.Pool) (*PostgresStore, error) {
	if pool == nil {
		return nil, fmt.Errorf("pool is required")
	}

	store := &PostgresStore{pool: pool}
	if err := store.ensureSchema(ctx); err != nil {
		return nil, fmt.Errorf("ensure schema: %w", err)
	}
	return store, nil
}

func (s *PostgresStore) ensureSchema(ctx context.Context) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS feeds (
	id SERIAL PRIMARY KEY,
	source_url TEXT UNIQUE NOT NULL,
	title TEXT,
	description TEXT,
	link TEXT,
	items JSONB NOT NULL DEFAULT '[]'::jsonb,
	fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

	_, err := s.pool.Exec(ctx, ddl)
	return err
}

// Save upserts the feed snapshot for the given URL.
func (s *PostgresStore) Save(ctx context.Context, entry *feed.Feed) error {
	if entry == nil {
		return fmt.Errorf("feed entry is nil")
	}

	serialized, err := json.Marshal(entry.Items)
	if err != nil {
		return fmt.Errorf("marshal items: %w", err)
	}

	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		return fmt.Errorf("feed source url is required")
	}

	if entry.FetchedAt.IsZero() {
		entry.FetchedAt = time.Now().UTC()
	}

	const query = `
INSERT INTO feeds (source_url, title, description, link, items, fetched_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (source_url)
DO UPDATE SET title = EXCLUDED.title,
              description = EXCLUDED.description,
              link = EXCLUDED.link,
              items = EXCLUDED.items,
              fetched_at = EXCLUDED.fetched_at;
`

	_, err = s.pool.Exec(ctx, query,
		sourceURL,
		entry.Title,
		entry.Description,
		entry.Link,
		serialized,
		entry.FetchedAt,
	)
	if err != nil {
		return fmt.Errorf("save feed: %w", err)
	}

	return nil
}

// ListRecent returns the most recent feeds.
func (s *PostgresStore) ListRecent(ctx context.Context, limit int) ([]feed.Summary, error) {
	if limit <= 0 {
		limit = 10
	}

	const query = `
SELECT source_url, title, description, link, fetched_at
FROM feeds
ORDER BY fetched_at DESC
LIMIT $1;
`

	rows, err := s.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("list feeds: %w", err)
	}
	defer rows.Close()

	var result []feed.Summary
	for rows.Next() {
		var summary feed.Summary
		if err := rows.Scan(&summary.SourceURL, &summary.Title, &summary.Description, &summary.Link, &summary.FetchedAt); err != nil {
			return nil, fmt.Errorf("scan feed: %w", err)
		}
		result = append(result, summary)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

// FindByURL returns the latest feed snapshot for a URL.
func (s *PostgresStore) FindByURL(ctx context.Context, url string) (*feed.Feed, error) {
	const query = `
SELECT source_url, title, description, link, items, fetched_at
FROM feeds
WHERE source_url = $1;
`

	lookupURL := strings.TrimSpace(url)
	if lookupURL == "" {
		return nil, nil
	}

	row := s.pool.QueryRow(ctx, query, lookupURL)

	var (
		sourceURL   string
		title       string
		description string
		link        string
		itemsRaw    []byte
		fetchedAt   time.Time
	)

	if err := row.Scan(&sourceURL, &title, &description, &link, &itemsRaw, &fetchedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find feed: %w", err)
	}

	var items []feed.Item
	if len(itemsRaw) > 0 {
		if err := json.Unmarshal(itemsRaw, &items); err != nil {
			return nil, fmt.Errorf("unmarshal items: %w", err)
		}
	}

	return &feed.Feed{
		SourceURL:   sourceURL,
		Title:       title,
		Description: description,
		Link:        link,
		Items:       items,
		FetchedAt:   fetchedAt,
	}, nil
}

// Clear removes all stored feeds.
func (s *PostgresStore) Clear(ctx context.Context) error {
	if _, err := s.pool.Exec(ctx, `TRUNCATE TABLE feeds;`); err != nil {
		return fmt.Errorf("clear feeds: %w", err)
	}
	return nil
}
