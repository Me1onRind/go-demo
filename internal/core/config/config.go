package config

import (
	"time"
)

var (
	LocalConfig  = &StaticConfig{}
	RemoteConfig = &DynamicConfig{}
)

type DynamicConfig struct {
	DBs   MysqlConfigs `yaml:"dbs"`
	Redis RedisConfig  `yaml:"redis"`
	Asynq AsynqConfig  `yaml:"asynd"`
}

type StaticConfig struct {
	Etcd EtcdConfig `yaml:"etcd"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type MysqlConfigs struct {
	ConfigDB MysqlResources `yaml:"config_db"`
	DB       MysqlResources `yaml:"db"`
}

type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idel_conns"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type AsynqConfig struct {
	Redis RedisConfig `yaml:"redis"`
}
