package async

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/redis"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
)

type RedisJob[T any] struct {
	jobBase[T]
}

func NewRedisJob[T any](name string, handler func(context.Context, *T) error) JobWorker {
	return &RedisJob[T]{
		jobBase: jobBase[T]{
			JobName: name,
			Handler: handler,
		},
	}
}

func (r *RedisJob[T]) Send(ctx context.Context, msgEntity []byte, sendParam *SendParam) error {
	jobCfg := gconfig.DynamicCfg.RedisJobConfigs[r.JobName]
	logger.CtxInfof(ctx, "Job %s send to redis, redis_name=%s, queue_key=%s", r.Name(), jobCfg.RedisLabel, jobCfg.QueueKey)

	redisClient := redis.GetRedisClient(jobCfg.RedisLabel)
	if err := redisClient.LPush(ctx, jobCfg.QueueKey, msgEntity).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisJob[T]) BackendType() JobBackendType {
	return RedisBackendJob
}
