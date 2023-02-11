package redis

import (
	"context"
	"errors"
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(cfg *configmd.RedisConfig) (*redis.Client, error) {
	opt := redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
	}
	client := redis.NewClient(&opt)
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	client.AddHook(&redisHook{})
	return client, nil
}

type startTimeKey struct{}

type redisHook struct {
}

func (h *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx = context.WithValue(ctx, startTimeKey{}, time.Now())
	return ctx, nil
}

func (h *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	var key any
	var params []any
	args := cmd.Args()
	if len(args) > 0 {
		key = args[1]
		if len(args) > 1 {
			params = args[2:]
		}
	}
	startTime := ctx.Value(startTimeKey{}).(time.Time)
	duration := time.Since(startTime)

	if cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil) {
		logger.CtxInfof(ctx, "opt:[%s],key:[%s],params:%v,cost:[%s],err:[%s]", cmd.Name(), key, params, duration, cmd.Err())
	} else {
		logger.CtxInfof(ctx, "opt:[%s],key:[%s],params:%v,cost:[%s]", cmd.Name(), key, params, duration)
	}
	return nil
}

func (h *redisHook) BeforeProcessPipeline(ctx context.Context, cmd []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *redisHook) AfterProcessPipeline(ctx context.Context, cmd []redis.Cmder) error {
	return nil
}
