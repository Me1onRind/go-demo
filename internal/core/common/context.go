package common

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

const (
	cKey = "cmtx"
)

func Init() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeDuration = zapcore.MillisDurationEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	f, err := os.OpenFile("./log/info.log", os.O_WRONLY|os.O_RDONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(f), zapcore.InfoLevel)
	logger = zap.New(core, zap.AddCaller())
}

type Context struct {
	context.Context

	Logger *zap.Logger
	Span   opentracing.Span
}

func NewContext(ctx context.Context) *Context {
	c := &Context{}
	c.Context = storeContext(ctx, c)
	c.Logger = logger
	return c
}

func GetContext(c context.Context) *Context {
	return c.Value(cKey).(*Context)
}

func storeContext(c context.Context, ctx *Context) context.Context {
	switch v := c.(type) {
	case *gin.Context:
		v.Set(cKey, ctx)
		return c
	default:
		return context.WithValue(c, cKey, ctx)
	}
}
