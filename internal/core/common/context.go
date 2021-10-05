package common

import (
	"context"

	"github.com/Me1onRind/go-demo/internal/core/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type contextKey string

const (
	cKey contextKey = "cmtx"
)

type Context struct {
	context.Context

	Logger *zap.Logger
	Span   opentracing.Span
}

func NewContext(ctx context.Context) *Context {
	c := &Context{}
	c.Context = storeContext(ctx, c)
	c.Logger = logger.Logger
	return c
}

func GetContext(c context.Context) *Context {
	return c.Value(cKey).(*Context)
}

func storeContext(c context.Context, ctx *Context) context.Context {
	switch v := c.(type) {
	case *gin.Context:
		v.Set(string(cKey), ctx)
		return c
	default:
		return context.WithValue(c, cKey, ctx)
	}
}

func ContextLogger(ctx context.Context) *zap.Logger {
	if commonCtx, ok := ctx.(*Context); ok {
		return commonCtx.Logger
	}

	return logger.StdoutLogger
}
