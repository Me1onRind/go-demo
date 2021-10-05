package main

import (
	"context"
	"fmt"

	task_constant "github.com/Me1onRind/go-demo/internal/constant/task"
	"github.com/Me1onRind/go-demo/internal/core/config"
	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/Me1onRind/go-demo/internal/dao/periodic_task_dao"
	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/initialize"
	"github.com/Me1onRind/go-demo/internal/lib/localcache"
	"github.com/hibiken/asynq"
)

func Init(ctx context.Context) {
	funcs := []func() error{
		initialize.InitLogger,
		initialize.InitLocalConfig("./conf"),
		initialize.InitEtcdClient,
		initialize.InitEtcdConfig(ctx, "/go-demo/config.yaml"),
		initialize.InitDB,
		initialize.InitLocalCache,
	}

	for _, v := range funcs {
		if err := v(); err != nil {
			panic(err)
		}
	}
}

func Close() {
	logger.StdoutLogger.Info("Process exit close")

	funcs := []func() error{
		initialize.CloseLogger,
		initialize.CloseEtcdClient,
		initialize.CloseDB,
	}

	for _, v := range funcs {
		if err := v(); err != nil {
			logger.StdoutLogger.Error("Process exit close")
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	Init(ctx)

	defer func() {
		cancel()
		Close()
	}()
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

	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//goroutine.Go(func() {
	//sig := <-sigCh
	//logger.Logger.Info("Catch exit signal", zap.Any("signal", sig))
	//})

	if err := scheduler.Run(); err != nil {
		panic(err)
	}
	//scheduler.Register(
}
