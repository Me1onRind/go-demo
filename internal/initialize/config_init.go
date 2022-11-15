package initialize

import (
	"io/ioutil"
	"os"

	"github.com/Me1onRind/go-demo/internal/global/config"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"gopkg.in/yaml.v2"
)

type InitHandle func() error

func InitEtcdConfig() {
}

func InitFileConfig(filepath string) InitHandle {
	return func() error {
		file, err := os.Open(filepath)
		if err != nil {
			logger.Errorf("Open file:[%s] failed, err:[%v]", filepath, err)
			return err
		}
		defer file.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Errorf("Read file:[%s] failed, err:[%v]", filepath, err)
			return err
		}

		if err := yaml.Unmarshal(content, config.LocalFileCfg); err != nil {
			logger.Errorf("Yaml unmarchal file:[%s] failed, content=%s, err:[%v]", filepath, content, err)
			return err
		}

		return nil
	}
}
