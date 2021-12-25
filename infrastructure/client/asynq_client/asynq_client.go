package asynq_client

import (
	"fmt"

	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/hibiken/asynq"
)

var (
	AsynqClient *asynq.Client
)

func InitAsynqClientByRedisAddr(addr string) {
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Network: "tcp",
		Addr:    addr,
	})
}

func InitAsynqClient(cfg *config.AsynqConfig) {
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})
}
