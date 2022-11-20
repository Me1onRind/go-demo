package configmd

type StaticConfig struct {
}

type DynamicConfig struct {
	IdentifyCode string `yaml:"identify_code"`
}

type LocalFileConfig struct {
	Etcd      EtcdConfig `yaml:"etcd"`
	DymCfgKey string     `yaml:"dym_cfg_key"`
}
