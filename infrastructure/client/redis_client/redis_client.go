package redis_client

import (
	"context"
	"time"

	"github.com/Me1onRind/go-demo/config"
	"github.com/Me1onRind/go-demo/infrastructure/logger"
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
	beginTime := ctx.Value(beginTimeKey{}).(time.Time)
	logger.CtxInfo(ctx, "Redis command excute done", zap.String("detail", cmd.String()), zap.Duration("cost", time.Since(beginTime)))
	if err := cmd.Err(); err != nil && err != redis.Nil {
		logger.CtxError(ctx, "Redis command error", zap.String("detail", cmd.String()), zap.Error(err))
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
		Addr:         cfg.Addr,
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
