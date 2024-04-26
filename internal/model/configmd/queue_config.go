package configmd

import (
	"time"
)

type KafkaJobConfig struct {
	KafkaName       string        `yaml:"kafka_name"`
	Topic           string        `yaml:"topic"`
	Addr            []string      `yaml:"addr"`
	ConsumerGroup   string        `yaml:"consumer_group"`
	ProducerTimeout time.Duration `yaml:"producer_timeout"`
}

type RedisJobConfig struct {
	RedisLabel string `yaml:"redis_label"`
	QueueKey   string `yaml:"queue_key"`
}
