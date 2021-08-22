package middleware

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/core/common"
	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
)

func GrpcContext() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = common.NewContext(ctx)
		return handler(ctx, req)
	}
}

func GinContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = common.NewContext(c)
		c.Next()
	}
}
