package feed

import (
	"context"

	"rssreader/internal/infra/httpclient"
)

// HTTPRepository downloads raw feed content via HTTP.
type HTTPRepository struct {
	client httpclient.Client
}

// NewHTTPRepository wires a new HTTPRepository instance.
func NewHTTPRepository(client httpclient.Client) *HTTPRepository {
	return &HTTPRepository{client: client}
}

// Fetch retrieves the feed bytes from the given URL.
func (r *HTTPRepository) Fetch(ctx context.Context, url string) ([]byte, error) {
	return httpclient.FetchBytes(ctx, r.client, url)
}
