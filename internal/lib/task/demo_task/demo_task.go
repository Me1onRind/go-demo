package demo_task

import (
	"fmt"

	"github.com/Me1onRind/go-demo/err_code"
	"github.com/Me1onRind/go-demo/infrastructure/client/asynq_client"
	"github.com/Me1onRind/go-demo/infrastructure/ctm_context"
	"github.com/Me1onRind/go-demo/internal/constant/task"
	"github.com/hibiken/asynq"
	json "github.com/json-iterator/go"
)

func SendDemoTask(demoTask *DemoTask) *err_code.Error {
	payload, err := json.Marshal(demoTask)
	if err != nil {
		return err_code.JsonDecodeError.WithErr(err)
	}
	_, err = asynq_client.AsynqClient.Enqueue(
		asynq.NewTask(task.TaskDemo, payload),
	)
	if err != nil {
		return err_code.AsyncTaskSendError.WithErr(err)
	}
	return nil
}

func HandleDemoTask(ctx *ctm_context.Context, task interface{}) *err_code.Error {
	demoTask := (task).(*DemoTask)
	fmt.Printf("demo task:%v\n", demoTask)
	return nil
}
