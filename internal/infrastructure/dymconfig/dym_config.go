package dymconfig

import (
	"context"
	"strings"

	"github.com/Me1onRind/go-demo/internal/global/config"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/etcd"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

func AssociateEtcd(ctx context.Context, cli etcd.Client, key string, cfgPointer any) error {
	body, err := cli.Get(ctx, key, config.LocalFileCfg.Etcd.ReadTimeout)
	if err != nil {
		return err
	}

	if err := unmarshalBySuffix(key, body, cfgPointer); err != nil {
		return err
	}

	return nil
}

func unmarshalBySuffix(key string, body []byte, cfgPointer any) error {
	var err error
	if strings.HasSuffix(key, ".json") {
		err = jsoniter.Unmarshal(body, cfgPointer)
	} else {
		err = yaml.Unmarshal(body, cfgPointer)
	}

	return err
}
