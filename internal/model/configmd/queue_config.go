package configmd

import (
	"fmt"
	"time"
)

type KafkaJobConfig struct {
	KafkaName       string        `yaml:"kafka_name"`
	Topic           string        `yaml:"topic"`
	Addr            []string      `yaml:"addr"`
	ConsumerGroup   string        `yaml:"consumer_group"`
	ProducerTimeout time.Duration `yaml:"producer_timeout"`
}

func (k *KafkaJobConfig) UniqueKey() string {
	return fmt.Sprintf("%s_%s", k.KafkaName, k.Topic)
}
