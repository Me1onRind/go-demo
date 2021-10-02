package asynq_client

import (
	"github.com/hibiken/asynq"
)

var (
	AsynqClient *asynq.Client
)

func InitAsynqClient(addr string) {
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Network: "tcp",
		Addr:    addr,
	})
}
