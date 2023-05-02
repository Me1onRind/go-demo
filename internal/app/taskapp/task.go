package taskapp

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/global/gconfig"
	"github.com/Me1onRind/go-demo/internal/infrastructure/async"
	"github.com/Me1onRind/go-demo/internal/infrastructure/goroutine"
	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/initialize"
	"github.com/Me1onRind/go-demo/internal/model/configmd"
	"github.com/Me1onRind/go-demo/internal/usecase/pinguc"
	"github.com/Shopify/sarama"
)

type TaskServer struct {
	JobManager  *async.JobManager
	PingUsecase *pinguc.PingUsecase
}

func NewTaskServer() *TaskServer {
	return &TaskServer{
		JobManager:  async.NewJobManager(),
		PingUsecase: pinguc.NewPingUsecase(),
	}
}

func (t *TaskServer) RegisterJob() *TaskServer {
	_ = t.JobManager.RegisterJob(async.NewKafkaJob("ping", t.PingUsecase.Ping))
	return t
}

func (t *TaskServer) Init() *TaskServer {
	initFuncs := []initialize.InitHandler{
		initialize.InitOpentracking("go-demo", "0.0.1"),
		initialize.InitFileConfig("./conf.yml"),
		//initialize.InitEtcdClient(),
		//initialize.InitDynamicConfig(),
		//initialize.InitMysqlClient(),
	}
	ctx := context.Background()
	ctx = logger.WithFields(ctx, logger.TraceIdKey, "main-goruntine")

	for _, f := range initFuncs {
		if err := f(ctx); err != nil {
			logger.CtxFatalf(ctx, "initialize fail, err:[%s]", err)
			panic(err)
		}
	}
	return t
}

func (t *TaskServer) Run() error {
	kafkaDuplicateConsumer := map[string]struct{}{}

	for jobName, job := range t.JobManager.GetAllJobs() {
		switch job.BackendType() {
		case async.KafkaBackendJob:
			jobCfg := gconfig.DynamicCfg.KafkaJobConfigs[jobName]
			if _, ok := kafkaDuplicateConsumer[jobName]; ok {
				continue
			}
			kafkaDuplicateConsumer[jobName] = struct{}{}
			_ = consumeKakfa(context.TODO(), jobCfg)
		}
	}
	return nil
}

func consumeKakfa(ctx context.Context, cfg configmd.KafkaJobConfig) error {
	kafkaCfg := sarama.NewConfig()
	client, err := sarama.NewConsumerGroup(cfg.Addr, cfg.ConsumerGroup, kafkaCfg)
	if err != nil {
		logger.CtxErrorf(ctx, "NewConsumerGroup failed, cause:[%s]", err)
		return err
	}
	goroutine.SafeGo(ctx, func() {
		for {
			if err := client.Consume(ctx, []string{cfg.Topic}, nil); err != nil {
				logger.CtxErrorf(ctx, "Kafka consume fail, addr:[%v], topic:[%s], cause:[%s]", cfg.Addr, cfg.Topic, err)
				continue
			}
		}
	})
	return nil
}
