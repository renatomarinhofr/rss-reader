package fetchfeed

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"rssreader/internal/domain/feed"
	"rssreader/internal/repository"
)

// UseCase orchestrates parsing an RSS feed from a given URL.
type UseCase struct {
	fetcher repository.FeedFetcher
	store   repository.FeedStore
	parser  feedParser
	clock   func() time.Time
}

type feedParser interface {
	ParseString(input string) (*gofeed.Feed, error)
}

// New creates a new UseCase instance.
func New(fetcher repository.FeedFetcher, store repository.FeedStore, clock func() time.Time) *UseCase {
	if clock == nil {
		clock = time.Now
	}
	return &UseCase{
		fetcher: fetcher,
		store:   store,
		parser:  gofeed.NewParser(),
		clock:   clock,
	}
}

// Execute downloads and parses the feed for the provided URL.
func (uc *UseCase) Execute(ctx context.Context, url string) (*feed.Feed, error) {
	trimmedURL := strings.TrimSpace(url)
	if trimmedURL == "" {
		return nil, errors.New("url is required")
	}

	raw, err := uc.fetcher.Fetch(ctx, trimmedURL)
	if err != nil {
		if uc.store != nil {
			cached, cacheErr := uc.store.FindByURL(ctx, trimmedURL)
			if cacheErr == nil && cached != nil {
				return cached, nil
			}
			if cacheErr != nil {
				return nil, fmt.Errorf("fetch feed: %v (fallback lookup failed: %w)", err, cacheErr)
			}
		}
		return nil, fmt.Errorf("fetch feed: %w", err)
	}

	parsed, err := uc.parser.ParseString(string(raw))
	if err != nil {
		return nil, fmt.Errorf("parse feed: %w", err)
	}

	fetchedAt := uc.clock()
	result := transformFeed(parsed, uc.clock)
	result.SourceURL = trimmedURL
	result.FetchedAt = fetchedAt

	if uc.store != nil {
		if err := uc.store.Save(ctx, result); err != nil {
			return nil, fmt.Errorf("save feed: %w", err)
		}
	}

	return result, nil
}

func transformFeed(parsed *gofeed.Feed, clock func() time.Time) *feed.Feed {
	if parsed == nil {
		return &feed.Feed{}
	}

	items := make([]feed.Item, 0, len(parsed.Items))
	for _, item := range parsed.Items {
		if item == nil {
			continue
		}

		published := resolvePublishedAt(item, clock)

		items = append(items, feed.Item{
			Title:       strings.TrimSpace(item.Title),
			Link:        strings.TrimSpace(item.Link),
			Description: strings.TrimSpace(item.Description),
			PublishedAt: published,
		})
	}

	return &feed.Feed{
		Title:       strings.TrimSpace(parsed.Title),
		Description: strings.TrimSpace(parsed.Description),
		Link:        strings.TrimSpace(parsed.Link),
		Items:       items,
	}
}

func resolvePublishedAt(item *gofeed.Item, clock func() time.Time) time.Time {
	switch {
	case item.PublishedParsed != nil:
		return *item.PublishedParsed
	case item.UpdatedParsed != nil:
		return *item.UpdatedParsed
	default:
		if clock != nil {
			return clock()
		}
		return time.Time{}
	}
}
