package worker

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisLogWriter struct {
	rdb     *redis.Client
	ctx     context.Context
	buildID string
}

func NewRedisLogWriter(ctx context.Context, rdb *redis.Client, buildID string) *RedisLogWriter {
	return &RedisLogWriter{
		rdb:     rdb,
		ctx:     ctx,
		buildID: buildID,
	}
}

func (w *RedisLogWriter) Write(p []byte) (n int, err error) {
	channel := fmt.Sprintf("logs:%s", w.buildID)
	if err := w.rdb.Publish(w.ctx, channel, string(p)).Err(); err != nil {
		return 0, err
	}
	return len(p), nil
}
