package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/Me1onRind/go-demo/infrastructure/logger"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GrpcRecover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if e := doRecover(ctx); e != nil {
				err = err_code.ServerInternalError.Withf("%v", e).GrpcErr()
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func GinRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := doRecover(c); e != nil {
				c.JSON(200, err_code.ServerInternalError.Withf("%v", e))
			}
		}()
		c.Next()
	}
}

func doRecover(ctx context.Context) interface{} {
	if e := recover(); e != nil {
		stack := string(debug.Stack())
		logger.CtxError(ctx, "server panic", zap.Any("panicErr", e), zap.Stack("stack"))
		logger.CtxError(ctx, stack)
		fmt.Println(stack)
		return e
	}
	return nil
}
