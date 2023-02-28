package gclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/Me1onRind/go-demo/internal/infrastructure/client/etcd"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/kafka"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
)

var (
	EtcdClient   etcd.Client
	kafkaClients = map[string]*kafka.KafkaClient{}
)

func GetKafkaClient(ctx context.Context, name string) (*kafka.KafkaClient, error) {
	if client, ok := kafkaClients[name]; ok {
		return client, nil
	}
	errMsg := fmt.Sprintf("Kafka client not found, name:[%s]", name)
	logger.CtxErrorf(ctx, errMsg)
	return nil, errors.New(errMsg)
}

func RegisterKafkaClient(name string, client *kafka.KafkaClient) error {
	if _, ok := kafkaClients[name]; ok {
		return fmt.Errorf("Register kafka client failed, name:[%s], case:[duplicate]", name)
	}
	kafkaClients[name] = client
	return nil
}

func CleanKafkaClient(name string) {
	delete(kafkaClients, name)
}
