package initialize

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gclient"
	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/etcd"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/mysql"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitEtcdClient() InitHandler {
	return func(ctx context.Context) error {
		var err error
		cfg := clientv3.Config{
			Endpoints:   gconfig.LocalFileCfg.Etcd.Endpoints,
			DialTimeout: gconfig.LocalFileCfg.Etcd.DialTimeout,
		}
		gclient.EtcdClient, err = etcd.NewEtcdClient(&cfg)
		return err
	}
}

func InitMysqlClient() InitHandler {
	return func(ctx context.Context) error {
		_, err := mysql.NewMysqlClusterClient(&gconfig.DynamicCfg.DefaultDB)
		if err != nil {
			return err
		}
		return nil
	}
}
