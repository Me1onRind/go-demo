package etcd

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func PullAndWatchConfig(ctx context.Context, cfg *configmd.EtcdConfig, key string, cfgPointer any) error {
	etcdCfg := clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
	}

	client, err := clientv3.New(etcdCfg)
	if err != nil {
		logger.CtxErrorf(ctx, "Etcd new client failed, err:[%v]", err)
		return err
	}

	resp, err := client.Get(ctx, key)
	if err != nil {
		logger.CtxErrorf(ctx, "Etcd client get failed, key=%s, err:[%v]", key, err)
		return err
	}

	fmt.Println(resp.Kvs)

	return nil
}
