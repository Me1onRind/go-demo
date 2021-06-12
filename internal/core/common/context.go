package common

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey struct{}

var (
	logger *zap.Logger
	cKey   = contextKey{}
)

func init() {
	//file, _ := os.OpenFile("./log/info.log", os.O_WRONLY|os.O_RDONLY|os.O_APPEND, 0755)
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	logger = zap.New(core, zap.AddCaller())
}

type Context struct {
	context.Context

	Logger *zap.Logger
}

func NewContext(ctx context.Context) *Context {
	c := &Context{}
	c.Context = storeContext(ctx, c)
	requestID, _ := uuid.NewRandom()
	c.Logger = logger.With(zap.String("request_id", requestID.String()))
	return c
}

func storeContext(c context.Context, ctx *Context) context.Context {
	return context.WithValue(c, cKey, ctx)
}

func GetContext(c context.Context) *Context {
	return c.Value(cKey).(*Context)
}
