package middleware

import (
	"net/http"

	"github.com/Me1onRind/go-demo/internal/infrastructure/goroutine"
	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				goroutine.LogPanicStack(mustGetGinExtractContext(c), err)
				c.Data(http.StatusInternalServerError, "text/plain", []byte("Server Internal Error"))
			}
		}()
		c.Next()
	}
}
