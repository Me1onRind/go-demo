package config

var (
	LocalConfig  = &StaticConfig{}
	RemoteConfig = &DynamicConfig{}
)

type DynamicConfig struct {
	ConfigDB MysqlConfig `yaml:"config_db"`
}

type StaticConfig struct {
	Etcd EtcdConfig `yaml:"etcd"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}
