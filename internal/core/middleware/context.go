package middleware

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"google.golang.org/grpc"
)

func GrpcContext() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		commonCtx := common.NewContext(ctx)
		ctx = common.StoreContext(ctx, commonCtx)
		return handler(ctx, req)
	}
}
