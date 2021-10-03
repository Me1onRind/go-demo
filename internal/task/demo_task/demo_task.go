package demo_task

import (
	"fmt"

	"github.com/Me1onRind/go-demo/internal/constant/task"
	"github.com/Me1onRind/go-demo/internal/core/client/asynq_client"
	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/Me1onRind/go-demo/internal/err_code"
	"github.com/hibiken/asynq"
	json "github.com/json-iterator/go"
)

func SendDemoTask(demoTask *DemoTask) *common.Error {
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

func HandleDemoTask(ctx *common.Context, task interface{}) *common.Error {
	demoTask := (task).(*DemoTask)
	fmt.Printf("demo task:%v\n", demoTask)
	return nil
}
