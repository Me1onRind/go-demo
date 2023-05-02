package async

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/internal/global/gclient"
	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/global/gerror"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/kafka"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	jsoniter "github.com/json-iterator/go"
)

const (
	KafkaJobNameKey = "job_name"
)

type KafkaJob[T any] struct {
	jobBase[T]
}

type KafkaJobEntity struct {
	JobName string `json:"job_name"`
	Content any    `json:"content"`
}

func NewKafkaJob[T any](name string, handler func(context.Context, *T) error) Job {
	return &KafkaJob[T]{
		jobBase: jobBase[T]{
			JobName: name,
			Handler: handler,
		},
	}
}

func (j *KafkaJob[T]) Name() string {
	return j.JobName
}

func (j *KafkaJob[T]) Send(ctx context.Context, protocol any, opts ...Option) error {
	if _, ok := protocol.(*T); !ok {
		errMsg := fmt.Sprintf("Job:[%s] send fail, protocol:[%+v] is not match register protocol", j.JobName, protocol)
		logger.CtxErrorf(ctx, errMsg)
		return gerror.SendJobError.With(errMsg)
	}

	sendParam := &SendParam{}
	for _, opt := range opts {
		opt(sendParam)
	}

	body, err := jsoniter.Marshal(protocol)
	if err != nil {
		return err
	}

	jobCfg := gconfig.DynamicCfg.KafkaJobConfigs[j.JobName]
	logger.CtxInfof(ctx, "Job:[%s] send to kafka:[%s], topic:[%s]", j.Name(), jobCfg.KafkaName, jobCfg.Topic)
	client, err := gclient.GetKafkaClient(ctx, jobCfg.KafkaName)
	if err != nil {
		return err
	}
	if _, _, err := client.SendMessage(ctx, jobCfg.Topic, body,
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
