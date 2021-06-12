package middleware

import (
	"context"
	"time"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GrpcLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		begin := time.Now()
		defer func() {
			commonCtx := common.GetContext(ctx)
			commonCtx.Logger.Info("access request", zap.Reflect("req", req), zap.Reflect("resp", resp),
				zap.String("method", info.FullMethod), zap.Duration("cost", time.Since(begin)),
			)
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func ClientLogger() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var err error
		begin := time.Now()
		defer func() {
			commonCtx := common.GetContext(ctx)
			commonCtx.Logger.Info("grpc request",
				zap.String("method", method), zap.Reflect("req", req), zap.Reflect("reply", reply),
				zap.Duration("cost", time.Since(begin)), zap.Error(err),
			)
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
