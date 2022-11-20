package configmd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

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

type DBCluster struct {
	Label  string     `yaml:"label"`
	Master DBConfig   `yaml:"master"`
	Slaves []DBConfig `yaml:"slaves"`
}

type DBConfig struct {
	DSN string `yaml:"dsn"`
}
