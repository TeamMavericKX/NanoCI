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
	"github.com/princetheprogrammerbtw/nanoci/internal/auth"
	"github.com/princetheprogrammerbtw/nanoci/internal/config"
	"github.com/princetheprogrammerbtw/nanoci/internal/db"
	"github.com/princetheprogrammerbtw/nanoci/internal/queue"
	"github.com/princetheprogrammerbtw/nanoci/internal/repository/postgres"
	"github.com/princetheprogrammerbtw/nanoci/internal/server/handlers"
	"github.com/princetheprogrammerbtw/nanoci/internal/server/logstream"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Load Config
	cfg, err := config.Load()
	if err != nil {
		zap.L().Fatal("failed to load config", zap.Error(err))
	}

	// Initialize DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DBURL)
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Initialize Redis
	opt, _ := redis.ParseURL(cfg.RedisURL)
	rdb := redis.NewClient(opt)
	defer rdb.Close()

	// Initialize Repositories
	userRepo := postgres.NewUserRepository(pool)
	projectRepo := postgres.NewProjectRepository(pool)
	buildRepo := postgres.NewBuildRepository(pool)
	secretRepo := postgres.NewSecretRepository(pool)

	// Initialize Queue
	q := queue.NewRedisQueue(rdb)

	// Initialize Services
	authService := auth.NewAuthService(cfg, userRepo)
	logManager := logstream.NewLogManager(rdb)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	webhookHandler := handlers.NewWebhookHandler(projectRepo, buildRepo, q)
	projectHandler := handlers.NewProjectHandler(projectRepo)
	buildHandler := handlers.NewBuildHandler(buildRepo)
	secretHandler := handlers.NewSecretHandler(secretRepo, cfg.EncryptionKey)

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

	r.Get("/ws/logs/{buildID}", func(w http.ResponseWriter, r *http.Request) {
		buildID := chi.URLParam(r, "buildID")
		logManager.HandleWS(w, r, buildID)
	})

	r.Route("/api/v1", func(r chi.Router) {
// ...
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", projectHandler.List)
			r.Post("/", projectHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", projectHandler.Get)
				r.Get("/builds", buildHandler.ListByProject)
				r.Get("/secrets", secretHandler.List)
				r.Post("/secrets", secretHandler.Create)
			})
		})
		r.Get("/builds/{id}", buildHandler.Get)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", authHandler.Login)
		r.Get("/callback", authHandler.Callback)
	})

	r.Post("/webhooks/github", webhookHandler.HandleGithub)

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
