package etcd

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/internal/global/config"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
)

func PullAndWatchConfig(ctx context.Context, key string, cfgPointer any) error {
	cfg := &config.LocalFileCfg.Etcd
	client, err := getEtcdClient(ctx, cfg.GetEtcdConfig())
	if err != nil {
		logger.CtxErrorf(ctx, "Get etcd client failed, cause:[%s]", err)
		return err
	}

	getCtx, cancel := context.WithTimeout(ctx, cfg.ReadTimeout)
	defer cancel()
	resp, err := client.Get(getCtx, key)
	if err != nil {
		logger.CtxErrorf(ctx, "Etcd client get value failed, cause:[%s], value=%s", err, key)
		return err
	}

	if len(resp.Kvs) == 0 {
		logger.CtxErrorf(ctx, "Can't not get value, key=%s", key)
		return fmt.Errorf("Can't not get value, key=%s", key)
	}

	item := resp.Kvs[0]
	if err := yaml.Unmarshal(item.Value, cfgPointer); err != nil {
		logger.CtxErrorf(ctx, "Unmarshal yaml failed, cause:[%s], value=%s", err, item.Value)
		return err
	}

	return nil
}

func getEtcdClient(ctx context.Context, cfg *clientv3.Config) (*clientv3.Client, error) {
	client, err := clientv3.New(*cfg)
	if err != nil {
		logger.CtxErrorf(ctx, "Etcd new client failed, err:[%s]", err)
		return nil, err
	}
	return client, nil
}
