package store

import (
	"context"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/go-redis/redis/v8"
)

var (
	RedisPool *redis.Client
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

func NewRedisPoolFromAddr(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    addr,
		DB:      0,
	})
	client.AddHook(&redisHook{})
	return client
}

func NewRedisPool(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  time.Millisecond * cfg.ReadTimeout,
		WriteTimeout: time.Millisecond * cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		IdleTimeout:  time.Second * cfg.IdleTimeout,
		DB:           0,
	})
	client.AddHook(&redisHook{})
	return client
}
