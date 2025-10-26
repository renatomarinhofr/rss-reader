package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server wraps the HTTP server configuration.
type Server struct {
	srv *http.Server
}

// NewServer configures an HTTP server using the provided handler registration function.
func NewServer(addr string, register func(mux *http.ServeMux)) *Server {
	mux := http.NewServeMux()
	register(mux)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &Server{srv: srv}
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	log.Printf("HTTP server listening on %s", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
