package configmd

type StaticConfig struct {
}

type DynamicConfig struct {
	IdentifyCode    string                    `yaml:"identify_code"`
	DefaultDB       DBCluster                 `yaml:"default_db"`
	KafkaJobConfigs map[string]KafkaJobConfig `yaml:"kafka_job_configs"` // key: jobName
	RedisJobConfigs map[string]RedisJobConfig `yaml:"redis_job_configs"`
}

type LocalFileConfig struct {
	Etcd      EtcdConfig `yaml:"etcd"`
	DymCfgKey string     `yaml:"dym_cfg_key"`
}
