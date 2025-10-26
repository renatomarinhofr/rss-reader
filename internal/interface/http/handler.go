package http

import (
	stdhttp "net/http"

	"rssreader/internal/usecase/clearfeeds"
	"rssreader/internal/usecase/fetchfeed"
	"rssreader/internal/usecase/listfeeds"
)

// Handler bundles HTTP handlers for the API surface.
type Handler struct {
	fetch *fetchfeed.UseCase
	list  *listfeeds.UseCase
	clear *clearfeeds.UseCase
}

// NewHandler wires dependencies.
func NewHandler(fetch *fetchfeed.UseCase, list *listfeeds.UseCase, clear *clearfeeds.UseCase) *Handler {
	return &Handler{fetch: fetch, list: list, clear: clear}
}

// Register mounts the routes on the provided ServeMux.
func (h *Handler) Register(mux *stdhttp.ServeMux) {
	mux.HandleFunc("/api/feed", h.getFeed)
	mux.HandleFunc("/api/feeds/recent", h.handleRecentFeeds)
}
