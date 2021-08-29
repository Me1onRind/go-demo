package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisHook struct {
}

func (r *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (r *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (r *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (r *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func NewRedisPool(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         addr,
		ReadTimeout:  time.Millisecond * 500,
		WriteTimeout: time.Millisecond * 500,
		PoolSize:     5,
		MinIdleConns: 0,
		IdleTimeout:  time.Second * 5,
		DB:           0,
	})
	client.AddHook(&redisHook{})
	return client
}
