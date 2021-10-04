package initialize

import (
	"github.com/Me1onRind/go-demo/internal/core/client/asynq_client"
	"github.com/Me1onRind/go-demo/internal/core/client/etcd_client"
	"github.com/Me1onRind/go-demo/internal/core/client/grpc_client"
	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/internal/lib/register"
)

func InitGrpcClients() error {
	resolver, err := register.GrpcResolvers()
	if err != nil {
		return err
	}
	if err := grpc_client.InitGoDemoClient(register.DialTarget("go-demo"), resolver); err != nil {
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

func InitEtcdClient() error {
	etcdConfig := &config.LocalConfig.Etcd
	return etcd_client.InitEtcdClient(etcdConfig)
}
