package configmd

import "time"

type StaticConfig struct {
}

type DynamicConfig struct {
	IdentifyCode string
}

type LocalFileConfig struct {
	Etcd EtcdConfig `yaml:"etcd"`
}

type EtcdConfig struct {
	Endpoints   []string      `yaml:"endpoints"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
}
