package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"rssreader/internal/domain/feed"
)

func (h *Handler) getFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	url := r.URL.Query().Get("url")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	feed, err := h.fetch.Execute(ctx, url)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, toResponse(feed))
}

func (h *Handler) handleRecentFeeds(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getRecentFeeds(w, r)
	case http.MethodDelete:
		h.clearRecentFeeds(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getRecentFeeds(w http.ResponseWriter, r *http.Request) {
	if h.list == nil {
		writeError(w, http.ErrNotSupported)
		return
	}

	ctx := r.Context()
	feeds, err := h.list.Execute(ctx, 10)
	if err != nil {
		writeError(w, err)
		return
	}

	response := make([]feedSummaryResponse, 0, len(feeds))
	for _, entry := range feeds {
		response = append(response, feedSummaryResponse{
			SourceURL:   entry.SourceURL,
			Title:       entry.Title,
			Description: entry.Description,
			Link:        entry.Link,
			FetchedAt:   entry.FetchedAt,
		})
	}

	writeJSON(w, recentFeedsResponse{Feeds: response})
}

func (h *Handler) clearRecentFeeds(w http.ResponseWriter, r *http.Request) {
	if h.clear == nil {
		writeError(w, http.ErrNotSupported)
		return
	}

	if err := h.clear.Execute(r.Context()); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type errorResponse struct {
	Error string `json:"error"`
}

type feedResponse struct {
	SourceURL   string         `json:"sourceUrl"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Link        string         `json:"link"`
	Items       []feedItemResp `json:"items"`
	FetchedAt   time.Time      `json:"fetchedAt"`
}

type feedItemResp struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
}

type recentFeedsResponse struct {
	Feeds []feedSummaryResponse `json:"feeds"`
}

type feedSummaryResponse struct {
	SourceURL   string    `json:"sourceUrl"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	FetchedAt   time.Time `json:"fetchedAt"`
}

func toResponse(f *feed.Feed) feedResponse {
	items := make([]feedItemResp, 0, len(f.Items))
	for _, item := range f.Items {
		items = append(items, feedItemResp{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			PublishedAt: item.PublishedAt,
		})
	}

	return feedResponse{
		SourceURL:   f.SourceURL,
		Title:       f.Title,
		Description: f.Description,
		Link:        f.Link,
		Items:       items,
		FetchedAt:   f.FetchedAt,
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		writeError(w, err)
	}
}
