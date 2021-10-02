package initialize

import (
	"github.com/Me1onRind/go-demo/internal/core/client/asynq_client"
	"github.com/Me1onRind/go-demo/internal/core/client/grpc_client"
)

func InitGrpcClients() error {
	if err := grpc_client.InitGoDemoClient(); err != nil {
		return err
	}
	return nil
}

func InitAsynqClient(addr string) func() error {
	return func() error {
		asynq_client.InitAsynqClient(addr)
		return nil
	}
}
