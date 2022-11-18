package initialize

import (
	"context"
	"io"
	"os"

	"github.com/Me1onRind/go-demo/internal/global/gclient"
	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/dymconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"gopkg.in/yaml.v2"
)

type InitHandler func(ctx context.Context) error

func InitFileConfig(filepath string) InitHandler {
	return func(ctx context.Context) error {
		file, err := os.Open(filepath)
		if err != nil {
			logger.Errorf("Open file:[%s] failed, err:[%v]", filepath, err)
			return err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			logger.Errorf("Read file:[%s] failed, err:[%v]", filepath, err)
			return err
		}
		logger.CtxInfof(ctx, "Read file content:\n%s", content)

		if err := yaml.Unmarshal(content, gconfig.LocalFileCfg); err != nil {
			logger.Errorf("Yaml unmarchal file:[%s] failed, content=%s, err:[%v]", filepath, content, err)
			return err
		}
		logger.CtxInfof(ctx, "Unmarshal to cfg:[%+v]", gconfig.LocalFileCfg)

		return nil
	}
}

func InitDynamicConfig() InitHandler {
	return func(ctx context.Context) error {
		return dymconfig.AssociateEtcd(ctx, gclient.EtcdClient, gconfig.LocalFileCfg.DymCfgKey, gconfig.DynamicCfg)
	}
}
