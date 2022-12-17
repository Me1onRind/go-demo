package dymconfig

import (
	"context"
	"fmt"
	"strings"

	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/etcd"
	"github.com/Me1onRind/go-demo/internal/infrastructure/goroutine"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	jsoniter "github.com/json-iterator/go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
)

func AssociateEtcd(ctx context.Context, cli etcd.Client, key string, cfgPointer any) error {
	body, err := cli.Get(ctx, key, gconfig.LocalFileCfg.Etcd.ReadTimeout)
	if err != nil {
		return err
	}

	if err := unmarshalBySuffix(ctx, key, body, cfgPointer); err != nil {
		return err
	}

	callback := func(ctx context.Context, event *clientv3.Event) {
		if event.Type == clientv3.EventTypePut {
			_ = unmarshalBySuffix(ctx, key, event.Kv.Value, cfgPointer)
		}
	}

	goroutine.SafeGo(ctx, func() {
		for {
			cli.Watch(ctx, key, callback)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	})

	return nil
}

func unmarshalBySuffix(ctx context.Context, key string, body []byte, cfgPointer any) error {
	var err error
	if strings.HasSuffix(key, ".json") {
		logger.CtxInfof(ctx, "Use json unmarshal, key:[%s]", key)
		err = jsoniter.Unmarshal(body, cfgPointer)
	} else if strings.HasSuffix(key, ".yaml") || strings.HasSuffix(key, ".yml") {
		logger.CtxInfof(ctx, "Use yaml unmarshal, key:[%s]", key)
		err = yaml.Unmarshal(body, cfgPointer)
	} else {
		err = fmt.Errorf("Key:[%s] has not support suffix", key)
	}

	if err != nil {
		logger.CtxErrorf(ctx, "Unmarshal yaml fail, key:[%s], body:\n%s", key, body)
		return err
	}

	logger.CtxInfof(ctx, "Unmarshal key:[%s] to cfg:[%+v] success, body:\n%s", key, cfgPointer, body)
	return nil
}
