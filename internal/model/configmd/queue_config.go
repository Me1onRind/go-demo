package configmd

import "time"

type KafkaConfig struct {
	Addr            []string      `yaml:"addr"`
	ProducerTimeout time.Duration `yaml:"producer_timeout"`
}
