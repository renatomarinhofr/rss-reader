package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"rssreader/internal/infra/database"
	feedRepo "rssreader/internal/infra/feed"
	"rssreader/internal/infra/httpclient"
	iface "rssreader/internal/interface/http"
	"rssreader/internal/usecase/clearfeeds"
	"rssreader/internal/usecase/fetchfeed"
	"rssreader/internal/usecase/listfeeds"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/rssreader?sslmode=disable"
	}

	pool, err := database.Connect(ctx, dsn, 5)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	store, err := feedRepo.NewPostgresStore(context.Background(), pool)
	if err != nil {
		log.Fatalf("failed to initialise feed store: %v", err)
	}

	client := httpclient.NewDefault(10 * time.Second)
	repository := feedRepo.NewHTTPRepository(client)
	fetchUseCase := fetchfeed.New(repository, store, time.Now)
	listUseCase := listfeeds.New(store)
	clearUseCase := clearfeeds.New(store)
	handler := iface.NewHandler(fetchUseCase, listUseCase, clearUseCase)

	distDir := filepath.Join("web", "dist")

	server := iface.NewServer(addr, func(mux *http.ServeMux) {
		handler.Register(mux)

		// Health check endpoint for readiness probes.
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

		if h := serveStatic(distDir); h != nil {
			mux.Handle("/", h)
		}
	})

	errs := make(chan error, 1)
	go func() {
		errs <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errs:
		if err != nil {
			log.Fatalf("server failed: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("signal received: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown gracefully: %v", err)
		}
	}
}

func serveStatic(distDir string) http.Handler {
	info, err := os.Stat(distDir)
	if err != nil || !info.IsDir() {
		return nil
	}

	fs := http.FileServer(http.Dir(distDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(distDir, strings.TrimPrefix(r.URL.Path, "/"))
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		// Fallback to index for SPA routes.
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	})
}
