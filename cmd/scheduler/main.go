package main

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/config"
	"github.com/Me1onRind/go-demo/infrastructure/ctm_context"
	"github.com/Me1onRind/go-demo/internal/dao/periodic_task_dao"
	"github.com/Me1onRind/go-demo/internal/lib/localcache"
	"github.com/hibiken/asynq"
)

func Init(ctx context.Context) {
	funcs := []func() error{}

	for _, v := range funcs {
		if err := v(); err != nil {
			panic(err)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//ctx := context.Background()
	Init(ctx)

	ctmCtx := ctm_context.NewContext(ctx)

	localcache.LoadCache(ctmCtx)

	asynqConfig := config.RemoteConfig.Asynq
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%d", asynqConfig.Redis.Host, asynqConfig.Redis.Port),
	}, &asynq.SchedulerOpts{
		Logger: logger.StdoutLogger.Sugar(),
	})

	periodicTaskDao := periodic_task_dao.NewPeriodicTaskDao()
	tasks := periodicTaskDao.ListAllTask()
	for _, task := range tasks {
		_, err := scheduler.Register(task.Cronspec, asynq.NewTask(task_constant.TaskDemo, []byte(`{"id":123,"name":"test scheduler"}`)))
		if err != nil {
			panic(err)
		}
	}

	if err := scheduler.Run(); err != nil {
		panic(err)
	}
}
