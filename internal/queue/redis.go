package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Job struct {
	BuildID string `json:"build_id"`
}

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{client: client}
}

func (q *RedisQueue) Enqueue(ctx context.Context, job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.client.LPush(ctx, "nanoci:jobs", data).Err()
}

func (q *RedisQueue) Close() error {
	return q.client.Close()
}
