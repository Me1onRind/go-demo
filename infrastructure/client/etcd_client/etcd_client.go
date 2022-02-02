package etcd_client

import (
	"time"

	//"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitEtcdClient(cfg *config.EtcdConfig) (*clientv3.Client, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:         cfg.Endpoints,
		DialTimeout:       time.Second * 5,
		DialKeepAliveTime: time.Second * 5,
	})
	return etcdClient, err
}
