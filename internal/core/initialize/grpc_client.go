package initialize

import "github.com/Me1onRind/go-demo/internal/core/client/grpc_client"

func InitGrpcClients() error {
	if err := grpc_client.InitGoDemoClient(); err != nil {
		return err
	}
	return nil
}
