package clearfeeds

import (
	"context"
	"errors"
	"fmt"

	"rssreader/internal/repository"
)

// UseCase removes stored feed snapshots.
type UseCase struct {
	store repository.FeedStore
}

// New constructs the use case with its dependencies.
func New(store repository.FeedStore) *UseCase {
	return &UseCase{store: store}
}

// Execute clears all stored feeds.
func (uc *UseCase) Execute(ctx context.Context) error {
	if uc.store == nil {
		return errors.New("feed store not configured")
	}
	if err := uc.store.Clear(ctx); err != nil {
		return fmt.Errorf("clear feeds: %w", err)
	}
	return nil
}
