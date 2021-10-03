package gateway

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/hibiken/asynq"
	json "github.com/json-iterator/go"
)

type TaskHandler func(c *common.Context, taskParam interface{}) (err *common.Error)

func JsonTask(handler TaskHandler, paramType interface{}) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var taskParam interface{}
		if paramType != nil {
			taskParam = parserProtocol(paramType)
		}

		if err := json.Unmarshal(t.Payload(), taskParam); err != nil {
			return err
		}

		newCtx := common.NewContext(ctx)
		err := handler(newCtx, taskParam)
		if err != nil {
			return err.GenError()
		}
		return nil
	}
}
