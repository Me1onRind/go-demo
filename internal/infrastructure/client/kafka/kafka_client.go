package kafka

import (
	"context"
	"time"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/Shopify/sarama"
)

type KafkaClient struct {
	producer sarama.SyncProducer
}

func NewKafkaClient(cfg configmd.KafkaConfig) (*KafkaClient, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Version = sarama.MaxVersion
	kafkaCfg.Producer.Return.Successes = true
	kafkaCfg.Producer.Timeout = cfg.ProducerTimeout
	p, err := sarama.NewSyncProducer(cfg.Addr, kafkaCfg)
	if err != nil {
		return nil, err
	}
	return &KafkaClient{
		producer: p,
	}, nil
}

func (k *KafkaClient) SendMessage(ctx context.Context, msg *sarama.ProducerMessage) error {
	startTime := time.Now()
	partitionId, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		logger.CtxErrorf(ctx, "Sender kafka message fail, topic:[%s], msg:[%s], cost:[%s], err:[%s]",
			msg.Topic, msg.Value, time.Since(startTime), err.Error())
		return err
	}
	logger.CtxErrorf(ctx, "Sender kafka message success, topic:[%s], metadata:[%v], headers:[%v], msg:[%s], partition_id:[%d], offset:[%d], cost:[%s]",
		msg.Topic, msg.Metadata, msg.Headers, msg.Value, partitionId, offset, time.Since(startTime))
	return nil
}
