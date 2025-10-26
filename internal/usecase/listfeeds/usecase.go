package listfeeds

import (
	"context"
	"errors"
	"fmt"

	"rssreader/internal/domain/feed"
	"rssreader/internal/repository"
)

// UseCase retrieves recent feed snapshots from storage.
type UseCase struct {
	store repository.FeedStore
}

// New constructs the use case with the required dependencies.
func New(store repository.FeedStore) *UseCase {
	return &UseCase{store: store}
}

// Execute returns a limited list of recent feeds.
func (uc *UseCase) Execute(ctx context.Context, limit int) ([]feed.Summary, error) {
	if uc.store == nil {
		return nil, errors.New("feed store not configured")
	}
	if limit <= 0 {
		limit = 10
	}

	feeds, err := uc.store.ListRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list feeds: %w", err)
	}

	return feeds, nil
}
