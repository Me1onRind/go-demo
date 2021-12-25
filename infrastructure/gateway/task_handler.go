package gateway

//import (
//"context"

//"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
//"github.com/Me1onRind/go-demo/internal/lib/err_code"
//"github.com/hibiken/asynq"
//json "github.com/json-iterator/go"
//)

//type TaskHandler func(c *ctm_context.Context, taskParam interface{}) (err *err_code.Error)

//func JsonTask(handler TaskHandler, paramType interface{}) asynq.HandlerFunc {
//return func(ctx context.Context, t *asynq.Task) error {
//var taskParam interface{}
//if paramType != nil {
//taskParam = parserProtocol(paramType)
//}

//if err := json.Unmarshal(t.Payload(), taskParam); err != nil {
//return err
//}

//newCtx := ctm_context.NewContext(ctx)
//err := handler(newCtx, taskParam)
//if err != nil {
//return err.GenError()
//}
//return nil
//}
//}
