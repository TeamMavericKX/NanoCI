package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/princetheprogrammerbtw/nanoci/internal/config"
	"github.com/princetheprogrammerbtw/nanoci/internal/db"
	"github.com/princetheprogrammerbtw/nanoci/internal/queue"
	"github.com/princetheprogrammerbtw/nanoci/internal/repository/postgres"
	"github.com/princetheprogrammerbtw/nanoci/internal/runner"
	"github.com/princetheprogrammerbtw/nanoci/internal/worker"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// ... (logger and config)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	cfg, err := config.Load()
	if err != nil {
		zap.L().Fatal("failed to load config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize DB
	pool, err := db.NewPool(ctx, cfg.DBURL)
	if err != nil {
		zap.L().Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Initialize Repositories
	buildRepo := postgres.NewBuildRepository(pool)
	projectRepo := postgres.NewProjectRepository(pool)
	secretRepo := postgres.NewSecretRepository(pool)

	// Initialize Runner
	dockerRunner, err := runner.NewDockerRunner()
	if err != nil {
		zap.L().Fatal("failed to initialize docker runner", zap.Error(err))
	}

	// Initialize Executor
	executor := worker.NewExecutor(buildRepo, projectRepo, secretRepo, dockerRunner, cfg.EncryptionKey)

	// Initialize Redis for polling
	opt, _ := redis.ParseURL(cfg.RedisURL)
	rdb := redis.NewClient(opt)
	defer rdb.Close()

	zap.L().Info("worker started, waiting for jobs...")

	for {
		select {
		case <-ctx.Done():
			zap.L().Info("worker shutting down")
			return
		default:
			result, err := rdb.BRPop(ctx, 5*time.Second, "nanoci:jobs").Result()
			if err == redis.Nil {
				continue
			}
			if err != nil {
				zap.L().Error("failed to pop job from redis", zap.Error(err))
				continue
			}

			var job queue.Job
			if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
				zap.L().Error("failed to unmarshal job", zap.Error(err))
				continue
			}

			zap.L().Info("processing job", zap.String("build_id", job.BuildID))
			
			if err := executor.Execute(ctx, job.BuildID); err != nil {
				zap.L().Error("execution failed", zap.String("build_id", job.BuildID), zap.Error(err))
			}
		}
	}
}
