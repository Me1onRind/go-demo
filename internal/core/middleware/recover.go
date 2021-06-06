package middleware

import (
	"context"
	"runtime/debug"
	//"fmt"
	//"io/ioutil"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func GrpcRecover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			commonCtx := common.GetContext(ctx)
			if err := recover(); err != nil {
				commonCtx.Logger.Error("server panic", zap.Any("panicErr", err))
				commonCtx.Logger.Sugar().Errorf("%s", debug.Stack())
			}
		}()
		return handler(ctx, req)
	}
}
