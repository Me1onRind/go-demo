package configmd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type StaticConfig struct {
}

type DynamicConfig struct {
	IdentifyCode string `yaml:"identify_code"`
}

type LocalFileConfig struct {
	Etcd      EtcdConfig `yaml:"etcd"`
	DymCfgKey string     `yaml:"dym_cfg_key"`
}

type EtcdConfig struct {
	Endpoints   []string      `yaml:"endpoints"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	ReadTimeout time.Duration `yaml:"read_timeout"`
}

func (e *EtcdConfig) GetReadTimeout() time.Duration {
	if e.ReadTimeout == 0 {
		return time.Second * 2
	}
	return e.ReadTimeout
}

func (e *EtcdConfig) GetEtcdConfig() *clientv3.Config {
	etcdCfg := clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: e.DialTimeout,
	}
	return &etcdCfg
}
