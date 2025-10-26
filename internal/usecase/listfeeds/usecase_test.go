package listfeeds_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"rssreader/internal/domain/feed"
	"rssreader/internal/usecase/listfeeds"
)

type storeStub struct {
	feeds []feed.Summary
	err   error
}

func (s storeStub) Save(ctx context.Context, entry *feed.Feed) error {
	return nil
}

func (s storeStub) ListRecent(ctx context.Context, limit int) ([]feed.Summary, error) {
	return s.feeds, s.err
}

func (s storeStub) FindByURL(ctx context.Context, url string) (*feed.Feed, error) {
	return nil, nil
}

func (s storeStub) Clear(ctx context.Context) error {
	return nil
}

func TestExecuteReturnsSummaries(t *testing.T) {
	expected := []feed.Summary{
		{
			SourceURL: "https://example.com",
			Title:     "Example",
			FetchedAt: time.Now(),
		},
	}

	usecase := listfeeds.New(storeStub{feeds: expected})

	result, err := usecase.Execute(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d entries, got %d", len(expected), len(result))
	}
}

func TestExecuteRequiresStore(t *testing.T) {
	usecase := listfeeds.New(nil)
	if _, err := usecase.Execute(context.Background(), 5); err == nil {
		t.Fatal("expected error when store is nil")
	}
}

func TestExecutePropagatesStoreError(t *testing.T) {
	usecase := listfeeds.New(storeStub{err: errors.New("db error")})
	if _, err := usecase.Execute(context.Background(), 5); err == nil {
		t.Fatal("expected error when store fails")
	}
}
