package async

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gclient"
	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/kafka"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
)

const (
	KafkaJobNameKey = "job_name"
)

type KafkaJob[T any] struct {
	jobBase[T]
}

func NewKafkaJob[T any](name string, handler func(context.Context, *T) error) JobWorker {
	return &KafkaJob[T]{
		jobBase: jobBase[T]{
			JobName: name,
			Handler: handler,
		},
	}
}

func (j *KafkaJob[T]) Send(ctx context.Context, msgEntity []byte, sendParam *SendParam) error {
	jobCfg := gconfig.DynamicCfg.KafkaJobConfigs[j.JobName]
	logger.CtxInfof(ctx, "Job %s send to kafka, kafka_name=%s, topic=%s", j.Name(), jobCfg.KafkaName, jobCfg.Topic)
	client, err := gclient.GetKafkaClient(jobCfg.KafkaName)
	if err != nil {
		return err
	}
	if _, _, err := client.SendMessage(ctx, jobCfg.Topic, msgEntity,
		kafka.PartitionKey(sendParam.Key),
		kafka.WithHeaers(map[string]string{
			KafkaJobNameKey: j.Name(),
		}),
	); err != nil {
		return err
	}
	return nil
}

func (j *KafkaJob[T]) BackendType() JobBackendType {
	return KafkaBackendJob
}
