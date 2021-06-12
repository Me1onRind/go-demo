package middleware

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GrpcRecover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			commonCtx := common.GetContext(ctx)
			if e := recover(); e != nil {
				commonCtx.Logger.Error("server panic", zap.Any("panicErr", e))
				commonCtx.Logger.Sugar().Errorf("%s", debug.Stack())
				err = errors.New(fmt.Sprintf("panic:%v", e))
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}
