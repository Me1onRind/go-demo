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

func CloseGrpcClients() error {
	return grpc_client.CloseGoDemoClient()
}

func InitAsynqClient() error {
	asynqConfig := &config.RemoteConfig.Asynq
	asynq_client.InitAsynqClient(asynqConfig)
	return nil
}

func CloseAsynqClient() error {
	return asynq_client.AsynqClient.Close()
}

func InitEtcdClient() error {
	etcdConfig := &config.LocalConfig.Etcd
	return etcd_client.InitEtcdClient(etcdConfig)
}

func CloseEtcdClient() error {
	return etcd_client.EtcdClient.Close()
}
