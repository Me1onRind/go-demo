package configmd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	DefaultDBLabel = "default"
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

func (d *DBCluster) GetLabel() string {
	if len(d.Label) == 0 {
		return DefaultDBLabel
	}
	return d.Label
}

type DBConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MinIdleConns int
	MaxConnAge   time.Duration
	IdleTimeout  time.Duration
	PoolTimeout  time.Duration
}
