package middleware

import (
	"context"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GrpcLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var resp interface{}
		var err error
		commonCtx := common.GetContext(ctx)
		begin := time.Now()
		defer func() {
			commonCtx.Logger.Info("access request", zap.Reflect("req", req), zap.Reflect("resp", resp),
				zap.String("method", info.FullMethod), zap.Duration("cost", time.Since(begin)),
			)
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}
