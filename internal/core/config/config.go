package config

import "time"

var (
	LocalConfig  = &StaticConfig{}
	RemoteConfig = &DynamicConfig{}
)

type DynamicConfig struct {
	DBs   MysqlConfigs `yaml:"dbs"`
	Redis RedisConfig  `yaml:"redis"`
}

type StaticConfig struct {
	Etcd EtcdConfig `yaml:"etcd"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type MysqlConfigs struct {
	ConfigDB MysqlConfig `yaml:"config_db"`
	DB       MysqlConfig `yaml:"db"`
}

type MysqlConfig struct {
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	DBName       string        `yaml:"dbname"`
	Timeout      string        `yaml:"timeout"`
	ReadTimeout  string        `yaml:"read_timeout"`
	WriteTimeout string        `yaml:"write_timeout"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	MaxIdleTime  time.Duration `yaml:"max_idle_time"`
	MaxLifetime  time.Duration `yaml:"max_lifetime"`
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
