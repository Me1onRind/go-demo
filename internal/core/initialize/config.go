package initialize

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/client/etcd_client"
	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/utils/env"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func InitLocalConfig(configDir string) func() error {
	return func() error {
		configPath := fmt.Sprintf("%s/%s.yaml", configDir, env.GetEnv("env", "test"))
		configFile, err := os.Open(configPath)
		if err != nil {
			return err
		}

		configBytes, err := ioutil.ReadAll(configFile)
		if err != nil {
			return err
		}

		logger.Logger.Info("Read config file", zap.ByteString("content", configBytes))

		if err := yaml.Unmarshal(configBytes, config.LocalConfig); err != nil {
			return err
		}

		return nil
	}
}

func InitEtcdConfig(key string) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		resp, err := etcd_client.EtcdClient.Get(ctx, key)
		if err != nil {
			return err
		}

		if len(resp.Kvs) == 0 {
			return fmt.Errorf("Can't find key:%s from etcd", key)
		}

		configBytes := resp.Kvs[0].Value
		if err := loadEtcdConfig(configBytes); err != nil {
			return err
		}

		listenEtcdConfigChange(key)

		return nil
	}
}

func loadEtcdConfig(configBytes []byte) error {
	logger.Logger.Info("Read config content from etcd", zap.ByteString("content", configBytes))
	if err := yaml.Unmarshal(configBytes, config.RemoteConfig); err != nil {
		return err
	}
	return nil
}

func listenEtcdConfigChange(key string) {
	go func() {
		for {
			rch := etcd_client.EtcdClient.Watch(context.TODO(), key) // context from main.go
			for wresp := range rch {
				logger.Logger.Info("Etcd config change")
				for _, ev := range wresp.Events {
					configBytes := ev.Kv.Value
					if err := loadEtcdConfig(configBytes); err != nil {
						logger.Logger.Error("Listen etcd config failed", zap.ByteString("content", configBytes),
							zap.String("key", key))
					}
				}
			}
		}
	}()
}
