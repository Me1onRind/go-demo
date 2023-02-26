package configmd

type StaticConfig struct {
}

type DynamicConfig struct {
	IdentifyCode    string                    `yaml:"identify_code"`
	DefaultDB       DBCluster                 `yaml:"default_db"`
	KafkaConfigs    []KafkaConfig             `yaml:"kafka_configs"`
	KafkaJobConfigs map[string]KafkaJobConfig `yaml:"kafka_job_configs"` // key: jobName
}

func (d *DynamicConfig) GetKafkaJobConfig(jobName string) *KafkaJobConfig {
	if v, ok := d.KafkaJobConfigs[jobName]; ok {
		return &v
	}
	return &KafkaJobConfig{}
}

type LocalFileConfig struct {
	Etcd      EtcdConfig `yaml:"etcd"`
	DymCfgKey string     `yaml:"dym_cfg_key"`
}
