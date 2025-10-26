package fetchfeed_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"rssreader/internal/domain/feed"
	"rssreader/internal/usecase/fetchfeed"
)

const sampleFeed = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Example Feed</title>
    <description>Feed description</description>
    <link>https://example.com</link>
    <item>
      <title>First item</title>
      <link>https://example.com/item1</link>
      <description>First description</description>
      <pubDate>Mon, 02 Jan 2006 15:04:05 -0700</pubDate>
    </item>
    <item>
      <title>Second item</title>
      <link>https://example.com/item2</link>
      <description>Second description</description>
    </item>
  </channel>
</rss>`

type fetcherStub struct {
	payload []byte
	err     error
}

func (f fetcherStub) Fetch(ctx context.Context, url string) ([]byte, error) {
	return f.payload, f.err
}

type storeStub struct {
	saved    []*feed.Feed
	saveErr  error
	findFeed *feed.Feed
	findErr  error
}

func (s *storeStub) Save(ctx context.Context, entry *feed.Feed) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.saved = append(s.saved, entry)
	return nil
}

func (s *storeStub) ListRecent(ctx context.Context, limit int) ([]feed.Summary, error) {
	return nil, nil
}

func (s *storeStub) FindByURL(ctx context.Context, url string) (*feed.Feed, error) {
	if s.findErr != nil {
		return nil, s.findErr
	}
	return s.findFeed, nil
}

func (s *storeStub) Clear(ctx context.Context) error {
	return nil
}

func TestExecuteReturnsParsedFeed(t *testing.T) {
	fetcher := fetcherStub{payload: []byte(sampleFeed)}
	now := func() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) }
	store := &storeStub{}

	uc := fetchfeed.New(fetcher, store, now)

	result, err := uc.Execute(context.Background(), "https://example.com/rss")
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if result.Title != "Example Feed" {
		t.Errorf("unexpected title: %q", result.Title)
	}

	if len(result.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result.Items))
	}

	if got, want := result.Items[0].PublishedAt, time.Date(2006, time.January, 2, 15, 4, 5, 0, time.FixedZone("", -7*3600)); !got.Equal(want) {
		t.Errorf("unexpected first published at: %v", got)
	}

	if got := result.Items[1].PublishedAt; got != now() {
		t.Errorf("expected fallback published at to use clock: %v", got)
	}

	if len(store.saved) != 1 {
		t.Fatalf("expected store to be called once, got %d", len(store.saved))
	}
	if savedURL := store.saved[0].SourceURL; savedURL != "https://example.com/rss" {
		t.Errorf("unexpected stored source url: %s", savedURL)
	}
}

func TestExecuteValidatesURL(t *testing.T) {
	uc := fetchfeed.New(fetcherStub{}, &storeStub{}, time.Now)

	if _, err := uc.Execute(context.Background(), "  "); err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestExecutePropagatesFetchError(t *testing.T) {
	expectedErr := errors.New("network down")
	fetcher := fetcherStub{err: expectedErr}
	uc := fetchfeed.New(fetcher, &storeStub{}, time.Now)

	_, err := uc.Execute(context.Background(), "https://example.com/rss")
	if err == nil || !errors.Is(err, expectedErr) {
		t.Fatalf("expected fetch error, got %v", err)
	}
}

func TestExecutePropagatesParseError(t *testing.T) {
	fetcher := fetcherStub{payload: []byte("not xml")}
	uc := fetchfeed.New(fetcher, &storeStub{}, time.Now)

	if _, err := uc.Execute(context.Background(), "https://example.com/rss"); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestExecuteFallsBackToCacheOnFetchFailure(t *testing.T) {
	expected := &feed.Feed{
		SourceURL: "https://example.com/rss",
		Title:     "Cached",
		FetchedAt: time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC),
	}

	store := &storeStub{findFeed: expected}
	fetcher := fetcherStub{err: errors.New("network down")}
	uc := fetchfeed.New(fetcher, store, time.Now)

	result, err := uc.Execute(context.Background(), expected.SourceURL)
	if err != nil {
		t.Fatalf("expected cached feed, got error: %v", err)
	}

	if result != expected {
		t.Fatalf("expected cached feed to be returned")
	}
}

func TestExecuteReturnsErrorWhenStoreFails(t *testing.T) {
	fetcher := fetcherStub{payload: []byte(sampleFeed)}
	store := &storeStub{saveErr: errors.New("db down")}
	uc := fetchfeed.New(fetcher, store, time.Now)

	if _, err := uc.Execute(context.Background(), "https://example.com/rss"); err == nil {
		t.Fatal("expected error when store fails")
	}
}
