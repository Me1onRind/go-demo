package unittest

import (
	"testing"

	"github.com/Me1onRind/go-demo/internal/global/gclient"
	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/kafka"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
)

func GetMockKafkaClient(t *testing.T, name string) *kafka.KafkaClient {
	gclient.CleanKafkaClient(name)
	kafkaClient := kafka.NewMockKafkaClient(t)
	_ = gclient.RegisterKafkaClient(name, kafkaClient)
	return kafkaClient
}

func SetKafkaJobConfig(jobname, kafkaName, topic string) {
	if gconfig.DynamicCfg.KafkaJobConfigs == nil {
		gconfig.DynamicCfg.KafkaJobConfigs = map[string]configmd.KafkaJobConfig{}
	}
	gconfig.DynamicCfg.KafkaJobConfigs[jobname] = configmd.KafkaJobConfig{
		KafkaName: kafkaName,
		Topic:     topic,
	}
}
