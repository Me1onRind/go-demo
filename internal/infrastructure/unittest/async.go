package unittest

import (
	"testing"

	"github.com/Me1onRind/go-demo/internal/global/gclient"
	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/kafka"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/redis"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func GetMockKafkaClient(t *testing.T, name string) *kafka.KafkaClient {
	gclient.CleanKafkaClient(name)
	kafkaClient := kafka.NewMockKafkaClient(t)
	_ = gclient.RegisterKafkaClient(name, kafkaClient)
	return kafkaClient
}

func GetMockRedis(t *testing.T, label string) *miniredis.Miniredis {
	redis.CleanRedisClient(label)
	redisServer := miniredis.RunT(t)
	cfg := configmd.RedisConfig{
		Addr: redisServer.Addr(),
	}

	client, err := redis.NewRedisPool(&cfg)
	if err != nil {
		panic(err)
	}

	err = redis.RegisterRedisClient(label, client)
	assert.Empty(t, err)

	return redisServer
}

func SetKafkaJobConfig(jobName, kafkaName, topic string) {
	if gconfig.DynamicCfg.KafkaJobConfigs == nil {
		gconfig.DynamicCfg.KafkaJobConfigs = map[string]configmd.KafkaJobConfig{}
	}
	gconfig.DynamicCfg.KafkaJobConfigs[jobName] = configmd.KafkaJobConfig{
		KafkaName: kafkaName,
		Topic:     topic,
	}
}

func SetRedisJobConfig(jobName, redisName, queueKey string) {
	if gconfig.DynamicCfg.RedisJobConfigs == nil {
		gconfig.DynamicCfg.RedisJobConfigs = map[string]configmd.RedisJobConfig{}
	}
	gconfig.DynamicCfg.RedisJobConfigs[jobName] = configmd.RedisJobConfig{
		RedisLabel: redisName,
		QueueKey:   queueKey,
	}
}
