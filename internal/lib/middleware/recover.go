package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/Me1onRind/go-demo/internal/lib/ctm_context"
	"github.com/Me1onRind/go-demo/internal/lib/err_code"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GrpcRecover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			commonCtx := ctx.(*ctm_context.Context)
			if e := recover(); e != nil {
				//commonCtx.Logger.Error("server panic", zap.Any("panicErr", e))
				commonCtx.Logger.Sugar().Errorf("%s", debug.Stack())
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
			commonCtx := ctm_context.GetCtmCtxFromGinCtx(c)
			if e := recover(); e != nil {
				commonCtx.Logger.Error("server panic", zap.Any("panicErr", e), zap.Stack("stack"))
				fmt.Println(string(debug.Stack()))
				//commonCtx.Logger.Sugar().Errorf("%s", debug.Stack())
				c.JSON(200, err_code.ServerInternalError.Withf("%v", e))
			}
		}()
		c.Next()
	}
}
