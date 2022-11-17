package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	ginReqCtxKey = "req-ctx"
)

func mustGetGinExtractContext(c *gin.Context) context.Context {
	value, _ := c.Get(ginReqCtxKey)
	if value == nil {
		setGinExtractContext(c, c.Request.Context())
	}
	return value.(context.Context)
}

func setGinExtractContext(c *gin.Context, ctx context.Context) {
	c.Set(ginReqCtxKey, ctx)
}
