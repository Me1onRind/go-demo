package gclient

import (
	"fmt"

	"github.com/Me1onRind/go-demo/internal/infrastructure/client/etcd"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/kafka"
)

var (
	EtcdClient   etcd.Client
	kafkaClients = map[string]*kafka.KafkaClient{}
)

func GetKafkaClient(name string) (*kafka.KafkaClient, error) {
	if client, ok := kafkaClients[name]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("Kafka client not found, name=%s", name)
}

func RegisterKafkaClient(name string, client *kafka.KafkaClient) error {
	if _, ok := kafkaClients[name]; ok {
		return fmt.Errorf("Register kafka client failed, cause duplicate, name=%s", name)
	}
	kafkaClients[name] = client
	return nil
}

func CleanKafkaClient(name string) {
	delete(kafkaClients, name)
}
