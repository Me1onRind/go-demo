package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	ginReqCtxKey = "req-ctx"
)

func ExtractRequestCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		setGinExtractContext(c, c.Request.Context())
	}
}

func mustGetGinExtractContext(c *gin.Context) context.Context {
	value, _ := c.Get(ginReqCtxKey)
	return value.(context.Context)
}

func setGinExtractContext(c *gin.Context, ctx context.Context) {
	c.Set(ginReqCtxKey, ctx)
}
