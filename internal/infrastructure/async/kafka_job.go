package async

import (
	"context"
	"errors"
	"fmt"

	"github.com/Me1onRind/go-demo/internal/global/gclient"
	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/client/kafka"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	jsoniter "github.com/json-iterator/go"
)

type KafkaJob[T any] struct {
	JobName  string
	Protocol *T
	Handler  func(context.Context, *T) error
}

type KafkaJobEntity struct {
	JobName string `json:"job_name"`
	Content any    `json:"content"`
}

func NewKafkaJob[T any](name, kafkaName, topic string, protocol T, handler func(*context.Context, *T) error) Job {
	return &KafkaJob[T]{
		JobName:  name,
		Protocol: &protocol,
	}
}

func (j *KafkaJob[T]) Name() string {
	return j.JobName
}

func (j *KafkaJob[T]) Send(ctx context.Context, protocol any, opts ...Option) error {
	if _, ok := protocol.(*T); !ok {
		errMsg := fmt.Sprintf("Job:[%s] send fail, protocol:[%+v] is not match register protocol", j.JobName, protocol)
		logger.CtxErrorf(ctx, errMsg)
		return errors.New(errMsg)
	}

	sendParam := &SendParam{}
	for _, opt := range opts {
		opt(sendParam)
	}

	jobEntity := &KafkaJobEntity{
		JobName: j.JobName,
		Content: protocol,
	}

	body, err := jsoniter.Marshal(jobEntity)
	if err != nil {
		return err
	}

	jobCfg := gconfig.DynamicCfg.GetKafkaJobConfig(j.JobName)
	logger.CtxInfof(ctx, "Job:[%s] send to kafka:[%s], topic:[%s]", j.Name, jobCfg.KafkaName, jobCfg.Topic)
	client, err := gclient.GetKafkaClient(ctx, jobCfg.KafkaName)
	if err != nil {
		return err
	}
	if _, _, err := client.SendMessage(ctx, jobCfg.Topic, body, kafka.PartitionKey(sendParam.Key)); err != nil {
		return err
	}
	return nil
}

func (j *KafkaJob[T]) Handle(ctx context.Context, protocol any) error {
	p, ok := protocol.(*T)
	if !ok {
		errMsg := fmt.Sprintf("Job:[%s] handle fail, protocol:[%+v] is not match register protocol", j.JobName, protocol)
		logger.CtxErrorf(ctx, errMsg)
		return errors.New(errMsg)
	}

	return j.Handler(ctx, p)
}
