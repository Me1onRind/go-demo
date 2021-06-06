package common

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey struct{}

var (
	logger *zap.Logger
	cKey   = contextKey{}
)

func init() {
	logger, _ = zap.NewDevelopment()
}

type Context struct {
	context.Context

	Logger *zap.Logger
}

func NewContext(ctx context.Context) *Context {
	c := &Context{}
	requestID, _ := uuid.NewRandom()
	c.Logger = logger.With(zap.String("request_id", requestID.String()))
	return c
}

func StoreContext(c context.Context, ctx *Context) context.Context {
	return context.WithValue(c, cKey, ctx)
}

func GetContext(c context.Context) *Context {
	return c.Value(cKey).(*Context)
}
