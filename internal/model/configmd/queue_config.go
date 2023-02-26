package configmd

import "time"

type KafkaConfig struct {
	Name            string        `yaml:"name"`
	Addr            []string      `yaml:"addr"`
	ProducerTimeout time.Duration `yaml:"producer_timeout"`
}

type KafkaJobConfig struct {
	KafkaName string `yaml:"kafka_name"`
	Topic     string `yaml:"topic"`
}
