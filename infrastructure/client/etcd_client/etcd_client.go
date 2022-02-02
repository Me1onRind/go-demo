package etcd_client

import (
	"time"

	//"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	EtcdClient *clientv3.Client
)

func InitEtcdClient() error {
	var err error
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:         config.LocalConfig.Etcd.Endpoints,
		DialTimeout:       time.Second * 5,
		DialKeepAliveTime: time.Second * 5,
	})
	return err
}
