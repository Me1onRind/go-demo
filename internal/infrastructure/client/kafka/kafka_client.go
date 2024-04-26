package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
)

type KafkaClient struct {
	producer     sarama.SyncProducer
	MockProducer *mocks.SyncProducer
}

func NewKafkaClient(cfg configmd.KafkaJobConfig) (*KafkaClient, error) {
	kafkaCfg := getBaseProducerConfig()
	kafkaCfg.Producer.Timeout = cfg.ProducerTimeout
	p, err := sarama.NewSyncProducer(cfg.Addr, kafkaCfg)
	if err != nil {
		return nil, err
	}
	return &KafkaClient{
		producer: p,
	}, nil
}

func NewMockKafkaClient(t *testing.T) *KafkaClient {
	mock := mocks.NewSyncProducer(t, mocks.NewTestConfig())
	mock.SetDefaultPartitions(4)
	return &KafkaClient{
		MockProducer: mock,
	}
}

func getBaseProducerConfig() *sarama.Config {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Version = sarama.MaxVersion
	kafkaCfg.Producer.Return.Successes = true
	return kafkaCfg
}

type Option func(*sarama.ProducerMessage)

func (k *KafkaClient) SendMessage(ctx context.Context, topic string, body []byte, opts ...Option) (int32, int64, error) {
	startTime := time.Now()

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(body),
		Timestamp: time.Now(),
	}

	for _, opt := range opts {
		opt(msg)
	}

	_, _, err := k.getProducer().SendMessage(msg)
	if err != nil {
		logger.CtxErrorf(ctx, "Sender kafka message fail, topic=%s, msg=%s, cost=%s, err=%s",
			msg.Topic, msg.Value, time.Since(startTime), err.Error())
		return 0, 0, err
	}
	logger.CtxInfof(ctx, "Sender kafka message success, topic=%s, headers=%s, key=%s, msg=%s, partition_id=%d, offset=%d, cost=%s",
		msg.Topic, msg.Headers, msg.Key, msg.Value, msg.Partition, msg.Offset, time.Since(startTime))
	return msg.Partition, msg.Offset, nil
}

func (k *KafkaClient) getProducer() sarama.SyncProducer {
	if k.producer != nil {
		return k.producer
	}
	return k.MockProducer
}

func PartitionKey(key string) Option {
	return func(msg *sarama.ProducerMessage) {
		if len(key) == 0 {
			return
		}
		msg.Key = sarama.StringEncoder(key)
	}
}

func WithHeaers(headers map[string]string) Option {
	return func(msg *sarama.ProducerMessage) {
		for key, value := range headers {
			msg.Headers = append(msg.Headers, sarama.RecordHeader{
				Key:   []byte(key),
				Value: []byte(value),
			})
		}
	}
}
