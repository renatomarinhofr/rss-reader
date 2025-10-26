package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client defines the minimal HTTP client interface we rely on.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewDefault returns an http.Client with sane defaults.
func NewDefault(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &http.Client{Timeout: timeout}
}

// FetchBytes performs a GET request and returns the response body bytes.
func FetchBytes(ctx context.Context, client Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		io.Copy(io.Discard, res.Body)
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
