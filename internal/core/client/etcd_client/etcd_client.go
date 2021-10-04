package etcd_client

import (
	"time"

	"github.com/Me1onRind/go-demo/internal/core/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	EtcdClient *clientv3.Client
)

func InitEtcdClient(cfg *config.EtcdConfig) error {
	var err error
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:         cfg.Endpoints,
		DialTimeout:       time.Second * 5,
		DialKeepAliveTime: time.Second * 5,
	})
	return err
}
