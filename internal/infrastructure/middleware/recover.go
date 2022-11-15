package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/Me1onRind/go-demo/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := mustGetGinExtractContext(c)
				stack := debug.Stack()
				logger.CtxErrorf(ctx, "Server internal error:[%s]\n%s", err, stack)
				c.Data(http.StatusInternalServerError, "text/plain", []byte("Server Internal Error"))
			}
		}()
		c.Next()
	}
}
