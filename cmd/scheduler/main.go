package main

import (
	"context"
	"fmt"

	"github.com/Me1onRind/go-demo/config"
	"github.com/Me1onRind/go-demo/global/logger_singleton"
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

	//ctmCtx := tm_context.NewContext(ctx)

	asynqConfig := config.RemoteConfig.Asynq
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%d", asynqConfig.Redis.Host, asynqConfig.Redis.Port),
	}, &asynq.SchedulerOpts{
		Logger: logger_singleton.Logger.Sugar(),
	})

	//periodicTaskDao := periodic_task_dao.NewPeriodicTaskDao()
	//tasks := periodicTaskDao.ListAllTask()
	//for _, task := range tasks {
	//_, err := scheduler.Register(task.Cronspec, asynq.NewTask(task_constant.TaskDemo, []byte(`{"id":123,"name":"test scheduler"}`)))
	//if err != nil {
	//panic(err)
	//}
	//}

	if err := scheduler.Run(); err != nil {
		panic(err)
	}
}
