package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/princetheprogrammerbtw/nanoci/internal/config"
	"go.uber.org/zap"
)

func main() {
	// Initialize Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Load Config
	cfg, err := config.Load()
	if err != nil {
		zap.L().Fatal("failed to load config", zap.Error(err))
	}

	// Setup Router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Server setup
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful Shutdown
	go func() {
		zap.L().Info("starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen and serve error", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	zap.L().Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("forced shutdown", zap.Error(err))
	}

	zap.L().Info("server exited gracefully")
}
