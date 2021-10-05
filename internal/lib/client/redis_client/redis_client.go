package redis_client

import (
	"context"
	"fmt"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	RedisPool *redis.Client
)

type beginTimeKey struct{}

type redisHook struct {
}

func (r *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	beginTime := time.Now()
	ctx = context.WithValue(ctx, beginTimeKey{}, beginTime)
	return ctx, nil
}

func (r *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	lg := ctm_context.ContextLogger(ctx)
	beginTime := ctx.Value(beginTimeKey{}).(time.Time)
	lg.Info("Redis command excute done", zap.String("detail", cmd.String()), zap.Duration("cost", time.Since(beginTime)))
	if err := cmd.Err(); err != nil {
		lg.Error("Redis command error", zap.String("detail", cmd.String()), zap.Error(err))
	}
	return nil
}

func (r *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (r *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func NewRedisClientFromAddr(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    addr,
		DB:      0,
	})
	client.AddHook(&redisHook{})
	return client
}

func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
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
