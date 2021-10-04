package main

import (
	task_constant "github.com/Me1onRind/go-demo/internal/constant/task"
	"github.com/Me1onRind/go-demo/internal/core/initialize"
	"github.com/Me1onRind/go-demo/internal/dao/periodic_task"
	"github.com/Me1onRind/go-demo/internal/lib/localcache"
	"github.com/hibiken/asynq"
)

func Init() {
	funcs := []func() error{
		initialize.InitLogger,
		initialize.InitLocalConfig("./conf"),
		initialize.InitEtcdClient,
		initialize.InitDB,
		initialize.InitLocalCache,
	}

	for _, v := range funcs {
		if err := v(); err != nil {
			panic(err)
		}
	}
}

func main() {
	Init()
	localcache.LoadConfigCache()
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{}, nil)
	periodicTaskDao := periodic_task.NewPeriodicTaskDao()
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
	//scheduler.Register(
}
