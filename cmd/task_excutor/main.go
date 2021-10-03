package main

import (
	"github.com/Me1onRind/go-demo/internal/constant/task"
	"github.com/Me1onRind/go-demo/internal/core/gateway"
	"github.com/Me1onRind/go-demo/internal/lib/task/demo_task"
	"github.com/hibiken/asynq"
)

//"github.com/Me1onRind/go-demo/internal/core/initialize"

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: "127.0.0.1:6379",
		},
		asynq.Config{},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TaskDemo, gateway.JsonTask(demo_task.HandleDemoTask, &demo_task.DemoTask{}))
	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}
