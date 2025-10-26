package clearfeeds_test

import (
	"context"
	"errors"
	"testing"

	"rssreader/internal/domain/feed"
	"rssreader/internal/usecase/clearfeeds"
)

type storeStub struct {
	clearErr error
}

func (s storeStub) Save(ctx context.Context, entry *feed.Feed) error {
	return nil
}

func (s storeStub) ListRecent(ctx context.Context, limit int) ([]feed.Summary, error) {
	return nil, nil
}

func (s storeStub) FindByURL(ctx context.Context, url string) (*feed.Feed, error) {
	return nil, nil
}

func (s storeStub) Clear(ctx context.Context) error {
	return s.clearErr
}

func TestExecuteClearsFeeds(t *testing.T) {
	usecase := clearfeeds.New(storeStub{})
	if err := usecase.Execute(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteRequiresStore(t *testing.T) {
	usecase := clearfeeds.New(nil)
	if err := usecase.Execute(context.Background()); err == nil {
		t.Fatal("expected error when store is nil")
	}
}

func TestExecutePropagatesError(t *testing.T) {
	usecase := clearfeeds.New(storeStub{clearErr: errors.New("db error")})
	if err := usecase.Execute(context.Background()); err == nil {
		t.Fatal("expected error when clear fails")
	}
}
